package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	figure "github.com/common-nighthawk/go-figure"

	"dunshell/internal/game"
)

var titleArt = figure.NewFigure(game.GameTitle, "small", true).Slicify()

func (m *Model) renderTopHeader(width int) string {
	return m.styles.HeaderBar.Copy().Width(width).Render(game.GameTitle)
}

func (m *Model) renderScreen(content string, width int, height int, hPos lipgloss.Position, vPos lipgloss.Position) string {
	header := m.renderTopHeader(width)
	availableHeight := max(1, height-lipgloss.Height(header))
	body := lipgloss.Place(width, availableHeight, hPos, vPos, content)
	return lipgloss.JoinVertical(lipgloss.Left, header, body)
}

func (m *Model) gameplayMetrics() (int, int, int, int, int) {
	width := max(88, m.width)
	height := max(30, m.height)
	headerHeight := lipgloss.Height(m.renderTopHeader(width))
	footerHeight := lipgloss.Height(m.styles.Footer.Render(m.help.View(m.keys)))
	availableHeight := max(22, height-headerHeight-footerHeight)
	logHeight := clamp(height/4, 6, 9)
	if availableHeight-logHeight < 14 {
		logHeight = max(6, availableHeight-14)
	}
	bodyHeight := max(14, availableHeight-logHeight)
	sidebarWidth := clamp(width/3, 34, 40)
	return width, height, bodyHeight, logHeight, sidebarWidth
}

