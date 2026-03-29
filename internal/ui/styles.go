package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"dunshell/internal/game"
)

type styles struct {
	App                     lipgloss.Style
	HeaderBar               lipgloss.Style
	Title                   lipgloss.Style
	Subtitle                lipgloss.Style
	Panel                   lipgloss.Style
	PanelActive             lipgloss.Style
	ModalPanel              lipgloss.Style
	PanelTitle              lipgloss.Style
	PanelTitleActive        lipgloss.Style
	PanelTitleModal         lipgloss.Style
	Muted                   lipgloss.Style
	Dim                     lipgloss.Style
	Accent                  lipgloss.Style
	AccentSoft              lipgloss.Style
	Gold                    lipgloss.Style
	Success                 lipgloss.Style
	Danger                  lipgloss.Style
	Info                    lipgloss.Style
	Warning                 lipgloss.Style
	PanelNote               lipgloss.Style
	Quantity                lipgloss.Style
	CompareBetter           lipgloss.Style
	CompareEqual            lipgloss.Style
	CompareWorse            lipgloss.Style
	Heal                    lipgloss.Style
	Focus                   lipgloss.Style
	Cure                    lipgloss.Style
	Ember                   lipgloss.Style
	Attack                  lipgloss.Style
	Defense                 lipgloss.Style
	Vitality                lipgloss.Style
	Sight                   lipgloss.Style
	MenuItem                lipgloss.Style
	MenuSelected            lipgloss.Style
	MenuCursor              lipgloss.Style
	ListSelected            lipgloss.Style
	Footer                  lipgloss.Style
	ModalChoice             lipgloss.Style
	ModalChoiceActive       lipgloss.Style
	Void                    lipgloss.Style
	TileFloorVisible        lipgloss.Style
	TileFloorSeen           lipgloss.Style
	TileFloorClearedVisible lipgloss.Style
	TileFloorClearedSeen    lipgloss.Style
	TileWallVisible         lipgloss.Style
	TileWallSeen            lipgloss.Style
	TileWallClearedVisible  lipgloss.Style
	TileWallClearedSeen     lipgloss.Style
	BossFloorVisible        lipgloss.Style
	BossFloorSeen           lipgloss.Style
}

