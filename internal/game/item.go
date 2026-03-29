package game

type ItemKind int

const (
	ItemKindConsumable ItemKind = iota
	ItemKindEquipment
	ItemKindKey
	ItemKindRelic
)

type EquipmentSlot int

const (
	SlotNone EquipmentSlot = iota
	SlotWeapon
	SlotArmor
	SlotCharm
)

func (s EquipmentSlot) Label() string {
	switch s {
	case SlotWeapon:
		return "Weapon"
	case SlotArmor:
		return "Armor"
	case SlotCharm:
		return "Charm"
	default:
		return "Pack"
	}
}

type Item struct {
	ID           string
	Name         string
	Kind         ItemKind
	Slot         EquipmentSlot
	Rarity       Rarity
	Glyph        rune
	ASCII        rune
	Tint         string
	Description  string
	AttackBonus  int
	DefenseBonus int
	MaxHPBonus   int
	SightBonus   int
	Heal         int
	PoisonCure   bool
	FireCure     bool
	FocusTurns   int
	FocusBonus   int
	EmberDamage  int
	PoisonResist int
	FireResist   int
	Price        int
	KeyTier      KeyTier
	MinFloor     int
	MaxFloor     int
	Weight       int
}

func (i Item) ActionLabel() string {
	switch i.Kind {
	case ItemKindConsumable:
		return "Use"
	case ItemKindEquipment:
		return "Equip"
	case ItemKindRelic:
		return "Claim"
	default:
		return "Take"
	}
}

func (i Item) DetailLine() string {
	switch i.Kind {
	case ItemKindConsumable:
		parts := make([]string, 0, 4)
		if i.Heal > 0 {
			parts = append(parts, "+"+itoa(i.Heal)+" HP")
		}
		if i.PoisonCure {
			parts = append(parts, "cures poison")
		}
		if i.FireCure {
			parts = append(parts, "douses fire")
		}
		if i.FocusTurns > 0 {
			parts = append(parts, "focus +"+itoa(i.FocusBonus))
		}
		if i.EmberDamage > 0 {
			parts = append(parts, "ember "+itoa(i.EmberDamage))
		}
		if i.PoisonResist > 0 {
			parts = append(parts, "poison ward +"+itoa(i.PoisonResist))
		}
		if i.FireResist > 0 {
			parts = append(parts, "fire ward +"+itoa(i.FireResist))
		}
		return joinParts(parts)
	case ItemKindEquipment:
		parts := make([]string, 0, 4)
		if i.AttackBonus > 0 {
			parts = append(parts, "ATK+"+itoa(i.AttackBonus))
		}
		if i.DefenseBonus > 0 {
			parts = append(parts, "DEF+"+itoa(i.DefenseBonus))
		}
		if i.MaxHPBonus > 0 {
			parts = append(parts, "HP+"+itoa(i.MaxHPBonus))
		}
		if i.SightBonus > 0 {
			parts = append(parts, "SIGHT+"+itoa(i.SightBonus))
		}
		if i.PoisonResist > 0 {
			parts = append(parts, "POISON+"+itoa(i.PoisonResist))
		}
		if i.FireResist > 0 {
			parts = append(parts, "FIRE+"+itoa(i.FireResist))
		}
		return joinParts(parts)
	case ItemKindKey:
		return i.KeyTier.Label() + " key"
	case ItemKindRelic:
		return "the abbey's buried crown"
	default:
		return i.Description
	}
}

func (i Item) IsUnique() bool {
	return i.Rarity == RarityUnique
}

type GroundItem struct {
	Pos       Position
	Item      Item
	RoomIndex int
}

type MerchantOffer struct {
	Item  Item
	Price int
	Sold  bool
}

