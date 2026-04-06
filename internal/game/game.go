package game

import "strings"

const (
	GameTitle   = "Dunshell"
	GameVersion = "0.4.0"
)

type Summary struct {
	Seed                 int64
	Floor                int
	Level                int
	Gold                 int
	Kills                int
	Turn                 int
	RecoveredRelic       bool
	PersistentDifficulty int
	Endless              bool
}

type NewGameOptions struct {
	Seed                 int64
	PersistentDifficulty int
	GodMode              bool
}

type InteractionKind int

const (
	InteractionNone InteractionKind = iota
	InteractionPickup
	InteractionOpenChest
	InteractionMerchant
	InteractionBossEntry
	InteractionDescend
)

func (k InteractionKind) Label() string {
	switch k {
	case InteractionPickup:
		return "Gather"
	case InteractionOpenChest:
		return "Open Chest"
	case InteractionMerchant:
		return "Trade"
	case InteractionBossEntry:
		return "Enter Lair"
	case InteractionDescend:
		return "Descend"
	default:
		return "Interact"
	}
}

type InteractionContext struct {
	Primary   InteractionKind
	Secondary []InteractionKind
}

type Game struct {
	Title                string
	Seed                 int64
	GodMode              bool
	Mode                 GameMode
	FloorIndex           int
	MaxFloors            int
	Turn                 int
	Log                  []string
	Player               *Player
	Floor                *Floor
	PersistentDifficulty int
	Endless              bool
	PendingRoutes        []RouteChoice
	VictoryRecorded      bool
	rng                  *RNG
	nextEnemyID          int
	nextChestID          int
	nextMerchantID       int
}

func New(options NewGameOptions) *Game {
	if options.Seed == 0 {
		options.Seed = 1
	}
	rng := NewRNG(options.Seed)
	weapon, armor, charm := StarterEquipment()
	player := &Player{
		Pos:         Position{},
		BaseMaxHP:   30,
		BaseAttack:  4,
		BaseDefense: 1,
		Level:       1,
		Inventory:   StarterInventory(),
		Equipment: Equipment{
			Weapon: &weapon,
			Armor:  &armor,
			Charm:  &charm,
		},
	}
	player.HP = player.MaxHP()

	game := &Game{
		Title:                GameTitle,
		Seed:                 options.Seed,
		GodMode:              options.GodMode,
		Mode:                 ModePlaying,
		FloorIndex:           1,
		MaxFloors:            20,
		Player:               player,
		PersistentDifficulty: options.PersistentDifficulty,
		rng:                  rng,
		Log:                  make([]string, 0, 256),
	}
	if game.GodMode {
		game.enableGodModeLoadout()
	}

	game.Floor = GenerateFloor(game.rng, game.FloorIndex, game.MaxFloors, game.PersistentDifficulty, FloorModifier{}, false, &game.nextEnemyID, &game.nextChestID, &game.nextMerchantID)
	game.Player.Pos = game.Floor.Entrance
	ComputeFOV(game.Floor, game.Player.Pos, game.Player.VisionRadius())
	if game.GodMode {
		game.AddLog("GOD MODE: the abbey opens for testing, not judgment.")
	}
	game.AddLog("The abbey doors seal behind you. The Cinder Crown waits twenty floors below.")
	if game.PersistentDifficulty > 0 {
		game.AddLog("The abbey remembers your last victory. Its hatred is keener now.")
	}
	game.AddLog(game.floorIntro())
	return game
}

func (g *Game) enableGodModeLoadout() {
	if g == nil || g.Player == nil {
		return
	}
	weapon, armor, charm := godModeEquipment()
	g.Player.Level = 20
	g.Player.XP = 0
	g.Player.BaseMaxHP = 180
	g.Player.BaseAttack = 30
	g.Player.BaseDefense = 18
	g.Player.Gold = 999
	g.Player.Keys = KeyRing{Bronze: 9, Silver: 9, Gold: 9}
	g.Player.Inventory = godModeInventory()
	g.Player.Statuses = nil
	g.Player.Equipment = Equipment{
		Weapon: &weapon,
		Armor:  &armor,
		Charm:  &charm,
	}
	g.Player.HP = g.Player.MaxHP()
}

func (g *Game) restoreGodModeState() {
	if g == nil || !g.GodMode || g.Player == nil {
		return
	}
	if g.Player.HP <= 0 {
		g.Player.HP = g.Player.MaxHP()
	}
	nextStatuses := g.Player.Statuses[:0]
	for _, status := range g.Player.Statuses {
		if status.Harmful() {
			continue
		}
		nextStatuses = append(nextStatuses, status)
	}
	g.Player.Statuses = nextStatuses
	g.Player.ClampHP()
}

func (g *Game) Summary() Summary {
	return Summary{
		Seed:                 g.Seed,
		Floor:                g.FloorIndex,
		Level:                g.Player.Level,
		Gold:                 g.Player.Gold,
		Kills:                g.Player.Kills,
		Turn:                 g.Turn,
		RecoveredRelic:       g.Player.HasRelic,
		PersistentDifficulty: g.PersistentDifficulty,
		Endless:              g.Endless,
	}
}

