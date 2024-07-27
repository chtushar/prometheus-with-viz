package app

import (
	"fmt"

	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/sparkline"
	"github.com/prometheus/prometheus/cmd/promviz/widgets"
)

type Widget struct {
	Type string `json:"type,omitempty"`
	Gauge *widgets.Gauge `json:"gauge,omitempty"`
}

type App struct {
	Terminal *tcell.Terminal
	Container *container.Container
	Widgets []Widget
}

type GridPos struct {
	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
	W int `json:"w,omitempty"`
}


type Graph struct {
	Query string `json:"query,omitempty"`
	Pos GridPos `json:"pos,omitempty"`
	Line *sparkline.SparkLine `json:"line,omitempty"`
}

func New() *App {
	t, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		fmt.Printf("tcell.New => %v", err)
	}
	
	c, err := container.New(t, container.ID("root"))
	if err != nil {
		fmt.Printf("container.New => %v", err)
	}
	
	return &App{
		Terminal: t,
		Container: c,
	}
}


func (a *App) AddWidget(w Widget) {
	switch w.Type {
	case "gauge":
		a.Widgets = append(a.Widgets, w)
	}
}
