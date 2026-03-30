package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"dunshell/internal/game"
)

var titleArt = []string{
	"  ▄▄▄▄▄▄                                ▄▄ ▄▄ ",
	" █▀██▀▀██                    █▄          ██ ██",
	"   ██   ██       ▄           ██          ██ ██",
	"   ██   ██ ██ ██ ████▄ ▄██▀█ ████▄ ▄█▀█▄ ██ ██",
	" ▄ ██   ██ ██ ██ ██ ██ ▀███▄ ██ ██ ██▄█▀ ██ ██",
	" ▀██▀███▀ ▄▀██▀█▄██ ▀██▄▄██▀▄██ ██▄▀█▄▄▄▄██▄██",
}

func (m *Model) renderScreen(content string, width int, height int, hPos lipgloss.Position, vPos lipgloss.Position) string {
	footer := m.renderFooter(width)
	availableHeight := max(1, height-lipgloss.Height(footer))
	body := lipgloss.Place(width, availableHeight, hPos, vPos, content)
	return lipgloss.JoinVertical(lipgloss.Left, body, footer)
}

func (m *Model) gameplayMetrics() (int, int, int, int, int) {
	width := max(100, m.width)
	height := max(34, m.height)
	footerHeight := lipgloss.Height(m.renderFooter(width))
	availableHeight := max(24, height-footerHeight)
	logHeight := clamp(height/4, 7, 10)
	if availableHeight-logHeight < 15 {
		logHeight = max(6, availableHeight-15)
	}
	bodyHeight := max(15, availableHeight-logHeight)
	sidebarWidth := clamp(width/3, 38, 46)
	return width, height, bodyHeight, logHeight, sidebarWidth
}