func (g *Game) Objective() string {
	if g.Player.HasRelic && g.Endless {
		return "Push deeper into the ashfall depths. The run no longer has a floor."
	}
	if g.Player.HasRelic {
		return "The crown is yours. Choose whether to descend into the endless dark."
	}
	if isBossFloor(g.FloorIndex, g.MaxFloors, g.Endless) {
		boss := g.Floor.ActiveBoss()
		if boss != nil {
			return "Break the floor's keeper and claim its reward chest."
		}
		return "Find the sealed boss room and enter when the floor is spent."
	}
	return "Descend to floor 20, survive the Ashen Prior, and claim the Cinder Crown."
}

func (g *Game) FloorLabel() string {
	label := "Floor " + itoa(g.FloorIndex) + "  " + g.Floor.Theme
	if g.Floor.Modifier.HasEffect() {
		label += "  ·  " + g.Floor.Modifier.Label()
	}
	if g.Endless && g.FloorIndex > g.MaxFloors {
		label += "  ·  Endless"
	}
	return label
}

func (g *Game) RouteChoices() []RouteChoice {
	return append([]RouteChoice(nil), g.PendingRoutes...)
}

func (g *Game) MovePlayer(dx int, dy int) bool {
	if g.Mode != ModePlaying {
		return false
	}

	target := g.Player.Pos.Offset(dx, dy)
	if !g.Floor.InBounds(target) {
		return false
	}

	if enemy := g.Floor.EnemyAt(target); enemy != nil {
		g.playerAttack(enemy)
		if g.Mode == ModePlaying {
			g.advanceTurn()
		}
		return true
	}

	tile := g.Floor.TileAt(target)
	switch tile {
	case TileWall:
		g.AddLog("Your shoulder meets cold stone.")
		return false
	case TileDoorClosed:
		g.Floor.OpenDoor(target)
		g.Player.Pos = target
		g.AddLog("You ease open an ironwood door.")
		g.advanceTurn()
		return true
	case TileBossGate:
		g.AddLog("A blood-locked gate bars the chamber. Press E when you are ready to be sealed in.")
		return false
	case TileBossSeal:
		g.AddLog("The gate will not yield until the keeper falls.")
		return false
	case TileFloor, TileDoorOpen, TileStairsDown:
		g.Player.Pos = target
		if tile == TileStairsDown {
			g.AddLog("The stair mouth yawns below you.")
		}
		g.advanceTurn()
		return true
	default:
		return false
	}
}

func (g *Game) WaitTurn() bool {
	if g.Mode != ModePlaying {
		return false
	}
	g.AddLog("You hold your ground and listen.")
	g.advanceTurn()
	return true
}

func (g *Game) Pickup() bool {
	if g.Mode != ModePlaying {
		return false
	}
	indices := g.Floor.ItemIndicesAt(g.Player.Pos)
	if len(indices) == 0 {
		g.AddLog("Nothing here but damp dust.")
		return false
	}
	for index := len(indices) - 1; index >= 0; index-- {
		item := g.Floor.RemoveItemAt(indices[index])
		g.gainItem(item, "You gather ")
	}
	if g.Mode == ModePlaying {
		g.advanceTurn()
	}
	return true
}

func (g *Game) BeginDescendSelection() bool {
	if g.Mode != ModePlaying {
		return false
	}
	if !g.CanDescendHere() {
		g.AddLog("No stair waits beneath your boots.")
		return false
	}
	if len(g.PendingRoutes) == 0 {
		nextFloor := g.FloorIndex + 1
		g.PendingRoutes = GenerateRouteChoices(g.rng, nextFloor, g.MaxFloors, g.Endless)
	}
	return true
}

func (g *Game) DescendWithRoute(index int) bool {
	if g.Mode != ModePlaying || !g.CanDescendHere() || len(g.PendingRoutes) == 0 {
		return false
	}
	index = clamp(index, 0, len(g.PendingRoutes)-1)
	choice := g.PendingRoutes[index]
	g.PendingRoutes = nil

	g.FloorIndex++
	g.Floor = GenerateFloor(g.rng, g.FloorIndex, g.MaxFloors, g.PersistentDifficulty, choice.Modifier, g.Endless, &g.nextEnemyID, &g.nextChestID, &g.nextMerchantID)
	g.Player.Pos = g.Floor.Entrance
	g.Player.HP = min(g.Player.MaxHP(), g.Player.HP+4)
	g.tickPlayerFloorStatuses()
	g.applyFloorArrival(choice.Modifier)
	ComputeFOV(g.Floor, g.Player.Pos, g.Player.VisionRadius())
	g.AddLog("You descend by way of " + strings.ToLower(choice.Title) + ".")
	g.AddLog(g.floorIntro())
	return true
}

func (g *Game) ContinueEndless() bool {
	if !g.Player.HasRelic {
		return false
	}
	g.Endless = true
	g.Mode = ModePlaying
	if g.FloorIndex <= g.MaxFloors {
		g.FloorIndex = g.MaxFloors + 1
	} else {
		g.FloorIndex++
	}
	choice := GenerateRouteChoices(g.rng, g.FloorIndex, g.MaxFloors, true)[g.rng.Intn(3)]
	g.Floor = GenerateFloor(g.rng, g.FloorIndex, g.MaxFloors, g.PersistentDifficulty, choice.Modifier, true, &g.nextEnemyID, &g.nextChestID, &g.nextMerchantID)
	g.Player.Pos = g.Floor.Entrance
	g.Player.HP = min(g.Player.MaxHP(), g.Player.HP+max(6, g.Player.MaxHP()/5))
	g.PendingRoutes = nil
	g.AddLog("You keep the crown and go deeper. The abbey opens a throat beneath the sanctum.")
	g.tickPlayerFloorStatuses()
	g.applyFloorArrival(choice.Modifier)
	ComputeFOV(g.Floor, g.Player.Pos, g.Player.VisionRadius())
	g.AddLog(g.floorIntro())
	return true
}

