package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	figure "github.com/common-nighthawk/go-figure"

	"dunshell/internal/game"
)

var titleArt = figure.NewFigure(game.GameTitle, "small", true).Slicify()

func (m Model) renderTopHeader(width int) string {
	return m.styles.HeaderBar.Copy().Width(width).Render(game.GameTitle)
}

func (m Model) renderScreen(content string, width int, height int, hPos lipgloss.Position, vPos lipgloss.Position) string {
	header := m.renderTopHeader(width)
	availableHeight := max(1, height-lipgloss.Height(header))
	body := lipgloss.Place(width, availableHeight, hPos, vPos, content)
	return lipgloss.JoinVertical(lipgloss.Left, header, body)
}

func (m Model) gameplayMetrics() (int, int, int, int, int) {
	width := max(80, m.width)
	height := max(28, m.height)
	headerHeight := lipgloss.Height(m.renderTopHeader(width))
	footerHeight := lipgloss.Height(m.styles.Footer.Render(m.help.View(m.keys)))
	availableHeight := max(20, height-headerHeight-footerHeight)
	logHeight := clamp(height/4, 7, 10)
	maxLogHeight := max(7, availableHeight-12)
	if logHeight > maxLogHeight {
		logHeight = maxLogHeight
	}
	bodyHeight := max(12, availableHeight-logHeight)
	if width >= 112 && bodyHeight%2 != 0 {
		bodyHeight--
	}
	sidebarWidth := 34
	return width, height, bodyHeight, logHeight, sidebarWidth
}

func (m Model) viewTitle() string {
	width := max(80, m.width)
	height := max(28, m.height)

	titleLines := make([]string, 0, len(titleArt))
	for _, line := range titleArt {
		titleLines = append(titleLines, m.styles.Title.Render(line))
	}

	menuLabels := []string{"New Run", "Help", "Quit"}
	menuLines := make([]string, 0, len(menuLabels))
	for index, label := range menuLabels {
		if index == m.menuIndex {
			menuLines = append(menuLines, m.styles.MenuCursor.Render(">")+" "+m.styles.MenuSelected.Render(label))
			continue
		}
		menuLines = append(menuLines, "  "+m.styles.MenuItem.Render(label))
	}

	copy := wrapText("Delve beneath the drowned abbey, outfight its hungry choir, and return with the Cinder Crown before the dark learns your name.", 66)
	seedLine := m.styles.Muted.Render("Seed: random each run")
	if m.hasLockedSeed {
		seedLine = m.styles.Muted.Render("Seed locked to " + fmt.Sprintf("%d", m.lockedSeed))
	}

	panelBody := strings.Join([]string{
		m.styles.Subtitle.Render("A terminal roguelike forged with Bubble Tea"),
		"",
		copy,
		"",
		strings.Join(menuLines, "\n"),
		"",
		seedLine,
		m.styles.Dim.Render("Controls: up/down or w/s; letters work in either case; Enter selects"),
	}, "\n")

	panel := m.styles.box("Start Menu", panelBody, 72, 0)
	content := lipgloss.JoinVertical(lipgloss.Center, strings.Join(titleLines, "\n"), "", panel)
	return m.renderScreen(content, width, height, lipgloss.Center, lipgloss.Center)
}

