package game

type StatusKind int

const (
	StatusPoison StatusKind = iota
	StatusFocus
)

type StatusEffect struct {
	Kind    StatusKind
	Turns   int
	Potency int
}

func (s StatusEffect) Label() string {
	switch s.Kind {
	case StatusPoison:
		return "Poison " + itoa(s.Turns)
	case StatusFocus:
		return "Focus " + itoa(s.Turns)
	default:
		return "Unknown"
	}
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
	total := 7
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

func (p *Player) RemoveStatus(kind StatusKind) bool {
	for index := range p.Statuses {
		if p.Statuses[index].Kind == kind {
			p.Statuses = append(p.Statuses[:index], p.Statuses[index+1:]...)
			return true
		}
	}
	return false
}

func (p *Player) ClampHP() {
	p.HP = clamp(p.HP, 0, p.MaxHP())
}

func (p *Player) GainXP(amount int) bool {
	p.XP += amount
	leveled := false
	for p.XP >= p.NextLevelXP() {
		p.Level++
		p.BaseMaxHP += 4
		p.BaseAttack++
		if p.Level%2 == 0 {
			p.BaseDefense++
		}
		p.HP = min(p.MaxHP(), p.HP+6)
		leveled = true
	}
	return leveled
}

func (p *Player) NextLevelXP() int {
	return p.Level * 18
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
}

func (e *Enemy) AttackPower() int {
	return e.Template.Attack
}

func (e *Enemy) DefensePower() int {
	return e.Template.Defense
}

func (e *Enemy) IsAlive() bool {
	return e.HP > 0
}
