package tui

import (
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
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

type timeAttackKeymap struct {
	pause  key.Binding
	back   key.Binding
	delete key.Binding
	commit key.Binding
}

type CharState int

const (
	IncorrectChar CharState = iota
	CorrectChar
	PendingChar
)

type detailedCharacter struct {
	character string
	charState CharState
}

type detailedWord struct {
	isCorrect           bool
	word                string
	overtypedCharacters string
	characters          []detailedCharacter
	characterIndex      int
}

type timeAttackModel struct {
	timeAttackKeymap  timeAttackKeymap
	timer             timer.Model
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
		return m, signalResultsMsg(&m)

	case tea.KeyPressMsg:
		if !m.started {
			buildInitialBuffer(&m)
			if isAlphabetMsg(msg) {
				addCharToBuffer(msg.Text, &m)
				m.started = true
			}
			return m, m.timer.Init()
		}
		if m.wordIndex == len(m.wordSlice) {
			return m, signalResultsMsg(&m)
		}
		switch {
		case key.Matches(msg, m.timeAttackKeymap.back):
			return m, signalMenuMsg()
		case key.Matches(msg, m.timeAttackKeymap.pause):
			return m, m.timer.Toggle()
		case key.Matches(msg, m.timeAttackKeymap.delete):
			deleteChar(&m)
			return m, nil
		case key.Matches(msg, m.timeAttackKeymap.commit):
			commitWord(&m)
			return m, nil
		}
		if msg.Text != "" {
			lastCorrect := addCharToBuffer(msg.Text, &m)
			if lastCorrect && m.wordIndex == len(m.wordSlice)-1 {
				commitWord(&m)
				return m, signalResultsMsg(&m)
			}
			return m, nil
		}
	}
	return m, nil
}

func (m timeAttackModel) View() tea.View {
	s := m.timer.View()
	s += "\n" + renderWords(&m)
	timeAttackView := tea.NewView(s)
	timeAttackView.AltScreen = true
	return timeAttackView
}

func buildsentence() []string {
	path := filepath.Join("./wordbanks", "google-10000-english-usa-no-swears.txt")
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
	for range 100 {
		randIndex := rand.IntN(len(stringSlice) + 1)
		sentenceSlice = append(sentenceSlice, stringSlice[randIndex])
	}
	return sentenceSlice
}

func buildInitialBuffer(m *timeAttackModel) {
	if m.wordIndex == len(m.wordSlice) {
		return
	}
	currWord := m.wordSlice[m.wordIndex]
	m.currentBuffer = detailedWord{
		isCorrect: false,
		word:      currWord,
	}
	for _, char := range currWord {
		detailedChar := detailedCharacter{
			character: string(char),
			charState: PendingChar,
		}
		m.currentBuffer.characters = append(m.currentBuffer.characters, detailedChar)
	}
}

func isAlphabetMsg(msg tea.KeyMsg) bool {
	return unicode.IsLetter(msg.Key().Code)
}

func addCharToBuffer(character string, m *timeAttackModel) bool {
	if m.currentBuffer.characterIndex == len(m.currentBuffer.characters) {
		m.currentBuffer.overtypedCharacters += character
		m.mistakes++
		return false
	}
	if m.currentBuffer.characters[m.currentBuffer.characterIndex].character == character {
		m.currentBuffer.characters[m.currentBuffer.characterIndex].charState = CorrectChar
		if m.currentBuffer.characterIndex == len(m.currentBuffer.characters)-1 {
			m.typedWords++
			m.currentBuffer.characterIndex++
			return true
		}
		m.typedWords++
	} else {
		m.currentBuffer.characters[m.currentBuffer.characterIndex].charState = IncorrectChar
		m.mistakes++
	}
	m.currentBuffer.characterIndex++
	return false
}

func commitWord(m *timeAttackModel) {
	if m.currentBuffer.characterIndex == 0 {
		//do nothing as word not started yet
		return
	} else {
		//set iscorrect before appending to finished words
		m.currentBuffer.isCorrect = true
		for _, char := range m.currentBuffer.characters {
			if char.charState == IncorrectChar || char.charState == PendingChar {
				m.currentBuffer.isCorrect = false
			}
		}
		if len(m.currentBuffer.overtypedCharacters) > 0 {
			m.currentBuffer.isCorrect = false
		}
		m.finishedWordSlice = append(m.finishedWordSlice, m.currentBuffer)
		m.wordIndex++
		if m.wordIndex == len(m.wordSlice) {
			return
		}
		if !m.currentBuffer.isCorrect {
			m.mistakes++
		}
		m.typedWords++
		buildInitialBuffer(m)
	}
}

func deleteChar(m *timeAttackModel) {
	//check if currently typing in a buffer
	if m.currentBuffer.characters[0].charState != PendingChar {
		//check if any overflow
		if len(m.currentBuffer.overtypedCharacters) > 0 {
			m.currentBuffer.overtypedCharacters = m.currentBuffer.overtypedCharacters[:len(m.currentBuffer.overtypedCharacters)-1]
			return
		}
		m.currentBuffer.characters[m.currentBuffer.characterIndex-1].charState = PendingChar
		m.currentBuffer.characterIndex--
	} else {
		//if not in a buffer, go into the finished word slices
		if len(m.finishedWordSlice) > 0 {
			if !m.finishedWordSlice[len(m.finishedWordSlice)-1].isCorrect {
				m.currentBuffer = m.finishedWordSlice[len(m.finishedWordSlice)-1]
				m.finishedWordSlice = m.finishedWordSlice[:len(m.finishedWordSlice)-1]
				m.wordIndex--
			}
		}
	}
}

func renderWords(m *timeAttackModel) string {
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
		wordSlice: buildsentence(),
		started:   false,
	}
}
