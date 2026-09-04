package tui

import (
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type menuKeymap struct {
	start key.Binding
}

type menuModel struct {
	menuKeymap   menuKeymap
	globalKeymap globalKeymap
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.menuKeymap.start):
			return m, signalStartTimedTyper(time.Second * 30)
		case key.Matches(msg, m.globalKeymap.quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menuModel) View() tea.View {
	s := "Press enter to start"
	menuView := tea.NewView(s)
	menuView.AltScreen = true
	return menuView
}
