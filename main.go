package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	tui "github.com/jeremyu25/ttyper/tui"
)

func main() {
	m := tui.CreateNewRouterModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
}
