package game

var cardinalDirections = []Position{
	{X: 0, Y: -1},
	{X: 1, Y: 0},
	{X: 0, Y: 1},
	{X: -1, Y: 0},
}

func NextStepToward(floor *Floor, start Position, goal Position, canOpenDoors bool, blocked map[Position]bool) (Position, bool) {
	if start.Equals(goal) {
		return start, false
	}

	type node struct {
		Pos Position
	}

	queue := []Position{start}
	visited := map[Position]bool{start: true}
	previous := map[Position]Position{}
	found := false

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current.Equals(goal) {
			found = true
			break
		}

		for _, dir := range cardinalDirections {
			next := current.Add(dir)
			if visited[next] || !floor.InBounds(next) {
				continue
			}
			if blocked[next] && !next.Equals(goal) {
				continue
			}
			if !floor.IsWalkableFor(next, canOpenDoors) && !next.Equals(goal) {
				continue
			}
			visited[next] = true
			previous[next] = current
			queue = append(queue, next)
		}
	}

	if !found {
		return start, false
	}

	step := goal
	for {
		parent, ok := previous[step]
		if !ok {
			return start, false
		}
		if parent.Equals(start) {
			return step, true
		}
		step = parent
	}
}

func StepAway(floor *Floor, start Position, from Position, canOpenDoors bool, blocked map[Position]bool) (Position, bool) {
	best := start
	bestScore := distance(start, from)
	for _, dir := range cardinalDirections {
		next := start.Add(dir)
		if !floor.InBounds(next) || blocked[next] {
			continue
		}
		if !floor.IsWalkableFor(next, canOpenDoors) {
			continue
		}
		score := distance(next, from)
		if score > bestScore {
			best = next
			bestScore = score
		}
	}

	if best.Equals(start) {
		return start, false
	}
	return best, true
}
