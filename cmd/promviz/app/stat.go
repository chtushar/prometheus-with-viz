package app

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/prometheus/prometheus/cmd/promviz/dashboard"
)

func RenderStat(
	title string,
	value string,
	gridPos dashboard.GridPos,
	viewport *viewport.Model,
) string {

	// padding := 2
	colWidth := viewport.Width / 24
	totalWidth := colWidth * gridPos.W

	// Combine gauge and value
	result := lipgloss.NewStyle().Width(totalWidth).Align(lipgloss.Center).Render(value)

	return result
}
