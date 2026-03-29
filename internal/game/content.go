package game

type EnemyTemplate struct {
	ID                 string
	Name               string
	Glyph              rune
	ASCII              rune
	Tint               string
	Description        string
	MaxHP              int
	Attack             int
	Defense            int
	Sight              int
	XPReward           int
	GoldMin            int
	GoldMax            int
	MinFloor           int
	MaxFloor           int
	Weight             int
	Behavior           Behavior
	PoisonChance       float64
	PoisonTurns        int
	GoldStealMax       int
	CanOpenDoors       bool
	BossTier           int
	BurstName          string
	BurstRange         int
	BurstDamage        int
	BurstCooldown      int
	BurstStatus        StatusKind
	BurstStatusTurns   int
	BurstStatusPotency int
	EnrageThreshold    int
	EnrageAttackBonus  int
}

type FloorTheme struct {
	Name  string
	Intro string
}

var floorThemes = []FloorTheme{
	{},
	{Name: "Tidal Crypt", Intro: "Salt damp clings to the first crypts. The abbey still smells faintly of floodwater and candles."},
	{Name: "Pilgrim Warrens", Intro: "Narrow stone throats force every footstep into a whisper."},
	{Name: "Salt Ossuary", Intro: "Bone dust and chapel salt crunch together under your boots."},
	{Name: "Bell Rookery", Intro: "Broken clappers sway in the dark, heavy as hanging ribs."},
	{Name: "Kennel Reliquary", Intro: "The hounds are gone. Their chains are not."},
	{Name: "Flooded Cloister", Intro: "Cold water sheets between flagstones and drinks the torchlight."},
	{Name: "Choir Trenches", Intro: "Collapsed pews and cut masonry turn the nave into a battlefield."},
	{Name: "Mire Archive", Intro: "Swollen vellum rots in shelves that lean like drowned men."},
	{Name: "Sepulcher Galleries", Intro: "Portrait niches stare from the walls with blank mineral eyes."},
	{Name: "Resonant Nave", Intro: "The whole floor hums with a note too low to be song."},
	{Name: "Black Refectory", Intro: "Long tables split the hall into lanes of hunger and steel."},
	{Name: "Ash Baths", Intro: "Warm soot drifts through the rooms like funeral steam."},
	{Name: "Censer Loom", Intro: "Braided chains and brass burners fill the dark with perfumed smoke."},
	{Name: "Thorn Processional", Intro: "Iron vines thread the walls where the abbey once staged its rites."},
	{Name: "Saint's Wake", Intro: "Candles keep burning here with no wax left to feed them."},
	{Name: "Char Vault", Intro: "The stones themselves have learned how to hold heat."},
	{Name: "Ember Catacombs", Intro: "Ash glows between bricks like the abbey is remembering fire."},
	{Name: "Sable Reliquary", Intro: "Sealed caskets line the floor in ranks older than the flood."},
	{Name: "Crown Approach", Intro: "Every hall leans inward now, as though listening for your breath."},
	{Name: "Throne of Ash", Intro: "The deepest sanctum waits behind prayer, ash, and a final kneeling shape."},
}