func (m Model) viewHelp() string {
	width := max(90, m.width)
	height := max(32, m.height)
	helpView := help.New()
	helpView.Width = min(width-8, 96)
	helpView.ShowAll = true

	legend := []string{
		m.styles.Player.Render("@") + " You",
		m.styles.TileWall.Render("#") + " Wall",
		m.styles.TileDoor.Render("+") + " Closed door",
		m.styles.TileStairs.Render(">") + " Stairs down",
		m.styles.colorGlyph('!', "#d16078", true) + " Consumable",
		m.styles.colorGlyph(')', "#f2c97d", true) + " Weapon",
		m.styles.colorGlyph('[', "#7ea8c7", true) + " Armor",
		m.styles.colorGlyph('=', "#f6db7d", true) + " Charm",
		m.styles.colorGlyph('&', "#ffb347", true) + " Cinder Crown",
	}

	tips := wrapText("Movement is turn-based. Bump enemies to attack, use doors to break line of sight, drink Sunbrew before committing to a hard fight, and do not sit on poison longer than you have to.", 72)
	body := strings.Join([]string{
		m.styles.Subtitle.Render("Controls"),
		helpView.View(m.keys),
		"",
		m.styles.Subtitle.Render("Legend"),
		strings.Join(legend, "\n"),
		"",
		m.styles.Subtitle.Render("Field Notes"),
		tips,
		"",
		m.styles.Dim.Render("Press Esc, Enter, or ? to return."),
	}, "\n")

	panel := m.styles.box("How To Survive", body, min(width-4, 100), 0)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m Model) viewGame() string {
	if m.game == nil {
		return ""
	}

	width, height, bodyHeight, logHeight, sidebarWidth := m.gameplayMetrics()
	footer := m.styles.Footer.Render(m.help.View(m.keys))

	var body string
	if width >= 112 {
		mapWidth := max(36, width-sidebarWidth-5)
		mapPanel := m.renderMapPanel(mapWidth, bodyHeight)
		sidePanel := m.renderSidebar(sidebarWidth, bodyHeight)
		body = lipgloss.JoinHorizontal(lipgloss.Top, mapPanel, " ", sidePanel)
	} else {
		mapPanel := m.renderMapPanel(width-2, max(10, bodyHeight-14))
		sidePanel := m.renderSidebar(width-2, 12)
		body = lipgloss.JoinVertical(lipgloss.Left, mapPanel, sidePanel)
	}

	logPanel := m.renderLog(width-2, logHeight)
	content := lipgloss.JoinVertical(lipgloss.Left, body, logPanel, footer)
	return m.renderScreen(content, width, height, lipgloss.Left, lipgloss.Top)
}

func (m Model) viewQuitPrompt() string {
	width := max(60, m.width)
	height := max(20, m.height)
	body := strings.Join([]string{
		wrapText("Leave this run? Your current descent will be lost.", 36),
		"",
		m.styles.Dim.Render("Press y/Y to quit, n/N or esc to return."),
	}, "\n")
	panel := m.styles.box("Quit Run", body, 44, 0)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m Model) viewOutcome(victory bool) string {
	width := max(80, m.width)
	height := max(28, m.height)
	summary := m.game.Summary()

	var title string
	var accent lipgloss.Style
	var copy string
	if victory {
		title = "Cinder Crown Claimed"
		accent = m.styles.Success
		copy = "You return from the abbey carrying a relic it kept for centuries. The bells above will never sound the same again."
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
		m.styles.Dim.Render("Press Enter for the title screen, n/N for a new run, q/Q to quit."),
	}

	panel := m.styles.box(title, strings.Join(stats, "\n"), 70, 0)
	return m.renderScreen(panel, width, height, lipgloss.Center, lipgloss.Center)
}

func (m Model) renderMapPanel(width int, height int) string {
	innerWidth := max(18, width-4)
	innerHeight := max(8, height-4)
	player := m.game.Player
	floor := m.game.Floor

	cameraX := clamp(player.Pos.X-innerWidth/2, 0, max(0, floor.Width-innerWidth))
	cameraY := clamp(player.Pos.Y-innerHeight/2, 0, max(0, floor.Height-innerHeight))

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

	rows := make([]string, 0, innerHeight)
	for y := 0; y < innerHeight; y++ {
		var row strings.Builder
		for x := 0; x < innerWidth; x++ {
			pos := game.Position{X: cameraX + x, Y: cameraY + y}
			row.WriteString(m.renderMapCell(pos, enemyMap, itemMap))
		}
		rows = append(rows, row.String())
	}

	body := strings.Join(rows, "\n")
	return m.styles.box("Dungeon", body, width, height)
}

