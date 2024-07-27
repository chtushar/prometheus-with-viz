package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/cmd/promviz/app"
	"github.com/prometheus/prometheus/cmd/promviz/widgets"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	a := app.New()
	defer a.Terminal.Clear()
	defer a.Terminal.Close()

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

	client, err := NewPrometheusClient("http://192.168.0.105:9090")
	if err != nil {
		fmt.Printf("Error creating Prometheus client: %v\n", err)
		os.Exit(1)
	}

	node := "192.168.0.105:9100"

	for _, panel := range dashboard.Panels {
		if panel.Type == PanelTypeGauge {
			for _, target := range panel.Targets {
				if target.Expr != "" {
					fmt.Printf("Querying for panel '%s' (ID: %d)\n", panel.Title, panel.ID)

					now := time.Now()

					query := strings.ReplaceAll(target.Expr, "$node", node)
					query = strings.ReplaceAll(query, "$__rate_interval", "1m")
					query = strings.ReplaceAll(query, "$job", "node-exporter")

					result, err := client.Query(ctx, query, now)
					if err != nil {
						fmt.Printf("Error querying Prometheus: %v\n", err)
						continue
					}

					vector, ok := result.(model.Vector)
					if !ok {
						panic("result not of vector in gauge")
					}

					sample := vector[0]

					fmt.Println(sample.Value)
					a.AddWidget(app.Widget{
						Type:  string(PanelTypeGauge),
						Gauge: widgets.NewGauge(query, float32(sample.Timestamp.Unix()*1000), float32(sample.Value*100)),
					})

					// Process and display the result
					// displayResult(result)
					// fmt.Println()
					// fmt.Println()
				}
			}
		}
	}

	builder := grid.New()

	// builder.Add(
	// 	grid.ColWidthPerc(70,
	// 		grid.Widget(line,
	// 			container.Border(linestyle.Light),
	// 			container.BorderTitle("Press Esc to quit"),
	// 		),
	// 	))

	for _, widget := range a.Widgets {
		switch widget.Type {
		case string(PanelTypeGauge):
			builder.Add(
				grid.RowHeightPerc(1,
					grid.Widget(widget.Gauge.G),
				),
			)
		}
	}

	gridOpts, err := builder.Build()

	if err != nil {
		panic(err)
	}

	if err := a.Container.Update("root", gridOpts...); err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key.String() == "KeyEsc" {
			file.Close()
			a.Terminal.Flush()
			a.Terminal.Close()
			cancel()
		}
	}

	if err := termdash.Run(ctx, a.Terminal, a.Container, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(1 * time.Second)); err != nil {
		fmt.Printf("termdash.Run => %v", err)
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