var enemyCatalog = []EnemyTemplate{
	{
		ID: "gutter_rat", Name: "Gutter Rat", Glyph: 'r', ASCII: 'r', Tint: "#b9a38c",
		Description: "Fast vermin that lose nerve when cornered.",
		MaxHP:       8, Attack: 4, Defense: 0, Sight: 6, XPReward: 9, GoldMin: 1, GoldMax: 4,
		MinFloor: 1, MaxFloor: 3, Weight: 18, Behavior: BehaviorSkittish,
	},
	{
		ID: "bone_beetle", Name: "Bone Beetle", Glyph: 'b', ASCII: 'b', Tint: "#8a9965",
		Description: "A shrine beetle wearing fragments of old bone as shell.",
		MaxHP:       14, Attack: 5, Defense: 2, Sight: 5, XPReward: 12, GoldMin: 2, GoldMax: 5,
		MinFloor: 1, MaxFloor: 5, Weight: 16, Behavior: BehaviorBrute,
	},
	{
		ID: "lantern_wisp", Name: "Lantern Wisp", Glyph: 'w', ASCII: 'w', Tint: "#f2d16b",
		Description: "Cold fire that settles in the lungs like poison.",
		MaxHP:       12, Attack: 6, Defense: 1, Sight: 8, XPReward: 14, GoldMin: 3, GoldMax: 6,
		MinFloor: 2, MaxFloor: 6, Weight: 14, Behavior: BehaviorProwler, PoisonChance: 0.24, PoisonTurns: 3,
	},
	{
		ID: "mire_hound", Name: "Mire Hound", Glyph: 'h', ASCII: 'h', Tint: "#6fab78",
		Description: "It scents blood through doors and damp halls alike.",
		MaxHP:       18, Attack: 7, Defense: 1, Sight: 9, XPReward: 18, GoldMin: 4, GoldMax: 8,
		MinFloor: 3, MaxFloor: 8, Weight: 13, Behavior: BehaviorHunter, CanOpenDoors: true,
	},
	{
		ID: "tomb_brigand", Name: "Tomb Brigand", Glyph: 't', ASCII: 't', Tint: "#d4a66d",
		Description: "A grave robber who still remembers the shine of coin.",
		MaxHP:       19, Attack: 6, Defense: 2, Sight: 8, XPReward: 20, GoldMin: 6, GoldMax: 11,
		MinFloor: 4, MaxFloor: 9, Weight: 12, Behavior: BehaviorCutpurse, GoldStealMax: 8, CanOpenDoors: true,
	},
	{
		ID: "cathedral_knight", Name: "Cathedral Knight", Glyph: 'K', ASCII: 'K', Tint: "#bac7d0",
		Description: "Slow, stubborn, and impossible to dissuade once roused.",
		MaxHP:       28, Attack: 9, Defense: 4, Sight: 6, XPReward: 28, GoldMin: 8, GoldMax: 14,
		MinFloor: 6, MaxFloor: 12, Weight: 11, Behavior: BehaviorSentinel, CanOpenDoors: true,
	},
	{
		ID: "censer_acolyte", Name: "Censer Acolyte", Glyph: 'c', ASCII: 'c', Tint: "#c69d73",
		Description: "A smoke-veiled celebrant that hurls burning incense ahead of itself.",
		MaxHP:       22, Attack: 7, Defense: 2, Sight: 8, XPReward: 26, GoldMin: 9, GoldMax: 15,
		MinFloor: 7, MaxFloor: 15, Weight: 10, Behavior: BehaviorCaster, BurstName: "censer flare", BurstRange: 5,
		BurstDamage: 8, BurstCooldown: 3, BurstStatus: StatusScorched, BurstStatusTurns: 2, BurstStatusPotency: 1,
	},
	{
		ID: "grave_sycophant", Name: "Grave Sycophant", Glyph: 's', ASCII: 's', Tint: "#9b8476",
		Description: "A broad-shouldered penitent that has forgotten how to stop charging.",
		MaxHP:       32, Attack: 10, Defense: 3, Sight: 7, XPReward: 30, GoldMin: 10, GoldMax: 16,
		MinFloor: 9, MaxFloor: 16, Weight: 10, Behavior: BehaviorBrute,
	},
	{
		ID: "ash_archer", Name: "Ash Archer", Glyph: 'a', ASCII: 'a', Tint: "#e39b66",
		Description: "Its bowstring is a strip of embered gut that never snaps.",
		MaxHP:       24, Attack: 9, Defense: 2, Sight: 9, XPReward: 32, GoldMin: 12, GoldMax: 18,
		MinFloor: 10, MaxFloor: 18, Weight: 10, Behavior: BehaviorCaster, BurstName: "ash bolt", BurstRange: 7,
		BurstDamage: 10, BurstCooldown: 2,
	},
	{
		ID: "reliquary_ogre", Name: "Reliquary Ogre", Glyph: 'O', ASCII: 'O', Tint: "#9f8b7a",
		Description: "A chained giant swollen on relic dust and old vows.",
		MaxHP:       40, Attack: 13, Defense: 4, Sight: 6, XPReward: 38, GoldMin: 14, GoldMax: 22,
		MinFloor: 12, MaxFloor: 20, Weight: 8, Behavior: BehaviorBrute, CanOpenDoors: true,
	},
	{
		ID: "drowned_abbot", Name: "Drowned Abbot", Glyph: 'A', ASCII: 'A', Tint: "#86a6c1",
		Description: "A drowned priest whose blessing still arrives as punishment.",
		MaxHP:       34, Attack: 11, Defense: 5, Sight: 8, XPReward: 40, GoldMin: 16, GoldMax: 24,
		MinFloor: 14, MaxFloor: 20, Weight: 7, Behavior: BehaviorCaster, CanOpenDoors: true,
		BurstName: "drowned litany", BurstRange: 6, BurstDamage: 12, BurstCooldown: 3,
	},
	{
		ID: "ember_seraph", Name: "Ember Seraph", Glyph: 'S', ASCII: 'S', Tint: "#ef7f45",
		Description: "A flayed choir thing made of feather-shadows and furnace light.",
		MaxHP:       30, Attack: 12, Defense: 3, Sight: 10, XPReward: 42, GoldMin: 16, GoldMax: 26,
		MinFloor: 16, MaxFloor: 20, Weight: 6, Behavior: BehaviorHunter,
		BurstName: "ember lance", BurstRange: 6, BurstDamage: 11, BurstCooldown: 3, BurstStatus: StatusScorched,
		BurstStatusTurns: 2, BurstStatusPotency: 2,
	},
	{
		ID: "houndmaster_vey", Name: "Houndmaster Vey", Glyph: '◆', ASCII: 'V', Tint: "#d26767",
		Description: "The reliquary's kennel saint, all chain, scar, and command.",
		MaxHP:       58, Attack: 11, Defense: 3, Sight: 10, XPReward: 70, GoldMin: 30, GoldMax: 40,
		MinFloor: 5, MaxFloor: 5, Weight: 1, Behavior: BehaviorBoss, CanOpenDoors: true, BossTier: 1,
		BurstName: "ruinous pounce", BurstRange: 4, BurstDamage: 10, BurstCooldown: 2,
	},
	{
		ID: "bell_archivist_oria", Name: "Bell Archivist Oria", Glyph: '◆', ASCII: 'O', Tint: "#cf6767",
		Description: "A keeper of drowned bells whose toll still strips the air raw.",
		MaxHP:       82, Attack: 14, Defense: 5, Sight: 11, XPReward: 100, GoldMin: 48, GoldMax: 60,
		MinFloor: 10, MaxFloor: 10, Weight: 1, Behavior: BehaviorBoss, BossTier: 2,
		BurstName: "chime of ruin", BurstRange: 7, BurstDamage: 12, BurstCooldown: 2, BurstStatus: StatusScorched,
		BurstStatusTurns: 2, BurstStatusPotency: 1,
	},
	{
		ID: "censer_matriarch", Name: "Censer Matriarch", Glyph: '◆', ASCII: 'M', Tint: "#df6464",
		Description: "A smoke-crowned mother of cinders who burns the room by breathing.",
		MaxHP:       104, Attack: 17, Defense: 6, Sight: 11, XPReward: 140, GoldMin: 70, GoldMax: 84,
		MinFloor: 15, MaxFloor: 15, Weight: 1, Behavior: BehaviorBoss, BossTier: 3,
		BurstName: "incense storm", BurstRange: 6, BurstDamage: 14, BurstCooldown: 2, BurstStatus: StatusPoison,
		BurstStatusTurns: 4, BurstStatusPotency: 2,
	},
	{
		ID: "ashen_prior", Name: "Ashen Prior", Glyph: '✠', ASCII: 'P', Tint: "#f05f5f",
		Description: "Last abbot of the ember rite, still kneeling over the crown.",
		MaxHP:       148, Attack: 20, Defense: 8, Sight: 12, XPReward: 220, GoldMin: 120, GoldMax: 150,
		MinFloor: 20, MaxFloor: 20, Weight: 1, Behavior: BehaviorBoss, CanOpenDoors: true, BossTier: 4,
		BurstName: "funeral litany", BurstRange: 8, BurstDamage: 18, BurstCooldown: 2, BurstStatus: StatusScorched,
		BurstStatusTurns: 3, BurstStatusPotency: 2, EnrageThreshold: 50, EnrageAttackBonus: 4,
	},
}

