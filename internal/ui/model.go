package ui

import (
	"errors"
	"hash/fnv"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"dunshell/internal/game"
)

type screen int

const (
	screenTitle screen = iota
	screenSeed
	screenHelp
	screenPlaying
	screenVictory
	screenGameOver
)

type overlay int

const (
	overlayNone overlay = iota
	overlayInventory
	overlayMerchant
	overlayDescend
	overlayRoute
	overlayChest
	overlayBoss
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
	titleMenuIndex       int
	seedMode             int
	overlay              overlay
	inventoryPane        inventoryPane
	inventoryPackCursor  int
	inventoryPackScroll  int
	inventoryEquipCursor int
	routeCursor          int
	merchantCursor       int
	descendChoice        int
	chestChoice          int
	bossChoice           int
	quitChoice           int
	outcomeMenuIndex     int
	lockedSeed           int64
	hasLockedSeed        bool
	styles               styles
	glyphs               glyphSet
	keys                 keyMap
	help                 help.Model
	seedInput            textinput.Model
	seedError            string
	profile              game.Profile
	storageError         string
	hasContinue          bool
	savedRun             *game.Game
	game                 *game.Game
	mapViewport          mapViewportState
}

func NewModel(seed int64, hasLockedSeed bool) *Model {
	keys := newKeyMap()
	helpModel := help.New()
	helpModel.ShowAll = false

	input := textinput.New()
	input.Placeholder = "numeric or text seed"
	input.CharLimit = 40
	input.Width = 28

	model := &Model{
		screen:        screenTitle,
		returnScreen:  screenTitle,
		lockedSeed:    seed,
		hasLockedSeed: hasLockedSeed,
		styles:        newStyles(),
		glyphs:        newGlyphSet(),
		keys:          keys,
		help:          helpModel,
		seedInput:     input,
	}
	model.loadPersistence()
	return model
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.screen == screenSeed && m.seedMode == 1 {
		m.seedInput, cmd = m.seedInput.Update(msg)
	}

	switch typed := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = typed.Width
		m.height = typed.Height
		m.help.Width = typed.Width
		return m, cmd
	case tea.KeyMsg:
		if key.Matches(typed, key.NewBinding(key.WithKeys("ctrl+c"))) {
			return m, tea.Quit
		}

		switch m.screen {
		case screenTitle:
			return m.updateTitle(typed)
		case screenSeed:
			return m.updateSeed(typed)
		case screenHelp:
			return m.updateHelp(typed)
		case screenPlaying:
			switch m.overlay {
			case overlayInventory:
				return m.updateInventory(typed)
			case overlayMerchant:
				return m.updateMerchant(typed)
			case overlayDescend:
				return m.updateDescendPrompt(typed)
			case overlayRoute:
				return m.updateRouteChoice(typed)
			case overlayChest:
				return m.updateChestPrompt(typed)
			case overlayBoss:
				return m.updateBossPrompt(typed)
			case overlayQuit:
				return m.updateQuitPrompt(typed)
			default:
				return m.updateGame(typed)
			}
		case screenVictory, screenGameOver:
			return m.updateOutcome(typed)
		}
	}
	return m, cmd
}

func (m *Model) View() string {
	switch m.screen {
	case screenTitle:
		return m.viewTitle()
	case screenSeed:
		return m.viewSeed()
	case screenHelp:
		return m.viewHelp()
	case screenPlaying:
		switch m.overlay {
		case overlayRoute:
			return m.viewRouteChoice()
		case overlayDescend:
			return m.viewDescendPrompt()
		case overlayChest:
			return m.viewChestPrompt()
		case overlayBoss:
			return m.viewBossPrompt()
		case overlayQuit:
			return m.viewQuitPrompt()
		default:
			return m.viewGame()
		}
	case screenVictory:
		return m.viewOutcome(true)
	case screenGameOver:
		return m.viewOutcome(false)
	default:
		return ""
	}
}

