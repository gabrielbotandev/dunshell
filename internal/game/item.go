package game

import "math/rand"

type ItemKind int

const (
	ItemKindConsumable ItemKind = iota
	ItemKindEquipment
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
	Glyph        rune
	Tint         string
	Description  string
	AttackBonus  int
	DefenseBonus int
	MaxHPBonus   int
	SightBonus   int
	Heal         int
	PoisonCure   bool
	FocusTurns   int
	FocusBonus   int
	EmberDamage  int
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
		if i.Heal > 0 {
			return "+" + itoa(i.Heal) + " HP"
		}
		if i.FocusTurns > 0 {
			return "focus +" + itoa(i.FocusBonus) + " for " + itoa(i.FocusTurns) + " turns"
		}
		if i.EmberDamage > 0 {
			return "ember bolt " + itoa(i.EmberDamage)
		}
		if i.PoisonCure {
			return "cures venom"
		}
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
		return joinParts(parts)
	case ItemKindRelic:
		return "the run-winning relic"
	}
	return i.Description
}

type GroundItem struct {
	Pos       Position
	Item      Item
	RoomIndex int
}

type weightedItem struct {
	ID     string
	Weight int
}

var itemCatalog = map[string]Item{
	"healing_salve": {
		ID: "healing_salve", Name: "Healing Salve", Kind: ItemKindConsumable,
		Glyph: '!', Tint: "#d16078", Description: "A resin salve sealed in waxed cloth.",
		Heal: 10,
	},
	"sunbrew_tonic": {
		ID: "sunbrew_tonic", Name: "Sunbrew Tonic", Kind: ItemKindConsumable,
		Glyph: '!', Tint: "#f29f38", Description: "Distilled emberroot that sharpens the next few swings.",
		FocusTurns: 3, FocusBonus: 2,
	},
	"antivenom": {
		ID: "antivenom", Name: "Antivenom Phial", Kind: ItemKindConsumable,
		Glyph: '!', Tint: "#6bbf59", Description: "Cuts the mire's poison and steadies the breath.",
		Heal: 4, PoisonCure: true,
	},
	"ember_flask": {
		ID: "ember_flask", Name: "Ember Flask", Kind: ItemKindConsumable,
		Glyph: '*', Tint: "#f08a24", Description: "A glass spark that leaps to the nearest threat.",
		EmberDamage: 8,
	},
	"hearth_knife": {
		ID: "hearth_knife", Name: "Hearth Knife", Kind: ItemKindEquipment, Slot: SlotWeapon,
		Glyph: ')', Tint: "#d7c4a3", Description: "A pilgrim's knife, short but honest.",
		AttackBonus: 1,
	},
	"houndfang_dirk": {
		ID: "houndfang_dirk", Name: "Houndfang Dirk", Kind: ItemKindEquipment, Slot: SlotWeapon,
		Glyph: ')', Tint: "#f2c97d", Description: "Light steel with a bite meant for close halls.",
		AttackBonus: 2,
	},
	"abbey_spear": {
		ID: "abbey_spear", Name: "Abbey Spear", Kind: ItemKindEquipment, Slot: SlotWeapon,
		Glyph: ')', Tint: "#f2d7a1", Description: "A ceremonial spear hardened by old rites.",
		AttackBonus: 3, MaxHPBonus: 2,
	},
	"sunsteel_blade": {
		ID: "sunsteel_blade", Name: "Sunsteel Blade", Kind: ItemKindEquipment, Slot: SlotWeapon,
		Glyph: ')', Tint: "#ffe7b3", Description: "Rare metal that keeps an edge even in wet dark.",
		AttackBonus: 4, SightBonus: 1,
	},
	"patched_coat": {
		ID: "patched_coat", Name: "Patched Coat", Kind: ItemKindEquipment, Slot: SlotArmor,
		Glyph: '[', Tint: "#86a0b1", Description: "Layered cloth and leather, enough to turn a rat's bite.",
		DefenseBonus: 1,
	},
	"pilgrim_mail": {
		ID: "pilgrim_mail", Name: "Pilgrim Mail", Kind: ItemKindEquipment, Slot: SlotArmor,
		Glyph: '[', Tint: "#7ea8c7", Description: "Ring-mail with prayer tags knotted into the collar.",
		DefenseBonus: 2, MaxHPBonus: 2,
	},
	"cathedral_plate": {
		ID: "cathedral_plate", Name: "Cathedral Plate", Kind: ItemKindEquipment, Slot: SlotArmor,
		Glyph: '[', Tint: "#b6cad8", Description: "Sacristy armor too heavy for processions and just right here.",
		DefenseBonus: 3, MaxHPBonus: 4,
	},
	"lantern_charm": {
		ID: "lantern_charm", Name: "Lantern Charm", Kind: ItemKindEquipment, Slot: SlotCharm,
		Glyph: '=', Tint: "#f6db7d", Description: "A pocket sigil that coaxes light from old oil.",
		SightBonus: 1,
	},
	"warding_bead": {
		ID: "warding_bead", Name: "Warding Bead", Kind: ItemKindEquipment, Slot: SlotCharm,
		Glyph: '=', Tint: "#7fc8f8", Description: "A cool bead that stills the hand before impact.",
		DefenseBonus: 1, SightBonus: 1,
	},
	"vow_sigil": {
		ID: "vow_sigil", Name: "Vow Sigil", Kind: ItemKindEquipment, Slot: SlotCharm,
		Glyph: '=', Tint: "#f09c66", Description: "An oath-brand that rewards bold steps.",
		AttackBonus: 1, MaxHPBonus: 2,
	},
	"cinder_crown": {
		ID: "cinder_crown", Name: "Cinder Crown", Kind: ItemKindRelic,
		Glyph: '&', Tint: "#ffb347", Description: "The buried ember relic the abbey could never keep.",
	},
}