func (m *Model) viewTitle() string {
	width := max(88, m.width)
	height := max(30, m.height)
	titleLines := make([]string, 0, len(titleArt))
	for _, line := range titleArt {
		titleLines = append(titleLines, m.styles.Title.Render(line))
	}
	options := m.titleMenuOptions()
	menuLines := make([]string, 0, len(options))
	for index, label := range options {
		line := "  " + m.styles.MenuItem.Render(label)
		if index == m.titleMenuIndex {
			line = m.renderSelectedText("› "+label, 30)
		}
		menuLines = append(menuLines, line)
	}
	seedLine := m.styles.Muted.Render("Seed flow: random or manual in the new-run screen")
	if m.hasLockedSeed {
		seedLine = m.styles.Info.Render("CLI seed locked to " + fmt.Sprintf("%d", m.lockedSeed))
	}
	modeLine := m.styles.Muted.Render("New runs start in standard mode")
	if m.cliGodMode {
		modeLine = m.styles.Warning.Render("CLI god mode armed for new runs")
	}
	continueLine := ""
	if m.hasContinue && m.savedRun != nil {
		continueLine = m.styles.AccentSoft.Render(fmt.Sprintf("Continue floor %d  ·  level %d  ·  seed %d", m.savedRun.FloorIndex, m.savedRun.Player.Level, m.savedRun.Seed))
		if m.savedRun.GodMode {
			continueLine += "  ·  " + m.styles.Warning.Render("GOD MODE")
		}
	}
	metaLine := m.styles.Warning.Render(fmt.Sprintf("Omen tier %d  ·  wins %d", m.profile.Difficulty, m.profile.Wins))
	copy := wrapText("A dark-fantasy terminal roguelike of sealed boss chambers, route-mapped descents, keyed reliquaries, merchants, escalating minibosses, and a crown the abbey still kills to keep.", 72)
	lines := []string{
		"",
		copy,
		"",
	}
	if continueLine != "" {
		lines = append(lines, continueLine, "")
	}
	lines = append(lines,
		strings.Join(menuLines, "\n"),
		"",
		seedLine,
		modeLine,
		metaLine,
		m.styles.PanelNote.Render("Arrows or W/S move through menus. Enter confirms. Settings now controls glyph mode and descend prompts."),
	)
	if m.storageError != "" {
		lines = append(lines, "", m.styles.Danger.Render("Storage: "+m.storageError))
	}
	panel := m.styles.focusBox("Start Menu", strings.Join(lines, "\n"), 80, 0)
	content := lipgloss.JoinVertical(lipgloss.Center, strings.Join(titleLines, "\n"), "", panel)
	return m.renderScreen(content, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) viewSeed() string {
	width := max(84, m.width)
	height := max(28, m.height)
	randomLabel := m.styles.ModalChoice.Render(" Random Seed ")
	manualLabel := m.styles.ModalChoice.Render(" Manual Seed ")
	if m.seedMode == 0 {
		randomLabel = m.styles.ModalChoiceActive.Render(" Random Seed ")
	} else {
		manualLabel = m.styles.ModalChoiceActive.Render(" Manual Seed ")
	}
	input := m.styles.Panel.Render(m.seedInput.View())
	if m.seedMode == 0 {
		input = m.styles.Panel.Copy().Foreground(lipgloss.Color("#766c64")).Render("randomized on start")
	}
	lines := []string{
		wrapText("Choose how the next descent is born. Manual seeds accept either numbers or text and are hashed into a replayable run seed.", 56),
		"",
		lipgloss.JoinHorizontal(lipgloss.Left, randomLabel, "  ", manualLabel),
		"",
		m.styles.Subtitle.Render("Seed"),
		input,
		"",
	}
	if m.cliGodMode {
		lines = append(lines, m.styles.Warning.Render("Developer run: GOD MODE will be active from floor 1."), "")
	}
	lines = append(lines,
		m.styles.PanelNote.Render("Left/Right switch mode • Enter starts • Esc returns"),
	)
	if m.seedError != "" {
		lines = append(lines, "", m.styles.Danger.Render(m.seedError))
	}
	panel := m.styles.focusBox("New Run", strings.Join(lines, "\n"), 64, 0)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) viewHelp() string {
	width := max(102, m.width)
	height := max(36, m.height)
	saveDir, err := game.SaveDirectory()
	saveLine := saveDir
	if err != nil {
		saveLine = "unavailable: " + err.Error()
	}
	legend := []string{
		m.styles.cellGlyph(m.styles.TileFloorVisible, m.glyphs.player(), "#84d06f", true) + " You",
		m.styles.TileWallVisible.Render("█") + " Wall",
		m.styles.TileWallClearedVisible.Render("█") + " Cleared room wall",
		m.styles.cellGlyph(m.styles.TileFloorVisible, m.glyphs.bossGate(), "#b55e5e", true) + " Boss gate",
		m.styles.cellGlyph(m.styles.TileFloorVisible, m.glyphs.stairs(), "#9bc1d8", true) + " Stairs down",
		m.styles.colorGlyph(m.glyphs.merchant(), "#d5b36d", true) + " Merchant",
		m.styles.colorGlyph(m.glyphs.chest(game.KeyBronze), game.KeyBronze.Tint(), true) + " Bronze chest",
		m.styles.colorGlyph(m.glyphs.chest(game.KeySilver), game.KeySilver.Tint(), true) + " Silver chest",
		m.styles.colorGlyph(m.glyphs.chest(game.KeyGold), game.KeyGold.Tint(), true) + " Gold chest",
		m.styles.colorGlyph(m.glyphs.symbol('⚿', 'k'), game.KeyGold.Tint(), true) + " Keys",
	}
	body := strings.Join([]string{
		m.renderHelpSection("Run Flow", []string{
			"Descend twenty floors, face minibosses on floors 5, 10, and 15, then break the Ashen Prior on floor 20.",
			"Stairs now lead into a route-choice map before the next floor generates.",
			"Victory unlocks endless mode and permanently increases the omen tier for future runs.",
		}),
		"",
		m.renderHelpSection("Core Controls", []string{
			"Arrows or W/A/S/D move. Walking into an enemy attacks it.",
			"E interacts with your tile or nearby boss gate. C drinks the weakest healing item.",
			"I opens inventory. P opens settings. Esc closes the active overlay. Q opens the safe quit prompt.",
		}),
		"",
		m.renderHelpSection("Overlays", []string{
			"Chest prompts ask for the key spend without revealing the contents first.",
			"Merchants stock five offers: healing, weapon, armor, charm, and one extra slot.",
			"Boss entry prompts warn that the room will seal until the keeper dies.",
			"Settings persists glyph mode, ASCII fallback, descend confirmation, and message log length.",
		}),
		"",
		m.renderHelpSection("Legend", legend),
		"",
		m.renderHelpSection("Field Notes", []string{
			"Rarity colors: Common, Uncommon, Rare, Legendary, Unique.",
			"Nerd Font is recommended for the full glyph set. DUNSHELL_ASCII=1 still forces ASCII as an override.",
			"Auto-save location: " + saveLine,
			m.styles.PanelNote.Render("Press Esc, Enter, or ? to return."),
		}),
	}, "\n")
	panel := m.styles.box("Field Guide", body, min(width-4, 108), 0)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) viewSettings() string {
	width := max(96, m.width)
	height := max(32, m.height)
	settings := m.settingsDraft.Normalized()
	previewGlyphs := newGlyphSet(settings)
	rows := []string{
		m.renderSettingLine(0, "Glyph Mode", settings.GlyphMode.Label()),
		m.renderSettingLine(1, "ASCII Fallback", boolLabel(settings.ASCIIFallback)),
		m.renderSettingLine(2, "Confirm Descend", boolLabel(settings.ConfirmBeforeDescend)),
		m.renderSettingLine(3, "Message Log", fmt.Sprintf("%d lines", settings.MessageLogLines)),
		"",
		m.renderSettingsActions(),
		m.styles.PanelNote.Render("Up/Down choose • Left/Right change • Enter confirms • Esc returns"),
	}
	preview := []string{
		m.styles.Subtitle.Render("Preview"),
		previewGlyphs.player() + " you   " + previewGlyphs.stairs() + " stair   " + previewGlyphs.bossGate() + " gate",
		previewGlyphs.merchant() + " merchant   " + previewGlyphs.chest(game.KeySilver) + " chest   " + previewGlyphs.symbol('✠', 'X') + " unique",
	}
	if description := m.settingsDescription(settings); description != "" {
		preview = append(preview, "", m.styles.Subtitle.Render("Selection"), wrapText(description, 40))
	}
	if previewGlyphs.ForcedASCII() {
		preview = append(preview, "", m.styles.Warning.Render("DUNSHELL_ASCII=1 is active and currently forces ASCII output."))
	}
	left := m.styles.focusBox("Settings", strings.Join(rows, "\n"), 42, 0)
	right := m.styles.box("Detail", strings.Join(preview, "\n"), 48, 0)
	content := lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right)
	if width < 112 {
		content = lipgloss.JoinVertical(lipgloss.Left, left, "", right)
	}
	return m.renderScreen(content, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) renderHelpSection(title string, lines []string) string {
	return strings.Join([]string{m.styles.Subtitle.Render(title), strings.Join(lines, "\n")}, "\n")
}

func (m *Model) viewGame() string {
	if m.game == nil {
		return ""
	}
	width, height, bodyHeight, logHeight, sidebarWidth := m.gameplayMetrics()
	m.mapViewport = mapViewportState{}
	bodyOriginY := 0
	var body string
	if width >= 120 {
		mapWidth := max(54, width-sidebarWidth-1)
		mapPanel := m.renderMapPanel(mapWidth, bodyHeight, 0, bodyOriginY)
		sidePanel := m.renderSidebar(sidebarWidth, bodyHeight)
		body = lipgloss.JoinHorizontal(lipgloss.Top, mapPanel, " ", sidePanel)
	} else {
		sidebarHeight := max(11, bodyHeight-1)
		mapHeight := max(10, bodyHeight-sidebarHeight)
		mapPanel := m.renderMapPanel(width, mapHeight, 0, bodyOriginY)
		sidePanel := m.renderSidebar(width, sidebarHeight)
		body = lipgloss.JoinVertical(lipgloss.Left, mapPanel, sidePanel)
	}
	logPanel := m.renderLog(width, logHeight)
	content := lipgloss.JoinVertical(lipgloss.Left, body, logPanel)
	return m.renderScreen(content, width, height, lipgloss.Left, lipgloss.Top)
}

func (m *Model) renderFooter(width int) string {
	left := m.styles.Footer.Render("? Help  ·  q Quit")
	right := m.styles.AccentSoft.Render(game.GameTitle + " " + game.GameVersion)
	if m.game != nil && m.game.GodMode {
		right += m.styles.Dim.Render("  ·  ") + m.styles.Warning.Render("GOD MODE")
	}
	innerWidth := max(1, width-2)
	if lipgloss.Width(left)+lipgloss.Width(right)+1 > innerWidth {
		rightPlain := game.GameTitle + " " + game.GameVersion
		if m.game != nil && m.game.GodMode {
			rightPlain += " · GOD MODE"
		}
		available := max(1, innerWidth-lipgloss.Width(left)-1)
		right = m.styles.AccentSoft.Render(truncateText(rightPlain, available))
	}
	spacer := strings.Repeat(" ", max(1, innerWidth-lipgloss.Width(left)-lipgloss.Width(right)))
	line := left + spacer + right
	return m.styles.Footer.Copy().
		Width(width).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color("#4d4039")).
		Render(line)
}