var routePool = []FloorModifier{
	{
		ID: "gilded_way", Title: "Gilded Way", Subtitle: "Coin-pitted halls", Summary: "More gold and a bronze key on arrival.",
		BonusGold: 0.35, GuaranteedKey: keyTierPtr(KeyBronze),
	},
	{
		ID: "brokers_lantern", Title: "Broker's Lantern", Subtitle: "A warm lamp beyond barred teeth", Summary: "A merchant waits on the next floor.",
		Merchant: true,
	},
	{
		ID: "pilgrims_rest", Title: "Pilgrim's Rest", Subtitle: "Quiet pews and ash-water", Summary: "Heal, cleanse, and face a gentler floor.",
		Rest: true, HealOnStart: 16, CleanseOnRest: true, EnemyBonus: -2,
	},
	{
		ID: "reliquary_breach", Title: "Reliquary Breach", Subtitle: "Sealed vaults split open", Summary: "Extra chest and richer loot.",
		LootBonus: 2, ExtraChests: 1, GuaranteedKey: keyTierPtr(KeySilver),
	},
	{
		ID: "ashen_hunt", Title: "Ashen Hunt", Subtitle: "Marked prey stalks the dark", Summary: "More elites, better drops, sharper combat pressure.",
		LootBonus: 1, EnemyBonus: 1, EliteChance: 0.3,
	},
	{
		ID: "cursed_procession", Title: "Cursed Procession", Subtitle: "The abbey sings back", Summary: "Harder floor, more gold, extra chest, best rewards.",
		BonusGold: 0.4, LootBonus: 2, EnemyBonus: 2, EliteChance: 0.45, ExtraChests: 1, GuaranteedKey: keyTierPtr(KeySilver), Cursed: true,
	},
}

