package tui

import (
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

var (
	sentenceStyle = lipgloss.NewStyle().
		Bold(true)
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
	sentence        string
	started         bool
}

func (m typedTimerModel) Init() tea.Cmd {
	return nil
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
		if !m.started {
			if msg.Text != "" {
				m.userInput += msg.Text
				m.started = true
				return m, m.timer.Init()
			}
		}
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
			if m.userInput == m.sentence {
				return m, signalResultsMsg(m.timer.Timeout)
			}
		}
	}
	return m, nil
}

func (m typedTimerModel) View() tea.View {
	s := m.timer.View()
	s += "\n" + m.sentence
	s += "\n" + m.userInput
	// s = sentenceStyle.Render(s)
	typedTimerView := tea.NewView(s)
	typedTimerView.AltScreen = true
	return typedTimerView
}

func CreateNewTimedTypeModel(timerDuration time.Duration) tea.Model {
	return typedTimerModel{
		timer: timer.New(timerDuration, timer.WithInterval(time.Millisecond)),
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
		sentence: buildsentence(),
		started:  false,
	}
}

func buildsentence() string {
	path := filepath.Join(".", "google-10000-english-usa-no-swears-long.txt")
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	datString := strings.TrimSpace(string(dat))
	stringSlice := strings.Fields(datString)
	sentenceSlice := make([]string, 0)
	for _ = range 10 {
		randIndex := rand.IntN(len(stringSlice) + 1)
		sentenceSlice = append(sentenceSlice, stringSlice[randIndex])
	}
	sentence := strings.Join(sentenceSlice, " ")
	return sentence
}
