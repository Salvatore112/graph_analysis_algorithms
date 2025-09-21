package multigraph_algo

import (
	"reflect"
	"testing"

	"github.com/Salvatore112/graph_analysis_algorithms/graphs"
)

func makeMultiGraph(pairs map[[2]string]int) *graphs.MultiGraph {
	g := &graphs.MultiGraph{Vertices: make(map[string]map[string]int)}
	add := func(u, v string, cnt int) {
		if g.Vertices[u] == nil {
			g.Vertices[u] = make(map[string]int)
		}
		g.Vertices[u][v] = cnt
	}
	for k, cnt := range pairs {
		u, v := k[0], k[1]
		add(u, v, cnt)
		add(v, u, cnt)
	}
	return g
}

func TestEdgeListDegreeMaxMu(t *testing.T) {
	g := makeMultiGraph(map[[2]string]int{
		{"a", "b"}: 2,
		{"a", "c"}: 1,
		{"b", "c"}: 1,
	})

	edges := EdgeList(g)
	if len(edges) != 4 {
		t.Fatalf("EdgeList: ожидается 4 физических ребра, получили %d", len(edges))
	}

	if d := Degree(g, "a"); d != 3 {
		t.Fatalf("Degree(a) expected 3, got %d", d)
	}
	if d := Degree(g, "b"); d != 3 {
		t.Fatalf("Degree(b) expected 3, got %d", d)
	}
	if d := Degree(g, "c"); d != 2 {
		t.Fatalf("Degree(c) expected 2, got %d", d)
	}

	if md := MaxDegree(g); md != 3 {
		t.Fatalf("MaxDegree expected 3, got %d", md)
	}

	if mu := Mu(g); mu != 2 {
		t.Fatalf("Mu expected 2, got %d", mu)
	}
}

func TestGreedyEdgeColoringValidity(t *testing.T) {
	g := makeMultiGraph(map[[2]string]int{
		{"a", "b"}: 2,
		{"a", "c"}: 1,
		{"b", "c"}: 1,
	})

	colorsUsed, edges, colors := GreedyEdgeColoring(g)

	if err := VerifyEdgeColoring(edges, colors); err != nil {
		t.Fatalf("GreedyEdgeColoring produced invalid coloring: %v", err)
	}

	_max := -1
	for _, c := range colors {
		if c > _max {
			_max = c
		}
	}
	if colorsUsed != _max+1 {
		t.Fatalf("Greedy returned colorsUsed=%d but max(color)+1=%d", colorsUsed, _max+1)
	}
}

func TestBipartiteEdgeColoringExample(t *testing.T) {
	g := makeMultiGraph(map[[2]string]int{
		{"a", "x"}: 2,
		{"a", "y"}: 1,
		{"b", "x"}: 1,
	})

	if _, ok := IsBipartite(g); !ok {
		t.Fatalf("expected graph to be bipartite")
	}

	k, edges, colors, err := BipartiteEdgeColoring(g)
	if err != nil {
		t.Fatalf("BipartiteEdgeColoring returned error: %v", err)
	}

	if k != MaxDegree(g) {
		t.Fatalf("BipartiteEdgeColoring used %d colors; expected Δ=%d", k, MaxDegree(g))
	}

	if err := VerifyEdgeColoring(edges, colors); err != nil {
		t.Fatalf("BipartiteEdgeColoring produced invalid coloring: %v", err)
	}
}

func TestBipartiteEdgeColoringNonBipartite(t *testing.T) {
	g := makeMultiGraph(map[[2]string]int{
		{"a", "b"}: 1,
		{"b", "c"}: 1,
		{"c", "a"}: 1,
	})

	if _, ok := IsBipartite(g); ok {
		t.Fatalf("triangle should not be bipartite")
	}

	_, _, _, err := BipartiteEdgeColoring(g)
	if err == nil {
		t.Fatalf("expected BipartiteEdgeColoring to fail on non-bipartite graph")
	}
}

