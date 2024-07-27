package widgets

// type LineChart struct {
// 	Query string
// 	Range time.Duration
// 	Step time.Duration
// 	Values []float64
// }

// func NewLineChart(query string, r time.Duration, step time.Duration) *LineChart {
// 	return &LineChart{
// 		Query: query,
// 		Range: r,
// 		Step: step,
// 	}
// }

// func (l *LineChart) AddWidget() error {
// 	lc, err := linechart.New(
// 		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
// 		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
// 		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// }