func (m *Model) viewDescendPrompt() string {
	width := max(68, m.width)
	height := max(24, m.height)
	completion := m.game.Floor.Completion()
	missing := make([]string, 0, 4)
	if !completion.FullyExplored() {
		missing = append(missing, m.styles.Info.Render(fmt.Sprintf("• %d tiles still hide in shadow.", completion.UnexploredTiles)))
	}
	if !completion.LootCollected() {
		missing = append(missing, m.styles.Gold.Render(fmt.Sprintf("• %d spoils still lie unopened.", completion.RemainingItems)))
	}
	if !completion.EnemiesCleared() {
		missing = append(missing, m.styles.Danger.Render(fmt.Sprintf("• %d foes still draw breath.", completion.RemainingEnemies)))
	}
	if completion.UnclearedRooms() > 0 {
		missing = append(missing, m.styles.AccentSoft.Render(fmt.Sprintf("• %d chambers remain uncleared.", completion.UnclearedRooms())))
	}
	copy := "The stair no longer drops you straight into the next floor. It opens a route map first."
	if completion.Complete() {
		copy = "This floor lies spent. The next step will reveal a route map and let you choose the shape of the descent."
	}
	body := []string{wrapText(copy, 44)}
	if len(missing) > 0 {
		body = append(body, "", strings.Join(missing, "\n"))
	}
	body = append(body, "", m.renderBinaryChoice(m.descendChoice, "Take the route map", "Stay"), m.styles.PanelNote.Render("Left/Right choose • Enter confirms • Esc returns"))
	panel := m.styles.modalBox("Stair Mouth", strings.Join(body, "\n"), 54)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) viewRouteChoice() string {
	width := m.width
	if width <= 0 {
		width = 96
	}
	height := m.height
	if height <= 0 {
		height = 30
	}
	routes := m.game.RouteChoices()
	if len(routes) == 0 {
		return m.viewGame()
	}
	selectedIndex := clamp(m.routeCursor, 0, len(routes)-1)
	selected := routes[selectedIndex]
	contentWidth := width - 2
	if contentWidth <= 0 {
		contentWidth = 94
	}
	if contentWidth > 118 {
		contentWidth = 118
	}

	titleText := selected.MapLabel + "  ·  choose the next descent"
	if selected.BossFloor {
		titleText = selected.MapLabel + "  ·  keeper ahead"
	}
	title := m.styles.Accent.Render(truncateText(titleText, contentWidth))
	help := m.styles.Footer.Render(truncateText("Up/Down choose branch • Left/Right also cycle • Enter descend • Esc return", contentWidth))

	graphInnerHeight := clamp(height-10, 11, 18)
	if contentWidth < 74 {
		graphInnerHeight = clamp(height-15, 8, 12)
	}
	graphPanelHeight := graphInnerHeight + 3

	var panels string
	if contentWidth >= 74 {
		graphWidth := clamp(contentWidth*58/100, 36, 60)
		if graphWidth > contentWidth-26 {
			graphWidth = max(28, contentWidth-26)
		}
		detailWidth := max(24, contentWidth-graphWidth-2)
		graphPanel := m.styles.focusBox("Route Graph", m.renderRouteGraph(routes, selectedIndex, max(12, graphWidth-4), graphInnerHeight), graphWidth, graphPanelHeight)

		var sidebar string
		if contentWidth >= 104 && height >= 32 && detailWidth >= 30 {
			detailPanelHeight := max(10, graphPanelHeight*3/5)
			summaryPanelHeight := max(8, graphPanelHeight-detailPanelHeight-1)
			detailPanel := m.styles.box("Selected Route", fixedPanelBody(m.routeDetailLines(selected, false), max(12, detailWidth-4), max(1, detailPanelHeight-3)), detailWidth, detailPanelHeight)
			summaryPanel := m.styles.box("Rewards And Omens", fixedPanelBody(m.routeSignalLines(selected), max(12, detailWidth-4), max(1, summaryPanelHeight-3)), detailWidth, summaryPanelHeight)
			sidebar = lipgloss.JoinVertical(lipgloss.Left, detailPanel, " ", summaryPanel)
		} else {
			sidebar = m.styles.box("Selected Route", fixedPanelBody(m.routeDetailLines(selected, true), max(12, detailWidth-4), max(1, graphPanelHeight-3)), detailWidth, graphPanelHeight)
		}
		panels = lipgloss.JoinHorizontal(lipgloss.Top, graphPanel, "  ", sidebar)
	} else {
		detailPanelHeight := clamp(height-graphPanelHeight-4, 7, 12)
		graphPanel := m.styles.focusBox("Route Graph", m.renderRouteGraph(routes, selectedIndex, max(12, contentWidth-4), graphInnerHeight), contentWidth, graphPanelHeight)
		detailPanel := m.styles.box("Selected Route", fixedPanelBody(m.routeDetailLines(selected, true), max(12, contentWidth-4), max(1, detailPanelHeight-3)), contentWidth, detailPanelHeight)
		panels = lipgloss.JoinVertical(lipgloss.Left, graphPanel, " ", detailPanel)
	}

	sections := []string{title}
	if contentWidth >= 74 {
		sections = append(sections, m.styles.PanelNote.Render(wrapText("Branch lines show how the current stair can split. Route text stays outside the graph.", max(24, contentWidth))))
	}
	sections = append(sections, "", panels, "", help)
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return m.renderScreen(content, width, height, lipgloss.Center, lipgloss.Top)
}