func TestHopcroftKarpSimple(t *testing.T) {
	adj := map[string][]string{
		"u1": {"v1", "v2"},
		"u2": {"v2"},
	}
	left := []string{"u1", "u2"}
	pairU, pairV := hopcroftKarp(adj, left)

	count := 0
	for _, v := range pairU {
		if v != "" {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("expected matching size 2, got %d; pairU: %v, pairV: %v", count, pairU, pairV)
	}

	for u, v := range pairU {
		if v == "" {
			continue
		}
		if pairV[v] != u {
			t.Fatalf("inconsistent matching: pairU[%s]=%s but pairV[%s]=%s", u, v, v, pairV[v])
		}
	}
}

func TestVerifyEdgeColoringDetectsUncolored(t *testing.T) {
	g := makeMultiGraph(map[[2]string]int{
		{"a", "b"}: 1,
		{"b", "c"}: 1,
	})

	edges := EdgeList(g)
	colors := make(map[int]int)
	if len(edges) == 0 {
		t.Fatalf("no edges")
	}
	colors[edges[0].ID] = 0

	if err := VerifyEdgeColoring(edges, colors); err == nil {
		t.Fatalf("VerifyEdgeColoring should detect uncolored edge and return error")
	}
}

func TestEdgeListDeterministic(t *testing.T) {
	g := makeMultiGraph(map[[2]string]int{
		{"m", "n"}: 1,
		{"a", "z"}: 2,
		{"b", "c"}: 1,
	})

	e1 := EdgeList(g)
	e2 := EdgeList(g)

	if !reflect.DeepEqual(e1, e2) {
		t.Fatalf("EdgeList should be deterministic but results differ\n%v\n%v", e1, e2)
	}
}

func TestExactEdgeColoring_Empty(t *testing.T) {
	g := makeMultiGraph(map[[2]string]int{})
	k, edges, colors := ExactEdgeColoring(g)
	if k != 0 {
		t.Fatalf("expected 0 colors for empty graph, got %d", k)
	}
	if err := VerifyEdgeColoring(edges, colors); err != nil {
		t.Fatalf("verification failed on empty graph: %v", err)
	}
}

func TestExactEdgeColoring_ParallelEdgesOnly(t *testing.T) {
	g := makeMultiGraph(map[[2]string]int{
		{"u", "v"}: 3,
	})
	k, edges, colors := ExactEdgeColoring(g)
	if k != 3 {
		t.Fatalf("expected 3 colors (Δ=3), got %d", k)
	}
	if err := VerifyEdgeColoring(edges, colors); err != nil {
		t.Fatalf("verification failed: %v", err)
	}
}

func TestExactEdgeColoring_BipartiteExample(t *testing.T) {
	g := makeMultiGraph(map[[2]string]int{
		{"a", "x"}: 2,
		{"a", "y"}: 1,
		{"b", "x"}: 1,
	})

	k, edges, colors := ExactEdgeColoring(g)
	if k != 3 {
		t.Fatalf("expected optimal 3 colors for bipartite multigraph (Δ=3), got %d", k)
	}
	if err := VerifyEdgeColoring(edges, colors); err != nil {
		t.Fatalf("verification failed: %v", err)
	}
}

func TestExactEdgeColoring_TriangleSimple(t *testing.T) {
	g := makeMultiGraph(map[[2]string]int{
		{"a", "b"}: 1,
		{"b", "c"}: 1,
		{"c", "a"}: 1,
	})
	k, edges, colors := ExactEdgeColoring(g)
	if k != 3 {
		t.Fatalf("expected 3 colors for K3 (Δ=2, class 2), got %d", k)
	}
	if err := VerifyEdgeColoring(edges, colors); err != nil {
		t.Fatalf("verification failed: %v", err)
	}
}

func TestExactEdgeColoring_ShannonTightDoubleK3(t *testing.T) {
	g := makeMultiGraph(map[[2]string]int{
		{"a", "b"}: 2,
		{"b", "c"}: 2,
		{"c", "a"}: 2,
	})
	k, edges, colors := ExactEdgeColoring(g)
	if k != 6 {
		t.Fatalf("expected 6 colors for doubled K3 (Δ=4, χ'=6), got %d", k)
	}
	if err := VerifyEdgeColoring(edges, colors); err != nil {
		t.Fatalf("verification failed: %v", err)
	}
}