var itemCatalog = map[string]Item{
	"healing_salve": {
		ID: "healing_salve", Name: "Healing Salve", Kind: ItemKindConsumable, Rarity: RarityCommon,
		Glyph: '✚', ASCII: '!', Tint: RarityCommon.Tint(), Description: "Waxed cloth packed with bitter resin.",
		Heal: 10, Price: 18, MinFloor: 1, MaxFloor: 99, Weight: 18,
	},
	"pilgrim_bandage": {
		ID: "pilgrim_bandage", Name: "Pilgrim Bandage", Kind: ItemKindConsumable, Rarity: RarityCommon,
		Glyph: '✚', ASCII: '!', Tint: RarityCommon.Tint(), Description: "A tighter wrap blessed for long descents.",
		Heal: 16, Price: 32, MinFloor: 4, MaxFloor: 99, Weight: 15,
	},
	"antivenom_phial": {
		ID: "antivenom_phial", Name: "Antivenom Phial", Kind: ItemKindConsumable, Rarity: RarityUncommon,
		Glyph: '✚', ASCII: '!', Tint: RarityUncommon.Tint(), Description: "Cuts mire poison and steadies the lungs.",
		Heal: 8, PoisonCure: true, Price: 28, MinFloor: 2, MaxFloor: 99, Weight: 12,
	},
	"dousing_salts": {
		ID: "dousing_salts", Name: "Dousing Salts", Kind: ItemKindConsumable, Rarity: RarityUncommon,
		Glyph: '✚', ASCII: '!', Tint: RarityUncommon.Tint(), Description: "Ash-cold salts that smother ember cling and sting.",
		Heal: 8, FireCure: true, Price: 30, MinFloor: 4, MaxFloor: 99, Weight: 10,
	},
	"sunbrew_tonic": {
		ID: "sunbrew_tonic", Name: "Sunbrew Tonic", Kind: ItemKindConsumable, Rarity: RarityUncommon,
		Glyph: '✚', ASCII: '!', Tint: RarityUncommon.Tint(), Description: "Distilled emberroot that sharpens the next exchange.",
		FocusTurns: 4, FocusBonus: 2, Price: 40, MinFloor: 3, MaxFloor: 99, Weight: 10,
	},
	"ember_phial": {
		ID: "ember_phial", Name: "Ember Phial", Kind: ItemKindConsumable, Rarity: RarityRare,
		Glyph: '✶', ASCII: '*', Tint: RarityRare.Tint(), Description: "A captive spark that hunts the nearest threat.",
		EmberDamage: 12, Price: 58, MinFloor: 5, MaxFloor: 99, Weight: 9,
	},
	"grave_amber": {
		ID: "grave_amber", Name: "Grave Amber", Kind: ItemKindConsumable, Rarity: RarityRare,
		Glyph: '✚', ASCII: '!', Tint: RarityRare.Tint(), Description: "Old sap burned down into a tar-black draught.",
		Heal: 24, Price: 72, MinFloor: 10, MaxFloor: 99, Weight: 6,
	},
	"saintroot_draught": {
		ID: "saintroot_draught", Name: "Saintroot Draught", Kind: ItemKindConsumable, Rarity: RarityLegendary,
		Glyph: '✚', ASCII: '!', Tint: RarityLegendary.Tint(), Description: "Restores torn flesh and breaks the chapel's poisons and cinders.",
		Heal: 30, PoisonCure: true, FireCure: true, Price: 96, MinFloor: 12, MaxFloor: 99, Weight: 3,
	},
	"hearth_knife": {
		ID: "hearth_knife", Name: "Hearth Knife", Kind: ItemKindEquipment, Slot: SlotWeapon, Rarity: RarityCommon,
		Glyph: '†', ASCII: ')', Tint: RarityCommon.Tint(), Description: "A pilgrim's short blade, honest and plain.",
		AttackBonus: 1, Price: 24, MinFloor: 1, MaxFloor: 6, Weight: 12,
	},
	"gravehook": {
		ID: "gravehook", Name: "Gravehook", Kind: ItemKindEquipment, Slot: SlotWeapon, Rarity: RarityCommon,
		Glyph: '†', ASCII: ')', Tint: RarityCommon.Tint(), Description: "Curved steel dragged from a robber's satchel.",
		AttackBonus: 2, Price: 34, MinFloor: 2, MaxFloor: 9, Weight: 11,
	},
	"chapel_pike": {
		ID: "chapel_pike", Name: "Chapel Pike", Kind: ItemKindEquipment, Slot: SlotWeapon, Rarity: RarityUncommon,
		Glyph: '†', ASCII: ')', Tint: RarityUncommon.Tint(), Description: "Ceremonial reach reforged for narrow halls.",
		AttackBonus: 3, MaxHPBonus: 2, Price: 52, MinFloor: 4, MaxFloor: 12, Weight: 9,
	},
	"lantern_falx": {
		ID: "lantern_falx", Name: "Lantern Falx", Kind: ItemKindEquipment, Slot: SlotWeapon, Rarity: RarityUncommon,
		Glyph: '†', ASCII: ')', Tint: RarityUncommon.Tint(), Description: "A crescent blade keen enough to work by candlelight.",
		AttackBonus: 3, SightBonus: 1, Price: 58, MinFloor: 6, MaxFloor: 14, Weight: 8,
	},
	"ossuary_blade": {
		ID: "ossuary_blade", Name: "Ossuary Blade", Kind: ItemKindEquipment, Slot: SlotWeapon, Rarity: RarityRare,
		Glyph: '†', ASCII: ')', Tint: RarityRare.Tint(), Description: "Bone-dusted steel with a brutal backswing.",
		AttackBonus: 5, DefenseBonus: 1, Price: 88, MinFloor: 8, MaxFloor: 18, Weight: 7,
	},
	"saintbreaker_maul": {
		ID: "saintbreaker_maul", Name: "Saintbreaker Maul", Kind: ItemKindEquipment, Slot: SlotWeapon, Rarity: RarityRare,
		Glyph: '†', ASCII: ')', Tint: RarityRare.Tint(), Description: "A reliquary hammer too heavy for processions.",
		AttackBonus: 6, MaxHPBonus: 4, Price: 108, MinFloor: 11, MaxFloor: 20, Weight: 6,
	},
	"sunsteel_blade": {
		ID: "sunsteel_blade", Name: "Sunsteel Blade", Kind: ItemKindEquipment, Slot: SlotWeapon, Rarity: RarityLegendary,
		Glyph: '†', ASCII: ')', Tint: RarityLegendary.Tint(), Description: "Rare steel that glows warm in soaked darkness.",
		AttackBonus: 7, SightBonus: 2, Price: 144, MinFloor: 14, MaxFloor: 99, Weight: 3,
	},
	"thorn_of_the_prior": {
		ID: "thorn_of_the_prior", Name: "Thorn of the Prior", Kind: ItemKindEquipment, Slot: SlotWeapon, Rarity: RarityUnique,
		Glyph: '✠', ASCII: ')', Tint: RarityUnique.Tint(), Description: "The ashen crozier split into a killing edge.",
		AttackBonus: 9, MaxHPBonus: 4, SightBonus: 2, Price: 0, MinFloor: 20, MaxFloor: 99, Weight: 1,
	},
	"patched_coat": {
		ID: "patched_coat", Name: "Patched Coat", Kind: ItemKindEquipment, Slot: SlotArmor, Rarity: RarityCommon,
		Glyph: '⛨', ASCII: '[', Tint: RarityCommon.Tint(), Description: "Layered cloth and leather, enough to turn a rat's bite.",
		DefenseBonus: 1, Price: 22, MinFloor: 1, MaxFloor: 6, Weight: 12,
	},
	"pilgrim_mail": {
		ID: "pilgrim_mail", Name: "Pilgrim Mail", Kind: ItemKindEquipment, Slot: SlotArmor, Rarity: RarityCommon,
		Glyph: '⛨', ASCII: '[', Tint: RarityCommon.Tint(), Description: "Prayer tags still knot the collar closed.",
		DefenseBonus: 2, MaxHPBonus: 3, Price: 40, MinFloor: 3, MaxFloor: 10, Weight: 10,
	},
	"gravewax_hauberk": {
		ID: "gravewax_hauberk", Name: "Gravewax Hauberk", Kind: ItemKindEquipment, Slot: SlotArmor, Rarity: RarityUncommon,
		Glyph: '⛨', ASCII: '[', Tint: RarityUncommon.Tint(), Description: "Wax-black rings that drink the first shock of a blow.",
		DefenseBonus: 3, MaxHPBonus: 4, PoisonResist: 1, Price: 60, MinFloor: 6, MaxFloor: 14, Weight: 8,
	},
	"floodplate_cuirass": {
		ID: "floodplate_cuirass", Name: "Floodplate Cuirass", Kind: ItemKindEquipment, Slot: SlotArmor, Rarity: RarityRare,
		Glyph: '⛨', ASCII: '[', Tint: RarityRare.Tint(), Description: "River-worn plate that still answers to the strap.",
		DefenseBonus: 4, MaxHPBonus: 6, Price: 92, MinFloor: 9, MaxFloor: 18, Weight: 7,
	},
	"warden_harness": {
		ID: "warden_harness", Name: "Warden Harness", Kind: ItemKindEquipment, Slot: SlotArmor, Rarity: RarityRare,
		Glyph: '⛨', ASCII: '[', Tint: RarityRare.Tint(), Description: "A shrine guard's rig threaded with watch-lamps.",
		DefenseBonus: 4, MaxHPBonus: 4, SightBonus: 1, Price: 100, MinFloor: 11, MaxFloor: 20, Weight: 6,
	},
	"cathedral_plate": {
		ID: "cathedral_plate", Name: "Cathedral Plate", Kind: ItemKindEquipment, Slot: SlotArmor, Rarity: RarityLegendary,
		Glyph: '⛨', ASCII: '[', Tint: RarityLegendary.Tint(), Description: "Sacristy armor too heavy for processions and perfect for war.",
		DefenseBonus: 5, MaxHPBonus: 8, Price: 138, MinFloor: 14, MaxFloor: 99, Weight: 3,
	},
	"ashen_vestments": {
		ID: "ashen_vestments", Name: "Ashen Vestments", Kind: ItemKindEquipment, Slot: SlotArmor, Rarity: RarityUnique,
		Glyph: '✠', ASCII: '[', Tint: RarityUnique.Tint(), Description: "The prior's funerary layers, still warm with condemned prayer.",
		DefenseBonus: 6, MaxHPBonus: 12, SightBonus: 2, FireResist: 2, Price: 0, MinFloor: 20, MaxFloor: 99, Weight: 1,
	},
	"lantern_charm": {
		ID: "lantern_charm", Name: "Lantern Charm", Kind: ItemKindEquipment, Slot: SlotCharm, Rarity: RarityCommon,
		Glyph: '◌', ASCII: '=', Tint: RarityCommon.Tint(), Description: "A pocket sigil that coaxes light from old oil.",
		SightBonus: 1, Price: 20, MinFloor: 1, MaxFloor: 8, Weight: 12,
	},
	"warding_bead": {
		ID: "warding_bead", Name: "Warding Bead", Kind: ItemKindEquipment, Slot: SlotCharm, Rarity: RarityCommon,
		Glyph: '◌', ASCII: '=', Tint: RarityCommon.Tint(), Description: "A cool bead that stills the hand before impact.",
		DefenseBonus: 1, FireResist: 1, Price: 20, MinFloor: 1, MaxFloor: 9, Weight: 11,
	},
	"vow_sigil": {
		ID: "vow_sigil", Name: "Vow Sigil", Kind: ItemKindEquipment, Slot: SlotCharm, Rarity: RarityUncommon,
		Glyph: '◌', ASCII: '=', Tint: RarityUncommon.Tint(), Description: "An oath-brand that rewards bold steps.",
		AttackBonus: 1, MaxHPBonus: 2, Price: 36, MinFloor: 3, MaxFloor: 12, Weight: 9,
	},
	"salt_reliquary": {
		ID: "salt_reliquary", Name: "Salt Reliquary", Kind: ItemKindEquipment, Slot: SlotCharm, Rarity: RarityUncommon,
		Glyph: '◌', ASCII: '=', Tint: RarityUncommon.Tint(), Description: "A dry charm that keeps rot and fear at arm's reach.",
		DefenseBonus: 2, SightBonus: 1, PoisonResist: 1, Price: 52, MinFloor: 6, MaxFloor: 14, Weight: 8,
	},
	"blackglass_eye": {
		ID: "blackglass_eye", Name: "Blackglass Eye", Kind: ItemKindEquipment, Slot: SlotCharm, Rarity: RarityRare,
		Glyph: '◌', ASCII: '=', Tint: RarityRare.Tint(), Description: "A polished eye that catches hostile movement before the foot does.",
		AttackBonus: 2, SightBonus: 2, Price: 82, MinFloor: 9, MaxFloor: 18, Weight: 6,
	},
	"choir_knot": {
		ID: "choir_knot", Name: "Choir Knot", Kind: ItemKindEquipment, Slot: SlotCharm, Rarity: RarityLegendary,
		Glyph: '◌', ASCII: '=', Tint: RarityLegendary.Tint(), Description: "A severed braid bound around a prayer bell.",
		AttackBonus: 2, DefenseBonus: 2, MaxHPBonus: 4, PoisonResist: 1, FireResist: 1, Price: 124, MinFloor: 14, MaxFloor: 99, Weight: 3,
	},
	"crownshard_rosary": {
		ID: "crownshard_rosary", Name: "Crownshard Rosary", Kind: ItemKindEquipment, Slot: SlotCharm, Rarity: RarityUnique,
		Glyph: '✠', ASCII: '=', Tint: RarityUnique.Tint(), Description: "A string of splinters chipped from the crown's first casing.",
		AttackBonus: 3, DefenseBonus: 3, MaxHPBonus: 6, SightBonus: 2, PoisonResist: 1, FireResist: 1, Price: 0, MinFloor: 20, MaxFloor: 99, Weight: 1,
	},
	"bronze_key": {
		ID: "bronze_key", Name: "Bronze Key", Kind: ItemKindKey, Rarity: RarityCommon,
		Glyph: '⚿', ASCII: 'k', Tint: KeyBronze.Tint(), Description: "Fits the low vaults and pilgrim caskets.",
		Price: 28, KeyTier: KeyBronze, MinFloor: 1, MaxFloor: 99, Weight: 1,
	},
	"silver_key": {
		ID: "silver_key", Name: "Silver Key", Kind: ItemKindKey, Rarity: RarityRare,
		Glyph: '⚿', ASCII: 'k', Tint: KeySilver.Tint(), Description: "Cut for choir locks and grave treasuries.",
		Price: 56, KeyTier: KeySilver, MinFloor: 4, MaxFloor: 99, Weight: 1,
	},
	"gold_key": {
		ID: "gold_key", Name: "Gold Key", Kind: ItemKindKey, Rarity: RarityLegendary,
		Glyph: '⚿', ASCII: 'k', Tint: KeyGold.Tint(), Description: "Reserved for reliquaries and crowned vaults.",
		Price: 108, KeyTier: KeyGold, MinFloor: 10, MaxFloor: 99, Weight: 1,
	},
	"cinder_crown": {
		ID: "cinder_crown", Name: "Cinder Crown", Kind: ItemKindRelic, Rarity: RarityLegendary,
		Glyph: '♛', ASCII: '&', Tint: "#ffb347", Description: "The ember relic the abbey drowned itself to hide.",
		MinFloor: 20, MaxFloor: 99, Weight: 1,
	},
}