func FloorThemeName(level int) string {
	if level > 20 {
		return "Ashfall Depths " + itoa(level-20)
	}
	if level < 1 || level >= len(floorThemes) {
		return "Forgotten Vaults"
	}
	return floorThemes[level].Name
}

func FloorIntro(level int) string {
	if level > 20 {
		return "The crown has been taken, but the dark below the abbey keeps going. It has started learning your name."
	}
	if level < 1 || level >= len(floorThemes) {
		return "The dark waits."
	}
	return floorThemes[level].Intro
}

func GenerateRouteChoices(rng *RNG, nextFloor int, maxFloors int, endless bool) []RouteChoice {
	choices := make([]RouteChoice, 0, 3)
	indices := rng.Perm(len(routePool))
	bossFloor := isBossFloor(nextFloor, maxFloors, endless)
	for _, index := range indices[:3] {
		modifier := routePool[index]
		reward := modifier.Summary
		risk := "Steady"
		if modifier.Cursed {
			risk = "Cursed"
		} else if modifier.EliteChance > 0 {
			risk = "Sharper"
		} else if modifier.Rest {
			risk = "Gentle"
		}
		choices = append(choices, RouteChoice{
			ID:        modifier.ID,
			Title:     modifier.Title,
			Subtitle:  modifier.Subtitle,
			Reward:    reward,
			Risk:      risk,
			MapLabel:  floorMapLabel(nextFloor, bossFloor),
			Modifier:  modifier,
			BossFloor: bossFloor,
		})
	}
	return choices
}

func RandomEnemyTemplate(rng *RNG, floor int) EnemyTemplate {
	candidates := make([]EnemyTemplate, 0, len(enemyCatalog))
	total := 0
	for _, template := range enemyCatalog {
		if template.BossTier > 0 {
			continue
		}
		if floor >= template.MinFloor && floor <= template.MaxFloor {
			candidates = append(candidates, template)
			total += template.Weight
		}
	}
	if len(candidates) == 0 {
		return enemyCatalog[0]
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

func BossTemplateForFloor(rng *RNG, floor int, maxFloors int, endless bool) EnemyTemplate {
	for _, template := range enemyCatalog {
		if template.MinFloor == floor && template.MaxFloor == floor && template.BossTier > 0 {
			return template
		}
	}
	endlessBosses := []string{"houndmaster_vey", "bell_archivist_oria", "censer_matriarch"}
	choice := endlessBosses[rng.Intn(len(endlessBosses))]
	for _, template := range enemyCatalog {
		if template.ID == choice {
			return template
		}
	}
	return enemyCatalog[len(enemyCatalog)-1]
}

func ScaleEnemyTemplate(template EnemyTemplate, floor int, persistentDifficulty int, elite bool, cursed bool) EnemyTemplate {
	scaled := template
	depthScale := max(0, floor-template.MinFloor)
	scaled.MaxHP += depthScale*2 + persistentDifficulty*3
	scaled.Attack += depthScale/2 + persistentDifficulty
	scaled.Defense += depthScale/4 + persistentDifficulty/2
	scaled.XPReward += depthScale*2 + persistentDifficulty*2
	scaled.GoldMin += depthScale / 2
	scaled.GoldMax += depthScale

	if cursed {
		scaled.MaxHP += 4 + floor/2
		scaled.Attack += 2
	}
	if scaled.BossTier > 0 {
		scaled.MaxHP += floor*2 + persistentDifficulty*8 + scaled.BossTier*10
		scaled.Attack += floor/3 + scaled.BossTier*2
		scaled.Defense += scaled.BossTier
		scaled.BurstDamage += scaled.BossTier*2 + floor/6
	}
	if elite {
		scaled.MaxHP += 10 + floor
		scaled.Attack += 2 + floor/6
		scaled.Defense += 1 + floor/10
		scaled.XPReward += 10 + floor/2
		scaled.GoldMin += 5 + floor/3
		scaled.GoldMax += 7 + floor/2
	}
	return scaled
}

func floorMapLabel(nextFloor int, bossFloor bool) string {
	label := "Floor " + itoa(nextFloor)
	if bossFloor {
		return label + "  BOSS"
	}
	return label
}

func isBossFloor(level int, maxFloors int, endless bool) bool {
	if level <= 0 {
		return false
	}
	if level <= maxFloors {
		return level%5 == 0
	}
	return endless && level%5 == 0
}

func keyTierPtr(tier KeyTier) *KeyTier {
	value := tier
	return &value
}
