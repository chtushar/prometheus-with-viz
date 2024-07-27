package main

import (
	"context"
	"fmt"
	"github.com/prometheus/common/model"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if len(os.Args) < 2 {
		fmt.Println("incorrect call, need json")
		os.Exit(1)
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
		return
	}

	defer file.Close()

	dashboard, err := ParseDashboardJson(file)
	if err != nil {
		panic(err)
	}

	client, err := NewPrometheusClient("http://prometheus.lan")
	if err != nil {
		fmt.Printf("Error creating Prometheus client: %v\n", err)
		os.Exit(1)
	}

	node := "anakin-rpi.lan:9100"

	for _, panel := range dashboard.Panels {
		if panel.Type == PanelTypeTimeseries || panel.Type == PanelTypeGauge || panel.Type == PanelTypeStat {
			for _, target := range panel.Targets {
				if target.Expr != "" {
					fmt.Printf("Querying for panel '%s' (ID: %d)\n", panel.Title, panel.ID)

					query := strings.ReplaceAll(target.Expr, "$node", node)
					query = strings.ReplaceAll(query, "$__rate_interval", "1m")
					query = strings.ReplaceAll(query, "$job", "node-exporter")

					result, err := client.Query(ctx, query, time.Now())
					if err != nil {
						fmt.Printf("Error querying Prometheus: %v\n", err)
						continue
					}

					// Process and display the result
					displayResult(result)
					fmt.Println()
					fmt.Println()
				}
			}
		}
	}

}

func displayResult(result model.Value) {
	switch result.Type() {
	case model.ValVector:
		vector := result.(model.Vector)
		for _, sample := range vector {
			fmt.Printf("Metric: %v, Value: %v, Timestamp: %v\n", sample.Metric, sample.Value, sample.Timestamp)
		}
	case model.ValScalar:
		scalar := result.(*model.Scalar)
		fmt.Printf("Scalar Value: %v, Timestamp: %v\n", scalar.Value, scalar.Timestamp)
	case model.ValMatrix:
		matrix := result.(model.Matrix)
		for _, stream := range matrix {
			fmt.Printf("Metric: %v\n", stream.Metric)
			for _, value := range stream.Values {
				fmt.Printf("  Value: %v, Timestamp: %v\n", value.Value, value.Timestamp)
			}
		}
	default:
		fmt.Printf("Unsupported result type: %v\n", result.Type())
	}
}
