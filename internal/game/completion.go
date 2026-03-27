package game

type RoomState struct {
	Index            int
	Room             Room
	ExploredTiles    int
	TotalTiles       int
	Explored         bool
	OpenDoors        int
	TotalDoors       int
	DoorsOpen        bool
	RemainingEnemies int
	EnemiesCleared   bool
	RemainingItems   int
	LootCollected    bool
	Cleared          bool
}

type FloorCompletion struct {
	UnexploredTiles  int
	TotalTiles       int
	RemainingItems   int
	RemainingEnemies int
	ClearedRooms     int
	TotalRooms       int
}

func (c FloorCompletion) FullyExplored() bool {
	return c.UnexploredTiles == 0
}

func (c FloorCompletion) LootCollected() bool {
	return c.RemainingItems == 0
}

func (c FloorCompletion) EnemiesCleared() bool {
	return c.RemainingEnemies == 0
}

func (c FloorCompletion) Complete() bool {
	return c.FullyExplored() && c.LootCollected() && c.EnemiesCleared()
}

func (c FloorCompletion) UnclearedRooms() int {
	return max(0, c.TotalRooms-c.ClearedRooms)
}

func (f *Floor) RoomIndexAt(pos Position) int {
	for index, room := range f.Rooms {
		if room.Contains(pos) {
			return index
		}
	}
	return -1
}

func (f *Floor) AdjacentRoomIndices(pos Position) []int {
	indices := make([]int, 0, 2)
	seen := make(map[int]bool, 2)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			roomIndex := f.RoomIndexAt(pos.Offset(dx, dy))
			if roomIndex < 0 || seen[roomIndex] {
				continue
			}
			seen[roomIndex] = true
			indices = append(indices, roomIndex)
		}
	}
	return indices
}

func (f *Floor) BindRoomDoors() {
	f.RoomDoors = make([][]Position, len(f.Rooms))
	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			pos := Position{X: x, Y: y}
			tile := f.TileAt(pos)
			if tile != TileDoorClosed && tile != TileDoorOpen {
				continue
			}
			for _, roomIndex := range f.AdjacentRoomIndices(pos) {
				if containsPosition(f.RoomDoors[roomIndex], pos) {
					continue
				}
				f.RoomDoors[roomIndex] = append(f.RoomDoors[roomIndex], pos)
			}
		}
	}
}

func (f *Floor) RoomStates() []RoomState {
	states := make([]RoomState, len(f.Rooms))
	for index := range f.Rooms {
		states[index] = f.RoomState(index)
	}
	return states
}

func (f *Floor) RoomState(roomIndex int) RoomState {
	if roomIndex < 0 || roomIndex >= len(f.Rooms) {
		return RoomState{}
	}

	room := f.Rooms[roomIndex]
	state := RoomState{
		Index: roomIndex,
		Room:  room,
	}

	for y := room.Y; y < room.Y+room.H; y++ {
		for x := room.X; x < room.X+room.W; x++ {
			pos := Position{X: x, Y: y}
			if f.TileAt(pos) == TileWall {
				continue
			}
			state.TotalTiles++
			if f.IsExplored(pos) {
				state.ExploredTiles++
			}
		}
	}

	state.TotalDoors = len(f.RoomDoors[roomIndex])
	for _, pos := range f.RoomDoors[roomIndex] {
		if f.TileAt(pos) == TileDoorOpen {
			state.OpenDoors++
		}
	}

	for _, enemy := range f.Enemies {
		if enemy.IsAlive() && enemy.HomeRoom == roomIndex {
			state.RemainingEnemies++
		}
	}

	for _, item := range f.Items {
		if item.RoomIndex == roomIndex {
			state.RemainingItems++
		}
	}

	state.Explored = state.TotalTiles > 0 && state.ExploredTiles == state.TotalTiles
	state.DoorsOpen = state.OpenDoors == state.TotalDoors
	state.EnemiesCleared = state.RemainingEnemies == 0
	state.LootCollected = state.RemainingItems == 0
	state.Cleared = state.Explored && state.DoorsOpen && state.EnemiesCleared && state.LootCollected
	return state
}

func (f *Floor) Completion() FloorCompletion {
	completion := FloorCompletion{
		RemainingItems: len(f.Items),
		TotalRooms:     len(f.Rooms),
	}

	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			if f.Tiles[y][x] == TileWall {
				continue
			}
			completion.TotalTiles++
			if !f.Explored[y][x] {
				completion.UnexploredTiles++
			}
		}
	}

	for _, enemy := range f.Enemies {
		if enemy.IsAlive() {
			completion.RemainingEnemies++
		}
	}

	for _, state := range f.RoomStates() {
		if state.Cleared {
			completion.ClearedRooms++
		}
	}

	return completion
}

func containsPosition(positions []Position, target Position) bool {
	for _, pos := range positions {
		if pos.Equals(target) {
			return true
		}
	}
	return false
}
