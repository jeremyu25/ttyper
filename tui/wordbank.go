package tui

import (
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	timeAttackSeconds = time.Second * 60
	timeAttackWords   = 100
	endlessWords      = 10
)

func buildsentence(wordsToGenerate int) []string {
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
	for range wordsToGenerate {
		randIndex := rand.IntN(len(stringSlice) + 1)
		sentenceSlice = append(sentenceSlice, stringSlice[randIndex])
	}
	return sentenceSlice
}