var floorLootTables = map[int][]weightedItem{
	1: {
		{ID: "healing_salve", Weight: 34},
		{ID: "sunbrew_tonic", Weight: 20},
		{ID: "antivenom", Weight: 16},
		{ID: "ember_flask", Weight: 10},
		{ID: "houndfang_dirk", Weight: 8},
		{ID: "pilgrim_mail", Weight: 7},
		{ID: "warding_bead", Weight: 5},
	},
	2: {
		{ID: "healing_salve", Weight: 24},
		{ID: "sunbrew_tonic", Weight: 18},
		{ID: "antivenom", Weight: 16},
		{ID: "ember_flask", Weight: 14},
		{ID: "houndfang_dirk", Weight: 10},
		{ID: "abbey_spear", Weight: 8},
		{ID: "pilgrim_mail", Weight: 8},
		{ID: "warding_bead", Weight: 8},
		{ID: "vow_sigil", Weight: 6},
	},
	3: {
		{ID: "healing_salve", Weight: 18},
		{ID: "sunbrew_tonic", Weight: 16},
		{ID: "antivenom", Weight: 14},
		{ID: "ember_flask", Weight: 16},
		{ID: "abbey_spear", Weight: 12},
		{ID: "sunsteel_blade", Weight: 7},
		{ID: "pilgrim_mail", Weight: 8},
		{ID: "cathedral_plate", Weight: 6},
		{ID: "vow_sigil", Weight: 8},
	},
	4: {
		{ID: "healing_salve", Weight: 16},
		{ID: "sunbrew_tonic", Weight: 16},
		{ID: "antivenom", Weight: 12},
		{ID: "ember_flask", Weight: 18},
		{ID: "abbey_spear", Weight: 10},
		{ID: "sunsteel_blade", Weight: 10},
		{ID: "cathedral_plate", Weight: 10},
		{ID: "warding_bead", Weight: 8},
		{ID: "vow_sigil", Weight: 8},
	},
	5: {
		{ID: "healing_salve", Weight: 14},
		{ID: "sunbrew_tonic", Weight: 16},
		{ID: "antivenom", Weight: 10},
		{ID: "ember_flask", Weight: 18},
		{ID: "sunsteel_blade", Weight: 12},
		{ID: "cathedral_plate", Weight: 12},
		{ID: "warding_bead", Weight: 8},
		{ID: "vow_sigil", Weight: 10},
	},
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

func RandomGroundItem(rng *rand.Rand, floor int) Item {
	return weightedRandomItem(rng, floorLootTables[clamp(floor, 1, 5)])
}

func RandomDropItem(rng *rand.Rand, floor int) (Item, bool) {
	chance := 0.18 + float64(floor)*0.04
	if rng.Float64() > chance {
		return Item{}, false
	}
	return RandomGroundItem(rng, floor), true
}

func weightedRandomItem(rng *rand.Rand, table []weightedItem) Item {
	total := 0
	for _, entry := range table {
		total += entry.Weight
	}
	if total == 0 {
		return ItemByID("healing_salve")
	}

	roll := rng.Intn(total)
	for _, entry := range table {
		roll -= entry.Weight
		if roll < 0 {
			return ItemByID(entry.ID)
		}
	}
	return ItemByID(table[len(table)-1].ID)
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
