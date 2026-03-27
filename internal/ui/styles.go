package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
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
}

func newStyles() styles {
	frame := lipgloss.Color("#55483f")
	frameActive := lipgloss.Color("#9bc1d8")
	accent := lipgloss.Color("#da9b53")
	accentSoft := lipgloss.Color("#e7c89a")
	muted := lipgloss.Color("#b6ab9d")
	dim := lipgloss.Color("#786f66")
	text := lipgloss.Color("#ece4d9")

	return styles{
		App: lipgloss.NewStyle().Foreground(text),
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
			Foreground(muted).
			Bold(true),
		Panel: lipgloss.NewStyle().
			Foreground(text).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(frame).
			Padding(0, 1),
		PanelActive: lipgloss.NewStyle().
			Foreground(text).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(frameActive).
			Padding(0, 1),
		ModalPanel: lipgloss.NewStyle().
			Foreground(text).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(frameActive).
			Padding(0, 1),
		PanelTitle: lipgloss.NewStyle().
			Foreground(accentSoft).
			Bold(true),
		PanelTitleActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d6e8f1")).
			Bold(true),
		PanelTitleModal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d6e8f1")).
			Bold(true),
		Muted: lipgloss.NewStyle().Foreground(muted),
		Dim:   lipgloss.NewStyle().Foreground(dim),
		Accent: lipgloss.NewStyle().
			Foreground(accent).
			Bold(true),
		AccentSoft: lipgloss.NewStyle().
			Foreground(accentSoft).
			Bold(true),
		Gold: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d9b861")).
			Bold(true),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#84b870")).
			Bold(true),
		Danger: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d17081")).
			Bold(true),
		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8fb7cf")).
			Bold(true),
		PanelNote: lipgloss.NewStyle().
			Foreground(dim),
		Quantity: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8fb7cf")).
			Bold(true),
		CompareBetter: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#84b870")).
			Bold(true),
		CompareEqual: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#72a1c1")).
			Bold(true),
		CompareWorse: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d17081")).
			Bold(true),
		Heal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#84b870")).
			Bold(true),
		Focus: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f0c77d")).
			Bold(true),
		Cure: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8fd2d8")).
			Bold(true),
		Ember: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ef9852")).
			Bold(true),
		Attack: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ef9852")).
			Bold(true),
		Defense: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#72a1c1")).
			Bold(true),
		Vitality: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#84b870")).
			Bold(true),
		Sight: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f0c77d")).
			Bold(true),
		MenuItem: lipgloss.NewStyle().
			Foreground(muted),
		MenuSelected: lipgloss.NewStyle().
			Foreground(text).
			Bold(true),
		MenuCursor: lipgloss.NewStyle().
			Foreground(accentSoft).
			Bold(true),
		ListSelected: lipgloss.NewStyle().
			Foreground(text).
			Bold(true).
			Underline(true),
		Footer: lipgloss.NewStyle().
			Foreground(muted),
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
		Void:                    lipgloss.NewStyle(),
		TileFloorVisible:        lipgloss.NewStyle(),
		TileFloorSeen:           lipgloss.NewStyle(),
		TileFloorClearedVisible: lipgloss.NewStyle(),
		TileFloorClearedSeen:    lipgloss.NewStyle(),
		TileWallVisible: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8f8579")),
		TileWallSeen: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#514b45")),
		TileWallClearedVisible: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#55708a")),
		TileWallClearedSeen: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#334658")),
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

func (s styles) colorGlyph(glyph rune, tint string, bold bool) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(tint))
	if bold {
		style = style.Bold(true)
	}
	return style.Render(string(glyph))
}

func (s styles) cellGlyph(base lipgloss.Style, glyph string, tint string, bold bool) string {
	style := base.Copy().Foreground(lipgloss.Color(tint))
	if bold {
		style = style.Bold(true)
	}
	return style.Render(glyph)
}
