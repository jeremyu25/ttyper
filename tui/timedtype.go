package tui

import (
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
)

type typeTimerKeymap struct {
	pause  key.Binding
	back   key.Binding
	delete key.Binding
}

type typedTimerModel struct {
	typeTimerKeymap typeTimerKeymap
	timer           timer.Model
	userInput       string
}

func (m typedTimerModel) Init() tea.Cmd {
	return m.timer.Init()
}

func (m typedTimerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		return m, tea.Quit

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.typeTimerKeymap.back):
			m.timer.Timeout = time.Second * 30
			m.userInput = ""
			return m, signalMenuMsg()
		case key.Matches(msg, m.typeTimerKeymap.pause):
			return m, m.timer.Toggle()
		case key.Matches(msg, m.typeTimerKeymap.delete):
			if len(m.userInput) > 0 {
				m.userInput = m.userInput[:len(m.userInput)-1]
			}
			return m, nil
		}
		if msg.Text != "" {
			m.userInput += msg.Text
			if m.userInput == sentence {
				return m, signalResultsMsg(m.timer.Timeout)
			}
		}
	}
	return m, nil
}

func (m typedTimerModel) View() tea.View {
	s := m.timer.View()
	s += "\n" + sentence
	s += "\n" + m.userInput
	typedTimerView := tea.NewView(s)
	typedTimerView.AltScreen = true
	return typedTimerView
}
