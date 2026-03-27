package game

import "math/rand"

type EnemyTemplate struct {
	ID           string
	Name         string
	Glyph        rune
	Tint         string
	Description  string
	MaxHP        int
	Attack       int
	Defense      int
	Sight        int
	XPReward     int
	GoldMin      int
	GoldMax      int
	MinFloor     int
	MaxFloor     int
	Weight       int
	Behavior     Behavior
	PoisonChance float64
	PoisonTurns  int
	GoldStealMax int
	CanOpenDoors bool
}

var floorThemes = []string{
	"",
	"Moss Crypts",
	"Candle Warrens",
	"Salt Ossuary",
	"Wisp Galleries",
	"Ember Sanctum",
}

var enemyCatalog = []EnemyTemplate{
	{
		ID: "gutter_rat", Name: "Gutter Rat", Glyph: 'r', Tint: "#b9a38c",
		Description: "Fast vermin that lose nerve when cornered.",
		MaxHP:       7, Attack: 3, Defense: 0, Sight: 5, XPReward: 6, GoldMin: 1, GoldMax: 4,
		MinFloor: 1, MaxFloor: 2, Weight: 26, Behavior: BehaviorSkittish,
	},
	{
		ID: "bone_beetle", Name: "Bone Beetle", Glyph: 'b', Tint: "#8a9965",
		Description: "A shrine beetle wearing fragments of old bone as shell.",
		MaxHP:       12, Attack: 4, Defense: 2, Sight: 5, XPReward: 9, GoldMin: 2, GoldMax: 5,
		MinFloor: 1, MaxFloor: 3, Weight: 22, Behavior: BehaviorBrute,
	},
	{
		ID: "lantern_wisp", Name: "Lantern Wisp", Glyph: 'w', Tint: "#f2d16b",
		Description: "Cold fire that brands the lungs like poison.",
		MaxHP:       10, Attack: 5, Defense: 1, Sight: 8, XPReward: 12, GoldMin: 3, GoldMax: 7,
		MinFloor: 2, MaxFloor: 4, Weight: 20, Behavior: BehaviorProwler, PoisonChance: 0.28, PoisonTurns: 3,
	},
	{
		ID: "mire_hound", Name: "Mire Hound", Glyph: 'h', Tint: "#6fab78",
		Description: "It scents blood through doors and damp halls alike.",
		MaxHP:       15, Attack: 6, Defense: 1, Sight: 9, XPReward: 14, GoldMin: 4, GoldMax: 9,
		MinFloor: 2, MaxFloor: 5, Weight: 18, Behavior: BehaviorHunter, CanOpenDoors: true,
	},
	{
		ID: "tomb_brigand", Name: "Tomb Brigand", Glyph: 't', Tint: "#d4a66d",
		Description: "A grave robber who still remembers the shine of coin.",
		MaxHP:       16, Attack: 5, Defense: 1, Sight: 7, XPReward: 15, GoldMin: 6, GoldMax: 11,
		MinFloor: 3, MaxFloor: 5, Weight: 16, Behavior: BehaviorCutpurse, GoldStealMax: 6, CanOpenDoors: true,
	},
	{
		ID: "cathedral_knight", Name: "Cathedral Knight", Glyph: 'K', Tint: "#bac7d0",
		Description: "Slow, stubborn, and impossible to dissuade once roused.",
		MaxHP:       24, Attack: 8, Defense: 3, Sight: 6, XPReward: 22, GoldMin: 8, GoldMax: 14,
		MinFloor: 4, MaxFloor: 5, Weight: 10, Behavior: BehaviorSentinel, CanOpenDoors: true,
	},
	{
		ID: "ashen_prior", Name: "Ashen Prior", Glyph: 'A', Tint: "#ef7f45",
		Description: "Last abbot of the ember rite, still kneeling over the crown.",
		MaxHP:       38, Attack: 10, Defense: 4, Sight: 10, XPReward: 40, GoldMin: 20, GoldMax: 30,
		MinFloor: 5, MaxFloor: 5, Weight: 1, Behavior: BehaviorBoss, PoisonChance: 0.35, PoisonTurns: 4, CanOpenDoors: true,
	},
}

func FloorTheme(level int) string {
	if level < 1 || level >= len(floorThemes) {
		return "Forgotten Vaults"
	}
	return floorThemes[level]
}

func RandomEnemyTemplate(rng *rand.Rand, floor int) EnemyTemplate {
	candidates := make([]EnemyTemplate, 0, len(enemyCatalog))
	total := 0
	for _, template := range enemyCatalog {
		if floor >= template.MinFloor && floor <= template.MaxFloor && template.ID != "ashen_prior" {
			candidates = append(candidates, template)
			total += template.Weight
		}
	}

	roll := rng.Intn(total)
	for _, template := range candidates {
		roll -= template.Weight
		if roll < 0 {
			return template
		}
	}

	return candidates[len(candidates)-1]
}

func BossTemplate() EnemyTemplate {
	for _, template := range enemyCatalog {
		if template.ID == "ashen_prior" {
			return template
		}
	}
	return enemyCatalog[len(enemyCatalog)-1]
}
