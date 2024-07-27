package dashboard

import (
	"fmt"

	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/sparkline"
)

type Dashboard struct {
	Terminal *tcell.Terminal
	Container *container.Container
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

func New() *Dashboard {
	t, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		fmt.Printf("tcell.New => %v", err)
	}
	
	c, err := container.New(t, container.ID("root"))
	if err != nil {
		fmt.Printf("container.New => %v", err)
	}
	
	return &Dashboard{
		Terminal: t,
		Container: c,
	}
}
