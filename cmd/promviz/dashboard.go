package main

import (
	"context"
	"fmt"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/cmd/promviz/dashboard"
	"time"
)

func renderDashboard(
	ctx context.Context,
	querier Querier,
	dashboardJson *dashboard.Dashboard,
) error {
	variableValues := map[string]string{
		"$node":            "anakin-rpi.lan:9100",
		"$job":             "node-exporter",
		"$__rate_interval": dashboardJson.Refresh,
	}

	now := time.Now()
	start := now.Add(-24 * time.Hour)
	end := now

	results := make(map[int]model.Value)

	for _, panel := range dashboardJson.Panels {
		switch panel.Type {
		case dashboard.PanelTypeGauge:
			data, err := querier.fetchGaugePanelData(ctx, panel, variableValues)
			if err != nil {
				return fmt.Errorf("failed to load panel %d", panel.ID)
			}

			results[panel.ID] = data

		case dashboard.PanelTypeStat:
			data, err := querier.fetchStatPanelData(ctx, panel, variableValues)
			if err != nil {
				return fmt.Errorf("failed to load panel %d", panel.ID)
			}

			results[panel.ID] = data

		case dashboard.PanelTypeTimeseries:
			data, err := querier.fetchTimeSeriesPanelData(ctx, panel, start, end, variableValues)
			if err != nil {
				return fmt.Errorf("failed to load panel %d", panel.ID)
			}

			// TODO: rendering timeseries??

			results[panel.ID] = *data[0]

			fmt.Println(data)

		default:
			fmt.Printf("unsupported panel type: %s\n\n", panel.Type)
		}
	}

	if err := runUI(dashboardJson, results); err != nil {
		return fmt.Errorf("failed to render dashboard: %w", err)
	}

	return nil
}
