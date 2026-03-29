package ui

import (
	"os"
	"strings"

	"dunshell/internal/game"
)

type glyphSet struct {
	ascii bool
}

func newGlyphSet() glyphSet {
	value := strings.TrimSpace(strings.ToLower(os.Getenv("DUNSHELL_ASCII")))
	return glyphSet{ascii: value == "1" || value == "true" || value == "yes"}
}

func (g glyphSet) symbol(primary rune, fallback rune) string {
	if g.ascii {
		return string(fallback)
	}
	return string(primary)
}

func (g glyphSet) player() string {
	if g.ascii {
		return "@"
	}
	return "◆"
}

func (g glyphSet) stairs() string {
	if g.ascii {
		return ">"
	}
	return "▾"
}

func (g glyphSet) merchant() string {
	if g.ascii {
		return "$"
	}
	return "⚖"
}

func (g glyphSet) bossGate() string {
	if g.ascii {
		return "X"
	}
	return "⛧"
}

func (g glyphSet) chest(tier game.KeyTier) string {
	if g.ascii {
		switch tier {
		case game.KeyBronze:
			return "b"
		case game.KeySilver:
			return "s"
		default:
			return "g"
		}
	}
	return "▣"
}

func (g glyphSet) roomMarker(kind game.RoomKind) string {
	if g.ascii {
		switch kind {
		case game.RoomMerchant:
			return "$"
		case game.RoomBoss:
			return "B"
		case game.RoomSanctuary:
			return "+"
		default:
			return "o"
		}
	}
	switch kind {
	case game.RoomMerchant:
		return "◌"
	case game.RoomBoss:
		return "✠"
	case game.RoomSanctuary:
		return "✚"
	default:
		return "•"
	}
}
