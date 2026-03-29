package game

type StatusKind int

const (
	StatusPoison StatusKind = iota
	StatusFocus
	StatusFire
)

const StatusScorched = StatusFire

type StatusEffect struct {
	Kind    StatusKind
	Turns   int
	Potency int
}

func (s StatusEffect) Label() string {
	suffix := itoa(s.Turns) + "t"
	switch s.Kind {
	case StatusPoison:
		return "Poison " + suffix + " x" + itoa(max(1, s.Potency))
	case StatusFocus:
		return "Focus +" + itoa(max(1, s.Potency)) + " " + suffix
	case StatusFire:
		return "Fire " + suffix + " x" + itoa(max(1, s.Potency))
	default:
		return "Unknown"
	}
}

func (s StatusEffect) ShortLabel() string {
	potency := ""
	if s.Potency > 1 {
		potency = "x" + itoa(s.Potency)
	}
	switch s.Kind {
	case StatusPoison:
		return "VENOM " + itoa(s.Turns) + potency
	case StatusFocus:
		return "FOCUS " + itoa(s.Turns)
	case StatusFire:
		return "FIRE " + itoa(s.Turns) + potency
	default:
		return "STATE"
	}
}

func (s StatusEffect) Harmful() bool {
	return s.Kind != StatusFocus
}

func (k StatusKind) Name() string {
	switch k {
	case StatusPoison:
		return "poison"
	case StatusFocus:
		return "focus"
	case StatusFire:
		return "fire"
	default:
		return "status"
	}
}

type KeyRing struct {
	Bronze int
	Silver int
	Gold   int
}

func (k *KeyRing) Add(tier KeyTier, amount int) {
	switch tier {
	case KeyBronze:
		k.Bronze += amount
	case KeySilver:
		k.Silver += amount
	case KeyGold:
		k.Gold += amount
	}
}

func (k KeyRing) Count(tier KeyTier) int {
	switch tier {
	case KeyBronze:
		return k.Bronze
	case KeySilver:
		return k.Silver
	case KeyGold:
		return k.Gold
	default:
		return 0
	}
}

func (k *KeyRing) Spend(tier KeyTier) bool {
	if k.Count(tier) <= 0 {
		return false
	}
	switch tier {
	case KeyBronze:
		k.Bronze--
	case KeySilver:
		k.Silver--
	case KeyGold:
		k.Gold--
	}
	return true
}

type Equipment struct {
	Weapon *Item
	Armor  *Item
	Charm  *Item
}

func (e Equipment) Slot(slot EquipmentSlot) *Item {
	switch slot {
	case SlotWeapon:
		return e.Weapon
	case SlotArmor:
		return e.Armor
	case SlotCharm:
		return e.Charm
	default:
		return nil
	}
}

type Player struct {
	Pos         Position
	HP          int
	BaseMaxHP   int
	BaseAttack  int
	BaseDefense int
	XP          int
	Level       int
	Gold        int
	Kills       int
	Inventory   []Item
	Equipment   Equipment
	Statuses    []StatusEffect
	Keys        KeyRing
	HasRelic    bool
}

func (p *Player) MaxHP() int {
	total := p.BaseMaxHP
	for _, item := range p.equippedItems() {
		total += item.MaxHPBonus
	}
	return total
}

func (p *Player) AttackPower() int {
	total := p.BaseAttack
	for _, item := range p.equippedItems() {
		total += item.AttackBonus
	}
	if status, ok := p.Status(StatusFocus); ok {
		total += status.Potency
	}
	return total
}

func (p *Player) DefensePower() int {
	total := p.BaseDefense
	for _, item := range p.equippedItems() {
		total += item.DefenseBonus
	}
	return total
}

func (p *Player) VisionRadius() int {
	total := 8
	for _, item := range p.equippedItems() {
		total += item.SightBonus
	}
	return total
}

func (p *Player) Status(kind StatusKind) (StatusEffect, bool) {
	for _, status := range p.Statuses {
		if status.Kind == kind {
			return status, true
		}
	}
	return StatusEffect{}, false
}

func (p *Player) HasStatus(kind StatusKind) bool {
	_, ok := p.Status(kind)
	return ok
}

func (p *Player) ApplyStatus(kind StatusKind, turns int, potency int) {
	for index := range p.Statuses {
		if p.Statuses[index].Kind == kind {
			p.Statuses[index].Turns = max(p.Statuses[index].Turns, turns)
			p.Statuses[index].Potency = max(p.Statuses[index].Potency, potency)
			return
		}
	}

	p.Statuses = append(p.Statuses, StatusEffect{Kind: kind, Turns: turns, Potency: potency})
}

func (p *Player) StatusResistance(kind StatusKind) int {
	resistance := 0
	for _, item := range p.equippedItems() {
		switch kind {
		case StatusPoison:
			resistance += item.PoisonResist
		case StatusFire:
			resistance += item.FireResist
		}
	}
	return resistance
}

func (p *Player) RemoveStatus(kind StatusKind) bool {
	for index := range p.Statuses {
		if p.Statuses[index].Kind == kind {
			p.Statuses = append(p.Statuses[:index], p.Statuses[index+1:]...)
			return true
		}
	}
	return false
}

func (p *Player) RemoveStatuses(kinds ...StatusKind) []StatusKind {
	removed := make([]StatusKind, 0, len(kinds))
	for _, kind := range kinds {
		if p.RemoveStatus(kind) {
			removed = append(removed, kind)
		}
	}
	return removed
}

func (p *Player) ClampHP() {
	p.HP = clamp(p.HP, 0, p.MaxHP())
}

func (p *Player) GainXP(amount int) bool {
	p.XP += amount
	leveled := false
	for p.XP >= p.NextLevelXP() {
		p.Level++
		p.BaseMaxHP += 5
		p.BaseAttack += 2
		if p.Level%2 == 0 {
			p.BaseDefense++
		}
		if p.Level%3 == 0 {
			p.BaseDefense++
		}
		p.HP = min(p.MaxHP(), p.HP+6)
		leveled = true
	}
	return leveled
}

func (p *Player) NextLevelXP() int {
	return 16 + p.Level*16 + (p.Level-1)*(p.Level-1)
}

func (p *Player) equippedItems() []Item {
	items := make([]Item, 0, 3)
	if p.Equipment.Weapon != nil {
		items = append(items, *p.Equipment.Weapon)
	}
	if p.Equipment.Armor != nil {
		items = append(items, *p.Equipment.Armor)
	}
	if p.Equipment.Charm != nil {
		items = append(items, *p.Equipment.Charm)
	}
	return items
}

type Enemy struct {
	ID              int
	Template        EnemyTemplate
	Pos             Position
	Home            Position
	HomeRoom        int
	HP              int
	State           AIState
	LastKnownPlayer Position
	HasLastKnown    bool
	Memory          int
	Cooldown        int
	Elite           bool
	Enraged         bool
}

func (e *Enemy) AttackPower() int {
	total := e.Template.Attack
	if e.Elite {
		total++
	}
	if e.Enraged {
		total += e.Template.EnrageAttackBonus
	}
	return total
}

func (e *Enemy) DefensePower() int {
	total := e.Template.Defense
	if e.Elite {
		total++
	}
	return total
}

func (e *Enemy) IsAlive() bool {
	return e.HP > 0
}

func (e *Enemy) DisplayName() string {
	if e.Elite && e.Template.BossTier == 0 {
		return "Elite " + e.Template.Name
	}
	return e.Template.Name
}
