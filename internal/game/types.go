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

type TileType int

const (
	TileWall TileType = iota
	TileFloor
	TileDoorClosed
	TileDoorOpen
	TileStairsDown
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
	BehaviorBoss
)

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
