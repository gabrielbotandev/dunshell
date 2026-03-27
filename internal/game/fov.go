package game

func ComputeFOV(floor *Floor, origin Position, radius int) {
	floor.ResetVisibility()
	for y := max(0, origin.Y-radius); y <= min(floor.Height-1, origin.Y+radius); y++ {
		for x := max(0, origin.X-radius); x <= min(floor.Width-1, origin.X+radius); x++ {
			pos := Position{X: x, Y: y}
			dx := pos.X - origin.X
			dy := pos.Y - origin.Y
			if dx*dx+dy*dy > radius*radius {
				continue
			}
			if hasLineOfSight(floor, origin, pos) {
				floor.MarkVisible(pos)
			}
		}
	}
}

func hasLineOfSight(floor *Floor, from Position, to Position) bool {
	line := bresenhamLine(from, to)
	if len(line) == 0 {
		return false
	}
	for index := 1; index < len(line)-1; index++ {
		if !floor.IsTransparent(line[index]) {
			return false
		}
	}
	return true
}

func bresenhamLine(from Position, to Position) []Position {
	points := make([]Position, 0, max(abs(to.X-from.X), abs(to.Y-from.Y))+1)

	x0 := from.X
	y0 := from.Y
	x1 := to.X
	y1 := to.Y

	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy

	for {
		points = append(points, Position{X: x0, Y: y0})
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}

	return points
}
