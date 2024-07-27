package main

import (
	"encoding/json"
	"io"
)

func ParseDashboardJson(r io.Reader) (*Dashboard, error) {
	d := Dashboard{}

	err := json.NewDecoder(r).Decode(&d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}
