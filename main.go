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
	var godMode bool
	flag.Int64Var(&seed, "seed", 0, "replay a specific run seed")
	flag.BoolVar(&godMode, "god", false, "start a developer run with invulnerability and endgame gear")
	flag.Parse()

	model := ui.NewModel(ui.StartupOptions{
		Seed:          seed,
		HasLockedSeed: seed != 0,
		GodMode:       godMode,
	})
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "dunshell: %v\n", err)
		os.Exit(1)
	}
}
