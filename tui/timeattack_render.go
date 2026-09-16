package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	correctCharStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#0cf249"))

	wrongCharStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e62f17"))
	correctCharWrongWordStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#0cf249")).
					Underline(true).
					UnderlineStyle(lipgloss.UnderlineCurly).
					UnderlineColor(lipgloss.
						Color("#e62f17"))
	wrongCharWrongWordStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e62f17")).
				Underline(true).
				UnderlineStyle(lipgloss.UnderlineCurly).
				UnderlineColor(lipgloss.
					Color("#e62f17"))
	pendingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#b7c2c4")).
			Faint(true)
	overtypedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e62f17")).
			UnderlineStyle(lipgloss.UnderlineCurly).
			UnderlineColor(lipgloss.
				Color("#e62f17")).
			Faint(true)
)

func renderWords(m *typingEngine) string {
	var renderedWordsCombined strings.Builder
	var renderedWordsSlice []string
	if !m.started {
		//if not started, just render all as pending
		renderedWordsSlice = renderPendingWords(m.wordSlice)
	} else {
		//render finished words
		renderedWordsSlice = append(renderedWordsSlice, renderFinishedWords(m.finishedWordSlice)...)
		//render the word at current index
		currWord := renderCurrWord(m.currentBuffer)
		renderedWordsSlice = append(renderedWordsSlice, currWord)
		//render all other pending words
		if m.wordIndex < len(m.wordSlice)-1 {
			pendingWords := m.wordSlice[m.wordIndex+1:]
			renderedWordsSlice = append(renderedWordsSlice, renderPendingWords(pendingWords)...)
		}
	}
	for i, word := range renderedWordsSlice {
		if (i%20 == 0) && (i > 0) {
			renderedWordsCombined.WriteString(word)
			renderedWordsCombined.WriteString("\n")
		} else {
			renderedWordsCombined.WriteString(word)
			renderedWordsCombined.WriteString(" ")
		}
	}
	return renderedWordsCombined.String()
}

func renderPendingWords(wordSlice []string) []string {
	var renderedWordsSlice []string
	for i := range len(wordSlice) {
		finishedWord := pendingStyle.Render(wordSlice[i])
		renderedWordsSlice = append(renderedWordsSlice, finishedWord)
	}
	return renderedWordsSlice
}

func renderFinishedWords(wordSlice []detailedWord) []string {
	var renderedWordsSlice []string
	for i := range len(wordSlice) {
		finishedWord := renderFinishedWord(wordSlice[i])
		renderedWordsSlice = append(renderedWordsSlice, finishedWord)

	}
	return renderedWordsSlice
}

func renderFinishedWord(word detailedWord) string {
	var finishedWord strings.Builder
	for _, char := range word.characters {
		if char.charState == CorrectChar {
			if word.isCorrect {
				finishedWord.WriteString(correctCharStyle.Render(char.character))
			} else {
				finishedWord.WriteString(correctCharWrongWordStyle.Render(char.character))
			}
		} else {
			if word.isCorrect {
				finishedWord.WriteString(wrongCharStyle.Render(char.character))
			} else {
				finishedWord.WriteString(wrongCharWrongWordStyle.Render(char.character))
			}
		}
	}
	if len(word.overtypedCharacters) > 0 {
		finishedWord.WriteString(overtypedStyle.Render(word.overtypedCharacters))
	}
	return finishedWord.String()
}

func renderCurrWord(word detailedWord) string {
	var currWord strings.Builder
	for _, char := range word.characters {
		switch char.charState {
		case CorrectChar:
			currWord.WriteString(correctCharStyle.Render(char.character))
		case PendingChar:
			currWord.WriteString(pendingStyle.Render(char.character))
		case IncorrectChar:
			currWord.WriteString(wrongCharStyle.Render(char.character))
		}
	}
	if len(word.overtypedCharacters) > 0 {
		currWord.WriteString(overtypedStyle.Render(word.overtypedCharacters))
	}
	return currWord.String()
}
