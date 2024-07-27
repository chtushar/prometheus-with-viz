package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prometheus/prometheus/cmd/promviz/dashboard"
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

type Model struct {
	viewport viewport.Model
	content  string
	ready   bool
	dashboard *dashboard.Dashboard
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
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
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
		pos := dashboard.GetNewGridPos(p.GridPos, m.viewport.Width)

		panelStyle.Width(pos.W / 3)
		panelStyle.Height(pos.H / 3)
		panelStyle.MarginTop(pos.Y * m.viewport.Height / 24)
		panelStyle.MarginLeft(pos.X * m.viewport.Width / 24)

		panel := panelStyle.Render(fmt.Sprintf("%s\n%d\n%d", p.Title, pos.W, len(m.dashboard.Panels)))

		content += panel
	}
	m.viewport.SetContent(content)

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}


func New(d *dashboard.Dashboard) *App {
	p := tea.NewProgram(&Model{
			dashboard: d,
		},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)
	return &App{
		Program: p,
	}
}