func (m *Model) updateTitle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	options := m.titleMenuOptions()
	switch {
	case key.Matches(msg, m.keys.Up):
		m.titleMenuIndex = (m.titleMenuIndex + len(options) - 1) % len(options)
	case key.Matches(msg, m.keys.Down):
		m.titleMenuIndex = (m.titleMenuIndex + 1) % len(options)
	case key.Matches(msg, m.keys.Help):
		m.returnScreen = screenTitle
		m.screen = screenHelp
	case key.Matches(msg, m.keys.Select):
		switch options[m.titleMenuIndex] {
		case "Continue":
			m.continueRun()
		case "New Run":
			m.beginNewRunFlow()
		case "Field Guide":
			m.returnScreen = screenTitle
			m.screen = screenHelp
		case "Quit":
			return m, tea.Quit
		}
	case key.Matches(msg, m.keys.Quit, m.keys.Back):
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) updateSeed(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.seedError = ""
		m.seedInput.Blur()
		m.screen = screenTitle
	case key.Matches(msg, m.keys.Left, m.keys.Up):
		m.seedMode = 0
		m.seedError = ""
		m.seedInput.Blur()
	case key.Matches(msg, m.keys.Right, m.keys.Down):
		m.seedMode = 1
		m.seedInput.Focus()
	case key.Matches(msg, m.keys.Select):
		if m.seedMode == 0 {
			m.startRun(time.Now().UTC().UnixNano())
			return m, nil
		}
		seedText := strings.TrimSpace(m.seedInput.Value())
		if seedText == "" {
			m.seedError = "Enter a seed or switch back to Random."
			return m, nil
		}
		m.startRun(hashSeed(seedText))
		return m, nil
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
	options := m.outcomeMenuOptions()
	switch {
	case key.Matches(msg, m.keys.Up):
		m.outcomeMenuIndex = (m.outcomeMenuIndex + len(options) - 1) % len(options)
	case key.Matches(msg, m.keys.Down):
		m.outcomeMenuIndex = (m.outcomeMenuIndex + 1) % len(options)
	case key.Matches(msg, m.keys.Select):
		switch options[m.outcomeMenuIndex] {
		case "Continue Endless":
			if m.game != nil && m.game.ContinueEndless() {
				m.screen = screenPlaying
				m.overlay = overlayNone
				m.persistRun()
			}
		case "New Run":
			m.endRunState()
			m.beginNewRunFlow()
		case "Title Screen":
			m.endRunState()
			m.screen = screenTitle
		case "Quit":
			m.endRunState()
			return m, tea.Quit
		}
	case strings.EqualFold(msg.String(), "n"):
		m.endRunState()
		m.beginNewRunFlow()
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
		m.quitChoice = 0
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
		case game.InteractionOpenChest:
			m.overlay = overlayChest
			m.chestChoice = 0
		case game.InteractionMerchant:
			m.overlay = overlayMerchant
			m.merchantCursor = 0
		case game.InteractionBossEntry:
			m.overlay = overlayBoss
			m.bossChoice = 0
		case game.InteractionDescend:
			m.overlay = overlayDescend
			m.descendChoice = 0
		default:
			m.game.AddLog("Nothing here answers your hand.")
		}
	}
	m.afterGameMutation()
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
	m.afterGameMutation()
	return m, nil
}

func (m *Model) updateMerchant(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	merchant, index := m.game.MerchantAtPlayer()
	if merchant == nil {
		m.overlay = overlayNone
		return m, nil
	}
	switch {
	case key.Matches(msg, m.keys.Back, m.keys.Inventory):
		m.overlay = overlayNone
	case key.Matches(msg, m.keys.Up):
		m.merchantCursor--
	case key.Matches(msg, m.keys.Down):
		m.merchantCursor++
	case key.Matches(msg, m.keys.Interact, m.keys.Select):
		m.merchantCursor = clamp(m.merchantCursor, 0, len(merchant.Offers)-1)
		m.game.BuyMerchantOffer(index, m.merchantCursor)
	}
	if merchant != nil {
		m.merchantCursor = clamp(m.merchantCursor, 0, max(0, len(merchant.Offers)-1))
	}
	m.afterGameMutation()
	return m, nil
}

