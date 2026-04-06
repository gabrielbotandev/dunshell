package game

const (
	floorWidth  = 96
	floorHeight = 38
)

func GenerateFloor(rng *RNG, level int, maxFloors int, persistentDifficulty int, modifier FloorModifier, endless bool, nextEnemyID *int, nextChestID *int, nextMerchantID *int) *Floor {
	bossFloor := isBossFloor(level, maxFloors, endless)
	for attempt := 0; attempt < 16; attempt++ {
		floor := NewFloor(level, floorWidth, floorHeight, modifier)
		targetRooms := 12 + rng.Intn(5)
		if bossFloor {
			targetRooms--
		}

		for roomAttempts := 0; roomAttempts < 320 && len(floor.Rooms) < targetRooms; roomAttempts++ {
			room := Room{
				X: 2 + rng.Intn(floor.Width-18),
				Y: 2 + rng.Intn(floor.Height-12),
				W: 7 + rng.Intn(8),
				H: 6 + rng.Intn(5),
			}

			if overlapsRoom(floor.Rooms, room) {
				continue
			}

			carveRoom(floor, room)
			if len(floor.Rooms) > 0 {
				connectRooms(floor, floor.Rooms[len(floor.Rooms)-1], room, rng)
				if len(floor.Rooms) > 2 && rng.Float64() < 0.28 {
					other := floor.Rooms[rng.Intn(len(floor.Rooms))]
					connectRooms(floor, other, room, rng)
				}
			}
			floor.Rooms = append(floor.Rooms, room)
			floor.RoomKinds = append(floor.RoomKinds, RoomNormal)
		}

		if len(floor.Rooms) < 8 {
			continue
		}

		floor.Entrance = floor.Rooms[0].Center()
		if modifier.Rest && len(floor.RoomKinds) > 0 {
			floor.RoomKinds[0] = RoomSanctuary
		}

		if bossFloor && !addBossRoom(rng, floor) {
			continue
		}

		stairsRoom := chooseStairsRoomIndex(floor, bossFloor)
		if level < maxFloors || endless {
			floor.Stairs = floor.Rooms[stairsRoom].Center()
			floor.SetTile(floor.Stairs, TileStairsDown)
		}

		floor.BindRoomDoors()
		if bossFloor && !configureBossGate(floor) {
			continue
		}

		populateFloor(rng, floor, level, maxFloors, persistentDifficulty, modifier, endless, nextEnemyID, nextChestID, nextMerchantID)
		return floor
	}

	fallback := NewFloor(level, floorWidth, floorHeight, modifier)
	mainRoom := Room{X: 3, Y: 4, W: floorWidth - 6, H: floorHeight - 8}
	if bossFloor {
		mainRoom = Room{X: 3, Y: 4, W: floorWidth - 34, H: floorHeight - 8}
	}
	carveRoom(fallback, mainRoom)
	fallback.Rooms = append(fallback.Rooms, mainRoom)
	fallback.RoomKinds = append(fallback.RoomKinds, RoomNormal)
	fallback.Entrance = mainRoom.Center()
	if bossFloor {
		bossRoom := Room{X: floorWidth - 27, Y: floorHeight/2 - 6, W: 21, H: 12}
		carveRoom(fallback, bossRoom)
		connectRooms(fallback, mainRoom, bossRoom, rng)
		fallback.Rooms = append(fallback.Rooms, bossRoom)
		fallback.RoomKinds = append(fallback.RoomKinds, RoomBoss)
	}
	if level < maxFloors || endless {
		fallback.Stairs = mainRoom.Center().Offset(10, 0)
		fallback.SetTile(fallback.Stairs, TileStairsDown)
	}
	fallback.BindRoomDoors()
	if bossFloor {
		configureBossGate(fallback)
	}
	populateFloor(rng, fallback, level, maxFloors, persistentDifficulty, modifier, endless, nextEnemyID, nextChestID, nextMerchantID)
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

func connectRooms(floor *Floor, a Room, b Room, rng *RNG) {
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

func addBossRoom(rng *RNG, floor *Floor) bool {
	anchorIndex := 0
	anchorDistance := -1
	for index, room := range floor.Rooms {
		d := distance(floor.Rooms[0].Center(), room.Center())
		if d > anchorDistance {
			anchorDistance = d
			anchorIndex = index
		}
	}
	anchor := floor.Rooms[anchorIndex]
	anchorCenter := anchor.Center()

	for attempt := 0; attempt < 80; attempt++ {
		room := Room{W: 12 + rng.Intn(5), H: 8 + rng.Intn(4)}
		switch rng.Intn(4) {
		case 0:
			room.X = anchor.X + anchor.W + 7 + rng.Intn(5)
			room.Y = anchorCenter.Y - room.H/2 + rng.Intn(3) - 1
		case 1:
			room.X = anchor.X - room.W - 7 - rng.Intn(5)
			room.Y = anchorCenter.Y - room.H/2 + rng.Intn(3) - 1
		case 2:
			room.X = anchorCenter.X - room.W/2 + rng.Intn(3) - 1
			room.Y = anchor.Y + anchor.H + 6 + rng.Intn(5)
		default:
			room.X = anchorCenter.X - room.W/2 + rng.Intn(3) - 1
			room.Y = anchor.Y - room.H - 6 - rng.Intn(5)
		}
		if room.X < 2 || room.Y < 2 || room.X+room.W >= floor.Width-2 || room.Y+room.H >= floor.Height-2 {
			continue
		}
		if overlapsRoom(floor.Rooms, room) {
			continue
		}
		carveRoom(floor, room)
		connectRooms(floor, anchor, room, rng)
		floor.Rooms = append(floor.Rooms, room)
		floor.RoomKinds = append(floor.RoomKinds, RoomBoss)
		return true
	}
	return false
}

func chooseStairsRoomIndex(floor *Floor, bossFloor bool) int {
	bestIndex := len(floor.Rooms) - 1
	bestDistance := -1
	for index, room := range floor.Rooms {
		if bossFloor && index < len(floor.RoomKinds) && floor.RoomKinds[index] == RoomBoss {
			continue
		}
		d := distance(floor.Entrance, room.Center())
		if d > bestDistance {
			bestDistance = d
			bestIndex = index
		}
	}
	return bestIndex
}

func configureBossGate(floor *Floor) bool {
	bossRoomIndex := -1
	for index, kind := range floor.RoomKinds {
		if kind == RoomBoss {
			bossRoomIndex = index
			break
		}
	}
	if bossRoomIndex < 0 {
		return false
	}
	doors := append([]Position(nil), floor.RoomDoors[bossRoomIndex]...)
	if len(doors) == 0 {
		return false
	}
	for index := 1; index < len(doors); index++ {
		floor.SetTile(doors[index], TileWall)
	}
	gate := doors[0]
	entry := gate
	bossRoom := floor.Rooms[bossRoomIndex]
	for _, dir := range cardinalDirections {
		next := gate.Add(dir)
		if bossRoom.Contains(next) && floor.IsWalkable(next) {
			entry = next
			break
		}
	}
	floor.SetTile(gate, TileBossGate)
	floor.Boss = &BossEncounter{RoomIndex: bossRoomIndex, Gate: gate, Entry: entry, RewardChestID: -1}
	floor.BindRoomDoors()
	return true
}

func populateFloor(rng *RNG, floor *Floor, level int, maxFloors int, persistentDifficulty int, modifier FloorModifier, endless bool, nextEnemyID *int, nextChestID *int, nextMerchantID *int) {
	occupied := map[Position]bool{floor.Entrance: true}
	if floor.Stairs.X >= 0 {
		occupied[floor.Stairs] = true
	}

	bossRoomIndex := -1
	if floor.Boss != nil {
		bossRoomIndex = floor.Boss.RoomIndex
		occupied[floor.Boss.Gate] = true
		occupied[floor.Boss.Entry] = true
	}

	populateLooseItems(rng, floor, occupied, modifier, bossRoomIndex)
	populateChests(rng, floor, occupied, level, modifier, nextChestID, bossRoomIndex)
	populateMerchants(rng, floor, occupied, level, modifier, nextMerchantID, bossRoomIndex)
	populateKeys(rng, floor, occupied, level, bossRoomIndex)
	populateEnemies(rng, floor, occupied, level, persistentDifficulty, modifier, nextEnemyID, bossRoomIndex)
	populateBoss(rng, floor, occupied, level, maxFloors, persistentDifficulty, modifier, endless, nextEnemyID, nextChestID)
}

func populateLooseItems(rng *RNG, floor *Floor, occupied map[Position]bool, modifier FloorModifier, bossRoomIndex int) {
	itemCount := 5 + floor.Level/3 + rng.Intn(3)
	if modifier.Rest {
		itemCount++
	}
	if modifier.Cursed {
		itemCount++
	}
	for count := 0; count < itemCount; count++ {
		pos, roomIndex := randomPlacableTile(rng, floor, occupied, bossRoomIndex)
		if pos.X < 0 {
			break
		}
		occupied[pos] = true
		floor.Items = append(floor.Items, GroundItem{
			Pos:       pos,
			Item:      RandomGroundItem(rng, floor.Level, modifier),
			RoomIndex: roomIndex,
		})
	}
}

func populateChests(rng *RNG, floor *Floor, occupied map[Position]bool, level int, modifier FloorModifier, nextChestID *int, bossRoomIndex int) {
	chestCount := 0
	if rng.Float64() < 0.78 {
		chestCount = 1
		if rng.Float64() < 0.34 {
			chestCount++
		}
	}
	chestCount += modifier.ExtraChests

	for count := 0; count < chestCount; count++ {
		pos, roomIndex := randomPlacableTile(rng, floor, occupied, bossRoomIndex)
		if pos.X < 0 {
			break
		}
		tier := rollChestTier(rng, level, modifier)
		floor.Chests = append(floor.Chests, Chest{
			ID:        *nextChestID,
			Pos:       pos,
			Tier:      tier,
			RoomIndex: roomIndex,
			Rewards:   GenerateChestRewards(rng, tier, level, modifier, false, false),
		})
		*nextChestID++
		occupied[pos] = true
	}
}

func populateMerchants(rng *RNG, floor *Floor, occupied map[Position]bool, level int, modifier FloorModifier, nextMerchantID *int, bossRoomIndex int) {
	spawnMerchant := modifier.Merchant || (level > 3 && rng.Float64() < 0.04)
	if !spawnMerchant {
		return
	}
	pos, roomIndex := randomPlacableTile(rng, floor, occupied, bossRoomIndex)
	if pos.X < 0 {
		return
	}
	if roomIndex >= 0 && roomIndex < len(floor.RoomKinds) {
		floor.RoomKinds[roomIndex] = RoomMerchant
	}
	merchant := Merchant{
		ID:        *nextMerchantID,
		Name:      merchantNameForFloor(level),
		Pos:       pos,
		RoomIndex: roomIndex,
		Offers:    GenerateMerchantOffers(rng, level),
	}
	*nextMerchantID++
	floor.Merchants = append(floor.Merchants, merchant)
	occupied[pos] = true
}

func populateKeys(rng *RNG, floor *Floor, occupied map[Position]bool, level int, bossRoomIndex int) {
	keyCount := 0
	if rng.Float64() < 0.88 {
		keyCount = 1
		if level >= 8 && rng.Float64() < 0.28 {
			keyCount++
		}
	}
	for count := 0; count < keyCount; count++ {
		pos, roomIndex := randomPlacableTile(rng, floor, occupied, bossRoomIndex)
		if pos.X < 0 {
			break
		}
		occupied[pos] = true
		floor.Items = append(floor.Items, GroundItem{
			Pos:       pos,
			Item:      RandomKeyReward(rng, level),
			RoomIndex: roomIndex,
		})
	}
}

func populateEnemies(rng *RNG, floor *Floor, occupied map[Position]bool, level int, persistentDifficulty int, modifier FloorModifier, nextEnemyID *int, bossRoomIndex int) {
	enemyCount := 9 + level + rng.Intn(4) + modifier.EnemyBonus - earlyFloorEnemyCountRelief(level)
	if modifier.Cursed {
		enemyCount += 2
	}
	for count := 0; count < enemyCount; count++ {
		pos, roomIndex := randomPlacableTile(rng, floor, occupied, bossRoomIndex)
		if pos.X < 0 {
			break
		}
		if distance(pos, floor.Entrance) < 7 {
			continue
		}
		elite := modifier.EliteChance > 0 && rng.Float64() < modifier.EliteChance
		template := ScaleEnemyTemplate(RandomEnemyTemplate(rng, level), level, persistentDifficulty, elite, modifier.Cursed)
		enemy := &Enemy{
			ID:       *nextEnemyID,
			Template: template,
			Level:    enemyLevelForEncounter(level, template, elite),
			Pos:      pos,
			Home:     pos,
			HomeRoom: roomIndex,
			HP:       template.MaxHP,
			State:    AIStateWander,
			Elite:    elite,
		}
		*nextEnemyID++
		floor.Enemies = append(floor.Enemies, enemy)
		occupied[pos] = true
	}
}

func earlyFloorEnemyCountRelief(level int) int {
	switch level {
	case 1:
		return 2
	case 2:
		return 1
	default:
		return 0
	}
}

func populateBoss(rng *RNG, floor *Floor, occupied map[Position]bool, level int, maxFloors int, persistentDifficulty int, modifier FloorModifier, endless bool, nextEnemyID *int, nextChestID *int) {
	if floor.Boss == nil {
		return
	}
	bossRoom := floor.Rooms[floor.Boss.RoomIndex]
	bossPos := bossRoom.Center()
	if occupied[bossPos] {
		bossPos = bossRoom.Center().Offset(-1, 0)
	}
	template := ScaleEnemyTemplate(BossTemplateForFloor(rng, level, maxFloors, endless), level, persistentDifficulty, false, modifier.Cursed)
	boss := &Enemy{
		ID:       *nextEnemyID,
		Template: template,
		Level:    enemyLevelForEncounter(level, template, false),
		Pos:      bossPos,
		Home:     bossPos,
		HomeRoom: floor.Boss.RoomIndex,
		HP:       template.MaxHP,
		State:    AIStateWander,
	}
	*nextEnemyID++
	floor.Enemies = append(floor.Enemies, boss)
	floor.Boss.BossID = boss.ID
	occupied[bossPos] = true

	chestTier := KeySilver
	if level >= 10 {
		chestTier = KeyGold
	}
	if level >= maxFloors {
		chestTier = KeyGold
	}
	chestPos := Position{X: bossRoom.X + bossRoom.W - 3, Y: bossRoom.Y + bossRoom.H - 3}
	if occupied[chestPos] || !floor.IsWalkable(chestPos) {
		chestPos = Position{X: bossRoom.X + 2, Y: bossRoom.Y + bossRoom.H - 3}
	}
	chest := Chest{
		ID:         *nextChestID,
		Pos:        chestPos,
		Tier:       chestTier,
		Locked:     true,
		BossReward: true,
		RoomIndex:  floor.Boss.RoomIndex,
		Rewards:    GenerateChestRewards(rng, chestTier, level, modifier, true, level == maxFloors),
	}
	*nextChestID++
	floor.Chests = append(floor.Chests, chest)
	floor.Boss.RewardChestID = chest.ID
	occupied[chestPos] = true
}

func randomPlacableTile(rng *RNG, floor *Floor, occupied map[Position]bool, excludedRoom int) (Position, int) {
	for attempt := 0; attempt < 240; attempt++ {
		roomIndex := rng.Intn(len(floor.Rooms))
		if roomIndex == excludedRoom {
			continue
		}
		room := floor.Rooms[roomIndex]
		pos := Position{
			X: room.X + 1 + rng.Intn(max(1, room.W-2)),
			Y: room.Y + 1 + rng.Intn(max(1, room.H-2)),
		}
		if occupied[pos] || !floor.IsWalkable(pos) {
			continue
		}
		return pos, roomIndex
	}
	return Position{X: -1, Y: -1}, -1
}

func rollChestTier(rng *RNG, level int, modifier FloorModifier) KeyTier {
	roll := rng.Intn(100)
	if modifier.Cursed {
		roll += 10
	}
	if modifier.LootBonus > 0 {
		roll += modifier.LootBonus * 6
	}
	switch {
	case level >= 12 && roll > 78:
		return KeyGold
	case level >= 5 && roll > 38:
		return KeySilver
	default:
		return KeyBronze
	}
}

func merchantNameForFloor(level int) string {
	switch {
	case level >= 15:
		return "Ash Broker"
	case level >= 8:
		return "Relic Factor"
	default:
		return "Pilgrim Merchant"
	}
}
