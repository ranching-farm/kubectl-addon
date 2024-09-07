<img src="https://ranching.farm/images/logo.svg" alt="ranching.farm logo"/>

# [ranching.farm](https://ranching.farm) kubectl Addon

[![Build Status](https://img.shields.io/github/workflow/status/ranching-farm/kubectl-addon/CI?style=flat-square)](https://github.com/ranching-farm/kubectl-addon/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ranching-farm/kubectl-addon?style=flat-square)](https://goreportcard.com/report/github.com/ranching-farm/kubectl-addon)
[![License](https://img.shields.io/github/license/ranching-farm/kubectl-addon?style=flat-square)](https://github.com/ranching-farm/kubectl-addon/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/ranching-farm/kubectl-addon?style=flat-square)](https://github.com/ranching-farm/kubectl-addon/releases)

## What is ranching.farm?

ranching.farm is a Kubernetes management platform that leverages AI to simplify cluster operations and troubleshooting. It's designed for DevOps engineers and SREs who want to streamline their Kubernetes workflows and get intelligent assistance for common tasks. It can also be used as a training platform and e-learning resource for beginners and intermediates alike.

We offer two ways of connecting your k8s deployment to our AI:

1) Use our kubernetes-agent to deploy it straight with one-line into your k8s cluster.
2) Use our kubectl-addon (**this repository**) to run it locally on your machine using your kubectl configuration.

This kubectl addon allows you to connect your Kubernetes cluster to the [ranching.farm](https://ranching.farm) platform. Once connected, you can interact with your cluster through our AI-powered interface.

The addon provides a secure connection between your local environment and our AI services, enabling real-time analysis and assistance while maintaining your cluster's security. This approach complements our kubernetes-agent option, offering flexibility in how you integrate with [ranching.farm](https://ranching.farm).

## 🚀 Quick Start

```bash
kubectl krew install ranching.farm
kubectl ranching.farm connect --cluster-id <cluster-id> --cluster-secret <cluster-secret>
```

## 🌟 Features

- **Easy Setup**: Get your cluster connected in seconds
- **Lightweight**: Minimal resource footprint
- **Cross-Platform**: Works on Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64)
- **AI-Powered Assistance**: Access to ranching.farm's AI Kubernetes management features
- **Inspection Friendly**: Audit all commands being executed in your cluster, so you don't have to trust us

## 🛠 Installation

### Using Krew

1. [Install Krew](https://krew.sigs.k8s.io/docs/user-guide/setup/install/)
2. Run:
   ```bash
   kubectl krew install ranching.farm
   ```

### Manual Installation

Download the latest release for your platform from the [Releases page](https://github.com/ranching-farm/kubectl-addon/releases).

## 📄 License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♀️ Support

- 🐛 [Issue Tracker](https://github.com/ranching-farm/kubectl-addon/issues)

---

<p align="center">Made with ❤️ by the ranching.farm team</p>