func newStyles() styles {
	frame := lipgloss.Color("#4d4039")
	frameActive := lipgloss.Color("#8fa7bf")
	frameDanger := lipgloss.Color("#9f5151")
	accent := lipgloss.Color("#d39d62")
	accentSoft := lipgloss.Color("#e6c79a")
	muted := lipgloss.Color("#b6ab9d")
	dim := lipgloss.Color("#766c64")
	text := lipgloss.Color("#f1e7db")

	return styles{
		App: lipgloss.NewStyle().Foreground(text),
		HeaderBar: lipgloss.NewStyle().
			Foreground(accentSoft).
			Background(lipgloss.Color("#120f0e")).
			Bold(true).
			Align(lipgloss.Center).
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(frame),
		Title:    lipgloss.NewStyle().Foreground(accent).Bold(true),
		Subtitle: lipgloss.NewStyle().Foreground(muted).Bold(true),
		Panel: lipgloss.NewStyle().
			Foreground(text).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(frame).
			Background(lipgloss.Color("#110f0e")).
			Padding(0, 1),
		PanelActive: lipgloss.NewStyle().
			Foreground(text).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(frameActive).
			Background(lipgloss.Color("#110f0e")).
			Padding(0, 1),
		ModalPanel: lipgloss.NewStyle().
			Foreground(text).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(frameDanger).
			Background(lipgloss.Color("#100d0d")).
			Padding(0, 1),
		PanelTitle:       lipgloss.NewStyle().Foreground(accentSoft).Bold(true),
		PanelTitleActive: lipgloss.NewStyle().Foreground(lipgloss.Color("#d8e8f6")).Bold(true),
		PanelTitleModal:  lipgloss.NewStyle().Foreground(lipgloss.Color("#f0d2d2")).Bold(true),
		Muted:            lipgloss.NewStyle().Foreground(muted),
		Dim:              lipgloss.NewStyle().Foreground(dim),
		Accent:           lipgloss.NewStyle().Foreground(accent).Bold(true),
		AccentSoft:       lipgloss.NewStyle().Foreground(accentSoft).Bold(true),
		Gold:             lipgloss.NewStyle().Foreground(lipgloss.Color("#dfbe62")).Bold(true),
		Success:          lipgloss.NewStyle().Foreground(lipgloss.Color("#88ba7a")).Bold(true),
		Danger:           lipgloss.NewStyle().Foreground(lipgloss.Color("#d47777")).Bold(true),
		Info:             lipgloss.NewStyle().Foreground(lipgloss.Color("#88b2ce")).Bold(true),
		Warning:          lipgloss.NewStyle().Foreground(lipgloss.Color("#d7a16d")).Bold(true),
		PanelNote:        lipgloss.NewStyle().Foreground(dim),
		Quantity:         lipgloss.NewStyle().Foreground(lipgloss.Color("#9cc2e0")).Bold(true),
		CompareBetter:    lipgloss.NewStyle().Foreground(lipgloss.Color("#88ba7a")).Bold(true),
		CompareEqual:     lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa7c5")).Bold(true),
		CompareWorse:     lipgloss.NewStyle().Foreground(lipgloss.Color("#d47777")).Bold(true),
		Heal:             lipgloss.NewStyle().Foreground(lipgloss.Color("#88ba7a")).Bold(true),
		Focus:            lipgloss.NewStyle().Foreground(lipgloss.Color("#e5c67c")).Bold(true),
		Cure:             lipgloss.NewStyle().Foreground(lipgloss.Color("#8dd7d8")).Bold(true),
		Ember:            lipgloss.NewStyle().Foreground(lipgloss.Color("#ef9852")).Bold(true),
		Attack:           lipgloss.NewStyle().Foreground(lipgloss.Color("#ef9852")).Bold(true),
		Defense:          lipgloss.NewStyle().Foreground(lipgloss.Color("#7da6c7")).Bold(true),
		Vitality:         lipgloss.NewStyle().Foreground(lipgloss.Color("#88ba7a")).Bold(true),
		Sight:            lipgloss.NewStyle().Foreground(lipgloss.Color("#efcb76")).Bold(true),
		MenuItem:         lipgloss.NewStyle().Foreground(muted),
		MenuSelected:     lipgloss.NewStyle().Foreground(text).Bold(true),
		MenuCursor:       lipgloss.NewStyle().Foreground(accentSoft).Bold(true),
		ListSelected:     lipgloss.NewStyle().Foreground(text).Bold(true).Underline(true),
		Footer:           lipgloss.NewStyle().Foreground(muted),
		ModalChoice: lipgloss.NewStyle().
			Foreground(muted).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(frame).
			Padding(0, 1),
		ModalChoiceActive: lipgloss.NewStyle().
			Foreground(text).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(frameActive).
			Bold(true).
			Padding(0, 1),
		Void:                    lipgloss.NewStyle().Background(lipgloss.Color("#080707")),
		TileFloorVisible:        lipgloss.NewStyle().Background(lipgloss.Color("#151211")),
		TileFloorSeen:           lipgloss.NewStyle().Background(lipgloss.Color("#0d0b0a")),
		TileFloorClearedVisible: lipgloss.NewStyle().Background(lipgloss.Color("#101722")),
		TileFloorClearedSeen:    lipgloss.NewStyle().Background(lipgloss.Color("#0a1017")),
		TileWallVisible:         lipgloss.NewStyle().Foreground(lipgloss.Color("#8e8276")).Background(lipgloss.Color("#110f0e")),
		TileWallSeen:            lipgloss.NewStyle().Foreground(lipgloss.Color("#514943")).Background(lipgloss.Color("#090808")),
		TileWallClearedVisible:  lipgloss.NewStyle().Foreground(lipgloss.Color("#5b7389")).Background(lipgloss.Color("#101722")),
		TileWallClearedSeen:     lipgloss.NewStyle().Foreground(lipgloss.Color("#334556")).Background(lipgloss.Color("#0a1017")),
		BossFloorVisible:        lipgloss.NewStyle().Background(lipgloss.Color("#1d0f10")),
		BossFloorSeen:           lipgloss.NewStyle().Background(lipgloss.Color("#130909")),
	}
}

func (s styles) box(title string, body string, width int, height int) string {
	return s.renderBox(s.Panel, s.PanelTitle, title, body, width, height)
}

func (s styles) focusBox(title string, body string, width int, height int) string {
	return s.renderBox(s.PanelActive, s.PanelTitleActive, title, body, width, height)
}

func (s styles) modalBox(title string, body string, width int) string {
	return s.renderBox(s.ModalPanel, s.PanelTitleModal, title, body, width, 0)
}

func (s styles) renderBox(boxStyle lipgloss.Style, titleStyle lipgloss.Style, title string, body string, width int, height int) string {
	innerWidth := max(12, width-4)
	header := titleStyle.Render(strings.ToUpper(title))
	contentStyle := lipgloss.NewStyle().Width(innerWidth)
	if height > 0 {
		contentStyle = contentStyle.Height(max(1, height-3))
	}
	content := contentStyle.Render(body)
	rendered := lipgloss.JoinVertical(lipgloss.Left, header, content)
	return boxStyle.Copy().Width(width).Render(rendered)
}

func (s styles) colorGlyph(glyph string, tint string, bold bool) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(tint))
	if bold {
		style = style.Bold(true)
	}
	return style.Render(glyph)
}

func (s styles) cellGlyph(base lipgloss.Style, glyph string, tint string, bold bool) string {
	style := base.Copy().Foreground(lipgloss.Color(tint))
	if bold {
		style = style.Bold(true)
	}
	return style.Render(glyph)
}

func (s styles) rarityStyle(rarity game.Rarity) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(rarity.Tint())).Bold(rarity >= game.RarityRare)
}

func (s styles) keyStyle(tier game.KeyTier) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(tier.Tint())).Bold(true)
}