var consumableIDs = []string{
	"healing_salve",
	"pilgrim_bandage",
	"antivenom_phial",
	"dousing_salts",
	"sunbrew_tonic",
	"ember_phial",
	"grave_amber",
	"saintroot_draught",
}

var weaponIDs = []string{
	"hearth_knife",
	"gravehook",
	"chapel_pike",
	"lantern_falx",
	"ossuary_blade",
	"saintbreaker_maul",
	"sunsteel_blade",
}

var armorIDs = []string{
	"patched_coat",
	"pilgrim_mail",
	"gravewax_hauberk",
	"floodplate_cuirass",
	"warden_harness",
	"cathedral_plate",
}

var charmIDs = []string{
	"lantern_charm",
	"warding_bead",
	"vow_sigil",
	"salt_reliquary",
	"blackglass_eye",
	"choir_knot",
}

var uniqueIDs = []string{
	"thorn_of_the_prior",
	"ashen_vestments",
	"crownshard_rosary",
}

func ItemByID(id string) Item {
	return itemCatalog[id]
}

func StarterEquipment() (Item, Item, Item) {
	return ItemByID("hearth_knife"), ItemByID("patched_coat"), ItemByID("lantern_charm")
}

func StarterInventory() []Item {
	return []Item{
		ItemByID("healing_salve"),
		ItemByID("healing_salve"),
		ItemByID("sunbrew_tonic"),
	}
}

