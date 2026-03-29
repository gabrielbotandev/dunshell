package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"dunshell/internal/game"
)

type routeMapNodeKind int

const (
	routeMapNodeCurrent routeMapNodeKind = iota
	routeMapNodeChoice
	routeMapNodeDestination
)

type routeMapNode struct {
	X          int
	Y          int
	Size       int
	Kind       routeMapNodeKind
	RouteIndex int
}

type routeMapLayout struct {
	Width           int
	Height          int
	NodeSize        int
	SelectedRoute   int
	CurrentNode     int
	DestinationNode int
	LeftSpineX      int
	RightSpineX     int
	Nodes           []routeMapNode
}

type routePoint struct {
	X int
	Y int
}

type routeCanvasCell struct {
	mask         uint8
	lineStyle    lipgloss.Style
	linePriority int
	hasGlyph     bool
	glyph        string
	glyphStyle   lipgloss.Style
}

const (
	routeNorth uint8 = 1 << iota
	routeEast
	routeSouth
	routeWest
)

func buildRouteMapLayout(routeCount int, selected int, width int, height int) routeMapLayout {
	width = max(1, width)
	height = max(1, height)
	selected = clamp(selected, 0, max(0, routeCount-1))

	nodeSize := 5
	gapY := 2
	if width < 46 || height < 15 {
		nodeSize = 3
		gapY = 1
	}

	currentX := 1
	if width <= 20 {
		currentX = 0
	}
	destinationX := max(currentX+nodeSize+4, width-nodeSize-1)
	if destinationX+nodeSize > width {
		destinationX = max(0, width-nodeSize)
	}

	routeX := width/2 - nodeSize/2
	minRouteX := currentX + nodeSize + 2
	maxRouteX := destinationX - nodeSize - 2
	if maxRouteX < minRouteX {
		routeX = max(currentX+nodeSize+1, destinationX-nodeSize-2)
	} else {
		routeX = clamp(routeX, minRouteX, maxRouteX)
	}
	if routeX+nodeSize > width {
		routeX = max(0, width-nodeSize)
	}

	currentY := max(0, (height-nodeSize)/2)
	destinationY := currentY
	neededHeight := routeCount * nodeSize
	if routeCount > 1 {
		neededHeight += (routeCount - 1) * gapY
	}
	top := max(0, (height-neededHeight)/2)
	if routeCount == 1 {
		top = currentY
	}

	nodes := make([]routeMapNode, 0, routeCount+2)
	nodes = append(nodes, routeMapNode{X: currentX, Y: currentY, Size: nodeSize, Kind: routeMapNodeCurrent, RouteIndex: -1})
	for index := 0; index < routeCount; index++ {
		y := top + index*(nodeSize+gapY)
		if y+nodeSize > height {
			y = max(0, height-nodeSize)
		}
		nodes = append(nodes, routeMapNode{X: routeX, Y: y, Size: nodeSize, Kind: routeMapNodeChoice, RouteIndex: index})
	}
	destinationNode := len(nodes)
	nodes = append(nodes, routeMapNode{X: destinationX, Y: destinationY, Size: nodeSize, Kind: routeMapNodeDestination, RouteIndex: -1})

	currentExitX := nodes[0].X + nodeSize - 1
	leftSpineX := max(currentExitX+1, currentExitX+(routeX-currentExitX)/2)
	if leftSpineX >= routeX {
		leftSpineX = max(currentExitX, routeX-1)
	}

	routeExitX := routeX + nodeSize - 1
	rightSpineX := max(routeExitX+1, routeExitX+(destinationX-routeExitX)/2)
	if rightSpineX >= destinationX {
		rightSpineX = max(routeExitX, destinationX-1)
	}

	return routeMapLayout{
		Width:           width,
		Height:          height,
		NodeSize:        nodeSize,
		SelectedRoute:   selected,
		CurrentNode:     0,
		DestinationNode: destinationNode,
		LeftSpineX:      clamp(leftSpineX, 0, max(0, width-1)),
		RightSpineX:     clamp(rightSpineX, 0, max(0, width-1)),
		Nodes:           nodes,
	}
}

