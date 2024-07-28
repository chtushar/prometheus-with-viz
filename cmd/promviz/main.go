package main

import (
	"context"
	"fmt"
	"os"

	"github.com/prometheus/prometheus/cmd/promviz/app"
	"github.com/prometheus/prometheus/cmd/promviz/querier"
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


	client, err := querier.NewPrometheusClient("http://192.168.0.105:9090/")
	if err != nil {
		fmt.Printf("Error creating Prometheus client: %v\n", err)
		os.Exit(1)
	}

	q := querier.New(client)

	// err = renderDashboard(ctx, q, dashboard)
	// if err != nil {
	// 	panic(err)
	// }

	a := app.New(ctx, dashboard, q)

	if _, err := a.Program.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}