func (m *Model) viewChestPrompt() string {
	width := max(72, m.width)
	height := max(26, m.height)
	chest, _ := m.game.ChestAtPlayer()
	if chest == nil {
		return m.viewGame()
	}
	keyCount := m.game.Player.Keys.Count(chest.Tier)
	state := m.styles.Success.Render("Ready")
	if chest.Locked {
		state = m.styles.Danger.Render("Locked until the boss falls")
	} else if keyCount == 0 {
		state = m.styles.Warning.Render("Missing matching key")
	}
	body := []string{
		m.styles.keyStyle(chest.Tier).Render(chest.Tier.Label()+" Chest") + "  ·  " + state,
		"",
		wrapText("Break the seal with a matching key. The reliquary keeps its contents hidden until you commit.", 48),
		"",
		m.styles.Muted.Render(fmt.Sprintf("%s keys in ring: %d", chest.Tier.Label(), keyCount)),
		"",
		m.renderBinaryChoice(m.chestChoice, "Open chest", "Leave"),
		m.styles.PanelNote.Render("Left/Right choose • Enter confirms • Esc returns"),
	}
	panel := m.styles.modalBox("Reliquary", strings.Join(body, "\n"), 58)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) viewBossPrompt() string {
	width := max(74, m.width)
	height := max(28, m.height)
	boss := m.game.BossPreview()
	if boss == nil {
		return m.viewGame()
	}
	body := []string{
		m.styles.Danger.Render(boss.Template.Name),
		boss.Template.Description,
		"",
		fmt.Sprintf("HP %d   ATK %d   DEF %d", boss.Template.MaxHP, boss.AttackPower(), boss.DefensePower()),
		"",
		wrapText("Entering will seal the chamber. You cannot leave until the keeper is dead. The reward chest will only unlock after victory.", 48),
		"",
		m.renderBinaryChoice(m.bossChoice, "Enter chamber", "Withdraw"),
		m.styles.PanelNote.Render("Left/Right choose • Enter confirms • Esc returns"),
	}
	panel := m.styles.modalBox("Boss Gate", strings.Join(body, "\n"), 58)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) viewQuitPrompt() string {
	width := max(60, m.width)
	height := max(22, m.height)
	body := strings.Join([]string{
		wrapText("Leave this run? Auto-save will keep the current descent in place for Continue on the title screen.", 42),
		"",
		m.renderBinaryChoice(m.quitChoice, "Quit to terminal", "Return"),
		m.styles.PanelNote.Render("Left/Right choose • Enter confirms • Esc stays"),
	}, "\n")
	panel := m.styles.modalBox("Safe Quit", body, 48)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) viewOutcome(victory bool) string {
	width := max(88, m.width)
	height := max(30, m.height)
	summary := m.game.Summary()
	title := "Run Extinguished"
	copy := "The abbey keeps another name in its stone throat. The next descent will remember where this one failed."
	accent := m.styles.Danger
	if victory {
		title = "Cinder Crown Claimed"
		copy = "You break the sanctum's last prayer and take the crown. The abbey opens further, as though the theft only pleased it."
		accent = m.styles.Success
	}
	lines := []string{
		accent.Render(title),
		"",
		wrapText(copy, 64),
		"",
	}
	if m.game.GodMode {
		lines = append(lines, "Mode        "+m.styles.Warning.Render("GOD MODE"), "")
	}
	lines = append(lines,
		"Seed        "+fmt.Sprintf("%d", summary.Seed),
		"Floor       "+fmt.Sprintf("%d", summary.Floor),
		"Level       "+fmt.Sprintf("%d", summary.Level),
		"Gold        "+fmt.Sprintf("%d", summary.Gold),
		"Kills       "+fmt.Sprintf("%d", summary.Kills),
		"Turns       "+fmt.Sprintf("%d", summary.Turn),
		"Omen Tier   "+fmt.Sprintf("%d", summary.PersistentDifficulty),
		"",
		m.renderOutcomeOptions(victory),
		"",
		m.styles.PanelNote.Render("Up/Down choose • Enter confirms • n starts a new run"),
	)
	panel := m.styles.modalBox(title, strings.Join(lines, "\n"), 72)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) renderOutcomeOptions(victory bool) string {
	options := m.outcomeMenuOptions()
	lines := make([]string, 0, len(options))
	for index, option := range options {
		line := "  " + option
		if index == m.outcomeMenuIndex {
			line = m.renderSelectedText("› "+option, 30)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (m *Model) renderMapPanel(width int, height int, originX int, originY int) string {
	innerWidth := max(24, width-4)
	innerHeight := max(9, height-4)
	player := m.game.Player
	floor := m.game.Floor
	cameraX := clamp(player.Pos.X-innerWidth/2, 0, max(0, floor.Width-innerWidth))
	cameraY := clamp(player.Pos.Y-innerHeight/2, 0, max(0, floor.Height-innerHeight))
	m.mapViewport = mapViewportState{Panel: rect{X: originX, Y: originY, W: width, H: height}, Content: rect{X: originX + 2, Y: originY + 2, W: innerWidth, H: innerHeight}, CameraX: cameraX, CameraY: cameraY}
	enemyMap := make(map[game.Position]*game.Enemy, len(floor.Enemies))
	for _, enemy := range floor.Enemies {
		if floor.IsVisible(enemy.Pos) {
			enemyMap[enemy.Pos] = enemy
		}
	}
	itemMap := make(map[game.Position]game.GroundItem, len(floor.Items))
	for _, item := range floor.Items {
		if floor.IsVisible(item.Pos) {
			itemMap[item.Pos] = item
		}
	}
	chestMap := make(map[game.Position]game.Chest, len(floor.Chests))
	for _, chest := range floor.Chests {
		if !chest.Opened && floor.IsVisible(chest.Pos) {
			chestMap[chest.Pos] = chest
		}
	}
	merchantMap := make(map[game.Position]game.Merchant, len(floor.Merchants))
	for _, merchant := range floor.Merchants {
		if floor.IsVisible(merchant.Pos) {
			merchantMap[merchant.Pos] = merchant
		}
	}
	roomStates := floor.RoomStates()
	rows := make([]string, innerHeight)
	for y := 0; y < innerHeight; y++ {
		cells := make([]string, innerWidth)
		for x := 0; x < innerWidth; x++ {
			pos := game.Position{X: cameraX + x, Y: cameraY + y}
			cells[x] = m.renderMapCell(pos, enemyMap, itemMap, chestMap, merchantMap, roomStates)
		}
		rows[y] = strings.Join(cells, "")
	}
	return m.styles.box("Dungeon", strings.Join(rows, "\n"), width, height)
}

func (m *Model) renderMapCell(pos game.Position, enemies map[game.Position]*game.Enemy, items map[game.Position]game.GroundItem, chests map[game.Position]game.Chest, merchants map[game.Position]game.Merchant, roomStates []game.RoomState) string {
	floor := m.game.Floor
	if !floor.InBounds(pos) {
		return " "
	}
	visible := floor.IsVisible(pos)
	explored := floor.IsExplored(pos)
	if !explored {
		return " "
	}
	if m.game.Player.Pos.Equals(pos) {
		base := m.floorCellStyle(pos, visible, roomStates)
		return m.styles.cellGlyph(base, m.glyphs.player(), "#84d06f", true)
	}
	if enemy, ok := enemies[pos]; ok {
		base := m.floorCellStyle(pos, visible, roomStates)
		if enemy.Template.BossTier > 0 {
			base = m.styles.BossFloorVisible
			if !visible {
				base = m.styles.BossFloorSeen
			}
		}
		return m.styles.cellGlyph(base, m.glyphs.symbol(enemy.Template.Glyph, enemy.Template.ASCII), enemy.Template.Tint, true)
	}
	if merchant, ok := merchants[pos]; ok {
		base := m.floorCellStyle(pos, visible, roomStates)
		_ = merchant
		return m.styles.cellGlyph(base, m.glyphs.merchant(), "#d9b861", true)
	}
	if chest, ok := chests[pos]; ok {
		base := m.floorCellStyle(pos, visible, roomStates)
		return m.styles.cellGlyph(base, m.glyphs.chest(chest.Tier), chest.Tier.Tint(), true)
	}
	if item, ok := items[pos]; ok {
		base := m.floorCellStyle(pos, visible, roomStates)
		return m.styles.cellGlyph(base, m.glyphs.symbol(item.Item.Glyph, item.Item.ASCII), item.Item.Tint, true)
	}
	tile := floor.TileAt(pos)
	switch tile {
	case game.TileWall:
		if m.wallClearedContext(pos, roomStates) {
			if visible {
				return m.styles.TileWallClearedVisible.Render("█")
			}
			return m.styles.TileWallClearedSeen.Render("▓")
		}
		if visible {
			return m.styles.TileWallVisible.Render("█")
		}
		return m.styles.TileWallSeen.Render("▓")
	case game.TileFloor:
		return m.floorCellStyle(pos, visible, roomStates).Render(m.floorGlyph(pos, visible, roomStates))
	case game.TileDoorClosed:
		base := m.doorCellStyle(pos, visible, roomStates)
		return m.styles.cellGlyph(base, m.glyphs.symbol('+', '+'), "#d4a66d", true)
	case game.TileDoorOpen:
		base := m.doorCellStyle(pos, visible, roomStates)
		return m.styles.cellGlyph(base, m.glyphs.symbol('/', '/'), "#9b7a61", false)
	case game.TileStairsDown:
		base := m.floorCellStyle(pos, visible, roomStates)
		tint := "#9bc1d8"
		if !visible {
			tint = "#5c7686"
		}
		return m.styles.cellGlyph(base, m.glyphs.stairs(), tint, true)
	case game.TileBossGate:
		base := m.floorCellStyle(pos, visible, roomStates)
		return m.styles.cellGlyph(base, m.glyphs.bossGate(), "#d47777", true)
	case game.TileBossSeal:
		base := m.styles.BossFloorVisible
		if !visible {
			base = m.styles.BossFloorSeen
		}
		return m.styles.cellGlyph(base, m.glyphs.bossGate(), "#7c3d3d", true)
	default:
		return " "
	}
}

func (m *Model) floorCellStyle(pos game.Position, visible bool, roomStates []game.RoomState) lipgloss.Style {
	roomIndex := m.game.Floor.RoomIndexAt(pos)
	if roomIndex >= 0 && roomIndex < len(roomStates) && roomStates[roomIndex].Kind == game.RoomBoss && !roomStates[roomIndex].Cleared {
		if visible {
			return m.styles.BossFloorVisible
		}
		return m.styles.BossFloorSeen
	}
	cleared := m.roomClearedContext(pos, roomStates)
	switch {
	case visible && cleared:
		return m.styles.TileFloorClearedVisible
	case visible:
		return m.styles.TileFloorVisible
	case cleared:
		return m.styles.TileFloorClearedSeen
	default:
		return m.styles.TileFloorSeen
	}
}

func (m *Model) doorCellStyle(pos game.Position, visible bool, roomStates []game.RoomState) lipgloss.Style {
	if m.wallClearedContext(pos, roomStates) {
		if visible {
			return m.styles.TileFloorClearedVisible
		}
		return m.styles.TileFloorClearedSeen
	}
	if visible {
		return m.styles.TileFloorVisible
	}
	return m.styles.TileFloorSeen
}

func (m *Model) roomClearedContext(pos game.Position, roomStates []game.RoomState) bool {
	roomIndex := m.game.Floor.RoomIndexAt(pos)
	return roomIndex >= 0 && roomStates[roomIndex].Cleared
}

func (m *Model) wallClearedContext(pos game.Position, roomStates []game.RoomState) bool {
	for _, roomIndex := range m.game.Floor.AdjacentRoomIndices(pos) {
		if roomStates[roomIndex].Cleared {
			return true
		}
	}
	return false
}

func (m *Model) renderSidebar(width int, height int) string {
	switch m.overlay {
	case overlayInventory:
		return m.renderInventorySidebar(width, height)
	case overlayMerchant:
		return m.renderMerchantSidebar(width, height)
	default:
		return m.renderStatusSidebar(width, height)
	}
}

func (m *Model) renderStatusSidebar(width int, height int) string {
	player := m.game.Player
	completion := m.game.Floor.Completion()
	stats := []string{
		m.renderStatBar("HP", player.HP, player.MaxHP(), width-4, lipgloss.Color("#c56b78"), lipgloss.Color("#2d1d23")),
		m.renderStatBar("XP", player.XP, player.NextLevelXP(), width-4, lipgloss.Color("#5d88aa"), lipgloss.Color("#18232d")),
		"",
		"Lvl   " + m.styles.AccentSoft.Render(fmt.Sprintf("%d", player.Level)) + "   Gold " + m.styles.Gold.Render(fmt.Sprintf("%d", player.Gold)),
		"ATK   " + m.styles.Attack.Render(fmt.Sprintf("%d", player.AttackPower())) + "   DEF  " + m.styles.Defense.Render(fmt.Sprintf("%d", player.DefensePower())),
		"Keys  " + m.renderKeyLine(),
	}
	if m.game.GodMode {
		stats = append(stats, "Mode  "+m.styles.Warning.Render("GOD MODE")+"  "+m.styles.Success.Render("Invulnerable"))
	}
	stats = append(stats,
		"State",
	)
	stats = append(stats, m.renderStatusSummary(width-4)...)
	stats = append(stats, m.renderQuickHealPreview())
	floorLines := []string{
		m.styles.Muted.Render(m.game.FloorLabel()),
		m.renderCompletionSummary(completion),
		"Map    " + m.renderBoolMetric(fmt.Sprintf("%d%%", m.game.Floor.ExploredPercent()), completion.FullyExplored()),
		"Rooms  " + m.renderBoolMetric(fmt.Sprintf("%d/%d", completion.ClearedRooms, completion.TotalRooms), completion.UnclearedRooms() == 0),
		"Foes   " + m.renderCountMetric(completion.RemainingEnemies),
		"Loot   " + m.renderCountMetric(completion.RemainingItems),
		"Omen   " + m.styles.Warning.Render(fmt.Sprintf("%d", m.profile.Difficulty)),
		"Under  " + m.styles.AccentSoft.Render(m.game.TileDescriptionUnderPlayer()),
		"",
		wrapText(m.game.Objective(), width-6),
	}
	body := strings.Join([]string{m.styles.Subtitle.Render("Vitals"), strings.Join(stats, "\n"), "", m.styles.Subtitle.Render("Floor"), strings.Join(floorLines, "\n"), "", m.renderBossSection(width - 4)}, "\n")
	return m.styles.box("Status", body, width, height)
}

func (m *Model) renderMerchantSidebar(width int, height int) string {
	merchant, _ := m.game.MerchantAtPlayer()
	if merchant == nil {
		return m.renderStatusSidebar(width, height)
	}
	lines := make([]string, 0, len(merchant.Offers)+4)
	for index, offer := range merchant.Offers {
		price := m.styles.Gold.Render(fmt.Sprintf("%dg", offer.Price))
		name := m.renderItemName(offer.Item)
		line := m.glyphs.symbol(offer.Item.Glyph, offer.Item.ASCII) + " " + name + " " + m.styles.Dim.Render("·") + " " + price
		if offer.Sold {
			line = m.styles.Dim.Render("sold · ") + line
		}
		if index == clamp(m.merchantCursor, 0, len(merchant.Offers)-1) {
			line = m.renderSelectedText(ansi.Strip("› "+line), width-4)
		} else {
			line = "  " + truncateText(line, width-6)
		}
		lines = append(lines, line)
	}
	selected := merchant.Offers[clamp(m.merchantCursor, 0, len(merchant.Offers)-1)]
	body := strings.Join([]string{
		m.styles.Subtitle.Render(merchant.Name),
		"Gold " + m.styles.Gold.Render(fmt.Sprintf("%d", m.game.Player.Gold)),
		"",
		strings.Join(lines, "\n"),
		"",
		m.styles.rarityStyle(selected.Item.Rarity).Render(selected.Item.Rarity.Label()) + "  " + selected.Item.Description,
		m.renderPackItemDetail(selected.Item),
		"",
		m.styles.PanelNote.Render("Up/Down choose • Enter buys • Esc returns"),
	}, "\n")
	return m.styles.focusBox("Merchant", body, width, height)
}

func (m *Model) renderBossSection(width int) string {
	boss := m.game.ActiveBoss()
	if boss == nil {
		return ""
	}
	return strings.Join([]string{
		m.styles.Subtitle.Render("Boss"),
		m.styles.Danger.Render(boss.Template.Name),
		m.renderStatBar("HP", boss.HP, boss.Template.MaxHP, width, lipgloss.Color("#d95e5e"), lipgloss.Color("#321718")),
		boss.Template.Description,
	}, "\n")
}

func (m *Model) renderLog(width int, height int) string {
	innerWidth := max(24, width-6)
	maxLines := max(1, height-3)
	lines := make([]string, 0, maxLines)
	start := max(0, len(m.game.Log)-m.profile.Settings.MessageLogLines)
	for _, entry := range m.game.Log[start:] {
		wrapped := wrapText(entry, innerWidth-2)
		parts := strings.Split(wrapped, "\n")
		for partIndex, part := range parts {
			prefix := "  "
			if partIndex == 0 {
				prefix = m.styles.AccentSoft.Render("• ")
			}
			lines = append(lines, prefix+part)
		}
	}
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}
	body := strings.Join(lines, "\n")
	if body == "" {
		body = m.styles.Dim.Render("The abbey listens.")
	}
	return m.styles.box("Whispers", body, width, height)
}

func (m *Model) renderInventorySidebar(width int, height int) string {
	packHeight, detailsHeight := m.inventoryPanelHeights(height)
	packPanel := m.renderPackPanel(width, packHeight)
	detailsPanel := m.renderInventoryDetailsPanel(width, detailsHeight)
	return lipgloss.JoinVertical(lipgloss.Left, packPanel, detailsPanel)
}

func (m *Model) renderPackPanel(width int, height int) string {
	stacks := m.inventoryStacks()
	if len(stacks) == 0 {
		return m.styles.focusBox("Pack", m.styles.Dim.Render("Your pack is empty."), width, height)
	}
	rows := m.inventoryRows(stacks)
	visibleRows := max(1, height-4)
	start := clamp(m.inventoryPackScroll, 0, max(0, len(rows)-visibleRows))
	end := min(len(rows), start+visibleRows)
	lines := make([]string, 0, end-start)
	for _, row := range rows[start:end] {
		if row.StackIndex == -1 {
			lines = append(lines, m.styles.Subtitle.Render(strings.ToUpper(row.Category)))
			continue
		}
		selected := m.inventoryPane == inventoryPanePack && row.StackIndex == clamp(m.inventoryPackCursor, 0, len(stacks)-1)
		lines = append(lines, m.renderPackLine(stacks[row.StackIndex], selected, width-4))
	}
	return m.styles.focusBox("Pack", strings.Join(lines, "\n"), width, height)
}

func (m *Model) renderInventoryDetailsPanel(width int, height int) string {
	lines := make([]string, 0, 6)
	for _, slot := range []game.EquipmentSlot{game.SlotWeapon, game.SlotArmor, game.SlotCharm} {
		selected := m.inventoryPane == inventoryPaneEquipped && m.selectedEquipmentSlot() == slot
		rendered, _ := m.renderEquippedSlot(slot, selected, true, width-4)
		lines = append(lines, rendered...)
	}
	return m.styles.box("Equipped", strings.Join(lines, "\n"), width, height)
}

func (m *Model) renderPackLine(stack inventoryStack, selected bool, width int) string {
	quantity := m.styles.Quantity.Render(fmt.Sprintf("x%d", stack.Count))
	left := m.glyphs.symbol(stack.Item.Glyph, stack.Item.ASCII) + " " + m.renderItemName(stack.Item)
	if detail := m.renderPackItemDetail(stack.Item); detail != "" {
		left += " " + m.styles.Dim.Render("·") + " " + detail
	}
	prefix := "  "
	if selected {
		prefix = "› "
	}
	available := max(1, width-lipgloss.Width(prefix)-lipgloss.Width(quantity)-1)
	left = truncateText(left, available)
	line := prefix + padRight(left, available) + " " + quantity
	if selected {
		return m.renderSelectedText(ansi.Strip(line), width)
	}
	return lipgloss.NewStyle().Width(width).Render(line)
}

func (m *Model) renderEquippedSlot(slot game.EquipmentSlot, selected bool, detailed bool, width int) ([]string, int) {
	label := slot.Label()
	item := m.equippedItemFor(slot)
	cursor := "  "
	if selected {
		cursor = "› "
	}
	if item == nil {
		line := cursor + m.styles.Muted.Render(label) + "  " + m.styles.Dim.Render("empty")
		if selected {
			line = m.renderSelectedText(ansi.Strip(line), width)
		}
		return []string{line}, 1
	}
	first := cursor + m.styles.Muted.Render(label) + "  " + m.glyphs.symbol(item.Glyph, item.ASCII) + " " + truncateText(m.renderItemName(*item), max(8, width-12))
	second := "  " + m.renderEquippedItemDetail(*item)
	if !detailed {
		if selected {
			return []string{m.renderSelectedText(ansi.Strip(first), width)}, 1
		}
		return []string{first}, 1
	}
	if selected {
		return []string{m.renderSelectedText(ansi.Strip(first), width), m.renderSelectedText(ansi.Strip(truncateText(second, width)), width)}, 2
	}
	return []string{first, truncateText(second, width)}, 2
}

func (m *Model) renderPackItemDetail(item game.Item) string {
	switch item.Kind {
	case game.ItemKindConsumable:
		return m.renderConsumableEffects(item)
	case game.ItemKindEquipment:
		return m.renderComparedEquipmentStats(item, m.equippedItemFor(item.Slot))
	case game.ItemKindRelic:
		return m.styles.Ember.Render("relic")
	default:
		return ""
	}
}

func (m *Model) renderEquippedItemDetail(item game.Item) string {
	if item.Kind == game.ItemKindEquipment {
		return m.renderEquipmentStats(item)
	}
	return m.renderConsumableEffects(item)
}

func (m *Model) renderConsumableEffects(item game.Item) string {
	parts := make([]string, 0, 4)
	if item.Heal > 0 {
		parts = append(parts, m.styles.Heal.Render(fmt.Sprintf("+%d HP", item.Heal)))
	}
	if item.PoisonCure {
		parts = append(parts, m.styles.Cure.Render("cure poison"))
	}
	if item.FireCure {
		parts = append(parts, m.styles.Ember.Render("quench fire"))
	}
	if item.FocusTurns > 0 {
		parts = append(parts, m.styles.Focus.Render(fmt.Sprintf("focus +%d", item.FocusBonus)))
		parts = append(parts, m.styles.Focus.Render(fmt.Sprintf("%dt", item.FocusTurns)))
	}
	if item.EmberDamage > 0 {
		parts = append(parts, m.styles.Ember.Render(fmt.Sprintf("ember %d", item.EmberDamage)))
	}
	if item.PoisonResist > 0 {
		parts = append(parts, m.styles.Cure.Render(fmt.Sprintf("venom ward %d", item.PoisonResist)))
	}
	if item.FireResist > 0 {
		parts = append(parts, m.styles.Ember.Render(fmt.Sprintf("fire ward %d", item.FireResist)))
	}
	return strings.Join(parts, " ")
}

func (m *Model) renderComparedEquipmentStats(item game.Item, equipped *game.Item) string {
	currentAttack, currentDefense, currentMaxHP, currentSight := 0, 0, 0, 0
	if equipped != nil {
		currentAttack = equipped.AttackBonus
		currentDefense = equipped.DefenseBonus
		currentMaxHP = equipped.MaxHPBonus
		currentSight = equipped.SightBonus
	}
	parts := make([]string, 0, 6)
	if item.AttackBonus > 0 {
		parts = append(parts, m.renderComparedToken("ATK", item.AttackBonus, currentAttack))
	}
	if item.DefenseBonus > 0 {
		parts = append(parts, m.renderComparedToken("DEF", item.DefenseBonus, currentDefense))
	}
	if item.MaxHPBonus > 0 {
		parts = append(parts, m.renderComparedToken("HP", item.MaxHPBonus, currentMaxHP))
	}
	if item.SightBonus > 0 {
		parts = append(parts, m.renderComparedToken("SIGHT", item.SightBonus, currentSight))
	}
	if item.PoisonResist > 0 || equippedPoisonResist(equipped) > 0 {
		parts = append(parts, m.renderComparedToken("POISON", item.PoisonResist, equippedPoisonResist(equipped)))
	}
	if item.FireResist > 0 || equippedFireResist(equipped) > 0 {
		parts = append(parts, m.renderComparedToken("FIRE", item.FireResist, equippedFireResist(equipped)))
	}
	return strings.Join(parts, " ")
}

func (m *Model) renderEquipmentStats(item game.Item) string {
	parts := make([]string, 0, 6)
	if item.AttackBonus > 0 {
		parts = append(parts, m.styles.Attack.Render(fmt.Sprintf("ATK+%d", item.AttackBonus)))
	}
	if item.DefenseBonus > 0 {
		parts = append(parts, m.styles.Defense.Render(fmt.Sprintf("DEF+%d", item.DefenseBonus)))
	}
	if item.MaxHPBonus > 0 {
		parts = append(parts, m.styles.Vitality.Render(fmt.Sprintf("HP+%d", item.MaxHPBonus)))
	}
	if item.SightBonus > 0 {
		parts = append(parts, m.styles.Sight.Render(fmt.Sprintf("SIGHT+%d", item.SightBonus)))
	}
	if item.PoisonResist > 0 {
		parts = append(parts, m.styles.Cure.Render(fmt.Sprintf("POISON+%d", item.PoisonResist)))
	}
	if item.FireResist > 0 {
		parts = append(parts, m.styles.Ember.Render(fmt.Sprintf("FIRE+%d", item.FireResist)))
	}
	return strings.Join(parts, " ")
}

func (m *Model) renderComparedToken(label string, candidate int, current int) string {
	style := m.styles.CompareEqual
	switch {
	case candidate > current:
		style = m.styles.CompareBetter
	case candidate < current:
		style = m.styles.CompareWorse
	}
	return style.Render(fmt.Sprintf("%s+%d", label, candidate))
}

func (m *Model) renderQuickHealPreview() string {
	item, count, ok := m.game.QuickHealPreview()
	if !ok {
		return "Quick C  " + m.styles.Dim.Render("none")
	}
	line := "Quick C  " + m.glyphs.symbol(item.Glyph, item.ASCII) + " " + m.renderItemName(item) + " " + m.styles.Quantity.Render(fmt.Sprintf("x%d", count))
	if detail := m.renderConsumableEffects(item); detail != "" {
		line += " " + m.styles.Dim.Render("·") + " " + detail
	}
	return truncateText(line, 34)
}

func (m *Model) renderStateLine(statuses string) string {
	if statuses == "Steady" {
		return m.styles.Success.Render(statuses)
	}
	return m.styles.Danger.Render(statuses)
}

func (m *Model) renderCompletionSummary(completion game.FloorCompletion) string {
	switch {
	case completion.Complete():
		return m.styles.Success.Render("The floor lies hushed.")
	case completion.UnclearedRooms() == 0:
		return m.styles.Info.Render("Every chamber is opened, but the floor still keeps traces.")
	case completion.UnclearedRooms() == 1:
		return m.styles.Info.Render("One chamber still resists your passing.")
	default:
		return m.styles.Info.Render(fmt.Sprintf("%d chambers still resist your passing.", completion.UnclearedRooms()))
	}
}

func (m *Model) renderBoolMetric(label string, complete bool) string {
	if complete {
		return m.styles.Success.Render(label)
	}
	return m.styles.AccentSoft.Render(label)
}

func (m *Model) renderCountMetric(count int) string {
	if count == 0 {
		return m.styles.Success.Render("clear")
	}
	return m.styles.Danger.Render(fmt.Sprintf("%d remain", count))
}

func (m *Model) renderSettingLine(index int, label string, value string) string {
	line := padRight(label, 18) + value
	if index == m.settingsCursor {
		return m.renderSelectedText("› "+ansi.Strip(line), 38)
	}
	return "  " + line
}

func (m *Model) settingsDescription(settings game.Settings) string {
	switch m.settingsCursor {
	case 0:
		return "Choose how the interface renders symbols. Auto follows the fallback toggle, Nerd Font prefers the full rune set, and ASCII forces plain terminal-safe glyphs."
	case 1:
		return "Mirrors the old DUNSHELL_ASCII=1 behavior inside the game UI. When Glyph Mode is Auto, turning this on falls back to plain ASCII without restarting the run."
	case 2:
		return "Controls whether pressing E on stairs asks for confirmation before opening the route map. Turning it off makes descent selection faster while keeping the route choice itself explicit."
	case 3:
		return "Sets how many recent whisper-log entries stay visible in the lower panel. Higher values preserve more combat and status context at the cost of density."
	default:
		return ""
	}
}

func (m *Model) routeStyle(route game.RouteChoice) lipgloss.Style {
	switch route.ID {
	case "gilded_way":
		return m.styles.Gold.Copy()
	case "brokers_lantern":
		return m.styles.Info.Copy()
	case "pilgrims_rest":
		return m.styles.Success.Copy()
	case "reliquary_breach":
		return m.styles.Accent.Copy()
	case "ashen_hunt":
		return m.styles.Warning.Copy()
	case "cursed_procession":
		return m.styles.Danger.Copy()
	default:
		return m.styles.AccentSoft.Copy()
	}
}

func routeAccentColor(route game.RouteChoice) lipgloss.Color {
	switch route.ID {
	case "gilded_way":
		return lipgloss.Color("#dfbe62")
	case "brokers_lantern":
		return lipgloss.Color("#88b2ce")
	case "pilgrims_rest":
		return lipgloss.Color("#88ba7a")
	case "reliquary_breach":
		return lipgloss.Color("#d39d62")
	case "ashen_hunt":
		return lipgloss.Color("#d7a16d")
	case "cursed_procession":
		return lipgloss.Color("#d47777")
	default:
		return lipgloss.Color("#8fa7bf")
	}
}

func (m *Model) floorGlyph(pos game.Position, visible bool, roomStates []game.RoomState) string {
	roomIndex := m.game.Floor.RoomIndexAt(pos)
	bossRoom := roomIndex >= 0 && roomIndex < len(roomStates) && roomStates[roomIndex].Kind == game.RoomBoss && !roomStates[roomIndex].Cleared
	if bossRoom {
		return m.glyphs.bossFloor(visible)
	}
	cleared := m.roomClearedContext(pos, roomStates)
	switch {
	case visible && cleared:
		return m.glyphs.floorClearedVisible()
	case visible:
		return m.glyphs.floorVisible()
	case cleared:
		return m.glyphs.floorClearedSeen()
	default:
		return m.glyphs.floorSeen()
	}
}

func (m *Model) renderStatusSummary(width int) []string {
	if len(m.game.Player.Statuses) == 0 {
		return []string{"  " + m.styles.Success.Render("Steady")}
	}
	lines := make([]string, 0, 3)
	current := "  "
	for _, status := range m.game.Player.Statuses {
		badge := m.renderStatusBadge(status)
		candidate := current
		if strings.TrimSpace(current) == "" {
			candidate = "  " + badge
		} else {
			candidate = current + " " + badge
		}
		if lipgloss.Width(candidate) > width && strings.TrimSpace(current) != "" {
			lines = append(lines, current)
			current = "  " + badge
			continue
		}
		current = candidate
	}
	if strings.TrimSpace(current) != "" {
		lines = append(lines, current)
	}
	return lines
}

func (m *Model) renderStatusBadge(status game.StatusEffect) string {
	label := "[" + status.ShortLabel() + "]"
	switch status.Kind {
	case game.StatusPoison:
		return m.styles.Cure.Render(label)
	case game.StatusFire:
		return m.styles.Ember.Render(label)
	case game.StatusFocus:
		return m.styles.Focus.Render(label)
	default:
		return m.styles.Warning.Render(label)
	}
}

func equippedPoisonResist(item *game.Item) int {
	if item == nil {
		return 0
	}
	return item.PoisonResist
}

func equippedFireResist(item *game.Item) int {
	if item == nil {
		return 0
	}
	return item.FireResist
}

func boolLabel(value bool) string {
	if value {
		return "On"
	}
	return "Off"
}

func (m *Model) renderSettingsActions() string {
	active := -1
	if m.settingsCursor >= 4 {
		active = m.settingsCursor - 4
	}
	apply := m.styles.ModalChoice.Render(" Apply ")
	back := m.styles.ModalChoice.Render(" Back ")
	if active == 0 {
		apply = m.styles.ModalChoiceActive.Render(" Apply ")
	}
	if active == 1 {
		back = m.styles.ModalChoiceActive.Render(" Back ")
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, apply, "  ", back)
}

func fixedPanelBody(lines []string, wrapWidth int, maxLines int) string {
	if maxLines <= 0 {
		return ""
	}
	out := make([]string, 0, maxLines)
	truncated := false
	for _, line := range lines {
		if len(out) >= maxLines {
			truncated = true
			break
		}
		if strings.TrimSpace(ansi.Strip(line)) == "" {
			out = append(out, "")
			continue
		}
		wrapped := wrapText(line, wrapWidth)
		for _, part := range strings.Split(wrapped, "\n") {
			if len(out) >= maxLines {
				truncated = true
				break
			}
			out = append(out, truncateText(part, wrapWidth))
		}
	}
	if truncated && len(out) > 0 {
		last := out[len(out)-1]
		out[len(out)-1] = truncateText(ansi.Strip(last), max(1, wrapWidth-1))
		if lipgloss.Width(out[len(out)-1]) >= wrapWidth-1 {
			out[len(out)-1] = truncateText(out[len(out)-1], wrapWidth)
		}
	}
	for len(out) < maxLines {
		out = append(out, "")
	}
	return strings.Join(out, "\n")
}

func routeWaitsPanelLines(route game.RouteChoice) []string {
	modifier := route.Modifier
	lines := make([]string, 0, 8)
	if modifier.BonusGold > 0 {
		lines = append(lines, "• more gold from kills")
	}
	if modifier.Merchant {
		lines = append(lines, "• merchant guaranteed")
	}
	if modifier.Rest {
		lines = append(lines, fmt.Sprintf("• recover %d HP", modifier.HealOnStart))
		if modifier.CleanseOnRest {
			lines = append(lines, "• poison and fire cleansed")
		}
	}
	if modifier.GuaranteedKey != nil {
		lines = append(lines, "• guaranteed "+modifier.GuaranteedKey.LowerLabel()+" key")
	}
	if modifier.ExtraChests > 0 {
		lines = append(lines, "• extra reliquary chest")
	}
	if modifier.LootBonus > 0 {
		lines = append(lines, "• stronger loot rolls")
	}
	if modifier.EnemyBonus > 0 {
		lines = append(lines, fmt.Sprintf("• %d extra foes", modifier.EnemyBonus))
	} else if modifier.EnemyBonus < 0 {
		lines = append(lines, "• fewer enemies")
	}
	if modifier.EliteChance > 0 {
		lines = append(lines, fmt.Sprintf("• %.0f%% extra elites", modifier.EliteChance*100))
	}
	if modifier.Cursed {
		lines = append(lines, "• cursed scaling")
	}
	if len(lines) == 0 {
		lines = append(lines, "• steady descent")
	}
	return lines
}

func (m *Model) renderStatBar(label string, current int, total int, width int, fill lipgloss.Color, track lipgloss.Color) string {
	if total <= 0 {
		total = 1
	}
	numeric := m.styles.Quantity.Render(fmt.Sprintf("%d/%d", current, total))
	barWidth := clamp(width-lipgloss.Width(label)-lipgloss.Width(numeric)-2, 10, 16)
	bar := compactBar(current, total, barWidth, fill, track)
	return lipgloss.NewStyle().Width(width).Render(label + " " + numeric + " " + bar)
}

func compactBar(current int, total int, width int, fill lipgloss.Color, track lipgloss.Color) string {
	if total <= 0 {
		total = 1
	}
	current = clamp(current, 0, total)
	filled := current * width / total
	if filled > width {
		filled = width
	}
	full := lipgloss.NewStyle().Foreground(fill).Render(strings.Repeat("█", filled))
	empty := lipgloss.NewStyle().Foreground(track).Render(strings.Repeat("░", width-filled))
	return full + empty
}

func (m *Model) renderSelectedText(text string, width int) string {
	return m.styles.ListSelected.Copy().Width(width).Render(text)
}

func (m *Model) renderBinaryChoice(choice int, yes string, no string) string {
	left := m.styles.ModalChoice.Render(" " + yes + " ")
	right := m.styles.ModalChoice.Render(" " + no + " ")
	if choice == 0 {
		left = m.styles.ModalChoiceActive.Render(" " + yes + " ")
	} else {
		right = m.styles.ModalChoiceActive.Render(" " + no + " ")
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, left, "  ", right)
}

func (m *Model) renderKeyLine() string {
	keys := m.game.Player.Keys
	parts := []string{
		m.styles.keyStyle(game.KeyBronze).Render(fmt.Sprintf("B:%d", keys.Bronze)),
		m.styles.keyStyle(game.KeySilver).Render(fmt.Sprintf("S:%d", keys.Silver)),
		m.styles.keyStyle(game.KeyGold).Render(fmt.Sprintf("G:%d", keys.Gold)),
	}
	return strings.Join(parts, "  ")
}

func (m *Model) renderRarityName(item game.Item) string {
	return m.styles.rarityStyle(item.Rarity).Render(item.Rarity.Label())
}

func (m *Model) renderItemName(item game.Item) string {
	return m.styles.rarityStyle(item.Rarity).Render(item.Name)
}

func padRight(text string, width int) string {
	padding := max(0, width-lipgloss.Width(text))
	return text + strings.Repeat(" ", padding)
}

func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}
	lines := make([]string, 0, len(words)/8+1)
	line := words[0]
	for _, word := range words[1:] {
		if lipgloss.Width(line)+1+lipgloss.Width(word) > width {
			lines = append(lines, line)
			line = word
			continue
		}
		line += " " + word
	}
	lines = append(lines, line)
	return strings.Join(lines, "\n")
}

func truncateText(text string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(strings.ReplaceAll(ansi.Strip(text), "\n", " "))
	if len(runes) <= width {
		return string(runes)
	}
	if width <= 3 {
		return ansi.Truncate(text, width, "")
	}
	return ansi.Truncate(text, width, "...")
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func clamp(value int, low int, high int) int {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}