func (g *Game) HasLootHere() bool {
	return len(g.Floor.ItemIndicesAt(g.Player.Pos)) > 0
}

func (g *Game) CanOpenChestHere() bool {
	chest, _ := g.Floor.ChestAt(g.Player.Pos)
	return chest != nil && !chest.Opened
}

func (g *Game) CanTradeHere() bool {
	merchant, _ := g.Floor.MerchantAt(g.Player.Pos)
	return merchant != nil
}

func (g *Game) CanEnterBossRoom() bool {
	return g.Floor.BossGateNearby(g.Player.Pos)
}

func (g *Game) CanDescendHere() bool {
	return g.Player.Pos.Equals(g.Floor.Stairs) && (g.FloorIndex < g.MaxFloors || g.Endless)
}

func (g *Game) InteractionContext() InteractionContext {
	context := InteractionContext{}
	switch {
	case g.HasLootHere():
		context.Primary = InteractionPickup
		if g.CanOpenChestHere() {
			context.Secondary = append(context.Secondary, InteractionOpenChest)
		}
		if g.CanTradeHere() {
			context.Secondary = append(context.Secondary, InteractionMerchant)
		}
		if g.CanEnterBossRoom() {
			context.Secondary = append(context.Secondary, InteractionBossEntry)
		}
		if g.CanDescendHere() {
			context.Secondary = append(context.Secondary, InteractionDescend)
		}
	case g.CanOpenChestHere():
		context.Primary = InteractionOpenChest
	case g.CanTradeHere():
		context.Primary = InteractionMerchant
	case g.CanEnterBossRoom():
		context.Primary = InteractionBossEntry
	case g.CanDescendHere():
		context.Primary = InteractionDescend
	}
	return context
}

func (g *Game) ChestAtPlayer() (*Chest, int) {
	return g.Floor.ChestAt(g.Player.Pos)
}

func (g *Game) MerchantAtPlayer() (*Merchant, int) {
	return g.Floor.MerchantAt(g.Player.Pos)
}

func (g *Game) BossPreview() *Enemy {
	if g.Floor.Boss == nil {
		return nil
	}
	return g.Floor.EnemyByID(g.Floor.Boss.BossID)
}

func (g *Game) ActiveBoss() *Enemy {
	if g.Floor == nil {
		return nil
	}
	return g.Floor.ActiveBoss()
}

func (g *Game) EnterBossRoom() bool {
	if g.Mode != ModePlaying || g.Floor.Boss == nil || !g.CanEnterBossRoom() {
		return false
	}
	g.Player.Pos = g.Floor.Boss.Entry
	g.Floor.Boss.Active = true
	for _, door := range g.bossRoomDoors() {
		g.Floor.SetTile(door, TileBossSeal)
	}
	g.AddLog("The gate slams shut behind you. The chamber belongs to the keeper now.")
	g.advanceTurn()
	return true
}

func (g *Game) OpenChest(index int) bool {
	if g.Mode != ModePlaying {
		return false
	}
	if index < 0 || index >= len(g.Floor.Chests) {
		return false
	}
	chest := &g.Floor.Chests[index]
	if chest.Opened {
		g.AddLog("The chest is already empty.")
		return false
	}
	if chest.Locked {
		g.AddLog("The chest stays sealed until the keeper is broken.")
		return false
	}
	if !g.Player.Keys.Spend(chest.Tier) {
		g.AddLog("You need a " + chest.Tier.LowerLabel() + " key for this chest.")
		return false
	}
	chest.Opened = true
	g.AddLog("You break the " + chest.Tier.LowerLabel() + " seal and lift the lid.")
	for _, reward := range chest.Rewards {
		if reward.Kind == RewardGold {
			g.Player.Gold += reward.Gold
			g.AddLog("You claim " + itoa(reward.Gold) + " gold.")
			continue
		}
		g.gainItem(reward.Item, "You claim ")
	}
	if g.Mode == ModePlaying {
		g.advanceTurn()
	}
	return true
}

func (g *Game) BuyMerchantOffer(merchantIndex int, offerIndex int) bool {
	if g.Mode != ModePlaying || merchantIndex < 0 || merchantIndex >= len(g.Floor.Merchants) {
		return false
	}
	merchant := &g.Floor.Merchants[merchantIndex]
	if offerIndex < 0 || offerIndex >= len(merchant.Offers) {
		return false
	}
	offer := &merchant.Offers[offerIndex]
	if offer.Sold {
		g.AddLog("That stock is already gone.")
		return false
	}
	if g.Player.Gold < offer.Price {
		g.AddLog("You do not have enough gold for " + offer.Item.Name + ".")
		return false
	}
	g.Player.Gold -= offer.Price
	offer.Sold = true
	g.gainItem(offer.Item, "You buy ")
	g.advanceTurn()
	return true
}

