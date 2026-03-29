package game

type Position struct {
	X int
	Y int
}

func (p Position) Add(other Position) Position {
	return Position{X: p.X + other.X, Y: p.Y + other.Y}
}

func (p Position) Offset(dx int, dy int) Position {
	return Position{X: p.X + dx, Y: p.Y + dy}
}

func (p Position) Equals(other Position) bool {
	return p.X == other.X && p.Y == other.Y
}

type Room struct {
	X int
	Y int
	W int
	H int
}

func (r Room) Center() Position {
	return Position{X: r.X + r.W/2, Y: r.Y + r.H/2}
}

func (r Room) Contains(pos Position) bool {
	return pos.X >= r.X && pos.X < r.X+r.W && pos.Y >= r.Y && pos.Y < r.Y+r.H
}

func (r Room) Intersects(other Room) bool {
	return r.X <= other.X+other.W && r.X+r.W >= other.X &&
		r.Y <= other.Y+other.H && r.Y+r.H >= other.Y
}

type RoomKind int

const (
	RoomNormal RoomKind = iota
	RoomMerchant
	RoomBoss
	RoomSanctuary
)

func (k RoomKind) Label() string {
	switch k {
	case RoomMerchant:
		return "merchant"
	case RoomBoss:
		return "boss room"
	case RoomSanctuary:
		return "sanctuary"
	default:
		return "chamber"
	}
}

type TileType int

const (
	TileWall TileType = iota
	TileFloor
	TileDoorClosed
	TileDoorOpen
	TileStairsDown
	TileBossGate
	TileBossSeal
)

func (t TileType) Glyph() rune {
	switch t {
	case TileWall:
		return '#'
	case TileFloor:
		return '.'
	case TileDoorClosed:
		return '+'
	case TileDoorOpen:
		return '/'
	case TileStairsDown:
		return '>'
	case TileBossGate:
		return 'X'
	case TileBossSeal:
		return '#'
	default:
		return ' '
	}
}

func (t TileType) Name() string {
	switch t {
	case TileWall:
		return "wall"
	case TileFloor:
		return "flagstone"
	case TileDoorClosed:
		return "sealed door"
	case TileDoorOpen:
		return "open door"
	case TileStairsDown:
		return "stair"
	case TileBossGate:
		return "blood-locked gate"
	case TileBossSeal:
		return "sealed gate"
	default:
		return "void"
	}
}

func (t TileType) Walkable() bool {
	return t == TileFloor || t == TileDoorOpen || t == TileStairsDown
}

func (t TileType) Transparent() bool {
	return t == TileFloor || t == TileDoorOpen || t == TileStairsDown
}

type GameMode int

const (
	ModePlaying GameMode = iota
	ModeWon
	ModeLost
)

type AIState int

const (
	AIStateWander AIState = iota
	AIStateChase
	AIStateAttack
)

type Behavior int

const (
	BehaviorSkittish Behavior = iota
	BehaviorBrute
	BehaviorProwler
	BehaviorHunter
	BehaviorCutpurse
	BehaviorSentinel
	BehaviorCaster
	BehaviorBoss
)

type Rarity int

const (
	RarityCommon Rarity = iota
	RarityUncommon
	RarityRare
	RarityLegendary
	RarityUnique
)

func (r Rarity) Label() string {
	switch r {
	case RarityCommon:
		return "Common"
	case RarityUncommon:
		return "Uncommon"
	case RarityRare:
		return "Rare"
	case RarityLegendary:
		return "Legendary"
	case RarityUnique:
		return "Unique"
	default:
		return "Unknown"
	}
}

func (r Rarity) Tint() string {
	switch r {
	case RarityCommon:
		return "#c9c1b4"
	case RarityUncommon:
		return "#84be7f"
	case RarityRare:
		return "#69a8d6"
	case RarityLegendary:
		return "#e6b15a"
	case RarityUnique:
		return "#d8708c"
	default:
		return "#c9c1b4"
	}
}

type KeyTier int

const (
	KeyBronze KeyTier = iota
	KeySilver
	KeyGold
)

func (k KeyTier) Label() string {
	switch k {
	case KeyBronze:
		return "Bronze"
	case KeySilver:
		return "Silver"
	case KeyGold:
		return "Gold"
	default:
		return "Unknown"
	}
}

func (k KeyTier) LowerLabel() string {
	switch k {
	case KeyBronze:
		return "bronze"
	case KeySilver:
		return "silver"
	case KeyGold:
		return "gold"
	default:
		return "unknown"
	}
}

func (k KeyTier) Tint() string {
	switch k {
	case KeyBronze:
		return "#b6844e"
	case KeySilver:
		return "#aebccc"
	case KeyGold:
		return "#dfbb58"
	default:
		return "#c9c1b4"
	}
}

type RewardKind int

const (
	RewardGold RewardKind = iota
	RewardItem
)

type FloorModifier struct {
	ID            string
	Title         string
	Subtitle      string
	Summary       string
	Merchant      bool
	Rest          bool
	HealOnStart   int
	CleanseOnRest bool
	BonusGold     float64
	LootBonus     int
	EnemyBonus    int
	EliteChance   float64
	ExtraChests   int
	GuaranteedKey *KeyTier
	Cursed        bool
}

func (m FloorModifier) Label() string {
	if m.Title == "" {
		return "Unmarked Descent"
	}
	return m.Title
}

func (m FloorModifier) HasEffect() bool {
	return m.ID != ""
}

type RouteChoice struct {
	ID        string
	Title     string
	Subtitle  string
	Reward    string
	Risk      string
	MapLabel  string
	Modifier  FloorModifier
	BossFloor bool
}

type ChestReward struct {
	Kind RewardKind
	Gold int
	Item Item
}

func (r ChestReward) Summary() string {
	if r.Kind == RewardGold {
		return itoa(r.Gold) + " gold"
	}
	return r.Item.Name
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func clamp(value int, low int, high int) int {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}

func distance(a Position, b Position) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}
