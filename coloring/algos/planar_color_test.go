package algos

import (
	"testing"

	"github.com/Salvatore112/graph_analysis_algorithms/graphs"
)

func makeK4() *graphs.BasicGraph {
	g := graphs.NewBasicGraph()
	vs := []string{"A", "B", "C", "D"}
	for i := 0; i < len(vs); i++ {
		for j := i + 1; j < len(vs); j++ {
			g.AddEdge(vs[i], vs[j])
		}
	}
	return g
}

func TestFiveColor_K4(t *testing.T) {
	g := makeK4()
	colors, err := FiveColorPlanar(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for v, neis := range g.Vertices {
		for _, u := range neis {
			if colors[v] == colors[u] {
				t.Fatalf("adjacent %s and %s share color %d", v, u, colors[v])
			}
		}
	}
}

func TestFourColor_K4(t *testing.T) {
	g := makeK4()
	colors, err := FourColorPlanar(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	used := map[int]struct{}{}
	for _, c := range colors {
		used[c] = struct{}{}
	}
	if len(used) != 4 {
		t.Fatalf("expected to use 4 colors on K4, used %d", len(used))
	}
}