func (g *Game) UseItem(index int) bool {
	if g.Mode != ModePlaying || index < 0 || index >= len(g.Player.Inventory) {
		return false
	}
	item := g.Player.Inventory[index]
	if item.Kind != ItemKindConsumable {
		g.AddLog(item.Name + " is not something you can drink or throw.")
		return false
	}

	consumed := false
	switch {
	case item.Heal > 0:
		before := g.Player.HP
		g.Player.HP = min(g.Player.MaxHP(), g.Player.HP+item.Heal)
		cured := g.consumeCurativeEffects(item)
		healed := g.Player.HP - before
		if healed == 0 && len(cured) == 0 {
			g.AddLog("The draught would be wasted right now.")
			return false
		}
		if len(cured) > 0 && healed > 0 {
			g.AddLog("You recover " + itoa(healed) + " HP and clear " + statusList(cured) + " with " + item.Name + ".")
		} else if len(cured) > 0 {
			g.AddLog("You clear " + statusList(cured) + " with " + item.Name + ".")
		} else {
			g.AddLog("You recover " + itoa(healed) + " HP with " + item.Name + ".")
		}
		consumed = true
	case item.FocusFloors > 0 || item.FocusTurns > 0:
		floors := item.FocusFloors
		if floors <= 0 {
			floors = 1
		}
		g.Player.ApplyStatusByFloor(StatusFocus, floors, item.FocusBonus)
		if floors == 1 {
			g.AddLog("Your blood runs hot. The tonic's edge will hold for this floor.")
		} else {
			g.AddLog("Your blood runs hot. The tonic's edge will hold for the next " + itoa(floors) + " floors.")
		}
		consumed = true
	case item.EmberDamage > 0:
		target := g.nearestVisibleEnemy()
		if target == nil {
			g.AddLog("The phial finds no mark.")
			return false
		}
		damage := item.EmberDamage + g.rng.Intn(4)
		target.HP -= damage
		g.AddLog("Ember arcs into " + target.DisplayName() + " for " + itoa(damage) + ".")
		if target.HP <= 0 {
			g.killEnemy(target, "The ember blast shatters "+target.DisplayName()+".")
		}
		consumed = true
	}

	if !consumed {
		return false
	}
	g.Player.Inventory = append(g.Player.Inventory[:index], g.Player.Inventory[index+1:]...)
	if g.Mode == ModePlaying {
		g.advanceTurn()
	}
	return true
}

func (g *Game) QuickHealPreview() (Item, int, bool) {
	_, item, count, ok := g.quickHealCandidate()
	return item, count, ok
}

func (g *Game) QuickHeal() bool {
	index, _, _, ok := g.quickHealCandidate()
	if !ok {
		g.AddLog("You have no healing consumables left.")
		return false
	}
	return g.UseItem(index)
}

func (g *Game) EquipItem(index int) bool {
	if g.Mode != ModePlaying || index < 0 || index >= len(g.Player.Inventory) {
		return false
	}
	item := g.Player.Inventory[index]
	if item.Kind != ItemKindEquipment {
		g.AddLog("That belongs in the pack, not on your body.")
		return false
	}

	oldMax := g.Player.MaxHP()
	g.Player.Inventory = append(g.Player.Inventory[:index], g.Player.Inventory[index+1:]...)
	switch item.Slot {
	case SlotWeapon:
		if g.Player.Equipment.Weapon != nil {
			g.Player.Inventory = append(g.Player.Inventory, *g.Player.Equipment.Weapon)
		}
		equipped := item
		g.Player.Equipment.Weapon = &equipped
	case SlotArmor:
		if g.Player.Equipment.Armor != nil {
			g.Player.Inventory = append(g.Player.Inventory, *g.Player.Equipment.Armor)
		}
		equipped := item
		g.Player.Equipment.Armor = &equipped
	case SlotCharm:
		if g.Player.Equipment.Charm != nil {
			g.Player.Inventory = append(g.Player.Inventory, *g.Player.Equipment.Charm)
		}
		equipped := item
		g.Player.Equipment.Charm = &equipped
	}
	newMax := g.Player.MaxHP()
	if newMax > oldMax {
		g.Player.HP += newMax - oldMax
	}
	g.Player.ClampHP()
	g.AddLog("You equip " + item.Name + ".")
	g.advanceTurn()
	return true
}

func (g *Game) Unequip(slot EquipmentSlot) bool {
	if g.Mode != ModePlaying {
		return false
	}
	var item *Item
	switch slot {
	case SlotWeapon:
		item = g.Player.Equipment.Weapon
		g.Player.Equipment.Weapon = nil
	case SlotArmor:
		item = g.Player.Equipment.Armor
		g.Player.Equipment.Armor = nil
	case SlotCharm:
		item = g.Player.Equipment.Charm
		g.Player.Equipment.Charm = nil
	}
	if item == nil {
		g.AddLog("That slot is already empty.")
		return false
	}
	g.Player.Inventory = append(g.Player.Inventory, *item)
	g.Player.ClampHP()
	g.AddLog("You stow " + item.Name + ".")
	g.advanceTurn()
	return true
}

func (g *Game) VisibleEnemies() []*Enemy {
	visible := make([]*Enemy, 0, len(g.Floor.Enemies))
	for _, enemy := range g.Floor.Enemies {
		if g.Floor.IsVisible(enemy.Pos) {
			visible = append(visible, enemy)
		}
	}
	return visible
}

func (g *Game) VisibleItems() []GroundItem {
	items := make([]GroundItem, 0, len(g.Floor.Items))
	for _, item := range g.Floor.Items {
		if g.Floor.IsVisible(item.Pos) {
			items = append(items, item)
		}
	}
	return items
}

