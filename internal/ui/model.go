package ui

import (
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"dunshell/internal/game"
)

type screen int

const (
	screenTitle screen = iota
	screenHelp
	screenPlaying
	screenVictory
	screenGameOver
)

type overlay int

const (
	overlayNone overlay = iota
	overlayInventory
	overlayDescend
	overlayQuit
)

type inventoryPane int

const (
	inventoryPanePack inventoryPane = iota
	inventoryPaneEquipped
)

type rect struct {
	X int
	Y int
	W int
	H int
}

type inventoryStack struct {
	Item       game.Item
	Count      int
	FirstIndex int
}

type inventoryRow struct {
	Category   string
	StackIndex int
}

type mapViewportState struct {
	Panel   rect
	Content rect
	CameraX int
	CameraY int
}

type descendPromptState struct {
	Choice int
}

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Wait      key.Binding
	QuickHeal key.Binding
	Interact  key.Binding
	Inventory key.Binding
	Help      key.Binding
	Quit      key.Binding
	Select    key.Binding
	Back      key.Binding
	Use       key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Up:        key.NewBinding(key.WithKeys("up", "w", "W"), key.WithHelp("up/w", "move up")),
		Down:      key.NewBinding(key.WithKeys("down", "s", "S"), key.WithHelp("down/s", "move down")),
		Left:      key.NewBinding(key.WithKeys("left", "a", "A"), key.WithHelp("left/a", "move left")),
		Right:     key.NewBinding(key.WithKeys("right", "d", "D"), key.WithHelp("right/d", "move right")),
		Wait:      key.NewBinding(key.WithKeys("."), key.WithHelp(".", "wait")),
		QuickHeal: key.NewBinding(key.WithKeys("c", "C"), key.WithHelp("c", "quick heal")),
		Interact:  key.NewBinding(key.WithKeys("e", "E"), key.WithHelp("e", "interact")),
		Inventory: key.NewBinding(key.WithKeys("i", "I"), key.WithHelp("i", "inventory")),
		Help:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:      key.NewBinding(key.WithKeys("q", "Q"), key.WithHelp("q", "quit safely")),
		Select:    key.NewBinding(key.WithKeys("enter", " "), key.WithHelp("enter", "confirm")),
		Back:      key.NewBinding(key.WithKeys("esc", "backspace"), key.WithHelp("esc", "back")),
		Use:       key.NewBinding(key.WithKeys("u", "U"), key.WithHelp("u", "use item")),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Up, k.Down, k.Interact, k.QuickHeal, k.Inventory, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Wait, k.Interact, k.QuickHeal, k.Inventory},
		{k.Use, k.Help, k.Quit, k.Back},
	}
}

type Model struct {
	width                int
	height               int
	screen               screen
	returnScreen         screen
	menuIndex            int
	overlay              overlay
	inventoryPane        inventoryPane
	inventoryPackCursor  int
	inventoryPackScroll  int
	inventoryEquipCursor int
	lockedSeed           int64
	hasLockedSeed        bool
	styles               styles
	keys                 keyMap
	help                 help.Model
	game                 *game.Game
	mapViewport          mapViewportState
	descendPrompt        descendPromptState
}

func NewModel(seed int64, hasLockedSeed bool) *Model {
	keys := newKeyMap()
	helpModel := help.New()
	helpModel.ShowAll = false

	return &Model{
		screen:        screenTitle,
		returnScreen:  screenTitle,
		lockedSeed:    seed,
		hasLockedSeed: hasLockedSeed,
		styles:        newStyles(),
		keys:          keys,
		help:          helpModel,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typed := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = typed.Width
		m.height = typed.Height
		m.help.Width = typed.Width
		return m, nil
	case tea.KeyMsg:
		if key.Matches(typed, key.NewBinding(key.WithKeys("ctrl+c"))) {
			return m, tea.Quit
		}

		switch m.screen {
		case screenTitle:
			return m.updateTitle(typed)
		case screenHelp:
			return m.updateHelp(typed)
		case screenPlaying:
			switch m.overlay {
			case overlayInventory:
				return m.updateInventory(typed)
			case overlayDescend:
				return m.updateDescendPrompt(typed)
			case overlayQuit:
				return m.updateQuitPrompt(typed)
			default:
				return m.updateGame(typed)
			}
		case screenVictory, screenGameOver:
			return m.updateOutcome(typed)
		}
	}
	return m, nil
}

