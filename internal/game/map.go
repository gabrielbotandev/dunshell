package game

type Floor struct {
	Level    int
	Theme    string
	Width    int
	Height   int
	Tiles    [][]TileType
	Visible  [][]bool
	Explored [][]bool
	Rooms    []Room
	Entrance Position
	Stairs   Position
	Items    []GroundItem
	Enemies  []*Enemy
}

func NewFloor(level int, width int, height int) *Floor {
	tiles := make([][]TileType, height)
	visible := make([][]bool, height)
	explored := make([][]bool, height)

	for y := 0; y < height; y++ {
		tiles[y] = make([]TileType, width)
		visible[y] = make([]bool, width)
		explored[y] = make([]bool, width)
		for x := 0; x < width; x++ {
			tiles[y][x] = TileWall
		}
	}

	return &Floor{
		Level:    level,
		Theme:    FloorTheme(level),
		Width:    width,
		Height:   height,
		Tiles:    tiles,
		Visible:  visible,
		Explored: explored,
		Stairs:   Position{X: -1, Y: -1},
	}
}

func (f *Floor) InBounds(pos Position) bool {
	return pos.X >= 0 && pos.X < f.Width && pos.Y >= 0 && pos.Y < f.Height
}

func (f *Floor) TileAt(pos Position) TileType {
	if !f.InBounds(pos) {
		return TileWall
	}
	return f.Tiles[pos.Y][pos.X]
}

func (f *Floor) SetTile(pos Position, tile TileType) {
	if f.InBounds(pos) {
		f.Tiles[pos.Y][pos.X] = tile
	}
}

func (f *Floor) IsTransparent(pos Position) bool {
	return f.TileAt(pos).Transparent()
}

func (f *Floor) IsWalkable(pos Position) bool {
	return f.TileAt(pos).Walkable()
}

func (f *Floor) IsWalkableFor(pos Position, canOpenDoors bool) bool {
	tile := f.TileAt(pos)
	return tile.Walkable() || (canOpenDoors && tile == TileDoorClosed)
}

func (f *Floor) OpenDoor(pos Position) bool {
	if f.TileAt(pos) != TileDoorClosed {
		return false
	}
	f.SetTile(pos, TileDoorOpen)
	return true
}

func (f *Floor) ResetVisibility() {
	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			f.Visible[y][x] = false
		}
	}
}

func (f *Floor) MarkVisible(pos Position) {
	if !f.InBounds(pos) {
		return
	}
	f.Visible[pos.Y][pos.X] = true
	f.Explored[pos.Y][pos.X] = true
}

func (f *Floor) IsVisible(pos Position) bool {
	return f.InBounds(pos) && f.Visible[pos.Y][pos.X]
}

func (f *Floor) IsExplored(pos Position) bool {
	return f.InBounds(pos) && f.Explored[pos.Y][pos.X]
}

func (f *Floor) EnemyAt(pos Position) *Enemy {
	for _, enemy := range f.Enemies {
		if enemy.IsAlive() && enemy.Pos.Equals(pos) {
			return enemy
		}
	}
	return nil
}

func (f *Floor) ItemIndicesAt(pos Position) []int {
	indices := make([]int, 0, 2)
	for index, item := range f.Items {
		if item.Pos.Equals(pos) {
			indices = append(indices, index)
		}
	}
	return indices
}

func (f *Floor) TopItemAt(pos Position) (GroundItem, bool) {
	for _, item := range f.Items {
		if item.Pos.Equals(pos) {
			return item, true
		}
	}
	return GroundItem{}, false
}

func (f *Floor) RemoveItemAt(index int) Item {
	item := f.Items[index].Item
	f.Items = append(f.Items[:index], f.Items[index+1:]...)
	return item
}

func (f *Floor) RemoveEnemyByID(id int) *Enemy {
	for index, enemy := range f.Enemies {
		if enemy.ID == id {
			f.Enemies = append(f.Enemies[:index], f.Enemies[index+1:]...)
			return enemy
		}
	}
	return nil
}

func (f *Floor) ExploredPercent() int {
	explored := 0
	walkable := 0
	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			if f.Tiles[y][x] != TileWall {
				walkable++
				if f.Explored[y][x] {
					explored++
				}
			}
		}
	}
	if walkable == 0 {
		return 0
	}
	return explored * 100 / walkable
}