func KeyItem(tier KeyTier) Item {
	switch tier {
	case KeyBronze:
		return ItemByID("bronze_key")
	case KeySilver:
		return ItemByID("silver_key")
	case KeyGold:
		return ItemByID("gold_key")
	default:
		return ItemByID("bronze_key")
	}
}

func RandomGroundItem(rng *RNG, floor int, modifier FloorModifier) Item {
	consumableChance := 0.44
	if floor >= 10 {
		consumableChance = 0.36
	}
	if modifier.Rest {
		consumableChance += 0.08
	}
	if rng.Float64() < consumableChance {
		return randomConsumable(rng, floor)
	}

	slots := []EquipmentSlot{SlotWeapon, SlotArmor, SlotCharm}
	return randomEquipmentForSlot(rng, floor, slots[rng.Intn(len(slots))], modifier.LootBonus)
}

func RandomDropItem(rng *RNG, floor int, modifier FloorModifier, elite bool) (Item, bool) {
	chance := 0.18 + float64(floor)*0.012 + modifier.BonusGold*0.05
	if elite {
		chance += 0.18
	}
	if modifier.Cursed {
		chance += 0.08
	}
	if chance > 0.72 {
		chance = 0.72
	}
	if rng.Float64() > chance {
		return Item{}, false
	}
	bonus := modifier.LootBonus
	if elite {
		bonus++
	}
	return RandomGroundItem(rng, floor, FloorModifier{LootBonus: bonus}), true
}