func (m *Model) updateDescendPrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Left, m.keys.Up):
		m.descendChoice = 0
	case key.Matches(msg, m.keys.Right, m.keys.Down):
		m.descendChoice = 1
	case strings.EqualFold(msg.String(), "y"):
		m.descendChoice = 0
		m.confirmDescend()
	case strings.EqualFold(msg.String(), "n"):
		m.overlay = overlayNone
	case key.Matches(msg, m.keys.Select, m.keys.Interact):
		if m.descendChoice == 0 {
			m.confirmDescend()
		} else {
			m.overlay = overlayNone
		}
	case key.Matches(msg, m.keys.Back):
		m.overlay = overlayNone
	}
	return m, nil
}

func (m *Model) updateRouteChoice(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	routes := m.game.RouteChoices()
	if len(routes) == 0 {
		m.overlay = overlayNone
		return m, nil
	}
	switch {
	case key.Matches(msg, m.keys.Up, m.keys.Left):
		m.routeCursor = (m.routeCursor + len(routes) - 1) % len(routes)
	case key.Matches(msg, m.keys.Down, m.keys.Right):
		m.routeCursor = (m.routeCursor + 1) % len(routes)
	case key.Matches(msg, m.keys.Select, m.keys.Interact):
		m.game.DescendWithRoute(m.routeCursor)
		m.overlay = overlayNone
		m.afterGameMutation()
	case key.Matches(msg, m.keys.Back):
		m.overlay = overlayNone
	}
	return m, nil
}

func (m *Model) updateChestPrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Left, m.keys.Up):
		m.chestChoice = 0
	case key.Matches(msg, m.keys.Right, m.keys.Down):
		m.chestChoice = 1
	case key.Matches(msg, m.keys.Select, m.keys.Interact):
		if m.chestChoice == 0 {
			if _, index := m.game.ChestAtPlayer(); index >= 0 && m.game.OpenChest(index) {
				m.overlay = overlayNone
				m.afterGameMutation()
			}
		} else {
			m.overlay = overlayNone
		}
	case key.Matches(msg, m.keys.Back):
		m.overlay = overlayNone
	}
	return m, nil
}

