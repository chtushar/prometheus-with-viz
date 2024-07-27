package main

import (
	"context"
	"fmt"
	"os"

	"github.com/prometheus/prometheus/cmd/promviz/app"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()


	if len(os.Args) < 2 {
		fmt.Println("incorrect call, need json")
		os.Exit(1)
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	dashboard, err := ParseDashboardJson(file)
	if err != nil {
		panic(err)
	}

	// client, err := NewPrometheusClient("http://192.168.0.105:9090")
	// if err != nil {
	// 	fmt.Printf("Error creating Prometheus client: %v\n", err)
	// 	os.Exit(1)
	// }

	// node := "192.168.0.105:9100"

	// for _, panel := range dashboard.Panels {
	// 	if panel.Type == PanelTypeGauge {
	// 		for _, target := range panel.Targets {
	// 			if target.Expr != "" {
	// 				fmt.Printf("Querying for panel '%s' (ID: %d)\n", panel.Title, panel.ID)

	// 				now := time.Now()

	// 				query := strings.ReplaceAll(target.Expr, "$node", node)
	// 				query = strings.ReplaceAll(query, "$__rate_interval", "1m")
	// 				query = strings.ReplaceAll(query, "$job", "node-exporter")

	// 				result, err := client.Query(ctx, query, now)
	// 				if err != nil {
	// 					fmt.Printf("Error querying Prometheus: %v\n", err)
	// 					continue
	// 				}

	// 				// vector, ok := result.(model.Vector)
	// 				// if !ok {
	// 				// 	panic("result not of vector in gauge")
	// 				// }

	// 				// sample := vector[0]

	// 				// Process and display the result
	// 				displayResult(result)
	// 				fmt.Println()
	// 				fmt.Println()
	// 			}
	// 		}
	// 	}
	// }

	a := app.New(dashboard)

	if _, err := a.Program.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}

// func displayResult(result model.Value) {
// 	switch result.Type() {
// 	case model.ValVector:
// 		vector := result.(model.Vector)
// 		for _, sample := range vector {
// 			fmt.Printf("Metric: %v, Value: %v, Timestamp: %v\n", sample.Metric, sample.Value, sample.Timestamp)
// 		}
// 	case model.ValScalar:
// 		scalar := result.(*model.Scalar)
// 		fmt.Printf("Scalar Value: %v, Timestamp: %v\n", scalar.Value, scalar.Timestamp)
// 	case model.ValMatrix:
// 		matrix := result.(model.Matrix)
// 		for _, stream := range matrix {
// 			fmt.Printf("Metric: %v\n", stream.Metric)
// 			for _, value := range stream.Values {
// 				fmt.Printf("  Value: %v, Timestamp: %v\n", value.Value, value.Timestamp)
// 			}
// 		}
// 	default:
// 		fmt.Printf("Unsupported result type: %v\n", result.Type())
// 	}
// }