func RandomKeyReward(rng *RNG, floor int) Item {
	roll := rng.Intn(100)
	switch {
	case floor >= 12 && roll > 82:
		return KeyItem(KeyGold)
	case floor >= 5 && roll > 45:
		return KeyItem(KeySilver)
	default:
		return KeyItem(KeyBronze)
	}
}

func GenerateChestRewards(rng *RNG, tier KeyTier, floor int, modifier FloorModifier, bossReward bool, finalBoss bool) []ChestReward {
	rewards := make([]ChestReward, 0, 5)
	baseGold := 18 + floor*6
	switch tier {
	case KeyBronze:
		baseGold += 12
	case KeySilver:
		baseGold += 26
	case KeyGold:
		baseGold += 48
	}
	if bossReward {
		baseGold += 32 + floor*4
	}
	rewards = append(rewards, ChestReward{Kind: RewardGold, Gold: baseGold})

	lootBonus := modifier.LootBonus
	if bossReward {
		lootBonus += 2
	}

	switch tier {
	case KeyBronze:
		rewards = append(rewards, ChestReward{Kind: RewardItem, Item: RandomGroundItem(rng, floor, FloorModifier{LootBonus: lootBonus})})
		if rng.Float64() < 0.35 {
			rewards = append(rewards, ChestReward{Kind: RewardItem, Item: randomConsumable(rng, floor)})
		}
	case KeySilver:
		rewards = append(rewards, ChestReward{Kind: RewardItem, Item: randomEquipmentForSlot(rng, floor, randomEquipmentSlot(rng), lootBonus+1)})
		rewards = append(rewards, ChestReward{Kind: RewardItem, Item: randomConsumable(rng, floor)})
		if floor >= 7 && rng.Float64() < 0.4 {
			rewards = append(rewards, ChestReward{Kind: RewardItem, Item: KeyItem(KeyBronze)})
		}
	case KeyGold:
		rewards = append(rewards,
			ChestReward{Kind: RewardItem, Item: randomEquipmentForSlot(rng, floor, randomEquipmentSlot(rng), lootBonus+2)},
			ChestReward{Kind: RewardItem, Item: randomEquipmentForSlot(rng, floor, randomEquipmentSlot(rng), lootBonus+1)},
		)
		if rng.Float64() < 0.55 {
			rewards = append(rewards, ChestReward{Kind: RewardItem, Item: randomConsumable(rng, floor)})
		}
	}

	if bossReward {
		rewards = append(rewards, ChestReward{Kind: RewardItem, Item: randomEquipmentForSlot(rng, floor+1, randomEquipmentSlot(rng), lootBonus+2)})
	}
	if finalBoss {
		rewards = append(rewards,
			ChestReward{Kind: RewardItem, Item: RandomUniqueItem(rng)},
			ChestReward{Kind: RewardItem, Item: ItemByID("cinder_crown")},
		)
	}

	return rewards
}

