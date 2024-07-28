package app

import (
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/cmd/promviz/dashboard"
	"github.com/prometheus/prometheus/cmd/promviz/querier"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type App struct {
	Program *tea.Program
}

type DashboardUpdate struct{}

type Model struct {
	ctx               context.Context
	viewport          viewport.Model
	ready             bool
	dashboard         *dashboard.Dashboard
	results           map[int]model.Value
	timeseriesResults map[int][]*querier.TimeSeries
	querier           *querier.Querier
}

func (m Model) fetchDataFromPrometheus(t time.Time) error {
	variableValues := map[string]string{
		"$node":            "anakin-rpi.lan:9100",
		"$job":             "node-exporter",
		"$__rate_interval": "5m",
	}

	now := t
	start := now.Add(-24 * time.Hour)
	end := now

	for _, p := range m.dashboard.Panels {
		switch p.Type {
		case dashboard.PanelTypeGauge:
			data, err := m.querier.FetchGaugePanelData(m.ctx, p, variableValues)
			if err != nil {
				return fmt.Errorf("failed to load panel %d", p.ID)
			}

			m.results[p.ID] = data

		case dashboard.PanelTypeStat:
			data, err := m.querier.FetchStatPanelData(m.ctx, p, variableValues)
			if err != nil {
				return fmt.Errorf("failed to load panel %d", p.ID)
			}

			m.results[p.ID] = data

		case dashboard.PanelTypeTimeseries:
			data, err := m.querier.FetchTimeSeriesPanelData(m.ctx, p, start, end, variableValues)
			if err != nil {
				return fmt.Errorf("failed to load panel %d", p.ID)
			}

			m.timeseriesResults[p.ID] = data
		}
	}

	return nil
}

func (m Model) refreshDashboardData(t time.Time) tea.Msg {
	err := m.fetchDataFromPrometheus(t)
	if err != nil {
		panic(err)
	}

	return DashboardUpdate{}
}

func (m Model) checkServer(init bool) tea.Cmd {
	if init {
		return func() tea.Msg {
			return m.refreshDashboardData(time.Now())
		}
	}

	parsedDuration, err := time.ParseDuration(m.dashboard.Refresh)
	if err != nil {
		panic(err)
	}

	return tea.Every(parsedDuration, m.refreshDashboardData)
}

func (m Model) headerView() string {
	title := titleStyle.Render("promviz")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m Model) renderPanel(panel *dashboard.Panel) string {
	switch panel.Type {
	case dashboard.PanelTypeRow:
		return lipgloss.NewStyle().Bold(true).Render(panel.Title)
	case dashboard.PanelTypeGauge:
		value := m.getValueForPanel(panel.ID)
		return RenderGauge(panel.Title, value, 100, panel.GridPos, &m.viewport)
	case dashboard.PanelTypeStat:
		value := m.getValueForPanel(panel.ID)
		return RenderStat(panel.Title, formatter(panel.FieldConfig.Defaults.Unit, value), panel.GridPos, &m.viewport)
	case dashboard.PanelTypeTimeseries:
		var series []*querier.TimeSeries
		ts, ok := m.timeseriesResults[panel.ID]
		if ok {
			series = ts
		}
		return RenderTimeSeries(panel.Title, series, panel.GridPos, &m.viewport, panel.FieldConfig.Defaults.Unit)
	default:
		return fmt.Sprintf("Unsupported panel type: %s", panel.Type)
	}
}

func (m Model) getValueForPanel(panelID int) float64 {
	if result, ok := m.results[panelID]; ok {
		vec, ok := result.(model.Vector)
		if ok && len(vec) > 0 {
			return float64(vec[0].Value)
		}
	}
	return 0
}

func (m Model) stylePanelBox(title string, content string, width int) string {
	panelWidth := (m.viewport.Width * width) / 24 // Assuming a 24-column grid

	box := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(panelWidth)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Bold(true)

	return box.Render(
		titleStyle.Render(title)+"\n"+
			content,
	) + "\n"
}

func (m Model) Init() tea.Cmd {
	return m.checkServer(true)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case DashboardUpdate:
		cmds = append(cmds, m.checkServer(false))
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		case "pgup":
			m.viewport.HalfViewUp()
		case "pgdown":
			m.viewport.HalfViewDown()
		case "home":
			m.viewport.GotoTop()
		case "end":
			m.viewport.GotoBottom()
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.HighPerformanceRendering = false
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			m.viewport.LineUp(3)
		case tea.MouseWheelDown:
			m.viewport.LineDown(3)
		}
	}

	vp, cmd := m.viewport.Update(msg)

	m.viewport = vp
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) renderContent() string {
	var content strings.Builder
	var currentRow int
	var rowContent strings.Builder

	type PanelWithString struct {
		Panel   *dashboard.Panel
		Content string
	}

	rowGroups := make(map[int][]PanelWithString)
	maxY := 0

	for _, p := range m.dashboard.Panels {
		if p.GridPos.Y > currentRow {
			content.WriteString(rowContent.String() + "\n")
			rowContent.Reset()
			currentRow = p.GridPos.Y
		}

		y := p.GridPos.Y
		if y > maxY {
			maxY = y
		}

		panelContent := m.renderPanel(p)
		if p.Type != dashboard.PanelTypeRow {
			panelContent = m.stylePanelBox(p.Title, panelContent, p.GridPos.W)
		}
		rowGroups[y] = append(rowGroups[y], PanelWithString{Panel: p, Content: panelContent})
		rowContent.WriteString(panelContent)
	}

	for y := range rowGroups {
		sort.Slice(rowGroups[y], func(i, j int) bool {
			return rowGroups[y][i].Panel.GridPos.X < rowGroups[y][j].Panel.GridPos.X
		})
	}

	rows := make([]string, maxY+1)

	for y := 0; y <= maxY; y++ {
		rowPanels := rowGroups[y]
		if len(rowPanels) == 0 {
			continue
		}
		var horizontalPanels []string
		for _, p := range rowPanels {
			horizontalPanels = append(horizontalPanels, p.Content)
		}

		rows[y] = lipgloss.JoinHorizontal(lipgloss.Top, horizontalPanels...)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	content := m.renderContent()

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), content, m.footerView())
}

func New(ctx context.Context, d *dashboard.Dashboard, q *querier.Querier) *App {
	d.Panels = dashboard.SortPanelsByPosition(d.Panels)

	p := tea.NewProgram(
		&Model{
			ctx:               ctx,
			dashboard:         d,
			querier:           q,
			results:           make(map[int]model.Value),
			timeseriesResults: make(map[int][]*querier.TimeSeries),
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithMouseAllMotion(),
	)
	return &App{
		Program: p,
	}
}

func formatter(unit string, value float64) string {
	if unit == "s" {
		t := time.Now().Add(-1 * time.Duration(value) * time.Second)
		return humanize.Time(t)
	}

	if unit == "bytes" {
		return humanize.IBytes(uint64(value))
	}

	return fmt.Sprintf("%.2f", value)
}
