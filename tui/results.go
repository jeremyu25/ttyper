package tui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	resultsStyle = lipgloss.NewStyle().
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#63ceff"))
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
	_, err := fmt.Fprintf(&resultStringBuilder, ` ______     ______     __   __     ______     ______     ______     ______   __  __     __         ______     ______   __     ______     __   __     ______
/\  ___\   /\  __ \   /\ "-.\ \   /\  ___\   /\  == \   /\  __ \   /\__  _\ /\ \/\ \   /\ \       /\  __ \   /\__  _\ /\ \   /\  __ \   /\ "-.\ \   /\  ___\
\ \ \____  \ \ \/\ \  \ \ \-.  \  \ \ \__ \  \ \  __<   \ \  __ \  \/_/\ \/ \ \ \_\ \  \ \ \____  \ \  __ \  \/_/\ \/ \ \ \  \ \ \/\ \  \ \ \-.  \  \ \___  \
 \ \_____\  \ \_____\  \ \_\\"\_\  \ \_____\  \ \_\ \_\  \ \_\ \_\    \ \_\  \ \_____\  \ \_____\  \ \_\ \_\    \ \_\  \ \_\  \ \_____\  \ \_\\"\_\  \/\_____\
  \/_____/   \/_____/   \/_/ \/_/   \/_____/   \/_/ /_/   \/_/\/_/     \/_/   \/_____/   \/_____/   \/_/\/_/     \/_/   \/_/   \/_____/   \/_/ \/_/   \/_____/

  		You typed out the sentence in %v

		Your WPM was %v

		You made %v mistakes in total, giving an accuracy of %.2f%% out of %v chars typed.

		You typed out %v words, with %v wrongly typed words.

		Press enter to continue...`, m.resultTime, wpm, m.wrongChars, accuracy, m.totalChars, (m.correctWords + m.wrongWords), m.wrongWords)

	if err != nil {
		panic(err)
	}
	resultString := resultStringBuilder.String()
	resultString = resultsStyle.Render(resultString)
	resultsView := tea.NewView(resultString)
	resultsView.AltScreen = true
	return resultsView
}

func CreateNewResultsModel(resultTime time.Duration, correctWords, wrongWords, totalChars, wrongChars int) resultsModel {
	return resultsModel{
		resultTime:   resultTime,
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