func GenerateMerchantOffers(rng *RNG, floor int) []MerchantOffer {
	offers := []MerchantOffer{
		{Item: pickMerchantHealingItem(rng, floor)},
		{Item: randomEquipmentForSlot(rng, floor+1, SlotWeapon, 1)},
		{Item: randomEquipmentForSlot(rng, floor+1, SlotArmor, 1)},
		{Item: randomEquipmentForSlot(rng, floor+1, SlotCharm, 1)},
	}

	switch rng.Intn(4) {
	case 0:
		offers = append(offers, MerchantOffer{Item: pickMerchantUtilityItem(rng, floor)})
	case 1:
		offers = append(offers, MerchantOffer{Item: KeyItem(merchantKeyTier(floor))})
	default:
		offers = append(offers, MerchantOffer{Item: randomEquipmentForSlot(rng, floor+1, randomEquipmentSlot(rng), 2)})
	}

	seen := map[string]bool{}
	for index := range offers {
		for seen[offers[index].Item.ID] {
			offers[index].Item = rerollMerchantItem(rng, floor, index)
		}
		seen[offers[index].Item.ID] = true
		offers[index].Price = offers[index].Item.Price
	}
	return offers
}

func RandomUniqueItem(rng *RNG) Item {
	return ItemByID(uniqueIDs[rng.Intn(len(uniqueIDs))])
}