func (m *Model) View() string {
	switch m.screen {
	case screenTitle:
		return m.viewTitle()
	case screenHelp:
		return m.viewHelp()
	case screenPlaying:
		if m.overlay == overlayQuit {
			return m.viewQuitPrompt()
		}
		return m.viewGame()
	case screenVictory:
		return m.viewOutcome(true)
	case screenGameOver:
		return m.viewOutcome(false)
	default:
		return ""
	}
}

func (m *Model) updateTitle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		m.menuIndex = (m.menuIndex + 2) % 3
	case key.Matches(msg, m.keys.Down):
		m.menuIndex = (m.menuIndex + 1) % 3
	case key.Matches(msg, m.keys.Help):
		m.returnScreen = screenTitle
		m.screen = screenHelp
	case key.Matches(msg, m.keys.Select):
		switch m.menuIndex {
		case 0:
			m.startRun()
		case 1:
			m.returnScreen = screenTitle
			m.screen = screenHelp
		case 2:
			return m, tea.Quit
		}
	case key.Matches(msg, m.keys.Quit, m.keys.Back):
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Back, m.keys.Help, m.keys.Select) {
		m.screen = m.returnScreen
		return m, nil
	}
	return m, nil
}

func (m *Model) updateOutcome(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Select):
		m.screen = screenTitle
		m.overlay = overlayNone
		m.game = nil
	case strings.EqualFold(msg.String(), "n"):
		m.startRun()
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) updateGame(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Help):
		m.returnScreen = screenPlaying
		m.screen = screenHelp
	case key.Matches(msg, m.keys.Inventory):
		m.overlay = overlayInventory
		m.inventoryPane = inventoryPanePack
		m.inventoryPackCursor = 0
		m.inventoryPackScroll = 0
		m.inventoryEquipCursor = 0
	case key.Matches(msg, m.keys.Quit):
		m.overlay = overlayQuit
	case key.Matches(msg, m.keys.Up):
		m.game.MovePlayer(0, -1)
	case key.Matches(msg, m.keys.Down):
		m.game.MovePlayer(0, 1)
	case key.Matches(msg, m.keys.Left):
		m.game.MovePlayer(-1, 0)
	case key.Matches(msg, m.keys.Right):
		m.game.MovePlayer(1, 0)
	case key.Matches(msg, m.keys.Wait):
		m.game.WaitTurn()
	case key.Matches(msg, m.keys.QuickHeal):
		m.game.QuickHeal()
	case key.Matches(msg, m.keys.Interact):
		switch m.game.InteractionContext().Primary {
		case game.InteractionPickup:
			m.game.Pickup()
		case game.InteractionDescend:
			m.openDescendPrompt()
		default:
			m.game.AddLog("Nothing here answers your hand.")
		}
	}

	m.syncOutcome()
	return m, nil
}

func (m *Model) updateInventory(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back, m.keys.Inventory):
		m.overlay = overlayNone
	case key.Matches(msg, m.keys.Left):
		m.inventoryPane = inventoryPanePack
	case key.Matches(msg, m.keys.Right):
		m.inventoryPane = inventoryPaneEquipped
	case key.Matches(msg, m.keys.Up):
		if m.inventoryPane == inventoryPanePack {
			m.inventoryPackCursor--
		} else {
			m.inventoryEquipCursor--
		}
	case key.Matches(msg, m.keys.Down):
		if m.inventoryPane == inventoryPanePack {
			m.inventoryPackCursor++
		} else {
			m.inventoryEquipCursor++
		}
	case key.Matches(msg, m.keys.Interact, m.keys.Select):
		m.performInventoryPrimaryAction()
	case key.Matches(msg, m.keys.Use):
		m.handleInventoryUse()
	}

	m.adjustInventoryScroll()
	m.syncOutcome()
	return m, nil
}

func (m *Model) updateDescendPrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Left, m.keys.Up):
		m.descendPrompt.Choice = 0
	case key.Matches(msg, m.keys.Right, m.keys.Down):
		m.descendPrompt.Choice = 1
	case strings.EqualFold(msg.String(), "y"):
		m.descendPrompt.Choice = 0
		m.confirmDescend()
	case strings.EqualFold(msg.String(), "n"):
		m.overlay = overlayNone
	case key.Matches(msg, m.keys.Select, m.keys.Interact):
		if m.descendPrompt.Choice == 0 {
			m.confirmDescend()
		} else {
			m.overlay = overlayNone
		}
	case key.Matches(msg, m.keys.Back):
		m.overlay = overlayNone
	}
	return m, nil
}

