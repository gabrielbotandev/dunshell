package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"dunshell/internal/ui"
)

func main() {
	var seed int64
	flag.Int64Var(&seed, "seed", 0, "replay a specific run seed")
	flag.Parse()

	model := ui.NewModel(seed, seed != 0)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "dunshell: %v\n", err)
		os.Exit(1)
	}
}