func RandomUniqueEnemyDrop(rng *RNG, floor int, elite bool, bossTier int) (Item, bool) {
	if bossTier >= 4 {
		return Item{}, false
	}

	chance := 0.000015 + float64(max(0, floor-1))*0.000002
	if elite {
		chance += 0.0001
	}
	if bossTier > 0 {
		chance += float64(bossTier) * 0.00035
	}
	if chance > 0.0012 {
		chance = 0.0012
	}
	if rng.Float64() > chance {
		return Item{}, false
	}
	return RandomUniqueItem(rng), true
}

func randomConsumable(rng *RNG, floor int) Item {
	table := eligibleItemPool(consumableIDs, floor, RarityCommon, RarityLegendary)
	return weightedRandomItem(rng, table)
}

func randomEquipmentForSlot(rng *RNG, floor int, slot EquipmentSlot, lootBonus int) Item {
	rarity := rollRarity(rng, floor, lootBonus)
	ids := weaponIDs
	switch slot {
	case SlotArmor:
		ids = armorIDs
	case SlotCharm:
		ids = charmIDs
	}

	for current := rarity; current >= RarityCommon; current-- {
		candidates := eligibleItemPool(ids, floor, current, current)
		if len(candidates) > 0 {
			return weightedRandomItem(rng, candidates)
		}
	}
	return weightedRandomItem(rng, eligibleItemPool(ids, floor, RarityCommon, RarityLegendary))
}

