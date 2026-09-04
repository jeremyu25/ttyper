package tui

import (
	"fmt"
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
	s := fmt.Sprintf("You typed out the sentence in %v", m.resultTime)
	resultsView := tea.NewView(s)
	resultsView.AltScreen = true
	return resultsView
}
