package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/prometheus/prometheus/cmd/promviz/dashboard"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prometheus/common/model"
)

var (
	titleStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

type dashboardModel struct {
	dashboard *dashboard.Dashboard
	results   map[int]model.Value
	viewport  viewport.Model
	ready     bool
}

func initialModel(dashboard *dashboard.Dashboard) dashboardModel {
	return dashboardModel{
		dashboard: dashboard,
		results:   make(map[int]model.Value),
	}
}

func (m dashboardModel) Init() tea.Cmd {
	return nil
}

func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.SetContent(m.renderDashboard())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}

	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m dashboardModel) View() string {
	if !m.ready {
		return "Initializing..."
	}
	return m.viewport.View()
}

func (m dashboardModel) renderDashboard() string {
	var b strings.Builder

	for _, panel := range m.dashboard.Panels {

		switch panel.Type {
		case dashboard.PanelTypeGauge:
			if result, ok := m.results[panel.ID]; ok {
				// Assuming the result is a single scalar value
				value := float64(result.(model.Vector)[0].Value)
				fmt.Fprintf(&b, "%s\n", renderGauge(panel.Title, value, 100, 20))
			} else {
				fmt.Fprintf(&b, "No data available\n")
			}

		case dashboard.PanelTypeTimeseries, dashboard.PanelTypeStat:
			fmt.Fprintf(&b, "%s\n", titleStyle.Render(panel.Title))
			fmt.Fprintf(&b, "%s\n", infoStyle.Render(fmt.Sprintf("Panel ID: %d, Type: %s", panel.ID, panel.Type)))

			if result, ok := m.results[panel.ID]; ok {
				fmt.Fprintf(&b, "%s\n", renderResult(result))
			} else {
				fmt.Fprintf(&b, "No data available\n")
			}
		}

		fmt.Fprintf(&b, "\n")
	}

	return b.String()
}

func renderResult(result model.Value) string {
	var b strings.Builder

	switch result.Type() {
	case model.ValVector:
		vector := result.(model.Vector)
		for _, sample := range vector {
			fmt.Fprintf(&b, "  %v: %v\n", sample.Metric, sample.Value)
		}
	case model.ValScalar:
		scalar := result.(*model.Scalar)
		fmt.Fprintf(&b, "  Value: %v\n", scalar.Value)
	case model.ValMatrix:
		matrix := result.(model.Matrix)
		for _, stream := range matrix {
			fmt.Fprintf(&b, "  %v:\n", stream.Metric)
			for _, value := range stream.Values {
				fmt.Fprintf(&b, "    %v: %v\n", value.Timestamp, value.Value)
			}
		}
	default:
		fmt.Fprintf(&b, "  Unsupported result type: %v\n", result.Type())
	}

	return b.String()
}

func runUI(dashboard *dashboard.Dashboard, results map[int]model.Value) error {
	m := initialModel(dashboard)
	m.results = results

	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	return err
}

func renderGauge(
	title string,
	value float64,
	max float64,
	width int,
) string {
	// Ensure value is between 0 and max
	value = math.Max(0, math.Min(value, max))
	percentage := value / max
	filledWidth := int(math.Round(percentage * float64(width)))

	// Ensure filledWidth is not negative and not greater than width
	filledWidth = int(math.Max(0, math.Min(float64(filledWidth), float64(width))))

	// Create the gauge
	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", width-filledWidth)
	gauge := filled + empty

	// Color styling
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	var coloredGauge string
	if percentage < 0.7 {
		coloredGauge = green.Render(gauge)
	} else if percentage < 0.9 {
		coloredGauge = yellow.Render(gauge)
	} else {
		coloredGauge = red.Render(gauge)
	}

	// Create the value display
	valueDisplay := fmt.Sprintf("%.1f%%", percentage*100)

	// Combine gauge and value
	result := fmt.Sprintf("╭%s╮\n│%s│\n│%s│\n╰%s╯",
		strings.Repeat("─", width+2),
		lipgloss.NewStyle().Width(width+2).Align(lipgloss.Center).Render(title),
		lipgloss.NewStyle().Width(width+2).Align(lipgloss.Center).Render(valueDisplay),
		strings.Repeat("─", width+2))

	return fmt.Sprintf("%s\n%s", result, coloredGauge)
}