func (m *Model) renderRouteGraph(routes []game.RouteChoice, selected int, width int, height int) string {
	layout := buildRouteMapLayout(len(routes), selected, width, height)
	canvas := make([][]routeCanvasCell, layout.Height)
	for y := range canvas {
		canvas[y] = make([]routeCanvasCell, layout.Width)
	}

	current := layout.Nodes[layout.CurrentNode]
	destination := layout.Nodes[layout.DestinationNode]
	mutedEdge := lipgloss.NewStyle().Foreground(lipgloss.Color("#5f5852"))

	for _, node := range layout.Nodes {
		if node.Kind != routeMapNodeChoice || node.RouteIndex < 0 || node.RouteIndex >= len(routes) {
			continue
		}
		route := routes[node.RouteIndex]
		edgeStyle := mutedEdge
		priority := 1
		if node.RouteIndex == layout.SelectedRoute {
			edgeStyle = lipgloss.NewStyle().Foreground(routeAccentColor(route)).Bold(true)
			priority = 2
		}

		entryY := node.Y + node.Size/2
		leftPath := []routePoint{
			{X: current.X + current.Size - 1, Y: current.Y + current.Size/2},
			{X: layout.LeftSpineX, Y: current.Y + current.Size/2},
			{X: layout.LeftSpineX, Y: entryY},
			{X: node.X, Y: entryY},
		}
		rightPath := []routePoint{
			{X: node.X + node.Size - 1, Y: entryY},
			{X: layout.RightSpineX, Y: entryY},
			{X: layout.RightSpineX, Y: destination.Y + destination.Size/2},
			{X: destination.X, Y: destination.Y + destination.Size/2},
		}
		drawRouteMapPath(canvas, leftPath, edgeStyle, priority)
		drawRouteMapPath(canvas, rightPath, edgeStyle, priority)
	}

	for _, node := range layout.Nodes {
		switch node.Kind {
		case routeMapNodeCurrent:
			m.drawRouteMapNode(canvas, node, nil, false, false)
		case routeMapNodeDestination:
			m.drawRouteMapNode(canvas, node, &routes[layout.SelectedRoute], false, false)
		case routeMapNodeChoice:
			route := routes[node.RouteIndex]
			m.drawRouteMapNode(canvas, node, &route, node.RouteIndex == layout.SelectedRoute, true)
			if node.RouteIndex == layout.SelectedRoute {
				pointerStyle := lipgloss.NewStyle().Foreground(routeAccentColor(route)).Bold(true)
				pointerX := max(0, node.X-2)
				if layout.NodeSize <= 3 {
					pointerX = max(0, node.X-1)
				}
				routeCanvasPut(canvas, pointerX, node.Y+node.Size/2, "›", pointerStyle)
			}
		}
	}

	rows := make([]string, 0, layout.Height)
	for y := 0; y < layout.Height; y++ {
		var builder strings.Builder
		for x := 0; x < layout.Width; x++ {
			cell := canvas[y][x]
			switch {
			case cell.hasGlyph:
				builder.WriteString(cell.glyphStyle.Render(cell.glyph))
			case cell.mask != 0:
				builder.WriteString(cell.lineStyle.Render(routeLineGlyph(cell.mask, m.glyphs.ascii)))
			default:
				builder.WriteByte(' ')
			}
		}
		rows = append(rows, builder.String())
	}
	return strings.Join(rows, "\n")
}

func drawRouteMapPath(canvas [][]routeCanvasCell, points []routePoint, style lipgloss.Style, priority int) {
	if len(points) < 2 {
		return
	}
	for index := 0; index < len(points)-1; index++ {
		start := points[index]
		end := points[index+1]
		switch {
		case start.X == end.X:
			routeCanvasConnectV(canvas, start.X, start.Y, end.Y, style, priority)
		case start.Y == end.Y:
			routeCanvasConnectH(canvas, start.Y, start.X, end.X, style, priority)
		}
	}
}

func routeCanvasConnectH(canvas [][]routeCanvasCell, y int, x0 int, x1 int, style lipgloss.Style, priority int) {
	if y < 0 || y >= len(canvas) {
		return
	}
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	for x := x0; x < x1; x++ {
		routeCanvasAddLine(canvas, x, y, routeEast, style, priority)
		routeCanvasAddLine(canvas, x+1, y, routeWest, style, priority)
	}
}

