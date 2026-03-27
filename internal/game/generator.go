package game

import "math/rand"

const (
	floorWidth  = 72
	floorHeight = 30
)

func GenerateFloor(rng *rand.Rand, level int, maxFloors int, nextEnemyID *int) *Floor {
	for attempt := 0; attempt < 8; attempt++ {
		floor := NewFloor(level, floorWidth, floorHeight)
		targetRooms := 10 + rng.Intn(4)

		for roomAttempts := 0; roomAttempts < 220 && len(floor.Rooms) < targetRooms; roomAttempts++ {
			room := Room{
				X: 2 + rng.Intn(floor.Width-14),
				Y: 2 + rng.Intn(floor.Height-10),
				W: 6 + rng.Intn(7),
				H: 5 + rng.Intn(5),
			}

			if overlapsRoom(floor.Rooms, room) {
				continue
			}

			carveRoom(floor, room)
			if len(floor.Rooms) > 0 {
				connectRooms(floor, floor.Rooms[len(floor.Rooms)-1], room, rng)
				if rng.Float64() < 0.35 {
					other := floor.Rooms[rng.Intn(len(floor.Rooms))]
					connectRooms(floor, other, room, rng)
				}
			}
			floor.Rooms = append(floor.Rooms, room)
		}

		if len(floor.Rooms) < 7 {
			continue
		}

		floor.Entrance = floor.Rooms[0].Center()
		if level < maxFloors {
			floor.Stairs = floor.Rooms[len(floor.Rooms)-1].Center()
			floor.SetTile(floor.Stairs, TileStairsDown)
		}
		floor.BindRoomDoors()

		populateFloor(rng, floor, maxFloors, nextEnemyID)
		return floor
	}

	fallback := NewFloor(level, floorWidth, floorHeight)
	room := Room{X: 4, Y: 4, W: floorWidth - 8, H: floorHeight - 8}
	carveRoom(fallback, room)
	fallback.Rooms = append(fallback.Rooms, room)
	fallback.Entrance = room.Center()
	if level < maxFloors {
		fallback.Stairs = room.Center().Offset(8, 0)
		fallback.SetTile(fallback.Stairs, TileStairsDown)
	}
	fallback.BindRoomDoors()
	populateFloor(rng, fallback, maxFloors, nextEnemyID)
	return fallback
}

func overlapsRoom(existing []Room, room Room) bool {
	for _, other := range existing {
		padded := Room{X: other.X - 1, Y: other.Y - 1, W: other.W + 2, H: other.H + 2}
		if padded.Intersects(room) {
			return true
		}
	}
	return false
}

func carveRoom(floor *Floor, room Room) {
	for y := room.Y; y < room.Y+room.H; y++ {
		for x := room.X; x < room.X+room.W; x++ {
			floor.Tiles[y][x] = TileFloor
		}
	}
}

func connectRooms(floor *Floor, a Room, b Room, rng *rand.Rand) {
	path := tunnelPath(a.Center(), b.Center(), rng.Intn(2) == 0)
	for _, pos := range path {
		floor.SetTile(pos, TileFloor)
	}

	for index := 1; index < len(path); index++ {
		if a.Contains(path[index-1]) && !a.Contains(path[index]) {
			if floor.TileAt(path[index]) == TileFloor {
				floor.SetTile(path[index], TileDoorClosed)
			}
		}
		if !b.Contains(path[index-1]) && b.Contains(path[index]) {
			if floor.TileAt(path[index-1]) == TileFloor {
				floor.SetTile(path[index-1], TileDoorClosed)
			}
		}
	}
}

func tunnelPath(start Position, end Position, horizontalFirst bool) []Position {
	path := make([]Position, 0, abs(end.X-start.X)+abs(end.Y-start.Y)+1)
	path = append(path, start)

	current := start
	if horizontalFirst {
		for current.X != end.X {
			if current.X < end.X {
				current = current.Offset(1, 0)
			} else {
				current = current.Offset(-1, 0)
			}
			path = append(path, current)
		}
		for current.Y != end.Y {
			if current.Y < end.Y {
				current = current.Offset(0, 1)
			} else {
				current = current.Offset(0, -1)
			}
			path = append(path, current)
		}
		return path
	}

	for current.Y != end.Y {
		if current.Y < end.Y {
			current = current.Offset(0, 1)
		} else {
			current = current.Offset(0, -1)
		}
		path = append(path, current)
	}
	for current.X != end.X {
		if current.X < end.X {
			current = current.Offset(1, 0)
		} else {
			current = current.Offset(-1, 0)
		}
		path = append(path, current)
	}

	return path
}

func populateFloor(rng *rand.Rand, floor *Floor, maxFloors int, nextEnemyID *int) {
	occupied := map[Position]bool{
		floor.Entrance: true,
	}
	if floor.Stairs.X >= 0 {
		occupied[floor.Stairs] = true
	}

	itemCount := 4 + rng.Intn(2) + floor.Level
	for count := 0; count < itemCount; count++ {
		pos := randomPlacableTile(rng, floor, occupied)
		if pos.X < 0 {
			break
		}
		occupied[pos] = true
		floor.Items = append(floor.Items, GroundItem{
			Pos:       pos,
			Item:      RandomGroundItem(rng, floor.Level),
			RoomIndex: floor.RoomIndexAt(pos),
		})
	}

	enemyCount := 7 + floor.Level*2 + rng.Intn(3)
	for count := 0; count < enemyCount; count++ {
		pos := randomPlacableTile(rng, floor, occupied)
		if pos.X < 0 {
			break
		}
		if distance(pos, floor.Entrance) < 8 {
			continue
		}

		template := RandomEnemyTemplate(rng, floor.Level)
		enemy := &Enemy{
			ID:       *nextEnemyID,
			Template: template,
			Pos:      pos,
			Home:     pos,
			HomeRoom: floor.RoomIndexAt(pos),
			HP:       template.MaxHP,
			State:    AIStateWander,
		}
		*nextEnemyID++
		floor.Enemies = append(floor.Enemies, enemy)
		occupied[pos] = true
	}

	if floor.Level == maxFloors {
		bossRoom := floor.Rooms[len(floor.Rooms)-1]
		relicPos := bossRoom.Center()
		floor.Items = append(floor.Items, GroundItem{
			Pos:       relicPos,
			Item:      ItemByID("cinder_crown"),
			RoomIndex: floor.RoomIndexAt(relicPos),
		})
		occupied[relicPos] = true

		bossPos := relicPos.Offset(-2, 0)
		if !floor.IsWalkable(bossPos) || occupied[bossPos] {
			bossPos = relicPos.Offset(0, -2)
		}
		template := BossTemplate()
		boss := &Enemy{
			ID:       *nextEnemyID,
			Template: template,
			Pos:      bossPos,
			Home:     bossPos,
			HomeRoom: floor.RoomIndexAt(bossPos),
			HP:       template.MaxHP,
			State:    AIStateWander,
		}
		*nextEnemyID++
		floor.Enemies = append(floor.Enemies, boss)
	}
}

func randomPlacableTile(rng *rand.Rand, floor *Floor, occupied map[Position]bool) Position {
	for attempt := 0; attempt < 200; attempt++ {
		room := floor.Rooms[rng.Intn(len(floor.Rooms))]
		pos := Position{
			X: room.X + 1 + rng.Intn(max(1, room.W-2)),
			Y: room.Y + 1 + rng.Intn(max(1, room.H-2)),
		}
		if occupied[pos] || !floor.IsWalkable(pos) {
			continue
		}
		return pos
	}
	return Position{X: -1, Y: -1}
}
