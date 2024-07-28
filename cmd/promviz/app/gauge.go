package app

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/prometheus/prometheus/cmd/promviz/dashboard"
)

func RenderGauge(
	title string,
	value float64,
	max float64,
	gridPos dashboard.GridPos,
	viewport *viewport.Model,
) string {

	// padding := 2
	colWidth := viewport.Width / 24
	totalWidth := colWidth * gridPos.W


	// Ensure value is between 0 and max
	value = math.Max(0, math.Min(value, max))
	percentage := value / max
	filledWidth := int(math.Round(percentage * float64(totalWidth)))

	// Ensure filledWidth is not negative and not greater than width
	filledWidth = int(math.Max(0, math.Min(float64(filledWidth), float64(totalWidth))))

	// Create the gauge
	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", totalWidth - filledWidth)
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
	result := fmt.Sprintf("╭%s╮\n│%s│\n│%s│\n│%s│\n╰%s╯\n",
		strings.Repeat("─", totalWidth),
		lipgloss.NewStyle().Width(totalWidth).Align(lipgloss.Center).Render(title),
		lipgloss.NewStyle().Width(totalWidth).Align(lipgloss.Center).Render(valueDisplay),
		lipgloss.NewStyle().Width(totalWidth).Align(lipgloss.Center).Render(coloredGauge),
		strings.Repeat("─", totalWidth))

	return result
}