func (m *Model) updateQuitPrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case strings.EqualFold(msg.String(), "y"):
		return m, tea.Quit
	case strings.EqualFold(msg.String(), "n"):
		m.overlay = overlayNone
	}
	if key.Matches(msg, m.keys.Back, m.keys.Quit) {
		m.overlay = overlayNone
	}
	return m, nil
}

func (m *Model) performInventoryPrimaryAction() {
	if m.inventoryPane == inventoryPaneEquipped {
		slot := m.selectedEquipmentSlot()
		m.game.Unequip(slot)
		return
	}

	item, ok := m.selectedPackItem()
	if !ok {
		return
	}

	switch item.Kind {
	case game.ItemKindEquipment:
		m.handleInventoryEquip()
	case game.ItemKindConsumable:
		m.handleInventoryUse()
	}
}

func (m *Model) handleInventoryEquip() {
	if m.inventoryPane == inventoryPaneEquipped {
		slot := m.selectedEquipmentSlot()
		m.game.Unequip(slot)
		return
	}

	index := m.selectedPackIndex()
	if index >= 0 {
		m.game.EquipItem(index)
	}
}

func (m *Model) handleInventoryUse() {
	if m.inventoryPane == inventoryPaneEquipped {
		return
	}
	index := m.selectedPackIndex()
	if index >= 0 {
		m.game.UseItem(index)
	}
}

func (m *Model) openDescendPrompt() {
	m.overlay = overlayDescend
	m.descendPrompt.Choice = 0
}

func (m *Model) confirmDescend() {
	m.overlay = overlayNone
	m.game.Descend()
	m.syncOutcome()
}

func (m *Model) startRun() {
	seed := m.lockedSeed
	if !m.hasLockedSeed {
		seed = time.Now().UTC().UnixNano()
	}

	m.game = game.New(seed)
	m.screen = screenPlaying
	m.overlay = overlayNone
	m.inventoryPane = inventoryPanePack
	m.inventoryPackCursor = 0
	m.inventoryPackScroll = 0
	m.inventoryEquipCursor = 0
	m.mapViewport = mapViewportState{}
	m.descendPrompt = descendPromptState{}
}

func (m *Model) syncOutcome() {
	if m.game == nil {
		return
	}

	switch m.game.Mode {
	case game.ModeWon:
		m.overlay = overlayNone
		m.screen = screenVictory
	case game.ModeLost:
		m.overlay = overlayNone
		m.screen = screenGameOver
	}
}

func (m *Model) selectedEquipmentSlot() game.EquipmentSlot {
	switch clamp(m.inventoryEquipCursor, 0, 2) {
	case 0:
		return game.SlotWeapon
	case 1:
		return game.SlotArmor
	case 2:
		return game.SlotCharm
	default:
		return game.SlotNone
	}
}

func (m *Model) selectedPackIndex() int {
	stack, ok := m.selectedPackStack()
	if !ok {
		return -1
	}
	return stack.FirstIndex
}

func (m *Model) selectedPackItem() (game.Item, bool) {
	stack, ok := m.selectedPackStack()
	if !ok {
		return game.Item{}, false
	}
	return stack.Item, true
}

func (m *Model) selectedPackStack() (inventoryStack, bool) {
	stacks := m.inventoryStacks()
	if len(stacks) == 0 {
		return inventoryStack{}, false
	}
	index := clamp(m.inventoryPackCursor, 0, len(stacks)-1)
	return stacks[index], true
}

func (m *Model) adjustInventoryScroll() {
	m.inventoryEquipCursor = clamp(m.inventoryEquipCursor, 0, 2)

	stacks := m.inventoryStacks()
	if len(stacks) == 0 {
		m.inventoryPackCursor = 0
		m.inventoryPackScroll = 0
		return
	}

	m.inventoryPackCursor = clamp(m.inventoryPackCursor, 0, len(stacks)-1)

	rows := m.inventoryRows(stacks)
	selectedLine := m.selectedPackLine(rows)
	visibleRows := m.inventoryVisiblePackRows()
	maxScroll := max(0, len(rows)-visibleRows)
	if selectedLine < m.inventoryPackScroll {
		m.inventoryPackScroll = selectedLine
	}
	if selectedLine >= m.inventoryPackScroll+visibleRows {
		m.inventoryPackScroll = selectedLine - visibleRows + 1
	}
	m.inventoryPackScroll = clamp(m.inventoryPackScroll, 0, maxScroll)
}

