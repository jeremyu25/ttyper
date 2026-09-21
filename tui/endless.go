package tui

import (
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/stopwatch"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type endlessKeymap struct {
	finish key.Binding
	back   key.Binding
	delete key.Binding
	commit key.Binding
}

type endlessModel struct {
	endlessKeymap endlessKeymap
	stopwatch     stopwatch.Model
	typingEngine  typingEngine
	help          help.Model
}

func (m endlessModel) Init() tea.Cmd {
	return nil
}

func (m endlessModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case stopwatch.TickMsg:
		var cmd tea.Cmd
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		return m, cmd

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.endlessKeymap.back):
			return m, signalMenuMsg()
		case key.Matches(msg, m.endlessKeymap.finish):
			return m, signalResultsMsg(&m.typingEngine, m.stopwatch.Elapsed())
		case key.Matches(msg, m.endlessKeymap.delete):
			deleteChar(&m.typingEngine)
			return m, nil
		case key.Matches(msg, m.endlessKeymap.commit):
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
				return m, m.stopwatch.Init()
			}
			lastCorrect := addCharToBuffer(msg.Text, &m.typingEngine)
			if lastCorrect && m.typingEngine.wordIndex == len(m.typingEngine.wordSlice)-1 {
				commitWord(&m.typingEngine)
				m.typingEngine.wordSlice = buildsentence(endlessWords)
				m.typingEngine.wordIndex = 0
				//we need to rebuild the buffer after resetting the words for endless
				buildInitialBuffer(&m.typingEngine)
				return m, nil
			}
			return m, nil
		}
	}
	return m, nil
}

func (m endlessModel) helpView() string {
	return "\n\n" + m.help.ShortHelpView([]key.Binding{
		m.endlessKeymap.finish,
		m.endlessKeymap.back,
	})
}

func (m endlessModel) View() tea.View {
	s := m.stopwatch.View() + "\n"
	s += "\n" + renderEndlessWords(&m.typingEngine)
	s += m.helpView()
	endlessView := tea.NewView(s)
	endlessView.AltScreen = true
	return endlessView
}

func CreateNewEndlessModel() tea.Model {
	model := endlessModel{
		stopwatch: stopwatch.New(stopwatch.WithInterval(time.Millisecond)),
		endlessKeymap: endlessKeymap{
			finish: key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "finish typing"),
			),
			back: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "return to menu"),
			),
			delete: key.NewBinding(
				key.WithKeys("backspace"),
			),
			commit: key.NewBinding(
				key.WithKeys("space"),
			),
		},
		typingEngine: typingEngine{
			wordSlice: buildsentence(endlessWords),
			started:   false,
		},
		help: help.New(),
	}
	model.help.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#8535fc")).Bold(true)
	model.help.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#8535fc"))
	model.help.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color("#8535fc")).Faint(true)
	return model
}
