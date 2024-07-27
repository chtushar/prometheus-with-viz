package widgets

import (
	"fmt"

	"github.com/mum4k/termdash/widgets/gauge"
)

type Gauge struct {
	Query     string
	Timestamp float32
	Value     float32
	G         *gauge.Gauge
}

func NewGauge(query string, timestamp float32, value float32) *Gauge {
	gauge, err := gauge.New()

	fmt.Println("guage value:", value)

	if err != nil {
		panic(err)
	}

	return &Gauge{
		Query:     query,
		Timestamp: timestamp,
		Value:     value,
		G:         gauge,
	}
}
