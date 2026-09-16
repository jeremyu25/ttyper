package tui

import (
	"unicode"

	tea "charm.land/bubbletea/v2"
)

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

func buildInitialBuffer(m *typingEngine) {
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

func addCharToBuffer(character string, m *typingEngine) bool {
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

func commitWord(m *typingEngine) {
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

func deleteChar(m *typingEngine) {
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
