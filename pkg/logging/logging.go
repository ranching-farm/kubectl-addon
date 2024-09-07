package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

func Log(message string, data interface{}, format string) {
	fmt.Println(message)

	switch format {
	case "json":
		jsonOutput, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("Error marshalling JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonOutput))

	case "yaml":
		yamlOutput, err := yaml.Marshal(data)
		if err != nil {
			fmt.Printf("Error marshalling YAML: %v\n", err)
			return
		}
		fmt.Println(string(yamlOutput))

	case "table":
		table := tablewriter.NewWriter(os.Stdout)
		v := reflect.ValueOf(data)
		if v.Kind() == reflect.Map {
			table.SetHeader([]string{"Key", "Value"})
			for _, key := range v.MapKeys() {
				table.Append([]string{fmt.Sprintf("%v", key.Interface()), fmt.Sprintf("%v", v.MapIndex(key).Interface())})
			}
		} else if v.Kind() == reflect.Slice {
			if v.Len() > 0 {
				firstElem := v.Index(0)
				if firstElem.Kind() == reflect.Map {
					headers := make([]string, 0)
					for _, key := range firstElem.MapKeys() {
						headers = append(headers, fmt.Sprintf("%v", key.Interface()))
					}
					table.SetHeader(headers)

					for i := 0; i < v.Len(); i++ {
						elem := v.Index(i)
						row := make([]string, 0)
						for _, header := range headers {
							row = append(row, fmt.Sprintf("%v", elem.MapIndex(reflect.ValueOf(header)).Interface()))
						}
						table.Append(row)
					}
				}
			}
		}
		table.Render()

	default:
		fmt.Printf("%+v\n", data)
	}
}
