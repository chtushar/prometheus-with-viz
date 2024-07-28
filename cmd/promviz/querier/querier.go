package querier

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/cmd/promviz/dashboard"
)

type Querier struct {
	client *PrometheusClient
}

func New(client *PrometheusClient) *Querier {
	return &Querier{
		client: client,
	}
}

func (q *Querier) FetchGaugePanelData(
	ctx context.Context,
	panel *dashboard.Panel,
	variables map[string]string,
) (model.Vector, error) {
	if len(panel.Targets) < 1 {
		return nil, fmt.Errorf("fetchGaugePanelData: can't fetch data without targets")
	}

	target := panel.Targets[0]
	query := replaceVariablesWithValues(target.Expr, variables)

	result, err := q.client.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("fetchGaugePanelData: failed to query: %w", err)
	}

	vector, ok := result.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("failed to cast %T to vector type", result)
	}

	return vector, nil
}

func (q *Querier) FetchStatPanelData(
	ctx context.Context,
	panel *dashboard.Panel,
	variables map[string]string,
) (model.Vector, error) {
	if len(panel.Targets) < 1 {
		return nil, fmt.Errorf("fetchStatPanelData: can't fetch data without targets")
	}

	target := panel.Targets[0]
	query := replaceVariablesWithValues(target.Expr, variables)

	result, err := q.client.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("fetchStatPanelData: failed to query: %w", err)
	}

	vector, ok := result.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("failed to cast %T to vector type", result)
	}

	return vector, nil
}

type TimeSeries struct {
	LegendFormat string
	Matrix       model.Matrix
}

func (q *Querier) FetchTimeSeriesPanelData(
	ctx context.Context,
	panel *dashboard.Panel,
	start, end time.Time,
	variables map[string]string,
) ([]*TimeSeries, error) {
	results := make([]*TimeSeries, 0)

	for _, target := range panel.Targets {
		if target.Hide {
			continue
		}

		query := replaceVariablesWithValues(target.Expr, variables)

		result, err := q.client.QueryRange(ctx, query, v1.Range{
			Start: start,
			End:   end,
			Step:  time.Duration(target.Step) * time.Second,
		})
		if err != nil {
			return nil, fmt.Errorf("fetchTimeSeriesPanelData: failed to fetch: %w", err)
		}

		matrix, ok := result.(model.Matrix)
		if !ok {
			return nil, fmt.Errorf("failed to cast %T to vector type", result)
		}

		results = append(results, &TimeSeries{
			LegendFormat: target.LegendFormat,
			Matrix:       matrix,
		})
	}

	return results, nil
}

func replaceVariablesWithValues(
	query string,
	variables map[string]string,
) string {
	for variable, value := range variables {
		query = strings.ReplaceAll(query, variable, value)
	}

	return query
}
