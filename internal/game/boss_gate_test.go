package game

import "testing"

func TestBossGateOpensAfterBossDefeat(t *testing.T) {
	var nextEnemyID int
	var nextChestID int
	var nextMerchantID int

	floor := GenerateFloor(NewRNG(5), 5, 20, 0, FloorModifier{}, false, &nextEnemyID, &nextChestID, &nextMerchantID)
	if floor.Boss == nil {
		t.Fatal("expected boss floor to include a boss encounter")
	}

	g := &Game{
		Mode:  ModePlaying,
		Floor: floor,
		Player: &Player{
			HP:         30,
			BaseMaxHP:  30,
			BaseAttack: 4,
			Level:      1,
		},
		rng: NewRNG(9),
		Log: make([]string, 0, 16),
	}

	boss := g.Floor.EnemyByID(g.Floor.Boss.BossID)
	if boss == nil {
		t.Fatal("expected boss enemy to be present")
	}

	g.Floor.Boss.Active = true
	for _, door := range g.bossRoomDoors() {
		g.Floor.SetTile(door, TileBossSeal)
	}
	if got := g.Floor.TileAt(g.Floor.Boss.Gate); got != TileBossSeal {
		t.Fatalf("expected boss gate to seal during encounter, got %v", got)
	}

	g.killEnemy(boss, "The keeper breaks.")

	if got := g.Floor.TileAt(g.Floor.Boss.Gate); got != TileDoorOpen {
		t.Fatalf("expected boss gate to open after defeat, got %v", got)
	}
	if !g.Floor.Boss.Cleared {
		t.Fatal("expected boss encounter to be marked cleared")
	}
	if g.Floor.Boss.Active {
		t.Fatal("expected boss encounter to be inactive after defeat")
	}
}