func routeCanvasConnectV(canvas [][]routeCanvasCell, x int, y0 int, y1 int, style lipgloss.Style, priority int) {
	if x < 0 {
		return
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	for y := y0; y < y1; y++ {
		routeCanvasAddLine(canvas, x, y, routeSouth, style, priority)
		routeCanvasAddLine(canvas, x, y+1, routeNorth, style, priority)
	}
}

func routeCanvasAddLine(canvas [][]routeCanvasCell, x int, y int, mask uint8, style lipgloss.Style, priority int) {
	if y < 0 || y >= len(canvas) || x < 0 || x >= len(canvas[y]) {
		return
	}
	canvas[y][x].mask |= mask
	if priority >= canvas[y][x].linePriority {
		canvas[y][x].lineStyle = style
		canvas[y][x].linePriority = priority
	}
}

func routeCanvasPut(canvas [][]routeCanvasCell, x int, y int, glyph string, style lipgloss.Style) {
	if y < 0 || y >= len(canvas) || x < 0 || x >= len(canvas[y]) {
		return
	}
	canvas[y][x].hasGlyph = true
	canvas[y][x].glyph = glyph
	canvas[y][x].glyphStyle = style
	canvas[y][x].mask = 0
	canvas[y][x].linePriority = 0
}

func routeLineGlyph(mask uint8, ascii bool) string {
	if ascii {
		switch mask {
		case routeNorth | routeSouth:
			return "|"
		case routeEast | routeWest:
			return "-"
		default:
			return "+"
		}
	}

	switch mask {
	case routeNorth | routeSouth:
		return "│"
	case routeEast | routeWest:
		return "─"
	case routeSouth | routeEast:
		return "┌"
	case routeSouth | routeWest:
		return "┐"
	case routeNorth | routeEast:
		return "└"
	case routeNorth | routeWest:
		return "┘"
	case routeNorth | routeSouth | routeEast:
		return "├"
	case routeNorth | routeSouth | routeWest:
		return "┤"
	case routeSouth | routeEast | routeWest:
		return "┬"
	case routeNorth | routeEast | routeWest:
		return "┴"
	case routeNorth | routeSouth | routeEast | routeWest:
		return "┼"
	case routeNorth:
		return "│"
	case routeEast:
		return "─"
	case routeSouth:
		return "│"
	case routeWest:
		return "─"
	default:
		return " "
	}
}

func (m *Model) drawRouteMapNode(canvas [][]routeCanvasCell, node routeMapNode, route *game.RouteChoice, selected bool, selectable bool) {
	borderColor := lipgloss.Color("#6b5d53")
	background := lipgloss.Color("#171311")
	iconColor := lipgloss.Color("#c9c1b4")
	badgeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#c9c1b4")).Background(background).Bold(true)
	glyph := m.glyphs.player()
	badge := ""

	switch node.Kind {
	case routeMapNodeCurrent:
		borderColor = lipgloss.Color("#8fa7bf")
		background = lipgloss.Color("#1b232d")
		iconColor = lipgloss.Color("#84d06f")
		glyph = m.glyphs.player()
	case routeMapNodeDestination:
		background = lipgloss.Color("#141214")
		if route != nil && route.BossFloor {
			borderColor = lipgloss.Color("#8f4b4e")
			iconColor = lipgloss.Color("#d47777")
			glyph = m.glyphs.bossGate()
		} else {
			borderColor = lipgloss.Color("#56606b")
			iconColor = lipgloss.Color("#8fb2ca")
			glyph = m.glyphs.stairs()
		}
	case routeMapNodeChoice:
		if route == nil {
			break
		}
		borderColor = routeAccentColor(*route)
		background = lipgloss.Color("#181311")
		iconColor = routeAccentColor(*route)
		glyph = m.routePrimaryGlyph(*route)
		badge = m.routeBadgeGlyph(*route)
		badgeStyle = m.routeBadgeStyle(*route).Background(background).Bold(true)
		if !selected {
			background = lipgloss.Color("#15110f")
			badgeStyle = m.routeBadgeStyle(*route).Background(background)
		}
		if selectable && selected {
			background = lipgloss.Color("#221915")
			badgeStyle = m.routeBadgeStyle(*route).Background(background).Bold(true)
		}
	}

	border := lipgloss.NewStyle().Foreground(borderColor).Background(background).Bold(selected || node.Kind == routeMapNodeCurrent)
	fill := lipgloss.NewStyle().Foreground(lipgloss.Color("#f1e7db")).Background(background)
	icon := lipgloss.NewStyle().Foreground(iconColor).Background(background).Bold(true)

	horizontal := "─"
	vertical := "│"
	topLeft := "┌"
	topRight := "┐"
	bottomLeft := "└"
	bottomRight := "┘"
	if m.glyphs.ascii {
		horizontal = "-"
		vertical = "|"
		topLeft = "+"
		topRight = "+"
		bottomLeft = "+"
		bottomRight = "+"
	}

	last := node.Size - 1
	for y := 0; y < node.Size; y++ {
		for x := 0; x < node.Size; x++ {
			switch {
			case y == 0 && x == 0:
				routeCanvasPut(canvas, node.X+x, node.Y+y, topLeft, border)
			case y == 0 && x == last:
				routeCanvasPut(canvas, node.X+x, node.Y+y, topRight, border)
			case y == last && x == 0:
				routeCanvasPut(canvas, node.X+x, node.Y+y, bottomLeft, border)
			case y == last && x == last:
				routeCanvasPut(canvas, node.X+x, node.Y+y, bottomRight, border)
			case y == 0 || y == last:
				routeCanvasPut(canvas, node.X+x, node.Y+y, horizontal, border)
			case x == 0 || x == last:
				routeCanvasPut(canvas, node.X+x, node.Y+y, vertical, border)
			default:
				routeCanvasPut(canvas, node.X+x, node.Y+y, " ", fill)
			}
		}
	}

	centerX := node.X + node.Size/2
	centerY := node.Y + node.Size/2
	if node.Size >= 5 && badge != "" {
		routeCanvasPut(canvas, centerX, node.Y+1, badge, badgeStyle)
	}
	routeCanvasPut(canvas, centerX, centerY, glyph, icon)
}

func (m *Model) routePrimaryGlyph(route game.RouteChoice) string {
	switch route.ID {
	case "gilded_way":
		return m.glyphs.symbol('◈', 'g')
	case "brokers_lantern":
		return m.glyphs.merchant()
	case "pilgrims_rest":
		return m.glyphs.symbol('✚', '+')
	case "reliquary_breach":
		return m.glyphs.chest(game.KeySilver)
	case "ashen_hunt":
		return m.glyphs.symbol('✠', '!')
	case "cursed_procession":
		return m.glyphs.symbol('☠', 'x')
	default:
		return m.glyphs.symbol('◇', 'o')
	}
}

func (m *Model) routeBadgeGlyph(route game.RouteChoice) string {
	switch {
	case route.Modifier.GuaranteedKey != nil:
		return m.glyphs.symbol('⚿', 'k')
	case route.Modifier.ExtraChests > 0:
		return m.glyphs.chest(game.KeyBronze)
	case route.Modifier.BonusGold > 0:
		return m.glyphs.symbol('¤', 'g')
	case route.Modifier.EliteChance > 0:
		return m.glyphs.symbol('✠', 'E')
	case route.Modifier.Cursed:
		return m.glyphs.symbol('☠', 'x')
	default:
		return ""
	}
}

func (m *Model) routeBadgeStyle(route game.RouteChoice) lipgloss.Style {
	switch {
	case route.Modifier.GuaranteedKey != nil:
		return m.styles.keyStyle(*route.Modifier.GuaranteedKey)
	case route.Modifier.ExtraChests > 0:
		return m.styles.Accent
	case route.Modifier.BonusGold > 0:
		return m.styles.Gold
	case route.Modifier.Merchant:
		return m.styles.Info
	case route.Modifier.Rest:
		return m.styles.Success
	case route.Modifier.Cursed:
		return m.styles.Danger
	case route.Modifier.EliteChance > 0:
		return m.styles.Warning
	default:
		return m.styles.Muted
	}
}

func routeTypeLabel(route game.RouteChoice) string {
	switch route.ID {
	case "gilded_way":
		return "Gold Route"
	case "brokers_lantern":
		return "Merchant Route"
	case "pilgrims_rest":
		return "Rest Route"
	case "reliquary_breach":
		return "Reliquary Route"
	case "ashen_hunt":
		return "Elite Route"
	case "cursed_procession":
		return "Cursed Route"
	default:
		return "Unmarked Route"
	}
}

func (m *Model) routeDestinationLine(route game.RouteChoice) string {
	nextFloor := m.game.FloorIndex + 1
	if route.BossFloor {
		return fmt.Sprintf("Floor %d is a boss floor. The keeper waits behind a sealed gate.", nextFloor)
	}
	return fmt.Sprintf("Floor %d is generated when you descend, carrying this branch's omen.", nextFloor)
}

func (m *Model) routePressureDetail(route game.RouteChoice) string {
	modifier := route.Modifier
	switch {
	case route.BossFloor && modifier.Cursed:
		return "Boss floor ahead with cursed scaling, added elites, and extra foes."
	case route.BossFloor:
		return "Boss floor ahead. The keeper defines the pressure more than the halls."
	case modifier.Cursed:
		return "Cursed scaling hardens the floor while extra enemies and elites crowd the path."
	case modifier.EliteChance > 0 && modifier.EnemyBonus > 0:
		return "Sharper combat pressure with extra foes and a stronger elite presence."
	case modifier.EliteChance > 0:
		return "Sharper combat pressure. Expect a stronger elite presence."
	case modifier.Rest:
		return "Gentler opening pressure with recovery and fewer enemies."
	default:
		return "Steady descent. The next floor keeps its normal pace."
	}
}

func (m *Model) routeDetailLines(route game.RouteChoice, includeSignals bool) []string {
	lines := []string{
		m.routeStyle(route).Bold(true).Render(m.routePrimaryGlyph(route) + " " + route.Title),
		m.styles.Muted.Render(route.Subtitle),
		"",
		m.styles.Subtitle.Render("Type"),
		routeTypeLabel(route),
		"",
		m.styles.Subtitle.Render("Destination"),
		m.routeDestinationLine(route),
		"",
		m.styles.Success.Render("Reward"),
		route.Reward,
		"",
		m.styles.Warning.Render("Pressure"),
		m.routePressureDetail(route),
	}
	if includeSignals {
		lines = append(lines, "", m.styles.AccentSoft.Render("Rewards And Omens"))
		lines = append(lines, m.routeSignalLines(route)...)
	}
	return lines
}

func (m *Model) routeSignalLines(route game.RouteChoice) []string {
	rewards := m.routeBenefitLines(route)
	dangers := m.routeDangerLines(route)
	lines := []string{m.styles.Success.Render("Rewards")}
	lines = append(lines, rewards...)
	lines = append(lines, "", m.styles.Warning.Render("Omens"))
	lines = append(lines, dangers...)
	return lines
}

func (m *Model) routeBenefitLines(route game.RouteChoice) []string {
	modifier := route.Modifier
	lines := make([]string, 0, 8)
	if modifier.BonusGold > 0 {
		lines = append(lines, m.styles.Gold.Render(m.glyphs.symbol('◈', 'g'))+" Richer gold from kills")
	}
	if modifier.Merchant {
		lines = append(lines, m.styles.Info.Render(m.glyphs.merchant())+" Merchant guaranteed")
	}
	if modifier.Rest {
		lines = append(lines, m.styles.Success.Render(m.glyphs.symbol('✚', '+'))+fmt.Sprintf(" Recover %d HP on arrival", modifier.HealOnStart))
		if modifier.CleanseOnRest {
			lines = append(lines, m.styles.Cure.Render(m.glyphs.symbol('✢', '*'))+" Cleanse poison and fire")
		}
	}
	if modifier.GuaranteedKey != nil {
		icon := m.styles.keyStyle(*modifier.GuaranteedKey).Render(m.glyphs.symbol('⚿', 'k'))
		lines = append(lines, icon+" "+modifier.GuaranteedKey.Label()+" key on arrival")
	}
	if modifier.ExtraChests > 0 {
		lines = append(lines, m.styles.Accent.Render(m.glyphs.chest(game.KeyBronze))+" Extra reliquary chest")
	}
	if modifier.LootBonus > 0 {
		lines = append(lines, m.styles.AccentSoft.Render(m.glyphs.symbol('✦', '*'))+" Stronger loot rolls")
	}
	if modifier.EnemyBonus < 0 {
		lines = append(lines, m.styles.Success.Render(m.glyphs.symbol('☽', '~'))+" Fewer enemies in the halls")
	}
	if len(lines) == 0 {
		lines = append(lines, m.styles.Muted.Render("No special reward beyond a steady descent."))
	}
	return lines
}

func (m *Model) routeDangerLines(route game.RouteChoice) []string {
	modifier := route.Modifier
	lines := make([]string, 0, 8)
	if route.BossFloor {
		lines = append(lines, m.styles.Danger.Render(m.glyphs.bossGate())+" Keeper chamber on the next floor")
	}
	if modifier.EnemyBonus > 0 {
		lines = append(lines, m.styles.Warning.Render(m.glyphs.symbol('⚔', '!'))+fmt.Sprintf(" %d extra foes stalk the floor", modifier.EnemyBonus))
	}
	if modifier.EliteChance > 0 {
		lines = append(lines, m.styles.Warning.Render(m.glyphs.symbol('✠', 'E'))+fmt.Sprintf(" %.0f%% more elite pressure", modifier.EliteChance*100))
	}
	if modifier.Cursed {
		lines = append(lines, m.styles.Danger.Render(m.glyphs.symbol('☠', 'x'))+" Cursed scaling hardens the floor")
	}
	if len(lines) == 0 {
		lines = append(lines, m.styles.Muted.Render("No added omen beyond the normal floor."))
	}
	return lines
}
