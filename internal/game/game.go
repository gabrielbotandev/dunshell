package game

import (
	"fmt"
	"math/rand"
	"strings"
)

const GameTitle = "Dunshell"

type Summary struct {
	Seed           int64
	Floor          int
	Level          int
	Gold           int
	Kills          int
	Turn           int
	RecoveredRelic bool
}

type Game struct {
	Title       string
	Seed        int64
	Mode        GameMode
	FloorIndex  int
	MaxFloors   int
	Turn        int
	Log         []string
	Player      *Player
	Floor       *Floor
	rng         *rand.Rand
	nextEnemyID int
}

func New(seed int64) *Game {
	rng := rand.New(rand.NewSource(seed))
	weapon, armor, charm := StarterEquipment()
	player := &Player{
		Pos:         Position{},
		BaseMaxHP:   24,
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
		Title:      GameTitle,
		Seed:       seed,
		Mode:       ModePlaying,
		FloorIndex: 1,
		MaxFloors:  5,
		Player:     player,
		rng:        rng,
		Log:        make([]string, 0, 128),
	}

	game.Floor = GenerateFloor(game.rng, game.FloorIndex, game.MaxFloors, &game.nextEnemyID)
	game.Player.Pos = game.Floor.Entrance
	ComputeFOV(game.Floor, game.Player.Pos, game.Player.VisionRadius())
	game.AddLog("The abbey doors seal behind you. Find the Cinder Crown in the Ember Sanctum.")
	game.AddLog(game.floorIntro())
	return game
}

func (g *Game) Summary() Summary {
	return Summary{
		Seed:           g.Seed,
		Floor:          g.FloorIndex,
		Level:          g.Player.Level,
		Gold:           g.Player.Gold,
		Kills:          g.Player.Kills,
		Turn:           g.Turn,
		RecoveredRelic: g.Player.HasRelic,
	}
}

func (g *Game) Objective() string {
	if g.FloorIndex == g.MaxFloors {
		if g.Player.HasRelic {
			return "Escape the sanctum in memory alone."
		}
		return "Claim the Cinder Crown."
	}
	return "Descend to floor " + itoa(g.MaxFloors) + " and claim the Cinder Crown."
}

func (g *Game) FloorLabel() string {
	return "Floor " + itoa(g.FloorIndex) + "  " + g.Floor.Theme
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
		if item.Kind == ItemKindRelic {
			g.Player.HasRelic = true
			g.AddLog("You lift the Cinder Crown. Emberlight floods the hall.")
			g.Mode = ModeWon
			continue
		}
		g.Player.Inventory = append(g.Player.Inventory, item)
		g.AddLog("You gather " + item.Name + ".")
	}

	if g.Mode == ModePlaying {
		g.advanceTurn()
	}
	return true
}

