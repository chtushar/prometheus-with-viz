package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
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

	client, err := NewPrometheusClient("http://prometheus.lan")
	if err != nil {
		fmt.Printf("Error creating Prometheus client: %v\n", err)
		os.Exit(1)
	}

	querier := Querier{client: client}

	err = renderDashboard(ctx, querier, dashboard)
	if err != nil {
		panic(err)
	}
}