func (m Model) renderMapCell(pos game.Position, enemies map[game.Position]*game.Enemy, items map[game.Position]game.GroundItem) string {
	floor := m.game.Floor
	if !floor.InBounds(pos) {
		return " "
	}
	if m.game.Player.Pos.Equals(pos) {
		return m.styles.Player.Render("@")
	}
	if enemy, ok := enemies[pos]; ok {
		return m.styles.colorGlyph(enemy.Template.Glyph, enemy.Template.Tint, true)
	}
	if item, ok := items[pos]; ok {
		return m.styles.colorGlyph(item.Item.Glyph, item.Item.Tint, true)
	}
	if !floor.IsExplored(pos) {
		return " "
	}

	tile := floor.TileAt(pos)
	visible := floor.IsVisible(pos)
	switch tile {
	case game.TileWall:
		if visible {
			return m.styles.TileWall.Render("#")
		}
		return m.styles.TileWallSeen.Render("#")
	case game.TileFloor:
		if visible {
			return m.styles.TileFloor.Render(".")
		}
		return m.styles.TileFloorSeen.Render(".")
	case game.TileDoorClosed:
		if visible {
			return m.styles.TileDoor.Render("+")
		}
		return m.styles.TileDoorSeen.Render("+")
	case game.TileDoorOpen:
		if visible {
			return m.styles.TileDoor.Render("/")
		}
		return m.styles.TileDoorSeen.Render("/")
	case game.TileStairsDown:
		return m.styles.TileStairs.Render(">")
	default:
		return " "
	}
}

func (m Model) renderSidebar(width int, height int) string {
	if m.overlay == overlayInventory {
		return m.renderInventorySidebar(width, height)
	}

	player := m.game.Player
	statuses := "Steady"
	if len(player.Statuses) > 0 {
		labels := make([]string, 0, len(player.Statuses))
		for _, status := range player.Statuses {
			labels = append(labels, status.Label())
		}
		statuses = strings.Join(labels, ", ")
	}

	stats := []string{
		"HP   " + hpBar(player.HP, player.MaxHP(), 20) + " " + fmt.Sprintf("%d/%d", player.HP, player.MaxHP()),
		"XP   " + xpBar(player.XP, player.NextLevelXP(), 20) + " " + fmt.Sprintf("%d/%d", player.XP, player.NextLevelXP()),
		"Lvl  " + fmt.Sprintf("%d", player.Level),
		"ATK  " + fmt.Sprintf("%d", player.AttackPower()),
		"DEF  " + fmt.Sprintf("%d", player.DefensePower()),
		"Gold " + m.styles.Gold.Render(fmt.Sprintf("%d", player.Gold)),
		"State " + statuses,
		m.renderQuickHealPreview(),
	}

	floorInfo := []string{
		m.game.FloorLabel(),
		"Explored " + fmt.Sprintf("%d%%", m.game.Floor.ExploredPercent()),
		"Visible foes " + fmt.Sprintf("%d", len(m.game.VisibleEnemies())),
		"Visible loot " + fmt.Sprintf("%d", len(m.game.VisibleItems())),
		"Underfoot " + m.game.TileDescriptionUnderPlayer(),
		"",
		wrapText(m.game.Objective(), width-6),
	}

	body := strings.Join([]string{
		m.styles.Subtitle.Render("Stats"),
		strings.Join(stats, "\n"),
		"",
		m.styles.Subtitle.Render("Floor"),
		strings.Join(floorInfo, "\n"),
	}, "\n")

	return m.styles.box("Status", body, width, height)
}

func (m Model) renderInventorySidebar(width int, height int) string {
	panelHeight := max(4, height/2)
	packPanel := m.renderPackPanel(width, panelHeight)
	equippedPanel := m.renderInventoryDetailsPanel(width, panelHeight)
	return lipgloss.JoinVertical(lipgloss.Left, packPanel, equippedPanel)
}