func randomEquipmentSlot(rng *RNG) EquipmentSlot {
	slots := []EquipmentSlot{SlotWeapon, SlotArmor, SlotCharm}
	return slots[rng.Intn(len(slots))]
}

func rollRarity(rng *RNG, floor int, lootBonus int) Rarity {
	common := 62
	uncommon := 26
	rare := 10
	legendary := 2
	switch {
	case floor >= 15:
		common, uncommon, rare, legendary = 18, 34, 32, 16
	case floor >= 10:
		common, uncommon, rare, legendary = 30, 34, 26, 10
	case floor >= 5:
		common, uncommon, rare, legendary = 44, 32, 19, 5
	}

	common -= lootBonus * 7
	uncommon += lootBonus * 2
	rare += lootBonus * 3
	legendary += lootBonus * 2
	if common < 8 {
		common = 8
	}

	roll := rng.Intn(common + uncommon + rare + legendary)
	roll -= common
	if roll < 0 {
		return RarityCommon
	}
	roll -= uncommon
	if roll < 0 {
		return RarityUncommon
	}
	roll -= rare
	if roll < 0 {
		return RarityRare
	}
	return RarityLegendary
}

func eligibleItemPool(ids []string, floor int, minRarity Rarity, maxRarity Rarity) []Item {
	pool := make([]Item, 0, len(ids))
	for _, id := range ids {
		item := ItemByID(id)
		if floor < item.MinFloor || floor > item.MaxFloor {
			continue
		}
		if item.Rarity < minRarity || item.Rarity > maxRarity {
			continue
		}
		pool = append(pool, item)
	}
	return pool
}

func weightedRandomItem(rng *RNG, items []Item) Item {
	if len(items) == 0 {
		return ItemByID("healing_salve")
	}
	total := 0
	for _, item := range items {
		weight := item.Weight
		if weight <= 0 {
			weight = 1
		}
		total += weight
	}
	roll := rng.Intn(total)
	for _, item := range items {
		weight := item.Weight
		if weight <= 0 {
			weight = 1
		}
		roll -= weight
		if roll < 0 {
			return item
		}
	}
	return items[len(items)-1]
}

func merchantKeyTier(floor int) KeyTier {
	switch {
	case floor >= 12:
		return KeyGold
	case floor >= 6:
		return KeySilver
	default:
		return KeyBronze
	}
}

func pickMerchantHealingItem(rng *RNG, floor int) Item {
	pool := []Item{ItemByID("healing_salve"), ItemByID("pilgrim_bandage")}
	if floor >= 10 {
		pool = append(pool, ItemByID("grave_amber"))
	}
	if floor >= 12 {
		pool = append(pool, ItemByID("saintroot_draught"))
	}
	return pool[rng.Intn(len(pool))]
}

func pickMerchantUtilityItem(rng *RNG, floor int) Item {
	pool := []Item{ItemByID("antivenom_phial"), ItemByID("dousing_salts"), ItemByID("sunbrew_tonic"), ItemByID("ember_phial")}
	if floor >= 10 {
		pool = append(pool, ItemByID("grave_amber"))
	}
	return pool[rng.Intn(len(pool))]
}

func rerollMerchantItem(rng *RNG, floor int, slotIndex int) Item {
	switch slotIndex {
	case 0:
		return pickMerchantHealingItem(rng, floor)
	case 1:
		return randomEquipmentForSlot(rng, floor+1, SlotWeapon, 1)
	case 2:
		return randomEquipmentForSlot(rng, floor+1, SlotArmor, 1)
	case 3:
		return randomEquipmentForSlot(rng, floor+1, SlotCharm, 1)
	default:
		if rng.Intn(2) == 0 {
			return pickMerchantUtilityItem(rng, floor)
		}
		return KeyItem(merchantKeyTier(floor))
	}
}

func joinParts(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	joined := parts[0]
	for index := 1; index < len(parts); index++ {
		joined += "  " + parts[index]
	}
	return joined
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}
	digits := ""
	for value > 0 {
		digit := value % 10
		digits = string(rune('0'+digit)) + digits
		value /= 10
	}
	return sign + digits
}
