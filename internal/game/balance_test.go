package game

import "testing"

func TestScaleEnemyTemplateKeepsFloorOneAndTwoBaseAttackForNormalEnemies(t *testing.T) {
	rat := enemyTemplateByID(t, "gutter_rat")
	if got := ScaleEnemyTemplate(rat, 1, 0, false, false).Attack; got != rat.Attack {
		t.Fatalf("floor 1 attack = %d, want %d", got, rat.Attack)
	}
	if got := ScaleEnemyTemplate(rat, 2, 0, false, false).Attack; got != rat.Attack {
		t.Fatalf("floor 2 attack = %d, want %d", got, rat.Attack)
	}

	wisp := enemyTemplateByID(t, "lantern_wisp")
	if got := ScaleEnemyTemplate(wisp, 2, 0, false, false).Attack; got != wisp.Attack {
		t.Fatalf("floor 2 wisp attack = %d, want %d", got, wisp.Attack)
	}
}

func TestScaleEnemyTemplatePreservesFloorThreeAttackRamp(t *testing.T) {
	rat := enemyTemplateByID(t, "gutter_rat")
	got := ScaleEnemyTemplate(rat, 3, 0, false, false).Attack
	want := rat.Attack + 2
	if got != want {
		t.Fatalf("floor 3 attack = %d, want %d", got, want)
	}
}

func TestScaleEnemyTemplateAddsGlobalEnemyHealth(t *testing.T) {
	rat := enemyTemplateByID(t, "gutter_rat")
	got := ScaleEnemyTemplate(rat, 1, 0, false, false).MaxHP
	want := rat.MaxHP + 2
	if got != want {
		t.Fatalf("floor 1 hp = %d, want %d", got, want)
	}
}

func TestEarlyFloorEnemyCountRelief(t *testing.T) {
	tests := []struct {
		level int
		want  int
	}{
		{level: 1, want: 2},
		{level: 2, want: 1},
		{level: 3, want: 0},
		{level: 8, want: 0},
	}

	for _, test := range tests {
		if got := earlyFloorEnemyCountRelief(test.level); got != test.want {
			t.Fatalf("relief for floor %d = %d, want %d", test.level, got, test.want)
		}
	}
}

func TestHealingSalveEarlyGameRecovery(t *testing.T) {
	if got := ItemByID("healing_salve").Heal; got != 12 {
		t.Fatalf("healing salve heal = %d, want 12", got)
	}
}

func TestStarterInventoryCarriesFiveHealingSalves(t *testing.T) {
	inventory := StarterInventory()
	salves := 0
	tonics := 0
	for _, item := range inventory {
		switch item.ID {
		case "healing_salve":
			salves++
		case "sunbrew_tonic":
			tonics++
		}
	}

	if salves != 5 {
		t.Fatalf("starter salves = %d, want 5", salves)
	}
	if tonics != 0 {
		t.Fatalf("starter tonics = %d, want 0", tonics)
	}
	if len(inventory) != 5 {
		t.Fatalf("starter inventory size = %d, want 5", len(inventory))
	}
}

func TestSunbrewTonicUsesFloorDuration(t *testing.T) {
	item := ItemByID("sunbrew_tonic")
	if item.FocusFloors != 1 {
		t.Fatalf("sunbrew focus floors = %d, want 1", item.FocusFloors)
	}
	if item.FocusTurns != 0 {
		t.Fatalf("sunbrew focus turns = %d, want 0", item.FocusTurns)
	}
}

func TestFocusFloorStatusesPersistAcrossTurnsAndExpireOnFloorChange(t *testing.T) {
	game := &Game{Player: &Player{}}
	game.Player.ApplyStatusByFloor(StatusFocus, 1, 2)

	game.tickPlayerStatuses()
	status, ok := game.Player.Status(StatusFocus)
	if !ok {
		t.Fatal("expected focus to persist across turn ticks")
	}
	if status.Floors != 1 {
		t.Fatalf("focus floors after turn tick = %d, want 1", status.Floors)
	}

	game.tickPlayerFloorStatuses()
	if _, ok := game.Player.Status(StatusFocus); ok {
		t.Fatal("expected focus to expire on floor change")
	}
}

func TestRouteChoiceModifiersGateMerchantRoute(t *testing.T) {
	withoutMerchant := routeChoiceModifiers(NewRNG(1), false)
	for _, modifier := range withoutMerchant {
		if modifier.Merchant {
			t.Fatal("did not expect merchant route when disabled")
		}
	}

	withMerchant := routeChoiceModifiers(NewRNG(1), true)
	merchantRoutes := 0
	for _, modifier := range withMerchant {
		if modifier.Merchant {
			merchantRoutes++
		}
	}
	if merchantRoutes != 1 {
		t.Fatalf("merchant routes = %d, want 1", merchantRoutes)
	}
	if len(withMerchant) != 3 {
		t.Fatalf("route choices with merchant = %d, want 3", len(withMerchant))
	}
}