func (m Model) renderLog(width int, height int) string {
	innerWidth := max(20, width-4)
	lines := make([]string, 0, height)
	for index := len(m.game.Log) - 1; index >= 0 && len(lines) < height-3; index-- {
		wrapped := wrapText(m.game.Log[index], innerWidth)
		parts := strings.Split(wrapped, "\n")
		for partIndex := len(parts) - 1; partIndex >= 0 && len(lines) < height-3; partIndex-- {
			lines = append(lines, parts[partIndex])
		}
	}

	for left, right := 0, len(lines)-1; left < right; left, right = left+1, right-1 {
		lines[left], lines[right] = lines[right], lines[left]
	}

	body := strings.Join(lines, "\n")
	if body == "" {
		body = "The abbey listens."
	}
	return m.styles.box("Event Log", body, width, height)
}

func (m Model) renderPackPanel(width int, height int) string {
	stacks := m.inventoryStacks()
	if len(stacks) == 0 {
		return m.styles.box("Pack", m.styles.Dim.Render("Your pack is empty."), width, height)
	}

	rows := m.inventoryRows(stacks)
	visibleRows := max(1, height-4)
	start := clamp(m.inventoryPackScroll, 0, max(0, len(rows)-visibleRows))
	end := min(len(rows), start+visibleRows)

	lines := make([]string, 0, end-start)
	for _, row := range rows[start:end] {
		if row.StackIndex == -1 {
			lines = append(lines, m.styles.Subtitle.Render(row.Category))
			continue
		}
		selected := m.inventoryPane == inventoryPanePack && row.StackIndex == clamp(m.inventoryPackCursor, 0, len(stacks)-1)
		lines = append(lines, m.renderPackLine(stacks[row.StackIndex], selected, width-4))
	}

	return m.styles.box("Pack", strings.Join(lines, "\n"), width, height)
}

func (m Model) renderInventoryDetailsPanel(width int, height int) string {
	lines := []string{
		m.renderEquippedSlot(0, "Weapon", m.game.Player.Equipment.Weapon, width-4),
		m.renderEquippedSlot(1, "Armor", m.game.Player.Equipment.Armor, width-4),
		m.renderEquippedSlot(2, "Charm", m.game.Player.Equipment.Charm, width-4),
	}
	return m.styles.box("Equipped", strings.Join(lines, "\n"), width, height)
}

func (m Model) renderPackLine(stack inventoryStack, selected bool, width int) string {
	cursor := "  "
	nameStyle := m.styles.App
	if selected {
		cursor = m.styles.MenuCursor.Render("> ")
		nameStyle = m.styles.AccentSoft
	}

	parts := []string{
		cursor,
		m.styles.colorGlyph(stack.Item.Glyph, stack.Item.Tint, true),
		" ",
		nameStyle.Render(stack.Item.Name),
		" ",
		m.styles.Quantity.Render(fmt.Sprintf("(%d)", stack.Count)),
	}

	if detail := m.renderPackItemDetail(stack.Item); detail != "" {
		parts = append(parts, " ", detail)
	}

	return lipgloss.NewStyle().Width(width).Render(strings.Join(parts, ""))
}

