package main

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/nshafer/phx"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ranching-farm/kubectl-addon/pkg/logging"
)

var (
	version      = "unknown" // This will be set by the build flag
	clusterId    string
	secret       string
	endpointURL  string
	kubeconfig   string
	outputFormat string
)

var rootCmd = &cobra.Command{
	Use:   "kubectl-ranching_farm",
	Short: "Connect to ranching.farm AI-powered Kubernetes management",
	Long: `This plugin allows you to connect your Kubernetes cluster to the ranching.farm
platform. Once connected, you can interact with your cluster through our
AI-powered interface for simplified management and troubleshooting.`,
	Run: func(cmd *cobra.Command, args []string) {
		// This is called when no subcommands are provided
		cmd.Help()
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to ranching.farm",
	Long: `Connect your Kubernetes cluster to the ranching.farm platform.
This command establishes a secure connection between your local environment
and our AI services, enabling real-time analysis and assistance.`,
	Run: func(cmd *cobra.Command, args []string) {
		if clusterId == "" || secret == "" {
			logging.Log("Error", "Cluster ID and Secret are required. Use --cluster-id and --secret flags.", outputFormat)
			cmd.Usage()
			os.Exit(1)
		}
		runMain()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ranching.farm",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ranching.farm version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(versionCmd)

	connectCmd.Flags().StringVarP(&clusterId, "cluster-id", "i", "", "Cluster ID (required)")
	connectCmd.Flags().StringVarP(&secret, "secret", "s", "", "Cluster Secret (required)")
	connectCmd.Flags().StringVarP(&endpointURL, "endpoint", "e", "wss://ranching.farm/socket/kubernetes/cluster", "WebSocket endpoint URL")
	connectCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "Path to kubeconfig file")
	connectCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (json|yaml|table)")

	connectCmd.MarkFlagRequired("cluster-id")
	connectCmd.MarkFlagRequired("secret")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runMain() {
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = os.Getenv("HOME") + "/.kube/config"
		}
	}

	logging.Log("Starting", "kubectl ranching.farm addon", outputFormat)

	endPoint, err := url.Parse(endpointURL)
	if err != nil {
		logging.Log("Error", fmt.Sprintf("Failed to parse WebSocket URL: %v", err), outputFormat)
		os.Exit(1)
	}

	socket := phx.NewSocket(endPoint)
	err = socket.Connect()
	if err != nil {
		logging.Log("Error", fmt.Sprintf("Failed to connect to socket: %v", err), outputFormat)
		os.Exit(1)
	}
	logging.Log("Success", "Connected to WebSocket", outputFormat)

	channel := socket.Channel(fmt.Sprintf("cluster:%s", clusterId), nil)
	join, err := channel.Join()
	if err != nil {
		logging.Log("Error", fmt.Sprintf("Failed to join channel: %v", err), outputFormat)
		os.Exit(1)
	}

	join.Receive("ok", func(response any) {
		logging.Log("Success", fmt.Sprintf("Joined channel: %s", channel.Topic()), outputFormat)
		sendClusterInfo(channel)
	})

	channel.On("cmd", func(payload any) {
		handleCommand(channel, payload)
	})

	logging.Log("Info", "Main loop started. Waiting for events...", outputFormat)
	select {} // Keep the program running
}

func sendClusterInfo(channel *phx.Channel) {
	logging.Log("Info", "Sending cluster info...", outputFormat)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logging.Log("Error", fmt.Sprintf("Failed to build config from kubeconfig: %v", err), outputFormat)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logging.Log("Error", fmt.Sprintf("Failed to create Kubernetes client: %v", err), outputFormat)
		return
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), v1.ListOptions{})
	if err != nil {
		logging.Log("Error", fmt.Sprintf("Failed to list nodes: %v", err), outputFormat)
		return
	}

	nodeInfo := make([]map[string]string, 0)
	for _, node := range nodes.Items {
		nodeInfo = append(nodeInfo, map[string]string{
			"name":   node.Name,
			"status": string(node.Status.Phase),
		})
	}

	clusterInfo := map[string]interface{}{
		"nodes": nodeInfo,
	}

	logging.Log("Info", fmt.Sprintf("Found %d nodes in the cluster", len(nodes.Items)), outputFormat)
	logging.Log("Node Info", nodeInfo, outputFormat)

	push, err := channel.Push("info", clusterInfo)
	if err != nil {
		logging.Log("Error", fmt.Sprintf("Failed to send cluster info: %v", err), outputFormat)
		return
	}

	push.Receive("ok", func(response any) {
		logging.Log("Success", "Cluster info sent successfully", outputFormat)
	})
}

func handleCommand(channel *phx.Channel, payload any) {
	logging.Log("Info", fmt.Sprintf("Received command payload: %v", payload), outputFormat)
	cmd, ok := payload.(map[string]interface{})
	if !ok {
		logging.Log("Error", "Invalid payload format", outputFormat)
		return
	}

	command, ok := cmd["command"].(string)
	if !ok {
		logging.Log("Error", "Invalid command format", outputFormat)
		return
	}

	args, ok := cmd["arguments"].(string)
	if !ok {
		logging.Log("Error", "Invalid args format", outputFormat)
		return
	}

	uuid, ok := cmd["uuid"].(string)
	if !ok {
		logging.Log("Error", "Invalid uuid format", outputFormat)
		return
	}

	output, err := executeCommand(command, args)
	if err != nil {
		logging.Log("Error", fmt.Sprintf("Error executing command: %v", err), outputFormat)
		output = fmt.Sprintf("Error: %v", err)
	} else {
		logging.Log("Success", "Command executed successfully", outputFormat)
	}

	response := map[string]interface{}{
		"command":   command,
		"arguments": args,
		"output":    output,
		"uuid":      uuid,
	}

	push, err := channel.Push("output", response)
	if err != nil {
		logging.Log("Error", fmt.Sprintf("Failed to send command output: %v", err), outputFormat)
		return
	}

	push.Receive("ok", func(response any) {
		logging.Log("Success", "Command output sent successfully", outputFormat)
	})
}

func executeCommand(command string, params string) (string, error) {
	logging.Log("Info", fmt.Sprintf("Executing command: %s %s", command, params), outputFormat)

	args := strings.Fields(params)
	cmd := exec.Command(command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return fmt.Sprintf("Error: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String()), err
	}

	logging.Log("Success", "Command executed successfully", outputFormat)
	return stdout.String(), nil
}
