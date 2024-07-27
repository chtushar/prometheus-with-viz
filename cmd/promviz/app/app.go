package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	Program *tea.Program
}

type Model struct {

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg: 
			switch msg.Type {
				case tea.KeyEsc:
				  	return m, tea.Quit
			  	case tea.KeyCtrlC:
					return m, tea.Quit
			}
	}
	return m, nil
}

func (m Model) View() string {
	return "Hello, Bubble Tea!"
}

func New() *App {
	p := tea.NewProgram(&Model{})
	return &App{
		Program: p,
	}
}