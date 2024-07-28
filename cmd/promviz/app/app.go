package app

import (
	"context"
	"fmt"
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

type TickMsg time.Time

type Model struct {
	ctx               context.Context
	viewport          viewport.Model
	ready             bool
	dashboard         *dashboard.Dashboard
	results           map[int]model.Value
	timeseriesResults map[int][]*querier.TimeSeries
	querier           *querier.Querier
}

func (m Model) checkServer() tea.Cmd {
	parsedDuration, err := time.ParseDuration(m.dashboard.Refresh)

	if err != nil {
		panic(err)
	}

	variableValues := map[string]string{
		"$node":            "192.168.0.105:9100",
		"$job":             "node-exporter",
		"$__rate_interval": "5m",
	}

	now := time.Now()
	start := now.Add(-24 * time.Hour)
	end := now

	return tea.Every(parsedDuration, func(t time.Time) tea.Msg {
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
		return TickMsg(t)
	})
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
		return "\n"
	case dashboard.PanelTypeGauge:
		value := m.getValueForPanel(panel.ID)
		return RenderGauge(panel.Title, value, 100, panel.GridPos, &m.viewport)
	case dashboard.PanelTypeStat:
		value := m.getValueForPanel(panel.ID)
		return RenderStat(panel.Title, fmt.Sprintf("%.2f", value), panel.GridPos, &m.viewport)
	case dashboard.PanelTypeTimeseries:
		var (
			series []*querier.TimeSeries
		)

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
	return m.checkServer()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case TickMsg:
		cmds = append(cmds, m.checkServer())

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
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
		if msg.Action == tea.MouseActionMotion {
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	var content strings.Builder
	var currentRow int
	var rowContent strings.Builder

	for _, p := range m.dashboard.Panels {

		if p.GridPos.Y > currentRow {
			content.WriteString(rowContent.String() + "\n")
			rowContent.Reset()
			currentRow = p.GridPos.Y
		}

		panelContent := m.renderPanel(p)
		if p.Type != dashboard.PanelTypeRow {
			panelContent = m.stylePanelBox(p.Title, panelContent, p.GridPos.W)
		}
		rowContent.WriteString(panelContent)
	}

	content.WriteString(rowContent.String())
	m.viewport.SetContent(content.String())

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func New(ctx context.Context, d *dashboard.Dashboard, q *querier.Querier) *App {
	d.Panels = dashboard.SortPanelsByPosition(d.Panels)

	p := tea.NewProgram(&Model{
		ctx:               ctx,
		dashboard:         d,
		querier:           q,
		results:           make(map[int]model.Value),
		timeseriesResults: make(map[int][]*querier.TimeSeries),
	},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)
	return &App{
		Program: p,
	}
}
