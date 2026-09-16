package tui

import (
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
)

type timeAttackKeymap struct {
	pause  key.Binding
	back   key.Binding
	delete key.Binding
	commit key.Binding
}

type timeAttackModel struct {
	timeAttackKeymap timeAttackKeymap
	timer            timer.Model
	typingEngine     typingEngine
}

type typingEngine struct {
	userInput         string
	wordSlice         []string
	finishedWordSlice []detailedWord
	currentBuffer     detailedWord
	started           bool
	typedWords        int
	mistakes          int
	wordIndex         int
}

func (m timeAttackModel) Init() tea.Cmd {
	return nil
}

func (m timeAttackModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		return m, signalResultsMsg(&m.typingEngine)

	case tea.KeyPressMsg:
		if !m.typingEngine.started {
			buildInitialBuffer(&m.typingEngine)
			if isAlphabetMsg(msg) {
				addCharToBuffer(msg.Text, &m.typingEngine)
				m.typingEngine.started = true
			}
			return m, m.timer.Init()
		}
		if m.typingEngine.wordIndex == len(m.typingEngine.wordSlice) {
			return m, signalResultsMsg(&m.typingEngine)
		}
		switch {
		case key.Matches(msg, m.timeAttackKeymap.back):
			return m, signalMenuMsg()
		case key.Matches(msg, m.timeAttackKeymap.pause):
			return m, m.timer.Toggle()
		case key.Matches(msg, m.timeAttackKeymap.delete):
			deleteChar(&m.typingEngine)
			return m, nil
		case key.Matches(msg, m.timeAttackKeymap.commit):
			commitWord(&m.typingEngine)
			return m, nil
		}
		if msg.Text != "" {
			lastCorrect := addCharToBuffer(msg.Text, &m.typingEngine)
			if lastCorrect && m.typingEngine.wordIndex == len(m.typingEngine.wordSlice)-1 {
				commitWord(&m.typingEngine)
				return m, signalResultsMsg(&m.typingEngine)
			}
			return m, nil
		}
	}
	return m, nil
}

func (m timeAttackModel) View() tea.View {
	s := m.timer.View()
	s += "\n" + renderWords(&m.typingEngine)
	timeAttackView := tea.NewView(s)
	timeAttackView.AltScreen = true
	return timeAttackView
}

func CreateNewTimeAttackModel(timerDuration time.Duration) tea.Model {
	return timeAttackModel{
		timer: timer.New(timerDuration, timer.WithInterval(time.Millisecond)),
		timeAttackKeymap: timeAttackKeymap{
			pause: key.NewBinding(
				key.WithKeys("esc"),
			),
			back: key.NewBinding(
				key.WithKeys("shift+tab"),
			),
			delete: key.NewBinding(
				key.WithKeys("backspace"),
			),
			commit: key.NewBinding(
				key.WithKeys("space"),
			),
		},
		typingEngine: typingEngine{
			wordSlice: buildsentence(),
			started:   false,
		},
	}
}
