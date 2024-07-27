package dashboard

import "math"


type GridPos struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}

type PanelSize int

const (
	SingleColumn PanelSize = iota
	DoubleColumn
	TripleColumn
)

func DeterminePanelSize(panel GridPos, totalColumns int) PanelSize {
	// Calculate the percentage of the total width that this panel occupies
	panelWidthPercentage := float64(panel.W) / float64(totalColumns)

	switch {
	case panelWidthPercentage <= 0.33:
		return SingleColumn
	case panelWidthPercentage <= 0.66:
		return DoubleColumn
	default:
		return TripleColumn
	}
}

func CalculatePanelWidth(size PanelSize, availableWidth int) int {
	switch size {
	case SingleColumn:
		return availableWidth / 3
	case DoubleColumn:
		return (availableWidth * 2) / 3
	case TripleColumn:
		return availableWidth
	default:
		return availableWidth // Default to full width if unknown size
	}
}

func AdaptPanelToWidth(panel GridPos, availableWidth int) int {
	totalColumns := 24 // Grafana typically uses a 24-column grid
	columnWidth := availableWidth / totalColumns

	size := DeterminePanelSize(panel, totalColumns)
	width := CalculatePanelWidth(size, availableWidth)

	// Adjust width to nearest column multiple
	columns := int(math.Round(float64(width) / float64(columnWidth)))
	adaptedWidth := columns * columnWidth

	return adaptedWidth
}

func GetNewGridPos(panel GridPos, availableWidth int) GridPos {
	adaptedWidth := AdaptPanelToWidth(panel, availableWidth)
	adaptedHeight := int(float64(adaptedWidth) * (float64(panel.H) / float64(panel.W)))

	return GridPos{
		H: adaptedHeight,
		W: adaptedWidth,
		X: panel.X,
		Y: panel.Y,
	}
}