func (m *Model) updateBossPrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Left, m.keys.Up):
		m.bossChoice = 0
	case key.Matches(msg, m.keys.Right, m.keys.Down):
		m.bossChoice = 1
	case key.Matches(msg, m.keys.Select, m.keys.Interact):
		if m.bossChoice == 0 {
			m.game.EnterBossRoom()
			m.overlay = overlayNone
			m.afterGameMutation()
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
	case key.Matches(msg, m.keys.Left, m.keys.Up):
		m.quitChoice = 0
	case key.Matches(msg, m.keys.Right, m.keys.Down):
		m.quitChoice = 1
	case strings.EqualFold(msg.String(), "y"):
		m.persistRun()
		return m, tea.Quit
	case strings.EqualFold(msg.String(), "n"):
		m.overlay = overlayNone
	case key.Matches(msg, m.keys.Select, m.keys.Interact):
		if m.quitChoice == 0 {
			m.persistRun()
			return m, tea.Quit
		}
		m.overlay = overlayNone
	case key.Matches(msg, m.keys.Back, m.keys.Quit):
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
		m.game.Unequip(m.selectedEquipmentSlot())
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

func (m *Model) confirmDescend() {
	if m.game.BeginDescendSelection() {
		m.overlay = overlayRoute
		m.routeCursor = 0
		m.persistRun()
		return
	}
	m.overlay = overlayNone
	m.afterGameMutation()
}

func (m *Model) beginNewRunFlow() {
	if m.hasLockedSeed {
		m.startRun(m.lockedSeed)
		return
	}
	m.screen = screenSeed
	m.seedMode = 0
	m.seedError = ""
	m.seedInput.SetValue("")
	m.seedInput.Blur()
	if m.seedMode == 1 {
		m.seedInput.Focus()
	}
}

func (m *Model) startRun(seed int64) {
	if seed == 0 {
		seed = time.Now().UTC().UnixNano()
	}
	m.game = game.New(seed, m.profile.Difficulty)
	m.screen = screenPlaying
	m.overlay = overlayNone
	m.inventoryPane = inventoryPanePack
	m.inventoryPackCursor = 0
	m.inventoryPackScroll = 0
	m.inventoryEquipCursor = 0
	m.routeCursor = 0
	m.merchantCursor = 0
	m.mapViewport = mapViewportState{}
	m.persistRun()
	m.hasContinue = true
	m.savedRun = m.game
}

func (m *Model) continueRun() {
	if m.savedRun == nil {
		return
	}
	m.game = m.savedRun
	m.screen = screenPlaying
	m.overlay = overlayNone
	m.inventoryPane = inventoryPanePack
	m.syncOutcome()
	if m.screen == screenPlaying {
		m.persistRun()
	}
}

func (m *Model) afterGameMutation() {
	m.syncOutcome()
	if m.screen == screenPlaying {
		m.persistRun()
	}
}

func (m *Model) syncOutcome() {
	if m.game == nil {
		return
	}
	switch m.game.Mode {
	case game.ModeWon:
		if !m.game.VictoryRecorded {
			m.profile.Wins++
			m.profile.Difficulty++
			m.game.VictoryRecorded = true
			m.saveProfile()
		}
		m.overlay = overlayNone
		m.screen = screenVictory
	case game.ModeLost:
		m.overlay = overlayNone
		m.screen = screenGameOver
		m.clearSavedRun()
	}
	if m.game != nil && (m.screen == screenVictory || m.screen == screenPlaying) {
		m.persistRun()
	}
}

func (m *Model) persistRun() {
	if m.game == nil {
		return
	}
	if err := game.SaveRun(m.game); err != nil {
		m.storageError = err.Error()
		return
	}
	m.savedRun = m.game
	m.hasContinue = true
	m.storageError = ""
}

func (m *Model) clearSavedRun() {
	_ = game.ClearRun()
	m.savedRun = nil
	m.hasContinue = false
}

func (m *Model) endRunState() {
	m.clearSavedRun()
	m.game = nil
	m.overlay = overlayNone
	m.outcomeMenuIndex = 0
	if m.screen == screenVictory || m.screen == screenGameOver {
		m.screen = screenTitle
	}
}

func (m *Model) loadPersistence() {
	profile, err := game.LoadProfile()
	if err == nil {
		m.profile = profile
	} else {
		m.storageError = err.Error()
	}
	run, err := game.LoadRun()
	if err == nil {
		m.savedRun = run
		m.hasContinue = true
		return
	}
	if !errors.Is(err, game.ErrNoRunSave) {
		m.storageError = err.Error()
	}
}

func (m *Model) saveProfile() {
	if err := game.SaveProfile(m.profile); err != nil {
		m.storageError = err.Error()
	}
}

func (m *Model) titleMenuOptions() []string {
	options := make([]string, 0, 4)
	if m.hasContinue {
		options = append(options, "Continue")
	}
	options = append(options, "New Run", "Field Guide", "Quit")
	return options
}

func (m *Model) outcomeMenuOptions() []string {
	if m.screen == screenVictory {
		return []string{"Continue Endless", "New Run", "Title Screen", "Quit"}
	}
	return []string{"New Run", "Title Screen", "Quit"}
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
	if width < 116 {
		sidebarHeight = max(10, bodyHeight-1)
	}
	packHeight, _ := m.inventoryPanelHeights(sidebarHeight)
	return max(1, packHeight-4)
}

func (m *Model) inventoryPanelHeights(sidebarHeight int) (int, int) {
	detailsHeight := 9
	if m.overlay == overlayMerchant {
		detailsHeight = 10
	}
	if detailsHeight > sidebarHeight-5 {
		detailsHeight = max(5, sidebarHeight-5)
	}
	packHeight := sidebarHeight - detailsHeight
	if packHeight < 5 {
		packHeight = 5
		detailsHeight = max(5, sidebarHeight-packHeight)
	}
	return packHeight, detailsHeight
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
		byID[item.ID] = &inventoryStack{Item: item, Count: 1, FirstIndex: index}
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
		if left.Item.Rarity != right.Item.Rarity {
			return left.Item.Rarity > right.Item.Rarity
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

func hashSeed(input string) int64 {
	if input == "" {
		return time.Now().UTC().UnixNano()
	}
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte(input))
	value := int64(hasher.Sum64())
	if value == 0 {
		return 1
	}
	return value
}
