package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func ParseDashboardJson(r io.Reader) (*Dashboard, error) {
	d := Dashboard{}

	err := json.NewDecoder(r).Decode(&d)
	if err != nil {
		return nil, err
	}

	fmt.Println(d)

	return &d, nil
}