func (m *Model) inventoryVisiblePackRows() int {
	width, _, bodyHeight, _, _ := m.gameplayMetrics()
	sidebarHeight := bodyHeight
	if width < 112 {
		sidebarHeight = max(10, bodyHeight-1)
	}
	packHeight, _ := m.inventoryPanelHeights(sidebarHeight)
	return max(1, packHeight-4)
}

func (m *Model) inventoryPanelHeights(sidebarHeight int) (int, int) {
	equippedHeight := 9
	switch {
	case sidebarHeight >= 19:
		equippedHeight = 10
	case sidebarHeight <= 12:
		equippedHeight = 7
	case sidebarHeight <= 15:
		equippedHeight = 8
	}

	if equippedHeight > sidebarHeight-5 {
		equippedHeight = max(5, sidebarHeight-5)
	}

	packHeight := sidebarHeight - equippedHeight
	if packHeight < 5 {
		packHeight = 5
		equippedHeight = max(5, sidebarHeight-packHeight)
	}

	return packHeight, equippedHeight
}

func (m *Model) inventoryStacks() []inventoryStack {
	if m.game == nil {
		return nil
	}

	byID := make(map[string]*inventoryStack)
	order := make([]string, 0, len(m.game.Player.Inventory))
	for index, item := range m.game.Player.Inventory {
		if stack, ok := byID[item.ID]; ok {
			stack.Count++
			continue
		}
		byID[item.ID] = &inventoryStack{
			Item:       item,
			Count:      1,
			FirstIndex: index,
		}
		order = append(order, item.ID)
	}

	stacks := make([]inventoryStack, 0, len(order))
	for _, id := range order {
		stacks = append(stacks, *byID[id])
	}

	sort.Slice(stacks, func(i int, j int) bool {
		left := stacks[i]
		right := stacks[j]
		leftOrder := inventoryCategoryOrder(left.Item)
		rightOrder := inventoryCategoryOrder(right.Item)
		if leftOrder != rightOrder {
			return leftOrder < rightOrder
		}
		return left.Item.Name < right.Item.Name
	})

	return stacks
}

func (m *Model) inventoryRows(stacks []inventoryStack) []inventoryRow {
	rows := make([]inventoryRow, 0, len(stacks)+4)
	lastCategory := ""
	for index, stack := range stacks {
		category := inventoryCategoryLabel(stack.Item)
		if category != lastCategory {
			rows = append(rows, inventoryRow{Category: category, StackIndex: -1})
			lastCategory = category
		}
		rows = append(rows, inventoryRow{StackIndex: index})
	}
	return rows
}

func (m *Model) selectedPackLine(rows []inventoryRow) int {
	packIndex := 0
	for lineIndex, row := range rows {
		if row.StackIndex == -1 {
			continue
		}
		if packIndex == m.inventoryPackCursor {
			return lineIndex
		}
		packIndex++
	}
	return 0
}

func (m *Model) equippedItemFor(slot game.EquipmentSlot) *game.Item {
	if m.game == nil {
		return nil
	}

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

func inventoryCategoryLabel(item game.Item) string {
	switch item.Kind {
	case game.ItemKindConsumable:
		return "Consumables"
	case game.ItemKindEquipment:
		switch item.Slot {
		case game.SlotWeapon:
			return "Weapons"
		case game.SlotArmor:
			return "Armor"
		case game.SlotCharm:
			return "Charms"
		}
	case game.ItemKindRelic:
		return "Relics"
	}
	return "Misc"
}

func inventoryCategoryOrder(item game.Item) int {
	switch item.Kind {
	case game.ItemKindConsumable:
		return 0
	case game.ItemKindEquipment:
		switch item.Slot {
		case game.SlotWeapon:
			return 1
		case game.SlotArmor:
			return 2
		case game.SlotCharm:
			return 3
		}
	case game.ItemKindRelic:
		return 4
	}
	return 5
}