func (g *Game) TileDescriptionUnderPlayer() string {
	if chest, _ := g.Floor.ChestAt(g.Player.Pos); chest != nil && !chest.Opened {
		return chest.Tier.LowerLabel() + " chest"
	}
	if merchant, _ := g.Floor.MerchantAt(g.Player.Pos); merchant != nil {
		return merchant.Name
	}
	tile := g.Floor.TileAt(g.Player.Pos).Name()
	if items := g.Floor.ItemIndicesAt(g.Player.Pos); len(items) > 0 {
		item := g.Floor.Items[items[0]].Item
		return tile + " with " + item.Name
	}
	return tile
}

func (g *Game) AddLog(message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		return
	}
	g.Log = append(g.Log, message)
	if len(g.Log) > 240 {
		g.Log = g.Log[len(g.Log)-240:]
	}
}

func (g *Game) floorIntro() string {
	intro := FloorIntro(g.FloorIndex)
	if g.Floor.Modifier.HasEffect() {
		intro += " Route omen: " + g.Floor.Modifier.Summary
	}
	return intro
}

func (g *Game) playerAttack(enemy *Enemy) {
	damage := g.damageRoll(g.Player.AttackPower(), enemy.DefensePower())
	enemy.HP -= damage
	g.AddLog("You strike " + enemy.DisplayName() + " for " + itoa(damage) + ".")
	if enemy.HP <= 0 {
		g.killEnemy(enemy, "You finish "+enemy.DisplayName()+".")
	}
}

func (g *Game) killEnemy(enemy *Enemy, deathMessage string) {
	g.Floor.RemoveEnemyByID(enemy.ID)
	g.Player.Kills++
	g.AddLog(deathMessage)

	gold := enemy.Template.GoldMin
	if enemy.Template.GoldMax > enemy.Template.GoldMin {
		gold += g.rng.Intn(enemy.Template.GoldMax - enemy.Template.GoldMin + 1)
	}
	if g.Floor.Modifier.BonusGold > 0 {
		gold += int(float64(gold) * g.Floor.Modifier.BonusGold)
	}
	g.Player.Gold += gold
	g.AddLog(enemy.DisplayName() + " drops " + itoa(gold) + " gold.")

	uniqueDropped := false
	if item, ok := RandomUniqueEnemyDrop(g.rng, g.FloorIndex, enemy.Elite, enemy.Template.BossTier); ok {
		g.Floor.Items = append(g.Floor.Items, GroundItem{
			Pos:       enemy.Pos,
			Item:      item,
			RoomIndex: g.Floor.RoomIndexAt(enemy.Pos),
		})
		g.AddLog(enemy.DisplayName() + " leaves behind a blood-rose gleam: " + item.Name + ".")
		uniqueDropped = true
	}

	bonusDrop := enemy.Elite || enemy.Template.BossTier > 0
	if !uniqueDropped {
		if item, ok := RandomDropItem(g.rng, g.FloorIndex, g.Floor.Modifier, bonusDrop); ok {
			g.Floor.Items = append(g.Floor.Items, GroundItem{
				Pos:       enemy.Pos,
				Item:      item,
				RoomIndex: g.Floor.RoomIndexAt(enemy.Pos),
			})
			g.AddLog(enemy.DisplayName() + " leaves behind " + item.Name + ".")
		}
	}
	if bonusDrop && g.rng.Float64() < 0.22 {
		key := RandomKeyReward(g.rng, g.FloorIndex)
		g.Floor.Items = append(g.Floor.Items, GroundItem{Pos: enemy.Pos, Item: key, RoomIndex: g.Floor.RoomIndexAt(enemy.Pos)})
		g.AddLog(enemy.DisplayName() + " spills a " + strings.ToLower(key.Name) + ".")
	}
	if g.Player.GainXP(enemyXPProgress(enemy, g.Player, g.FloorIndex)) {
		g.AddLog("You rise to level " + itoa(g.Player.Level) + ".")
	}

	if g.Floor.Boss != nil && enemy.ID == g.Floor.Boss.BossID {
		g.Floor.Boss.Cleared = true
		g.Floor.Boss.Active = false
		for _, door := range g.bossRoomDoors() {
			g.Floor.SetTile(door, TileDoorOpen)
		}
		if chestIndex := g.findChestIndex(g.Floor.Boss.RewardChestID); chestIndex >= 0 {
			chest := &g.Floor.Chests[chestIndex]
			chest.Locked = false
			key := KeyItem(chest.Tier)
			g.gainItem(key, "The keeper drops ")
			g.AddLog("The boss chest unlocks with a deep iron click.")
		}
	}
}

func (g *Game) bossRoomDoors() []Position {
	if g == nil || g.Floor == nil || g.Floor.Boss == nil {
		return nil
	}
	doors := append([]Position(nil), g.Floor.RoomDoors[g.Floor.Boss.RoomIndex]...)
	if !containsPosition(doors, g.Floor.Boss.Gate) {
		doors = append(doors, g.Floor.Boss.Gate)
	}
	return doors
}

func (g *Game) advanceTurn() {
	if g.Mode != ModePlaying {
		return
	}
	g.Turn++
	g.runEnemyTurns()
	if g.Mode != ModePlaying {
		return
	}
	g.tickPlayerStatuses()
	if g.Player.HP <= 0 {
		g.Mode = ModeLost
		g.AddLog("The abbey finally claims you.")
		return
	}
	ComputeFOV(g.Floor, g.Player.Pos, g.Player.VisionRadius())
}

