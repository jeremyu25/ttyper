package main

import (
	"fmt"
	"os"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
)

const timeout = time.Second * 100

const sentence = "the quick brown fox jumps over the fence"

type AppState int

const (
	Menu AppState = iota
	TypeTimer
)

type model struct {
	timer           timer.Model
	globalKeymap    globalKeymap
	menuKeymap      menuKeymap
	typeTimerKeymap typeTimerKeymap
	quitting        bool
	currentState    AppState
	userInput       string
}

type globalKeymap struct {
	quit key.Binding
}

type menuKeymap struct {
	start key.Binding
}

type typeTimerKeymap struct {
	pause  key.Binding
	back   key.Binding
	delete key.Binding
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.globalKeymap.quit):
			m.quitting = true
			return m, tea.Quit
		}
	}
	switch m.currentState {
	case Menu:
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch {
			case key.Matches(msg, m.menuKeymap.start):
				m.currentState = TypeTimer
				m.timer = timer.New(timeout, timer.WithInterval(time.Millisecond))
				return m, m.timer.Init()
			}
		}

	case TypeTimer:
		switch msg := msg.(type) {

		case timer.StartStopMsg:
			var cmd tea.Cmd
			m.timer, cmd = m.timer.Update(msg)
			return m, cmd

		case timer.TickMsg:
			var cmd tea.Cmd
			m.timer, cmd = m.timer.Update(msg)
			return m, cmd

		case timer.TimeoutMsg:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyPressMsg:
			switch {
			case key.Matches(msg, m.typeTimerKeymap.back):
				m.currentState = Menu
				m.timer.Timeout = timeout
				m.userInput = ""
				return m, nil
			case key.Matches(msg, m.typeTimerKeymap.pause):
				return m, m.timer.Toggle()
			case key.Matches(msg, m.typeTimerKeymap.delete):
				if len(m.userInput) > 0 {
					m.userInput = m.userInput[:len(m.userInput)-1]
				}
			}
			if msg.Text != "" {
				m.userInput += msg.Text
			}
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	switch m.currentState {
	case Menu:
		s := "Press enter to start"
		return tea.NewView(s)
	case TypeTimer:
		s := m.timer.View()
		s += "\n" + sentence
		s += "\n" + m.userInput
		return tea.NewView(s)
	default:
		return tea.NewView("hi")
	}
}

func main() {
	m := model{
		globalKeymap: globalKeymap{
			quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
		menuKeymap: menuKeymap{
			start: key.NewBinding(
				key.WithKeys("enter", "space"),
				key.WithHelp("enter/space", "begin"),
			),
		},
		typeTimerKeymap: typeTimerKeymap{
			pause: key.NewBinding(
				key.WithKeys("esc"),
			),
			back: key.NewBinding(
				key.WithKeys("shift+tab"),
			),
			delete: key.NewBinding(
				key.WithKeys("backspace"),
			),
		},
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
}
