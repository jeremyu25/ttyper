package tui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
	itemStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	boldItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	modes         = []string{
		"Timed Typer",
		"Pied Piper",
	}
)

type menuKeymap struct {
	start key.Binding
	up    key.Binding
	down  key.Binding
}

type menuModel struct {
	menuKeymap   menuKeymap
	globalKeymap globalKeymap
	cursor       int
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
		case key.Matches(msg, m.menuKeymap.down):
			if m.cursor >= len(modes) {
				return m, nil
			} else {
				m.cursor++
				return m, nil
			}
		case key.Matches(msg, m.menuKeymap.up):
			if m.cursor == 0 {
				return m, nil
			} else {
				m.cursor--
				return m, nil
			}
		case key.Matches(msg, m.globalKeymap.quit):
			return m, tea.Quit

		}
	}
	return m, nil
}

func (m menuModel) View() tea.View {
	var menuStringBuilder strings.Builder
	menuStringBuilder.WriteString(titleStyle.Render("Welcome to TTyper! Select a mode below to begin:"))
	for i, mode := range modes {
		if m.cursor == i {
			menuStringBuilder.WriteString(boldItemStyle.Render(fmt.Sprintf("\n> %v", mode)))
		} else {
			menuStringBuilder.WriteString(itemStyle.Render(fmt.Sprintf("\n  %v", mode)))
		}

	}
	menuString := menuStringBuilder.String()
	menuView := tea.NewView(menuString)
	menuView.AltScreen = true
	return menuView
}

func NewMenuModel() menuModel {
	return menuModel{
		menuKeymap: menuKeymap{
			start: key.NewBinding(
				key.WithKeys("enter", "space"),
				key.WithHelp("enter/space", "begin"),
			),
			up: key.NewBinding(
				key.WithKeys("up"),
			),
			down: key.NewBinding(
				key.WithKeys("down"),
			),
		},
		globalKeymap: globalKeymap{
			quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("q", "quit"),
			)},
		cursor: 0,
	}
}