func (g *Game) runEnemyTurns() {
	for _, enemy := range append([]*Enemy(nil), g.Floor.Enemies...) {
		if g.Mode != ModePlaying {
			return
		}
		if !enemy.IsAlive() {
			continue
		}
		if g.Floor.Boss != nil && enemy.ID == g.Floor.Boss.BossID && !g.Floor.Boss.Active {
			continue
		}
		if enemy.Cooldown > 0 {
			enemy.Cooldown--
		}
		if enemy.Template.EnrageThreshold > 0 && !enemy.Enraged && enemy.HP*100 <= enemy.Template.MaxHP*enemy.Template.EnrageThreshold {
			enemy.Enraged = true
			if g.Floor.IsVisible(enemy.Pos) {
				g.AddLog(enemy.DisplayName() + " surges into a hotter fury.")
			}
		}

		dist := distance(enemy.Pos, g.Player.Pos)
		playerVisible := dist <= enemy.Template.Sight && hasLineOfSight(g.Floor, enemy.Pos, g.Player.Pos)
		if playerVisible {
			enemy.State = AIStateChase
			enemy.LastKnownPlayer = g.Player.Pos
			enemy.HasLastKnown = true
			enemy.Memory = 5
		}

		if playerVisible && dist > 1 && enemy.Template.BurstRange > 0 && dist <= enemy.Template.BurstRange && enemy.Cooldown == 0 {
			g.enemyBurst(enemy)
			continue
		}
		if dist == 1 {
			enemy.State = AIStateAttack
			g.enemyAttack(enemy)
			continue
		}

		switch enemy.Template.Behavior {
		case BehaviorSkittish:
			if enemy.HP <= enemy.Template.MaxHP/2 && dist <= 3 {
				if g.enemyStepAway(enemy) {
					continue
				}
			}
		case BehaviorSentinel:
			if !playerVisible && distance(enemy.Pos, enemy.Home) > 3 {
				if g.enemyMoveToward(enemy, enemy.Home) {
					continue
				}
			}
		}

		if enemy.HasLastKnown {
			if g.enemyMoveToward(enemy, enemy.LastKnownPlayer) {
				enemy.Memory--
				if enemy.Pos.Equals(enemy.LastKnownPlayer) || enemy.Memory <= 0 {
					enemy.HasLastKnown = false
				}
				continue
			}
			enemy.HasLastKnown = false
		}

		enemy.State = AIStateWander
		g.enemyWander(enemy)
	}
}

func (g *Game) enemyMoveToward(enemy *Enemy, goal Position) bool {
	blocked := g.blockedTiles(enemy.ID)
	step, ok := NextStepToward(g.Floor, enemy.Pos, goal, enemy.Template.CanOpenDoors, blocked)
	if !ok || step.Equals(enemy.Pos) {
		return false
	}
	return g.moveEnemy(enemy, step)
}

func (g *Game) enemyStepAway(enemy *Enemy) bool {
	step, ok := StepAway(g.Floor, enemy.Pos, g.Player.Pos, enemy.Template.CanOpenDoors, g.blockedTiles(enemy.ID))
	if !ok {
		return false
	}
	return g.moveEnemy(enemy, step)
}

func (g *Game) enemyWander(enemy *Enemy) bool {
	directions := g.rng.Perm(len(cardinalDirections))
	blocked := g.blockedTiles(enemy.ID)
	for _, index := range directions {
		next := enemy.Pos.Add(cardinalDirections[index])
		if distance(next, enemy.Home) > 6 && enemy.Template.Behavior != BehaviorHunter && enemy.Template.Behavior != BehaviorProwler && enemy.Template.Behavior != BehaviorBoss {
			continue
		}
		if !g.Floor.InBounds(next) || blocked[next] {
			continue
		}
		if !g.Floor.IsWalkableFor(next, enemy.Template.CanOpenDoors) {
			continue
		}
		return g.moveEnemy(enemy, next)
	}
	return false
}

func (g *Game) moveEnemy(enemy *Enemy, next Position) bool {
	if next.Equals(g.Player.Pos) {
		return false
	}
	if g.Floor.TileAt(next) == TileDoorClosed && enemy.Template.CanOpenDoors {
		g.Floor.OpenDoor(next)
	}
	if !g.Floor.IsWalkable(next) && g.Floor.TileAt(next) != TileDoorOpen {
		return false
	}
	enemy.Pos = next
	return true
}

func (g *Game) enemyAttack(enemy *Enemy) {
	damage := g.damageRoll(enemy.AttackPower(), g.Player.DefensePower())
	if enemy.Template.Behavior == BehaviorHunter && enemy.State == AIStateChase {
		damage++
	}
	applied, blocked := g.applyPlayerDamage(damage)
	if blocked {
		g.AddLog(enemy.DisplayName() + " strikes, but god mode turns the blow aside.")
	} else {
		g.AddLog(enemy.DisplayName() + " hits you for " + itoa(applied) + ".")
	}
	if enemy.Template.GoldStealMax > 0 && g.Player.Gold > 0 {
		stolen := min(g.Player.Gold, 1+g.rng.Intn(enemy.Template.GoldStealMax))
		g.Player.Gold -= stolen
		g.AddLog(enemy.DisplayName() + " palms " + itoa(stolen) + " of your gold.")
	}
	if kind, chance, turns, potency := enemyAttackStatus(enemy); chance > 0 && g.rng.Float64() < chance {
		g.applyPlayerStatus(kind, turns, potency, enemy.DisplayName())
	}
	if g.Player.HP <= 0 {
		g.Player.HP = 0
		g.Mode = ModeLost
		g.AddLog("You collapse beneath " + enemy.DisplayName() + ".")
	}
}

