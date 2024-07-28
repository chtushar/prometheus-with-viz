package widgets

type PanelType string

const (
	PanelTypeRow        PanelType = "row"
	PanelTypeGauge      PanelType = "gauge"
	PanelTypeStat       PanelType = "stat"
	PanelTypeTimeseries PanelType = "timeseries"
	PanelTypeBargauge   PanelType = "bargauge"
	PanelTypeUnknown    PanelType = "unknown"
)
