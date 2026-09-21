package tui

import (
	"embed"
	"io/fs"
	"math/rand/v2"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

const (
	timeAttackSeconds = time.Second * 60
	timeAttackWords   = 100
	endlessWords      = 10
)

//go:embed wordbanks/*.txt
var wordbankFS embed.FS

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
	wordbanks, err := getAllFilenames(&wordbankFS)
	if err != nil {
		return nil, err
	}
	model.wordbanks = wordbanks
	return model, nil
}

func buildsentence(wordsToGenerate int, selectedWordbank string) []string {
	dat, err := wordbankFS.ReadFile(selectedWordbank)
	if err != nil {
		panic(err)
	}
	datString := strings.TrimSpace(string(dat))
	stringSlice := strings.Fields(datString)
	sentenceSlice := make([]string, 0)
	for range wordsToGenerate {
		randIndex := rand.IntN(len(stringSlice))
		sentenceSlice = append(sentenceSlice, stringSlice[randIndex])
	}
	return sentenceSlice
}

func getAllFilenames(efs *embed.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}