func (g *Game) enemyBurst(enemy *Enemy) {
	damage := g.damageRoll(enemy.Template.BurstDamage, max(0, g.Player.DefensePower()/2))
	applied, blocked := g.applyPlayerDamage(damage)
	enemy.Cooldown = enemy.Template.BurstCooldown
	label := enemy.Template.BurstName
	if label == "" {
		label = "ranged strike"
	}
	if blocked {
		g.AddLog(enemy.DisplayName() + " uses " + label + ", but god mode scatters it harmlessly.")
	} else {
		g.AddLog(enemy.DisplayName() + " uses " + label + " for " + itoa(applied) + ".")
	}
	if enemy.Template.BurstStatusTurns > 0 {
		g.applyPlayerStatus(enemy.Template.BurstStatus, enemy.Template.BurstStatusTurns, max(1, enemy.Template.BurstStatusPotency), enemy.DisplayName()+"'s "+label)
	}
	if g.Player.HP <= 0 {
		g.Player.HP = 0
		g.Mode = ModeLost
		g.AddLog("You are broken by " + enemy.DisplayName() + ".")
	}
}

func (g *Game) tickPlayerStatuses() {
	nextStatuses := make([]StatusEffect, 0, len(g.Player.Statuses))
	for _, status := range g.Player.Statuses {
		if g.GodMode && status.Harmful() {
			g.AddLog("God mode burns away the " + status.Kind.Name() + ".")
			continue
		}
		switch status.Kind {
		case StatusPoison:
			g.Player.HP -= status.Potency
			g.AddLog("Poison bites for " + itoa(status.Potency) + ".")
		case StatusFire:
			g.Player.HP -= max(1, status.Potency)
			g.AddLog("Fire eats you for " + itoa(max(1, status.Potency)) + ".")
		}
		if status.Turns > 0 {
			status.Turns--
		}
		if status.Turns > 0 || status.Floors > 0 {
			nextStatuses = append(nextStatuses, status)
		} else {
			switch status.Kind {
			case StatusFocus:
				g.AddLog("The tonic's edge fades.")
			case StatusPoison:
				g.AddLog("The poison thins out.")
			case StatusFire:
				g.AddLog("The fire gutters out.")
			}
		}
	}
	g.Player.Statuses = nextStatuses
	g.Player.ClampHP()
}

func (g *Game) tickPlayerFloorStatuses() {
	nextStatuses := make([]StatusEffect, 0, len(g.Player.Statuses))
	for _, status := range g.Player.Statuses {
		if status.Floors <= 0 {
			nextStatuses = append(nextStatuses, status)
			continue
		}

		status.Floors--
		if status.Floors > 0 {
			nextStatuses = append(nextStatuses, status)
			continue
		}

		switch status.Kind {
		case StatusFocus:
			g.AddLog("The tonic's edge fades as a new floor begins.")
		default:
			g.AddLog("The " + status.Kind.Name() + " finally falls away.")
		}
	}
	g.Player.Statuses = nextStatuses
}

func (g *Game) applyPlayerDamage(amount int) (int, bool) {
	if amount <= 0 || g.Player == nil {
		return 0, false
	}
	if g.GodMode {
		return 0, true
	}
	g.Player.HP = max(0, g.Player.HP-amount)
	return amount, false
}

func (g *Game) nearestVisibleEnemy() *Enemy {
	var best *Enemy
	bestDistance := 999
	for _, enemy := range g.Floor.Enemies {
		if !g.Floor.IsVisible(enemy.Pos) {
			continue
		}
		currentDistance := distance(g.Player.Pos, enemy.Pos)
		if currentDistance < bestDistance {
			bestDistance = currentDistance
			best = enemy
		}
	}
	return best
}

func (g *Game) quickHealCandidate() (int, Item, int, bool) {
	bestIndex := -1
	bestItem := Item{}
	for index, item := range g.Player.Inventory {
		if item.Kind != ItemKindConsumable || item.Heal <= 0 {
			continue
		}
		utility := restorativeUtilityRank(item, g.Player)
		bestUtility := restorativeUtilityRank(bestItem, g.Player)
		if bestIndex == -1 || utility > bestUtility || (utility == bestUtility && (item.Heal < bestItem.Heal || (item.Heal == bestItem.Heal && item.Price < bestItem.Price))) {
			bestIndex = index
			bestItem = item
		}
	}
	if bestIndex == -1 {
		return -1, Item{}, 0, false
	}
	count := 0
	for _, item := range g.Player.Inventory {
		if item.ID == bestItem.ID {
			count++
		}
	}
	return bestIndex, bestItem, count, true
}

func restorativeUtilityRank(item Item, player *Player) int {
	rank := 0
	if player != nil {
		if player.HasStatus(StatusPoison) && item.PoisonCure {
			rank += 4
		}
		if player.HasStatus(StatusFire) && item.FireCure {
			rank += 4
		}
	}
	if item.PoisonCure {
		rank++
	}
	if item.FireCure {
		rank++
	}
	return rank
}