func (m Model) renderEquippedSlot(index int, label string, item *game.Item, width int) string {
	cursor := "  "
	labelStyle := m.styles.Muted
	nameStyle := m.styles.App
	if m.inventoryPane == inventoryPaneEquipped && index == clamp(m.inventoryEquipCursor, 0, 2) {
		cursor = m.styles.MenuCursor.Render("> ")
		labelStyle = m.styles.AccentSoft
		nameStyle = m.styles.AccentSoft
	}

	if item == nil {
		line := cursor + labelStyle.Render(label+": ") + m.styles.Dim.Render("empty")
		return lipgloss.NewStyle().Width(width).Render(line)
	}

	prefix := strings.Join([]string{
		cursor,
		labelStyle.Render(label + ": "),
		m.styles.colorGlyph(item.Glyph, item.Tint, true),
		" ",
	}, "")
	nameWidth := max(1, width-lipgloss.Width(prefix))
	name := truncateText(item.Name, nameWidth)
	line := prefix + nameStyle.Render(name)

	detail := m.renderEquippedItemDetail(*item)
	if detail != "" {
		candidate := line + " " + detail
		if lipgloss.Width(candidate) <= width {
			line = candidate
		}
	}
	return lipgloss.NewStyle().Width(width).Render(line)
}

func (m Model) renderPackItemDetail(item game.Item) string {
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

func (m Model) renderEquippedItemDetail(item game.Item) string {
	switch item.Kind {
	case game.ItemKindEquipment:
		return m.renderEquipmentStats(item)
	case game.ItemKindConsumable:
		return m.renderConsumableEffects(item)
	default:
		return ""
	}
}

func (m Model) renderConsumableEffects(item game.Item) string {
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

func (m Model) renderComparedEquipmentStats(item game.Item, equipped *game.Item) string {
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

func (m Model) renderEquipmentStats(item game.Item) string {
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

func (m Model) renderComparedToken(label string, candidate int, current int) string {
	style := m.styles.CompareEqual
	switch {
	case candidate > current:
		style = m.styles.CompareBetter
	case candidate < current:
		style = m.styles.CompareWorse
	}
	return style.Render(fmt.Sprintf("%s+%d", label, candidate))
}

func (m Model) equippedItemFor(slot game.EquipmentSlot) *game.Item {
	switch slot {
	case game.SlotWeapon:
		return m.game.Player.Equipment.Weapon
	case game.SlotArmor:
		return m.game.Player.Equipment.Armor
	case game.SlotCharm:
		return m.game.Player.Equipment.Charm
	default:
		return nil
	}
}

func (m Model) renderQuickHealPreview() string {
	item, count, ok := m.game.QuickHealPreview()
	if !ok {
		return "Heal C  " + m.styles.Dim.Render("none")
	}

	parts := []string{
		m.styles.AccentSoft.Render("Heal C"),
		"  ",
		m.styles.colorGlyph(item.Glyph, item.Tint, true),
		" ",
		lipgloss.NewStyle().Foreground(lipgloss.Color(item.Tint)).Render(item.Name),
		" ",
		m.styles.Quantity.Render(fmt.Sprintf("(%d)", count)),
	}
	if detail := m.renderConsumableEffects(item); detail != "" {
		parts = append(parts, " ", detail)
	}
	return strings.Join(parts, "")
}

func renderEquipLine(label string, item *game.Item) string {
	if item == nil {
		return label + "  none"
	}
	detail := item.DetailLine()
	if detail == "" {
		return label + "  " + item.Name
	}
	return label + "  " + item.Name + "  " + detail
}

func hpBar(current int, total int, width int) string {
	return progressBar(current, total, width, lipgloss.Color("#d16078"))
}

func xpBar(current int, total int, width int) string {
	return progressBar(current, total, width, lipgloss.Color("#5fa8d3"))
}

func progressBar(current int, total int, width int, color lipgloss.Color) string {
	if total <= 0 {
		total = 1
	}
	filled := current * width / total
	if filled > width {
		filled = width
	}
	full := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("#", filled))
	empty := lipgloss.NewStyle().Foreground(lipgloss.Color("#4a433c")).Render(strings.Repeat("-", width-filled))
	return full + empty
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
	runes := []rune(strings.ReplaceAll(text, "\n", " "))
	if len(runes) <= width {
		return string(runes)
	}
	if width <= 3 {
		return string(runes[:width])
	}
	return string(runes[:width-3]) + "..."
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
