package app

import (
	"github.com/NimbleMarkets/ntcharts/canvas/runes"
	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/prometheus/prometheus/cmd/promviz/dashboard"
	"github.com/prometheus/prometheus/cmd/promviz/querier"
)

var defaultStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63")) // purple

var graphLineStyle1 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("4")) // blue

var graphLineStyle2 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("10")) // green

var axisStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("3")) // yellow

var labelStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("6")) // cyan

func RenderTimeSeries(
	title string,
	series []*querier.TimeSeries,
	gridPos dashboard.GridPos,
	viewport *viewport.Model,
	unit string,
) string {
	colWidth := viewport.Width / 24
	colHeight := viewport.Height / 24
	totalWidth := colWidth * gridPos.W
	totalHeight := colHeight * gridPos.H

	if len(series) == 0 {
		return ""
	}

	t1 := timeserieslinechart.New(totalWidth, totalHeight)

	t1.AxisStyle = axisStyle
	t1.LabelStyle = labelStyle
	t1.XLabelFormatter = timeserieslinechart.HourTimeLabelFormatter()
	t1.UpdateHandler = timeserieslinechart.SecondUpdateHandler(1)
	t1.SetStyle(graphLineStyle1)
	t1.SetLineStyle(runes.ThinLineStyle)

	// Process the data
	var minValue, maxValue float64
	var minTime, maxTime int64

	ts := series[0]

	var values []timeserieslinechart.TimePoint

	for _, v := range ts.Matrix[0].Values {
		value := float64(v.Value)

		if unit == "percentunit" {
			value = value * 100
		}

		timestamp := v.Timestamp

		values = append(values, timeserieslinechart.TimePoint{
			Time:  v.Timestamp.Time(),
			Value: value,
		})

		t1.Push(timeserieslinechart.TimePoint{
			Time:  v.Timestamp.Time(),
			Value: value,
		})

		// Update min/max values
		if value < minValue || minValue == 0 {
			minValue = value
		}

		if value > maxValue {
			maxValue = value
		}

		if timestamp.Unix() < minTime || minTime == 0 {
			minTime = timestamp.Unix()
		}
		if timestamp.Unix() > maxTime {
			maxTime = timestamp.Unix()
		}
	}

	// Set Y and X ranges
	padding := (maxValue - minValue) * 0.1
	t1.SetYRange(minValue-padding, maxValue+padding)
	t1.SetViewYRange(minValue-padding, maxValue+padding)
	t1.SetXRange(float64(minTime), float64(maxTime))

	t1.DrawAll()

	// Combine title and chart
	result := defaultStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			t1.View(),
		),
	)

	return result
}
