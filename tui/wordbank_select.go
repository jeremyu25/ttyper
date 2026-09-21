package tui

import (
	"math/rand/v2"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/jeremyu25/ttyper/wordbanks"
)

const (
	timeAttackSeconds = time.Second * 60
	timeAttackWords   = 100
	endlessWords      = 10
)

type changeWordbankMsg struct {
	wordbank string
}

type wordbankKeymap struct {
	up    key.Binding
	down  key.Binding
	enter key.Binding
}

type wordbankModel struct {
	wordbankKeymap   wordbankKeymap
	selectedWordBank string
	wordbanks        []string
	cursor           int
}

func signalWordbankChangeMsg(wordbank string) tea.Cmd {
	return func() tea.Msg {
		return changeWordbankMsg{
			wordbank: wordbank,
		}
	}
}

func (m wordbankModel) Init() tea.Cmd {
	return nil
}

func (m wordbankModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.wordbankKeymap.down):
			if m.cursor >= len(m.wordbanks)-1 {
				return m, nil
			} else {
				m.cursor++
				return m, nil
			}
		case key.Matches(msg, m.wordbankKeymap.up):
			if m.cursor == 0 {
				return m, nil
			} else {
				m.cursor--
				return m, nil
			}
		case key.Matches(msg, m.wordbankKeymap.enter):
			return m, tea.Batch(signalWordbankChangeMsg(m.wordbanks[m.cursor]), signalMenuMsg())
		}
	}
	return m, nil
}
func (m wordbankModel) View() tea.View {
	var wordbankString strings.Builder
	wordbankString.WriteString("Select the wordbank you would like to use for typing below:\n\n")
	for i, wordbank := range m.wordbanks {
		if m.cursor == i {
			wordbankString.WriteString(boldmodeStyle.Render(wordbank))
		} else {
			wordbankString.WriteString(modeStyle.Render(wordbank))
		}
		wordbankString.WriteString("\n\n")
	}
	view := tea.NewView(wordbankString.String())
	view.AltScreen = true
	return view
}

func CreateNewWordbankModel() (tea.Model, error) {
	model := wordbankModel{
		wordbankKeymap: wordbankKeymap{
			up: key.NewBinding(
				key.WithKeys("up"),
			),
			down: key.NewBinding(
				key.WithKeys("down"),
			),
			enter: key.NewBinding(
				key.WithKeys("enter"),
			),
		},
	}
	wordbanks, err := wordbanks.GetAllFilenames()
	if err != nil {
		return nil, err
	}
	model.wordbanks = wordbanks
	return model, nil
}

func randomPickWords(wordsToGenerate int, selectedWordbank string) []string {
	stringSlice := wordbanks.CreateWordSlice(selectedWordbank)
	sentenceSlice := make([]string, 0)
	for range wordsToGenerate {
		randIndex := rand.IntN(len(stringSlice))
		sentenceSlice = append(sentenceSlice, stringSlice[randIndex])
	}
	return sentenceSlice
}
