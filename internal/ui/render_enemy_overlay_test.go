package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"dunshell/internal/game"
)

func TestEnemyOverlayCellsIncludeBarAndLevel(t *testing.T) {
	model := &Model{styles: newStyles()}
	enemy := &game.Enemy{
		Level: 3,
		HP:    6,
		Template: game.EnemyTemplate{
			MaxHP: 10,
		},
	}

	plain := ansi.Strip(strings.Join(model.enemyOverlayCells(enemy), ""))
	if !strings.Contains(plain, "·3") {
		t.Fatalf("overlay %q does not include enemy level", plain)
	}
	if len([]rune(plain)) < 5 {
		t.Fatalf("overlay %q is too short to include a visible bar", plain)
	}
}
