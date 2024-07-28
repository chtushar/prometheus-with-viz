package main

import (
	"encoding/json"
	"io"

	"github.com/prometheus/prometheus/cmd/promviz/dashboard"
)

func ParseDashboardJson(r io.Reader) (*dashboard.Dashboard, error) {
	d := dashboard.Dashboard{}

	err := json.NewDecoder(r).Decode(&d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}
