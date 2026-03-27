package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type styles struct {
	App           lipgloss.Style
	HeaderBar     lipgloss.Style
	Title         lipgloss.Style
	Subtitle      lipgloss.Style
	Panel         lipgloss.Style
	PanelTitle    lipgloss.Style
	Muted         lipgloss.Style
	Dim           lipgloss.Style
	Accent        lipgloss.Style
	AccentSoft    lipgloss.Style
	Gold          lipgloss.Style
	Success       lipgloss.Style
	Danger        lipgloss.Style
	Info          lipgloss.Style
	Quantity      lipgloss.Style
	CompareBetter lipgloss.Style
	CompareEqual  lipgloss.Style
	CompareWorse  lipgloss.Style
	Heal          lipgloss.Style
	Focus         lipgloss.Style
	Cure          lipgloss.Style
	Ember         lipgloss.Style
	Attack        lipgloss.Style
	Defense       lipgloss.Style
	Vitality      lipgloss.Style
	Sight         lipgloss.Style
	MenuItem      lipgloss.Style
	MenuSelected  lipgloss.Style
	MenuCursor    lipgloss.Style
	Footer        lipgloss.Style
	Player        lipgloss.Style
	TileWall      lipgloss.Style
	TileWallSeen  lipgloss.Style
	TileFloor     lipgloss.Style
	TileFloorSeen lipgloss.Style
	TileDoor      lipgloss.Style
	TileDoorSeen  lipgloss.Style
	TileStairs    lipgloss.Style
}

func newStyles() styles {
	frame := lipgloss.Color("#5f4a3a")
	accent := lipgloss.Color("#f08a24")
	accentSoft := lipgloss.Color("#f2c078")
	muted := lipgloss.Color("#b4a89b")
	dim := lipgloss.Color("#7a7168")

	return styles{
		App: lipgloss.NewStyle().Foreground(lipgloss.Color("#efe7dc")),
		HeaderBar: lipgloss.NewStyle().
			Foreground(accentSoft).
			Bold(true).
			Align(lipgloss.Center).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(frame),
		Title: lipgloss.NewStyle().
			Foreground(accent).
			Bold(true),
		Subtitle: lipgloss.NewStyle().
			Foreground(muted),
		Panel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(frame).
			Padding(0, 1),
		PanelTitle: lipgloss.NewStyle().
			Foreground(accentSoft).
			Bold(true),
		Muted: lipgloss.NewStyle().Foreground(muted),
		Dim:   lipgloss.NewStyle().Foreground(dim),
		Accent: lipgloss.NewStyle().
			Foreground(accent).
			Bold(true),
		AccentSoft: lipgloss.NewStyle().Foreground(accentSoft),
		Gold: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d9b44a")).
			Bold(true),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#79c267")).
			Bold(true),
		Danger: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d16078")).
			Bold(true),
		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7fc8f8")),
		Quantity: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7fc8f8")).
			Bold(true),
		CompareBetter: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#79c267")).
			Bold(true),
		CompareEqual: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5fa8d3")).
			Bold(true),
		CompareWorse: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d16078")).
			Bold(true),
		Heal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#79c267")).
			Bold(true),
		Focus: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f2c078")).
			Bold(true),
		Cure: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7fc8f8")).
			Bold(true),
		Ember: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f08a24")).
			Bold(true),
		Attack: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f08a24")).
			Bold(true),
		Defense: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5fa8d3")).
			Bold(true),
		Vitality: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#79c267")).
			Bold(true),
		Sight: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d9b44a")).
			Bold(true),
		MenuItem: lipgloss.NewStyle().
			Foreground(muted).
			Padding(0, 1),
		MenuSelected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fdf4e8")).
			Background(accent).
			Bold(true).
			Padding(0, 1),
		MenuCursor: lipgloss.NewStyle().
			Foreground(accentSoft).
			Bold(true),
		Footer: lipgloss.NewStyle().
			Foreground(muted),
		Player: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fff4db")).
			Bold(true),
		TileWall:      lipgloss.NewStyle().Foreground(lipgloss.Color("#93857a")),
		TileWallSeen:  lipgloss.NewStyle().Foreground(lipgloss.Color("#5d554e")),
		TileFloor:     lipgloss.NewStyle().Foreground(lipgloss.Color("#d5c6b8")),
		TileFloorSeen: lipgloss.NewStyle().Foreground(lipgloss.Color("#776d64")),
		TileDoor:      lipgloss.NewStyle().Foreground(lipgloss.Color("#d4a66d")).Bold(true),
		TileDoorSeen:  lipgloss.NewStyle().Foreground(lipgloss.Color("#7f6750")),
		TileStairs:    lipgloss.NewStyle().Foreground(lipgloss.Color("#93c6d6")).Bold(true),
	}
}

func (s styles) box(title string, body string, width int, height int) string {
	innerWidth := max(12, width-4)
	header := s.PanelTitle.Render(strings.ToUpper(title))
	contentStyle := lipgloss.NewStyle().Width(innerWidth)
	if height > 0 {
		contentStyle = contentStyle.Height(max(1, height-3))
	}
	content := contentStyle.Render(body)
	rendered := lipgloss.JoinVertical(lipgloss.Left, header, content)
	return s.Panel.Copy().Width(width).Render(rendered)
}

func (s styles) colorGlyph(glyph rune, tint string, bold bool) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(tint))
	if bold {
		style = style.Bold(true)
	}
	return style.Render(string(glyph))
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