func TestEnemyXPAwardReducesKillXP(t *testing.T) {
	ratTemplate := ScaleEnemyTemplate(enemyTemplateByID(t, "gutter_rat"), 1, 0, false, false)
	beetleTemplate := ScaleEnemyTemplate(enemyTemplateByID(t, "bone_beetle"), 1, 0, false, false)
	player := testPlayer(1)

	rat := &Enemy{Template: ratTemplate, Level: enemyLevelForEncounter(1, ratTemplate, false), HP: ratTemplate.MaxHP}
	beetle := &Enemy{Template: beetleTemplate, Level: enemyLevelForEncounter(1, beetleTemplate, false), HP: beetleTemplate.MaxHP}

	ratXP := enemyXPProgress(rat, player, 1)
	beetleXP := enemyXPProgress(beetle, player, 1)
	if ratXP >= beetleXP {
		t.Fatalf("expected beetle xp (%d) to exceed rat xp (%d)", beetleXP, ratXP)
	}
}

func TestEquipmentPoolsOfferMoreChoiceByFloorFour(t *testing.T) {
	if got := len(eligibleItemPool(weaponIDs, 4, RarityCommon, RarityLegendary)); got < 5 {
		t.Fatalf("weapon choices on floor 4 = %d, want at least 5", got)
	}
	if got := len(eligibleItemPool(armorIDs, 4, RarityCommon, RarityLegendary)); got < 4 {
		t.Fatalf("armor choices on floor 4 = %d, want at least 4", got)
	}
	if got := len(eligibleItemPool(charmIDs, 4, RarityCommon, RarityLegendary)); got < 5 {
		t.Fatalf("charm choices on floor 4 = %d, want at least 5", got)
	}
}

func TestEnemyLevelTracksFloor(t *testing.T) {
	rat := enemyTemplateByID(t, "gutter_rat")
	if got := enemyLevelForEncounter(1, rat, false); got != 1 {
		t.Fatalf("floor 1 rat level = %d, want 1", got)
	}
	if got := enemyLevelForEncounter(2, rat, false); got != 2 {
		t.Fatalf("floor 2 rat level = %d, want 2", got)
	}
}

func TestEnemyXPProgressFallsAsPlayerLevelClimbs(t *testing.T) {
	template := ScaleEnemyTemplate(enemyTemplateByID(t, "bone_beetle"), 2, 0, false, false)
	enemy := &Enemy{Template: template, Level: enemyLevelForEncounter(2, template, false), HP: template.MaxHP}

	lowLevel := enemyXPProgress(enemy, testPlayer(1), 2)
	highLevel := enemyXPProgress(enemy, testPlayer(5), 2)
	if lowLevel <= highLevel {
		t.Fatalf("expected level 1 xp (%d) to exceed level 5 xp (%d)", lowLevel, highLevel)
	}
}

func TestGainXPUsesPerLevelProgressionRollover(t *testing.T) {
	player := testPlayer(1)
	player.XP = 9990
	leveled := player.GainXP(25)
	if !leveled {
		t.Fatal("expected player to level")
	}
	if player.Level != 2 {
		t.Fatalf("player level = %d, want 2", player.Level)
	}
	if player.XP != 15 {
		t.Fatalf("player xp remainder = %d, want 15", player.XP)
	}
}

func TestHydrateProgressionStateMigratesLegacyXPAndEnemyLevels(t *testing.T) {
	template := enemyTemplateByID(t, "gutter_rat")
	game := &Game{
		Player: &Player{Level: 2, XP: 40},
		Floor:  &Floor{Level: 2, Enemies: []*Enemy{{Template: template}}},
	}

	hydrateProgressionState(game)

	if got := game.Floor.Enemies[0].Level; got != 2 {
		t.Fatalf("enemy level after hydrate = %d, want 2", got)
	}
	if game.Player.XP <= 0 || game.Player.XP >= experiencePerLevel {
		t.Fatalf("migrated xp = %d, want value within current level progress", game.Player.XP)
	}
}

func enemyTemplateByID(t *testing.T, id string) EnemyTemplate {
	t.Helper()
	for _, template := range enemyCatalog {
		if template.ID == id {
			return template
		}
	}
	t.Fatalf("enemy template %q not found", id)
	return EnemyTemplate{}
}

func testPlayer(level int) *Player {
	weapon, armor, charm := StarterEquipment()
	return &Player{
		BaseMaxHP:   30 + max(0, level-1)*5,
		BaseAttack:  4 + max(0, level-1)*2,
		BaseDefense: 1 + max(0, level/2),
		Level:       level,
		Equipment: Equipment{
			Weapon: &weapon,
			Armor:  &armor,
			Charm:  &charm,
		},
	}
}
