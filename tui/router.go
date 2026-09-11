package tui

import (
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type routerModel struct {
	globalKeymap globalKeymap
	currentState tea.Model
}

type globalKeymap struct {
	quit key.Binding
}

type timeAttackTransitionMsg struct {
	timerDuration time.Duration
}

type menuTransitionMsg struct{}

type resultsTransitionMsg struct {
	resultTime time.Duration
}

func signalMenuMsg() tea.Cmd {
	return func() tea.Msg {
		return menuTransitionMsg{}
	}
}

func signalStartTimeAttack(timerDuration time.Duration) tea.Cmd {
	return func() tea.Msg {
		return timeAttackTransitionMsg{
			timerDuration: timerDuration,
		}
	}
}
func signalResultsMsg(resultTime time.Duration) tea.Cmd {
	return func() tea.Msg {
		return resultsTransitionMsg{
			resultTime: resultTime,
		}
	}
}

func (m routerModel) Init() tea.Cmd {
	return nil
}

func (m routerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timeAttackTransitionMsg:
		m.currentState = CreateNewTimeAttackModel(msg.timerDuration)
		return m, m.currentState.Init()

	case menuTransitionMsg:
		m.currentState = menuModel{
			menuKeymap: menuKeymap{
				start: key.NewBinding(
					key.WithKeys("enter", "space"),
					key.WithHelp("enter/space", "begin"),
				),
			},
			globalKeymap: globalKeymap{
				quit: key.NewBinding(
					key.WithKeys("ctrl+c"),
					key.WithHelp("q", "quit"),
				),
			},
		}
		return m, m.currentState.Init()

	case resultsTransitionMsg:
		m.currentState = resultsModel{
			resultTime: msg.resultTime,
			resultsKeymap: resultsKeymap{
				restart: key.NewBinding(
					key.WithKeys("enter"),
					key.WithHelp("Press enter to resume", "enter"),
				),
			},
		}
		return m, m.currentState.Init()
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.globalKeymap.quit):
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.currentState, cmd = m.currentState.Update(msg)
	return m, cmd
}

func (m routerModel) View() tea.View {
	return m.currentState.View()
}

func CreateNewRouterModel() tea.Model {
	return routerModel{
		currentState: NewMenuModel(),
		globalKeymap: globalKeymap{
			quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
	}
}
