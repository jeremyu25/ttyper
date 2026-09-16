package tui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type resultsKeymap struct {
	restart key.Binding
}

type resultsModel struct {
	resultsKeymap resultsKeymap
	resultTime    time.Duration
	correctWords  int
	wrongWords    int
	totalChars    int
	wrongChars    int
}

func (m resultsModel) Init() tea.Cmd {
	return nil
}

func (m resultsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.resultsKeymap.restart):
			return m, signalMenuMsg()
		}

	}
	return m, nil
}

func (m resultsModel) View() tea.View {
	wpm := float64((m.correctWords + m.wrongWords)) / m.resultTime.Minutes()
	accuracy := (1 - (float64(m.wrongChars) / float64(m.totalChars))) * 100
	var resultStringBuilder strings.Builder
	_, err := fmt.Fprintf(&resultStringBuilder, `You typed out the sentence in %v
		Your WPM was %v
		You made %v mistakes in total, giving an accuracy of %.2f%% out of %v chars typed.
		You typed out %v words, with %v wrongly typed words.`, m.resultTime, wpm, m.wrongChars, accuracy, m.totalChars, (m.correctWords + m.wrongWords), m.wrongWords)
	if err != nil {
		panic(err)
	}
	resultString := resultStringBuilder.String()
	resultsView := tea.NewView(resultString)
	resultsView.AltScreen = true
	return resultsView
}

func CreateNewResultsModel(resultTime time.Duration, correctWords, wrongWords, totalChars, wrongChars int) resultsModel {
	return resultsModel{
		resultTime:   time.Second * 60,
		correctWords: correctWords,
		wrongWords:   wrongWords,
		totalChars:   totalChars,
		wrongChars:   wrongChars,
		resultsKeymap: resultsKeymap{
			restart: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("Press enter to resume", "enter"),
			),
		},
	}
}