func (g *Game) Descend() bool {
	if g.Mode != ModePlaying {
		return false
	}
	if !g.Player.Pos.Equals(g.Floor.Stairs) {
		g.AddLog("No stair waits beneath your boots.")
		return false
	}
	if g.FloorIndex >= g.MaxFloors {
		g.AddLog("This is the deepest hall.")
		return false
	}

	g.FloorIndex++
	g.Floor = GenerateFloor(g.rng, g.FloorIndex, g.MaxFloors, &g.nextEnemyID)
	g.Player.Pos = g.Floor.Entrance
	g.Player.HP = min(g.Player.MaxHP(), g.Player.HP+5)
	ComputeFOV(g.Floor, g.Player.Pos, g.Player.VisionRadius())
	g.AddLog("You descend into " + strings.ToLower(g.Floor.Theme) + ".")
	g.AddLog(g.floorIntro())
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
		if item.PoisonCure {
			g.Player.RemoveStatus(StatusPoison)
		}
		healed := g.Player.HP - before
		if healed == 0 && !item.PoisonCure {
			g.AddLog("The salve would be wasted right now.")
			return false
		}
		if item.PoisonCure {
			g.AddLog("You steady your breath with " + item.Name + ".")
		} else {
			g.AddLog("You patch yourself up for " + itoa(healed) + " HP.")
		}
		consumed = true
	case item.FocusTurns > 0:
		g.Player.ApplyStatus(StatusFocus, item.FocusTurns, item.FocusBonus)
		g.AddLog("Sunbrew heats your blood. Your strikes sharpen.")
		consumed = true
	case item.EmberDamage > 0:
		target := g.nearestVisibleEnemy()
		if target == nil {
			g.AddLog("The flask has no mark to chase.")
			return false
		}
		damage := item.EmberDamage + g.rng.Intn(3)
		target.HP -= damage
		g.AddLog("Ember arcs into " + target.Template.Name + " for " + itoa(damage) + ".")
		if target.HP <= 0 {
			g.killEnemy(target, "The ember blast shatters "+target.Template.Name+".")
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

	oldMax := g.Player.MaxHP()
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
	if g.Player.MaxHP() < oldMax {
		g.Player.ClampHP()
	}
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
	tile := g.Floor.TileAt(g.Player.Pos).Name()
	if items := g.Floor.ItemIndicesAt(g.Player.Pos); len(items) > 0 {
		item := g.Floor.Items[items[0]].Item
		return tile + " with " + item.Name
	}
	return tile
}

func (g *Game) AddLog(message string) {
	if message == "" {
		return
	}
	g.Log = append(g.Log, message)
	if len(g.Log) > 180 {
		g.Log = g.Log[len(g.Log)-180:]
	}
}

func (g *Game) floorIntro() string {
	switch g.FloorIndex {
	case 1:
		return "Moss climbs the crypt walls. Old candles still smell faintly sweet."
	case 2:
		return "The warrens breathe warm drafts through narrow stone throats."
	case 3:
		return "Salt and bone crackle underfoot in the ossuary."
	case 4:
		return "Wisps drift like choir lights through the galleries."
	case 5:
		return "Heat gathers below. The sanctum keeps its relic behind prayer and ash."
	default:
		return "The dark waits."
	}
}

func (g *Game) playerAttack(enemy *Enemy) {
	damage := g.damageRoll(g.Player.AttackPower(), enemy.DefensePower())
	enemy.HP -= damage
	g.AddLog("You strike " + enemy.Template.Name + " for " + itoa(damage) + ".")
	if enemy.HP <= 0 {
		g.killEnemy(enemy, "You finish "+enemy.Template.Name+".")
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
	g.Player.Gold += gold
	g.AddLog(enemy.Template.Name + " drops " + itoa(gold) + " gold.")

	if item, ok := RandomDropItem(g.rng, g.FloorIndex); ok {
		g.Floor.Items = append(g.Floor.Items, GroundItem{Pos: enemy.Pos, Item: item})
		g.AddLog(enemy.Template.Name + " leaves behind " + item.Name + ".")
	}

	if g.Player.GainXP(enemy.Template.XPReward) {
		g.AddLog("You rise to level " + itoa(g.Player.Level) + ".")
	}
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

		if distance(enemy.Pos, g.Player.Pos) == 1 {
			enemy.State = AIStateAttack
			g.enemyAttack(enemy)
			continue
		}

		playerVisible := distance(enemy.Pos, g.Player.Pos) <= enemy.Template.Sight && hasLineOfSight(g.Floor, enemy.Pos, g.Player.Pos)
		if playerVisible {
			enemy.State = AIStateChase
			enemy.LastKnownPlayer = g.Player.Pos
			enemy.HasLastKnown = true
			enemy.Memory = 4
		}

		switch enemy.Template.Behavior {
		case BehaviorSkittish:
			if enemy.HP <= enemy.Template.MaxHP/2 && distance(enemy.Pos, g.Player.Pos) <= 3 {
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
	for _, index := range directions {
		next := enemy.Pos.Add(cardinalDirections[index])
		if distance(next, enemy.Home) > 6 && enemy.Template.Behavior != BehaviorHunter && enemy.Template.Behavior != BehaviorProwler {
			continue
		}
		if !g.Floor.InBounds(next) || g.blockedTiles(enemy.ID)[next] {
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
	g.Player.HP -= damage
	g.AddLog(enemy.Template.Name + " hits you for " + itoa(damage) + ".")

	if enemy.Template.GoldStealMax > 0 && g.Player.Gold > 0 {
		stolen := min(g.Player.Gold, 1+g.rng.Intn(enemy.Template.GoldStealMax))
		g.Player.Gold -= stolen
		g.AddLog(enemy.Template.Name + " palms " + itoa(stolen) + " of your gold.")
	}

	if enemy.Template.PoisonChance > 0 && g.rng.Float64() < enemy.Template.PoisonChance {
		g.Player.ApplyStatus(StatusPoison, enemy.Template.PoisonTurns, 1)
		g.AddLog(enemy.Template.Name + " leaves a burning poison in the wound.")
	}

	if g.Player.HP <= 0 {
		g.Player.HP = 0
		g.Mode = ModeLost
		g.AddLog("You collapse beneath " + enemy.Template.Name + ".")
	}
}

func (g *Game) tickPlayerStatuses() {
	nextStatuses := make([]StatusEffect, 0, len(g.Player.Statuses))
	for _, status := range g.Player.Statuses {
		switch status.Kind {
		case StatusPoison:
			g.Player.HP -= status.Potency
			g.AddLog("Poison burns for " + itoa(status.Potency) + ".")
		}
		status.Turns--
		if status.Turns > 0 {
			nextStatuses = append(nextStatuses, status)
		} else if status.Kind == StatusFocus {
			g.AddLog("The sunbrew edge fades.")
		}
	}
	g.Player.Statuses = nextStatuses
	g.Player.ClampHP()
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
		if bestIndex == -1 ||
			item.Heal < bestItem.Heal ||
			(item.Heal == bestItem.Heal && poisonUtilityRank(item) < poisonUtilityRank(bestItem)) {
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

func poisonUtilityRank(item Item) int {
	if item.PoisonCure {
		return 1
	}
	return 0
}

func (g *Game) blockedTiles(exceptEnemyID int) map[Position]bool {
	blocked := make(map[Position]bool, len(g.Floor.Enemies))
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
	rolledAttack := attack + g.rng.Intn(4)
	rolledDefense := defense
	if defense > 0 {
		rolledDefense += g.rng.Intn(defense + 1)
	}
	return max(1, rolledAttack-rolledDefense/2)
}

func (g *Game) DebugState() string {
	return fmt.Sprintf("floor=%d turn=%d enemies=%d items=%d", g.FloorIndex, g.Turn, len(g.Floor.Enemies), len(g.Floor.Items))
}