func (m *Model) viewTitle() string {
	width := max(80, m.width)
	height := max(28, m.height)

	titleLines := make([]string, 0, len(titleArt))
	for _, line := range titleArt {
		titleLines = append(titleLines, m.styles.Title.Render(line))
	}

	menuLabels := []string{"New Run", "Field Guide", "Quit"}
	menuLines := make([]string, 0, len(menuLabels))
	for index, label := range menuLabels {
		line := "  " + m.styles.MenuItem.Render(label)
		if index == m.menuIndex {
			line = m.renderSelectedText("› "+label, 30)
		}
		menuLines = append(menuLines, line)
	}

	copy := wrapText("Descend beneath the drowned abbey, strip its halls bare, and carry the Cinder Crown back through whatever still remembers your footsteps.", 68)
	seedLine := m.styles.Muted.Render("Seed: random each run")
	if m.hasLockedSeed {
		seedLine = m.styles.Muted.Render("Seed locked to " + fmt.Sprintf("%d", m.lockedSeed))
	}

	panelBody := strings.Join([]string{
		m.styles.Subtitle.Render("A polished terminal roguelike carved in Bubble Tea"),
		"",
		copy,
		"",
		strings.Join(menuLines, "\n"),
		"",
		seedLine,
		m.styles.PanelNote.Render("Arrows or W/S move through menus. Enter confirms."),
	}, "\n")

	panel := m.styles.focusBox("Start Menu", panelBody, 74, 0)
	content := lipgloss.JoinVertical(lipgloss.Center, strings.Join(titleLines, "\n"), "", panel)
	return m.renderScreen(content, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) viewHelp() string {
	width := max(96, m.width)
	height := max(34, m.height)

	legend := []string{
		m.styles.cellGlyph(m.styles.TileFloorVisible, "■", "#84d06f", true) + " You",
		m.styles.TileWallVisible.Render("█") + " Wall",
		m.styles.TileWallClearedVisible.Render("█") + " Cleared chamber wall",
		m.styles.cellGlyph(m.styles.TileFloorVisible, "+", "#d4a66d", true) + " Closed door",
		m.styles.cellGlyph(m.styles.TileFloorVisible, "/", "#9f7f61", false) + " Open door",
		m.styles.cellGlyph(m.styles.TileFloorVisible, "▾", "#9bc1d8", true) + " Stairs down",
		m.styles.colorGlyph('!', "#d16078", true) + " Consumable",
		m.styles.colorGlyph(')', "#f2c97d", true) + " Weapon",
		m.styles.colorGlyph('[', "#7ea8c7", true) + " Armor",
		m.styles.colorGlyph('=', "#f6db7d", true) + " Charm",
	}

	body := strings.Join([]string{
		m.renderHelpSection("Movement", []string{
			"Arrows or W/A/S/D move one tile at a time.",
			"Walking into an enemy attacks it.",
			"Press . to wait and let the dungeon move.",
		}),
		"",
		m.renderHelpSection("Actions", []string{
			"E acts on your tile. Loot is gathered before stairs are offered.",
			"C drinks the weakest healing item in your pack.",
			"I opens the inventory sidebar. Esc closes the active overlay.",
			"Q opens a safe quit prompt.",
		}),
		"",
		m.renderHelpSection("Inventory", []string{
			"Left or right switches between Pack and Equipped.",
			"E or Enter performs the primary action on the selected item.",
			"U directly uses a consumable from the pack.",
			"Equipment is removed by selecting it in the Equipped pane and pressing E or Enter.",
		}),
		"",
		m.renderHelpSection("Legend", legend),
		"",
		m.renderHelpSection("Field Notes", []string{
			wrapText("Blue chamber walls mark rooms that have been fully laid bare: explored, opened, cleansed, and picked clean. The stair prompt tells you what the floor still keeps from you.", 78),
			m.styles.PanelNote.Render("Press Esc, Enter, or ? to return."),
		}),
	}, "\n")

	panel := m.styles.box("Field Guide", body, min(width-4, 104), 0)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) renderHelpSection(title string, lines []string) string {
	return strings.Join([]string{
		m.styles.Subtitle.Render(title),
		strings.Join(lines, "\n"),
	}, "\n")
}

func (m *Model) viewGame() string {
	if m.game == nil {
		return ""
	}

	width, height, bodyHeight, logHeight, sidebarWidth := m.gameplayMetrics()
	footer := m.styles.Footer.Render(m.help.View(m.keys))
	headerHeight := lipgloss.Height(m.renderTopHeader(width))

	m.mapViewport = mapViewportState{}

	bodyOriginY := headerHeight
	var body string
	if width >= 112 {
		mapWidth := max(46, width-sidebarWidth-1)
		mapPanel := m.renderMapPanel(mapWidth, bodyHeight, 0, bodyOriginY)
		sideX := mapWidth + 1
		sidePanel := m.renderSidebar(sidebarWidth, bodyHeight, sideX, bodyOriginY)
		body = lipgloss.JoinHorizontal(lipgloss.Top, mapPanel, " ", sidePanel)
	} else {
		sidebarHeight := max(10, bodyHeight-1)
		mapHeight := max(9, bodyHeight-sidebarHeight)
		if mapHeight+sidebarHeight > bodyHeight {
			sidebarHeight = max(8, bodyHeight-mapHeight)
		}
		mapPanel := m.renderMapPanel(width, mapHeight, 0, bodyOriginY)
		sidePanel := m.renderSidebar(width, sidebarHeight, 0, bodyOriginY+mapHeight)
		body = lipgloss.JoinVertical(lipgloss.Left, mapPanel, sidePanel)
	}

	logY := bodyOriginY + lipgloss.Height(body)
	logPanel := m.renderLog(width, logHeight, 0, logY)
	content := lipgloss.JoinVertical(lipgloss.Left, body, logPanel, footer)
	return m.renderScreen(content, width, height, lipgloss.Left, lipgloss.Top)
}

func (m *Model) viewQuitPrompt() string {
	width := max(60, m.width)
	height := max(20, m.height)
	body := strings.Join([]string{
		wrapText("Leave this run? The abbey will keep the shape of your descent, but not the life inside it.", 40),
		"",
		m.styles.PanelNote.Render("Press y to quit, n or Esc to return."),
	}, "\n")
	panel := m.styles.modalBox("Quit Run", body, 46)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) viewOutcome(victory bool) string {
	width := max(82, m.width)
	height := max(28, m.height)
	summary := m.game.Summary()

	var title string
	var accent lipgloss.Style
	var copy string
	if victory {
		title = "Cinder Crown Claimed"
		accent = m.styles.Success
		copy = "You return carrying a relic the abbey kept for centuries. The bells above will never ring the same way again."
	} else {
		title = "Run Extinguished"
		accent = m.styles.Danger
		copy = "The abbey keeps another name in its stone throat. The next descent will remember where this one failed."
	}

	stats := []string{
		accent.Render(title),
		"",
		wrapText(copy, 62),
		"",
		"Seed      " + fmt.Sprintf("%d", summary.Seed),
		"Floor     " + fmt.Sprintf("%d / %d", summary.Floor, m.game.MaxFloors),
		"Level     " + fmt.Sprintf("%d", summary.Level),
		"Gold      " + fmt.Sprintf("%d", summary.Gold),
		"Kills     " + fmt.Sprintf("%d", summary.Kills),
		"Turns     " + fmt.Sprintf("%d", summary.Turn),
		"",
		m.styles.PanelNote.Render("Press Enter for the title screen, n for a new run, q to quit."),
	}

	panel := m.styles.modalBox(title, strings.Join(stats, "\n"), 70)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m *Model) renderMapPanel(width int, height int, originX int, originY int) string {
	innerWidth := max(22, width-4)
	innerHeight := max(8, height-4)
	player := m.game.Player
	floor := m.game.Floor

	cameraX := clamp(player.Pos.X-innerWidth/2, 0, max(0, floor.Width-innerWidth))
	cameraY := clamp(player.Pos.Y-innerHeight/2, 0, max(0, floor.Height-innerHeight))

	m.mapViewport = mapViewportState{
		Panel:   rect{X: originX, Y: originY, W: width, H: height},
		Content: rect{X: originX + 2, Y: originY + 2, W: innerWidth, H: innerHeight},
		CameraX: cameraX,
		CameraY: cameraY,
	}

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

	roomStates := floor.RoomStates()
	cells := make([][]string, innerHeight)
	for y := 0; y < innerHeight; y++ {
		cells[y] = make([]string, innerWidth)
		for x := 0; x < innerWidth; x++ {
			pos := game.Position{X: cameraX + x, Y: cameraY + y}
			cells[y][x] = m.renderMapCell(pos, enemyMap, itemMap, roomStates)
		}
	}

	rows := m.renderCellRows(cells)
	if m.overlay == overlayDescend {
		rows = m.renderDescendPrompt(cells)
	}

	body := strings.Join(rows, "\n")
	if m.overlay == overlayDescend {
		return m.styles.focusBox("Dungeon", body, width, height)
	}
	return m.styles.box("Dungeon", body, width, height)
}

func (m *Model) renderCellRows(cells [][]string) []string {
	rows := make([]string, len(cells))
	for y := range cells {
		rows[y] = strings.Join(cells[y], "")
	}
	return rows
}

func (m *Model) renderMapCell(pos game.Position, enemies map[game.Position]*game.Enemy, items map[game.Position]game.GroundItem, roomStates []game.RoomState) string {
	floor := m.game.Floor
	if !floor.InBounds(pos) {
		return " "
	}

	visible := floor.IsVisible(pos)
	explored := floor.IsExplored(pos)
	if !explored {
		return m.styles.Void.Render(" ")
	}

	if m.game.Player.Pos.Equals(pos) {
		base := m.floorCellStyle(pos, visible, roomStates)
		return m.styles.cellGlyph(base, "■", "#84d06f", true)
	}

	if enemy, ok := enemies[pos]; ok {
		base := m.floorCellStyle(pos, visible, roomStates)
		return m.styles.cellGlyph(base, string(enemy.Template.Glyph), enemy.Template.Tint, true)
	}

	if item, ok := items[pos]; ok {
		base := m.floorCellStyle(pos, visible, roomStates)
		return m.styles.cellGlyph(base, string(item.Item.Glyph), item.Item.Tint, true)
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
		return m.floorCellStyle(pos, visible, roomStates).Render(" ")
	case game.TileDoorClosed:
		base := m.doorCellStyle(pos, visible, roomStates)
		return m.styles.cellGlyph(base, "+", "#d4a66d", true)
	case game.TileDoorOpen:
		base := m.doorCellStyle(pos, visible, roomStates)
		return m.styles.cellGlyph(base, "/", "#8b715c", false)
	case game.TileStairsDown:
		base := m.floorCellStyle(pos, visible, roomStates)
		tint := "#9bc1d8"
		if !visible {
			tint = "#5c7686"
		}
		return m.styles.cellGlyph(base, "▾", tint, true)
	default:
		return " "
	}
}

func (m *Model) floorCellStyle(pos game.Position, visible bool, roomStates []game.RoomState) lipgloss.Style {
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

func (m *Model) renderDescendPrompt(cells [][]string) []string {
	rows := m.renderCellRows(cells)
	completion := m.game.Floor.Completion()
	promptWidth := min(46, max(24, m.mapViewport.Content.W-2))
	if promptWidth > m.mapViewport.Content.W {
		promptWidth = m.mapViewport.Content.W
	}
	promptText := m.renderDescendPromptText(completion, max(18, promptWidth-4))
	prompt := m.styles.modalBox("Stair Mouth", promptText, promptWidth)
	lines := strings.Split(prompt, "\n")

	stairsX := m.game.Player.Pos.X - m.mapViewport.CameraX
	stairsY := m.game.Player.Pos.Y - m.mapViewport.CameraY
	promptHeight := len(lines)

	x := clamp(stairsX-promptWidth/2, 0, max(0, m.mapViewport.Content.W-promptWidth))
	y := stairsY - promptHeight - 1
	if y < 0 {
		y = min(max(0, stairsY+1), max(0, m.mapViewport.Content.H-promptHeight))
	}
	y = clamp(y, 0, max(0, m.mapViewport.Content.H-promptHeight))

	for lineIndex, line := range lines {
		rowIndex := y + lineIndex
		if rowIndex < 0 || rowIndex >= len(rows) {
			continue
		}
		left := strings.Join(cells[rowIndex][:x], "")
		rightStart := min(len(cells[rowIndex]), x+promptWidth)
		right := strings.Join(cells[rowIndex][rightStart:], "")
		rows[rowIndex] = left + line + right
	}
	return rows
}

func (m *Model) renderDescendPromptText(completion game.FloorCompletion, width int) string {
	options := m.renderDescendOptions()
	if completion.Complete() {
		lines := []string{
			m.styles.Success.Render(wrapText("The hall behind you lies cold and emptied. Only the deeper ash still keeps a voice.", width)),
			"",
			m.styles.PanelNote.Render(wrapText("Do you step into the next dark?", width)),
			"",
			options,
		}
		return strings.Join(lines, "\n")
	}

	missing := make([]string, 0, 4)
	if !completion.FullyExplored() {
		missing = append(missing, m.styles.Info.Render(fmt.Sprintf("• %d tiles still hide in shadow.", completion.UnexploredTiles)))
	}
	if !completion.LootCollected() {
		missing = append(missing, m.styles.Gold.Render(fmt.Sprintf("• %d spoils still glimmer below.", completion.RemainingItems)))
	}
	if !completion.EnemiesCleared() {
		missing = append(missing, m.styles.Danger.Render(fmt.Sprintf("• %d foes still draw breath.", completion.RemainingEnemies)))
	}
	if completion.UnclearedRooms() > 0 {
		missing = append(missing, m.styles.AccentSoft.Render(fmt.Sprintf("• %d chambers remain unpurged.", completion.UnclearedRooms())))
	}

	lines := []string{
		m.styles.Danger.Render(wrapText("The stair offers passage, but this floor still keeps unfinished hungers.", width)),
		"",
		strings.Join(missing, "\n"),
		"",
		options,
	}
	return strings.Join(lines, "\n")
}

func (m *Model) renderDescendOptions() string {
	yes := " Yes "
	no := " No "
	if m.descendPrompt.Choice == 0 {
		yes = m.styles.ModalChoiceActive.Render(yes)
		no = m.styles.ModalChoice.Render(no)
	} else {
		yes = m.styles.ModalChoice.Render(yes)
		no = m.styles.ModalChoiceActive.Render(no)
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, yes, "  ", no) + "\n" + m.styles.PanelNote.Render("Left/Right choose • Enter confirms • Esc stays")
}

func (m *Model) renderSidebar(width int, height int, originX int, originY int) string {
	if m.overlay == overlayInventory {
		return m.renderInventorySidebar(width, height, originX, originY)
	}
	return m.renderStatusSidebar(width, height)
}

func (m *Model) renderStatusSidebar(width int, height int) string {
	player := m.game.Player
	statuses := "Steady"
	if len(player.Statuses) > 0 {
		labels := make([]string, 0, len(player.Statuses))
		for _, status := range player.Statuses {
			labels = append(labels, status.Label())
		}
		statuses = strings.Join(labels, ", ")
	}

	completion := m.game.Floor.Completion()
	stats := []string{
		m.renderStatBar("HP", player.HP, player.MaxHP(), width-4, lipgloss.Color("#c56b78"), lipgloss.Color("#2d1d23")),
		m.renderStatBar("XP", player.XP, player.NextLevelXP(), width-4, lipgloss.Color("#5d88aa"), lipgloss.Color("#18232d")),
		"",
		"Lvl   " + m.styles.AccentSoft.Render(fmt.Sprintf("%d", player.Level)) + "   Gold " + m.styles.Gold.Render(fmt.Sprintf("%d", player.Gold)),
		"ATK   " + m.styles.Attack.Render(fmt.Sprintf("%d", player.AttackPower())) + "   DEF  " + m.styles.Defense.Render(fmt.Sprintf("%d", player.DefensePower())),
		"State " + m.renderStateLine(statuses),
		m.renderQuickHealPreview(),
	}

	floorLines := []string{
		m.styles.Muted.Render(m.game.FloorLabel()),
		m.renderCompletionSummary(completion),
		"Map   " + m.renderBoolMetric(fmt.Sprintf("%d%%", m.game.Floor.ExploredPercent()), completion.FullyExplored()),
		"Rooms " + m.renderBoolMetric(fmt.Sprintf("%d/%d", completion.ClearedRooms, completion.TotalRooms), completion.UnclearedRooms() == 0),
		"Foes  " + m.renderCountMetric(completion.RemainingEnemies),
		"Loot  " + m.renderCountMetric(completion.RemainingItems),
		"Under " + m.styles.AccentSoft.Render(m.game.TileDescriptionUnderPlayer()),
		"",
		m.renderInteractionHint(),
		"",
		wrapText(m.game.Objective(), width-6),
	}

	body := strings.Join([]string{
		m.styles.Subtitle.Render("Vitals"),
		strings.Join(stats, "\n"),
		"",
		m.styles.Subtitle.Render("Floor"),
		strings.Join(floorLines, "\n"),
	}, "\n")

	return m.styles.box("Status", body, width, height)
}

func (m *Model) renderInventorySidebar(width int, height int, originX int, originY int) string {
	_ = originX
	_ = originY
	packHeight, equippedHeight := m.inventoryPanelHeights(height)

	packPanel := m.renderPackPanel(width, packHeight)
	lowerPanel := m.renderInventoryDetailsPanel(width, equippedHeight)
	return lipgloss.JoinVertical(lipgloss.Left, packPanel, lowerPanel)
}

func (m *Model) renderLog(width int, height int, originX int, originY int) string {
	_ = originX
	_ = originY
	innerWidth := max(22, width-6)
	maxLines := max(1, height-3)
	lines := make([]string, 0, maxLines)

	start := max(0, len(m.game.Log)-24)
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

func (m *Model) renderPackPanel(width int, height int) string {
	stacks := m.inventoryStacks()
	if len(stacks) == 0 {
		box := m.styles.box
		if m.inventoryPane == inventoryPanePack {
			box = m.styles.focusBox
		}
		return box("Pack", m.styles.Dim.Render("Your pack is empty."), width, height)
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

	box := m.styles.box
	if m.inventoryPane == inventoryPanePack {
		box = m.styles.focusBox
	}
	return box("Pack", strings.Join(lines, "\n"), width, height)
}

func (m *Model) renderInventoryDetailsPanel(width int, height int) string {
	contentLines := max(1, height-3)
	detailed := contentLines >= 6

	lines := make([]string, 0, contentLines)
	for _, slot := range []game.EquipmentSlot{game.SlotWeapon, game.SlotArmor, game.SlotCharm} {
		selected := m.inventoryPane == inventoryPaneEquipped && m.selectedEquipmentSlot() == slot
		rendered, rowHeight := m.renderEquippedSlot(slot, selected, detailed, width-4)
		lines = append(lines, rendered...)
		_ = rowHeight
	}

	box := m.styles.box
	if m.inventoryPane == inventoryPaneEquipped {
		box = m.styles.focusBox
	}
	return box("Equipped", strings.Join(lines, "\n"), width, height)
}

func (m *Model) renderPackLine(stack inventoryStack, selected bool, width int) string {
	quantity := m.styles.Quantity.Render(fmt.Sprintf("x%d", stack.Count))
	left := m.styles.colorGlyph(stack.Item.Glyph, stack.Item.Tint, true) + " " + stack.Item.Name
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
		if detailed {
			second := "  " + m.styles.Dim.Render("No rite is bound to this slot.")
			if selected {
				return []string{
					m.renderSelectedText(ansi.Strip(line), width),
					m.renderSelectedText(ansi.Strip(second), width),
				}, 2
			}
			return []string{line, second}, 2
		}
		return []string{line}, 1
	}

	first := cursor + m.styles.Muted.Render(label) + "  " + m.styles.colorGlyph(item.Glyph, item.Tint, true) + " " + truncateText(item.Name, max(8, width-12))
	if !detailed {
		if selected {
			return []string{
				m.renderSelectedText(ansi.Strip(first), width),
				m.renderSelectedText("  "+ansi.Strip(truncateText(m.renderEquippedItemDetail(*item), width-2)), width),
			}, 2
		}
		return []string{first}, 1
	}

	second := "  " + m.renderEquippedItemDetail(*item)
	if selected {
		return []string{
			m.renderSelectedText(ansi.Strip(first), width),
			m.renderSelectedText(ansi.Strip(truncateText(second, width)), width),
		}, 2
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
	switch item.Kind {
	case game.ItemKindEquipment:
		return m.renderEquipmentStats(item)
	case game.ItemKindConsumable:
		return m.renderConsumableEffects(item)
	default:
		return ""
	}
}

func (m *Model) renderConsumableEffects(item game.Item) string {
	parts := make([]string, 0, 4)
	if item.Heal > 0 {
		parts = append(parts, m.styles.Heal.Render(fmt.Sprintf("+%d HP", item.Heal)))
	}
	if item.PoisonCure {
		parts = append(parts, m.styles.Cure.Render("cure"))
	}
	if item.FocusTurns > 0 {
		parts = append(parts, m.styles.Focus.Render(fmt.Sprintf("focus +%d", item.FocusBonus)))
		parts = append(parts, m.styles.Focus.Render(fmt.Sprintf("%dt", item.FocusTurns)))
	}
	if item.EmberDamage > 0 {
		parts = append(parts, m.styles.Ember.Render(fmt.Sprintf("ember %d", item.EmberDamage)))
	}
	return strings.Join(parts, " ")
}

func (m *Model) renderComparedEquipmentStats(item game.Item, equipped *game.Item) string {
	currentAttack := 0
	currentDefense := 0
	currentMaxHP := 0
	currentSight := 0
	if equipped != nil {
		currentAttack = equipped.AttackBonus
		currentDefense = equipped.DefenseBonus
		currentMaxHP = equipped.MaxHPBonus
		currentSight = equipped.SightBonus
	}

	parts := make([]string, 0, 4)
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
	return strings.Join(parts, " ")
}

func (m *Model) renderEquipmentStats(item game.Item) string {
	parts := make([]string, 0, 4)
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

	line := "Quick C  " + m.styles.colorGlyph(item.Glyph, item.Tint, true) + " " + item.Name + " " + m.styles.Quantity.Render(fmt.Sprintf("x%d", count))
	if detail := m.renderConsumableEffects(item); detail != "" {
		line += " " + m.styles.Dim.Render("·") + " " + detail
	}
	return truncateText(line, 30)
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

func (m *Model) renderInteractionHint() string {
	context := m.game.InteractionContext()
	switch context.Primary {
	case game.InteractionPickup:
		if len(context.Secondary) > 0 {
			return m.styles.Accent.Render("E gathers the loot first; the stair waits after.")
		}
		return m.styles.Accent.Render("E gathers what lies beneath your boots.")
	case game.InteractionDescend:
		return m.styles.Info.Render("E asks the stair if you are ready to go lower.")
	default:
		return m.styles.Dim.Render("E has nothing to answer here.")
	}
}

func (m *Model) renderStatBar(label string, current int, total int, width int, fill lipgloss.Color, track lipgloss.Color) string {
	if total <= 0 {
		total = 1
	}
	numeric := m.styles.Quantity.Render(fmt.Sprintf("%d/%d", current, total))
	barWidth := clamp(width-lipgloss.Width(label)-lipgloss.Width(numeric)-2, 8, 12)
	bar := compactBar(current, total, barWidth, fill, track)
	line := label + " " + numeric + " " + bar
	return lipgloss.NewStyle().Width(width).Render(line)
}

func compactBar(current int, total int, width int, fill lipgloss.Color, track lipgloss.Color) string {
	if total <= 0 {
		total = 1
	}
	if current < 0 {
		current = 0
	}
	if current > total {
		current = total
	}

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
