package ui

import (
	"os"
	"strings"

	"dunshell/internal/game"
)

type glyphSet struct {
	ascii    bool
	envASCII bool
}

func newGlyphSet(settings game.Settings) glyphSet {
	settings = settings.Normalized()
	value := strings.TrimSpace(strings.ToLower(os.Getenv("DUNSHELL_ASCII")))
	envASCII := value == "1" || value == "true" || value == "yes"
	ascii := envASCII
	if !envASCII {
		switch settings.GlyphMode {
		case game.GlyphModeASCII:
			ascii = true
		case game.GlyphModeNerd:
			ascii = false
		default:
			ascii = settings.ASCIIFallback
		}
	}
	return glyphSet{ascii: ascii, envASCII: envASCII}
}

func (g glyphSet) ForcedASCII() bool {
	return g.envASCII
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

func (g glyphSet) floorVisible() string {
	if g.ascii {
		return "."
	}
	return "·"
}

func (g glyphSet) floorSeen() string {
	if g.ascii {
		return ","
	}
	return "ˑ"
}

func (g glyphSet) floorClearedVisible() string {
	if g.ascii {
		return ":"
	}
	return "•"
}

func (g glyphSet) floorClearedSeen() string {
	if g.ascii {
		return "."
	}
	return "·"
}

func (g glyphSet) bossFloor(visible bool) string {
	if g.ascii {
		if visible {
			return ";"
		}
		return ":"
	}
	if visible {
		return "▪"
	}
	return "▫"
}
