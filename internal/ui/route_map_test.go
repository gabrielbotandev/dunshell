package ui

import "testing"

func TestBuildRouteMapLayoutKeepsNodesInBounds(t *testing.T) {
	layout := buildRouteMapLayout(3, 1, 52, 17)

	if len(layout.Nodes) != 5 {
		t.Fatalf("expected 5 nodes, got %d", len(layout.Nodes))
	}
	if layout.SelectedRoute != 1 {
		t.Fatalf("expected selected route 1, got %d", layout.SelectedRoute)
	}
	for index, node := range layout.Nodes {
		if node.X < 0 || node.Y < 0 {
			t.Fatalf("node %d out of bounds at negative position %+v", index, node)
		}
		if node.X+node.Size > layout.Width {
			t.Fatalf("node %d exceeds layout width: %+v in %d", index, node, layout.Width)
		}
		if node.Y+node.Size > layout.Height {
			t.Fatalf("node %d exceeds layout height: %+v in %d", index, node, layout.Height)
		}
	}
}

func TestBuildRouteMapLayoutCompactsForTightPanels(t *testing.T) {
	layout := buildRouteMapLayout(3, 7, 40, 14)

	if layout.NodeSize != 3 {
		t.Fatalf("expected compact node size 3, got %d", layout.NodeSize)
	}
	if layout.SelectedRoute != 2 {
		t.Fatalf("expected selection clamped to last route, got %d", layout.SelectedRoute)
	}
	if layout.Nodes[1].Y >= layout.Nodes[2].Y || layout.Nodes[2].Y >= layout.Nodes[3].Y {
		t.Fatalf("expected route nodes to stay vertically ordered: %+v", layout.Nodes[1:4])
	}
	if layout.LeftSpineX >= layout.Nodes[1].X {
		t.Fatalf("expected left spine before route column, got spine %d and route x %d", layout.LeftSpineX, layout.Nodes[1].X)
	}
	if layout.RightSpineX <= layout.Nodes[1].X+layout.NodeSize-1 {
		t.Fatalf("expected right spine after route column, got spine %d", layout.RightSpineX)
	}
}
