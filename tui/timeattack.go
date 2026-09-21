package tui

import (
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type timeAttackKeymap struct {
	back   key.Binding
	delete key.Binding
	commit key.Binding
}

type timeAttackModel struct {
	timeAttackKeymap timeAttackKeymap
	timer            timer.Model
	typingEngine     typingEngine
	help             help.Model
	selectedWordbank string
}

func (m timeAttackModel) Init() tea.Cmd {
	return nil
}

func (m timeAttackModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		return m, signalResultsMsg(&m.typingEngine, timeAttackSeconds)

	case tea.KeyPressMsg:
		if m.typingEngine.wordIndex == len(m.typingEngine.wordSlice) {
			return m, signalResultsMsg(&m.typingEngine, timeAttackSeconds)
		}
		switch {
		case key.Matches(msg, m.timeAttackKeymap.back):
			return m, signalMenuMsg()
		case key.Matches(msg, m.timeAttackKeymap.delete):
			deleteChar(&m.typingEngine)
			return m, nil
		case key.Matches(msg, m.timeAttackKeymap.commit):
			commitWord(&m.typingEngine)
			return m, nil
		}
		if msg.Text != "" {
			if !m.typingEngine.started {
				buildInitialBuffer(&m.typingEngine)
				if isAlphabetMsg(msg) {
					addCharToBuffer(msg.Text, &m.typingEngine)
					m.typingEngine.started = true
				}
				return m, m.timer.Init()
			}
			lastCorrect := addCharToBuffer(msg.Text, &m.typingEngine)
			if lastCorrect && m.typingEngine.wordIndex == len(m.typingEngine.wordSlice)-1 {
				commitWord(&m.typingEngine)
				return m, signalResultsMsg(&m.typingEngine, timeAttackSeconds)
			}
			return m, nil
		}
	}
	return m, nil
}

func (m timeAttackModel) helpView() string {
	return "\n\n" + m.help.ShortHelpView([]key.Binding{
		m.timeAttackKeymap.back,
		m.timeAttackKeymap.delete,
	})
}

func (m timeAttackModel) View() tea.View {
	s := m.timer.View() + "\n"
	s += "\n" + renderWords(&m.typingEngine)
	s += m.helpView()
	timeAttackView := tea.NewView(s)
	timeAttackView.AltScreen = true
	return timeAttackView
}

func CreateNewTimeAttackModel(timerDuration time.Duration, selectedWordbank string) tea.Model {
	model := timeAttackModel{
		timer: timer.New(timerDuration, timer.WithInterval(time.Millisecond)),
		timeAttackKeymap: timeAttackKeymap{
			back: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "return to menu"),
			),
			delete: key.NewBinding(
				key.WithKeys("backspace"),
				key.WithHelp("backspace", "delete character"),
			),
			commit: key.NewBinding(
				key.WithKeys("space"),
			),
		},
		typingEngine: typingEngine{
			wordSlice: buildsentence(timeAttackWords, selectedWordbank),
			started:   false,
		},
		help:             help.New(),
		selectedWordbank: selectedWordbank,
	}
	model.help.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#8535fc")).Bold(true)
	model.help.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#8535fc"))
	model.help.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color("#8535fc")).Faint(true)
	return model
}