func (g *Game) blockedTiles(exceptEnemyID int) map[Position]bool {
	blocked := make(map[Position]bool, len(g.Floor.Enemies)+1)
	for _, enemy := range g.Floor.Enemies {
		if enemy.ID == exceptEnemyID {
			continue
		}
		blocked[enemy.Pos] = true
	}
	blocked[g.Player.Pos] = true
	return blocked
}

func (g *Game) damageRoll(attack int, defense int) int {
	rolledAttack := attack + g.rng.Intn(max(3, 3+attack/3))
	rolledDefense := defense
	if defense > 0 {
		rolledDefense += g.rng.Intn(2 + defense)
	}
	return max(1, rolledAttack-rolledDefense/2)
}

func (g *Game) gainItem(item Item, prefix string) {
	switch item.Kind {
	case ItemKindKey:
		g.Player.Keys.Add(item.KeyTier, 1)
		g.AddLog(prefix + strings.ToLower(item.Name) + ".")
	case ItemKindRelic:
		g.Player.HasRelic = true
		g.AddLog("You lift the Cinder Crown. Emberlight floods the hall.")
		g.Mode = ModeWon
	default:
		g.Player.Inventory = append(g.Player.Inventory, item)
		g.AddLog(prefix + item.Name + ".")
	}
}

func (g *Game) applyFloorArrival(modifier FloorModifier) {
	if modifier.GuaranteedKey != nil {
		g.gainItem(KeyItem(*modifier.GuaranteedKey), "Route gift: ")
	}
	if modifier.Rest {
		g.Player.HP = min(g.Player.MaxHP(), g.Player.HP+modifier.HealOnStart)
		cleansed := []StatusKind(nil)
		if modifier.CleanseOnRest {
			cleansed = g.Player.RemoveStatuses(StatusPoison, StatusFire)
		}
		g.AddLog("A brief sanctuary steadies your hands before the next halls.")
		if len(cleansed) > 0 {
			g.AddLog("The rest washes away " + statusList(cleansed) + ".")
		}
	}
}

func enemyAttackStatus(enemy *Enemy) (StatusKind, float64, int, int) {
	if enemy == nil {
		return StatusPoison, 0, 0, 0
	}
	if enemy.Template.AttackStatusChance > 0 && enemy.Template.AttackStatusTurns > 0 {
		return enemy.Template.AttackStatus, enemy.Template.AttackStatusChance, enemy.Template.AttackStatusTurns, max(1, enemy.Template.AttackStatusPotency)
	}
	if enemy.Template.PoisonChance > 0 && enemy.Template.PoisonTurns > 0 {
		return StatusPoison, enemy.Template.PoisonChance, enemy.Template.PoisonTurns, 1
	}
	return StatusPoison, 0, 0, 0
}

func (g *Game) applyPlayerStatus(kind StatusKind, turns int, potency int, source string) bool {
	if turns <= 0 || potency <= 0 {
		return false
	}
	if g.GodMode && kind != StatusFocus {
		g.AddLog(source + " cannot afflict you through god mode.")
		return false
	}

	resistance := g.Player.StatusResistance(kind)
	originalTurns := turns
	for resistance > 0 && turns > 1 {
		turns--
		resistance--
	}
	for resistance > 0 && potency > 0 {
		potency--
		resistance--
	}
	if potency <= 0 {
		g.AddLog("Your wards blunt the " + kind.Name() + " from " + source + ".")
		return false
	}

	before, hadBefore := g.Player.Status(kind)
	g.Player.ApplyStatus(kind, turns, potency)
	after, _ := g.Player.Status(kind)
	if turns < originalTurns {
		g.AddLog("Your wards blunt the " + kind.Name() + ".")
	}
	if hadBefore && after.Turns == before.Turns && after.Potency == before.Potency {
		return true
	}

	switch kind {
	case StatusPoison:
		if hadBefore {
			g.AddLog(source + " deepens the poison on you (" + itoa(after.Turns) + "t).")
		} else {
			g.AddLog(source + " poisons you (" + itoa(after.Turns) + "t).")
		}
	case StatusFire:
		if hadBefore {
			g.AddLog(source + " feeds the fire on you (" + itoa(after.Turns) + "t).")
		} else {
			g.AddLog(source + " sets you ablaze (" + itoa(after.Turns) + "t).")
		}
	}
	return true
}

func (g *Game) consumeCurativeEffects(item Item) []StatusKind {
	toRemove := make([]StatusKind, 0, 2)
	if item.PoisonCure {
		toRemove = append(toRemove, StatusPoison)
	}
	if item.FireCure {
		toRemove = append(toRemove, StatusFire)
	}
	return g.Player.RemoveStatuses(toRemove...)
}

func statusList(statuses []StatusKind) string {
	if len(statuses) == 0 {
		return ""
	}
	parts := make([]string, 0, len(statuses))
	for _, status := range statuses {
		parts = append(parts, status.Name())
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1]
}

func (g *Game) findChestIndex(id int) int {
	for index, chest := range g.Floor.Chests {
		if chest.ID == id {
			return index
		}
	}
	return -1
}

func (g *Game) DebugState() string {
	return "floor=" + itoa(g.FloorIndex) + " turn=" + itoa(g.Turn) + " enemies=" + itoa(len(g.Floor.Enemies)) + " items=" + itoa(len(g.Floor.Items))
}
