package main

import "github.com/prometheus/prometheus/promql"

type Querier struct {
	engine *promql.Engine
}
