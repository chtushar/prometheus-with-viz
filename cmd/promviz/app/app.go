package app

import (
	"context"
	"fmt"
	"math"
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

	panelStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "│"
		b.Right = "│"
		b.Top = "─"
		b.Bottom = "─"
		return lipgloss.NewStyle().BorderStyle(b).Padding(1, 1)
	}()
)

type App struct {
	Program *tea.Program
}

type TickMsg time.Time

type Model struct {
	ctx       context.Context
	viewport  viewport.Model
	ready     bool
	dashboard *dashboard.Dashboard
	results   map[int]model.Value
	querier   *querier.Querier
}

func (m Model) checkServer() tea.Cmd {
	parsedDuration, err := time.ParseDuration(m.dashboard.Refresh)

	if err != nil {
		panic(err)
	}

	variableValues := map[string]string{
		"$node":            "anakin-rpi.lan:9100",
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
				_, err := m.querier.FetchTimeSeriesPanelData(m.ctx, p, start, end, variableValues)
				if err != nil {
					return fmt.Errorf("failed to load panel %d", p.ID)
				}

				// TODO: rendering timeseries??

				// results[panel.ID] = *data[0]

				// fmt.Println(data)
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

	content := ""

	for _, p := range m.dashboard.Panels {
		var panel string

		result := m.results[p.ID]

		switch p.Type {
		case dashboard.PanelTypeGauge:
			value := math.NaN()

			vec, ok := result.(model.Vector)
			if ok && len(vec) != 0 {
				value = float64(vec[0].Value)
			}

			panel = RenderGauge(p.Title, value, 100, p.GridPos, &m.viewport)

		case dashboard.PanelTypeStat:
			value := math.NaN()

			vec, ok := result.(model.Vector)
			if ok && len(vec) != 0 {
				value = float64(vec[0].Value)
			}

			str := fmt.Sprintf("%.2f %s", value, p.FieldConfig.Defaults.Unit)

			panel = RenderStat(p.Title, str, p.GridPos, &m.viewport)
		}

		content += panel
	}
	m.viewport.SetContent(content)

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func New(ctx context.Context, d *dashboard.Dashboard, querier *querier.Querier) *App {
	p := tea.NewProgram(&Model{
		ctx:       ctx,
		dashboard: d,
		querier:   querier,
		results:   make(map[int]model.Value),
	},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)
	return &App{
		Program: p,
	}
}
