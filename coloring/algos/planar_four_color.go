package algos

import (
	"errors"
	"slices"
	"sort"

	"github.com/Salvatore112/graph_analysis_algorithms/graphs"
)

func FourColorPlanar(g *graphs.BasicGraph) (map[string]int, error) {
	adj := cloneAdj(g.Vertices)
	verts := make([]string, 0, len(adj))
	for v := range adj {
		verts = append(verts, v)
	}
	sort.Slice(verts, func(i, j int) bool { return len(adj[verts[i]]) > len(adj[verts[j]]) })

	colors := map[string]int{}
	uncolored := map[string]struct{}{}
	for _, v := range verts {
		uncolored[v] = struct{}{}
	}

	ok := dsaturColor(adj, colors, uncolored, 4)
	if !ok {
		return nil, errors.New("не удалось 4-раскрасить граф (возможно, он не планарный)")
	}
	return colors, nil
}

type satInfo struct {
	v          string
	saturation int
	degree     int
}

func dsaturColor(adj Adj, colors map[string]int, uncolored map[string]struct{}, k int) bool {
	if len(uncolored) == 0 {
		return true
	}

	best := pickByDSATUR(adj, colors, uncolored)
	v := best.v

	used := map[int]struct{}{}
	for u := range adj[v] {
		if c, ok := colors[u]; ok {
			used[c] = struct{}{}
		}
	}
	cands := make([]int, 0, k)
	for c := 0; c < k; c++ {
		if _, ok := used[c]; !ok {
			cands = append(cands, c)
		}
	}
	slices.Sort(cands)

	delete(uncolored, v)
	for _, c := range cands {
		colors[v] = c
		if dsaturColor(adj, colors, uncolored, k) {
			return true
		}
		delete(colors, v)
	}
	uncolored[v] = struct{}{}
	return false
}

func pickByDSATUR(adj Adj, colors map[string]int, uncolored map[string]struct{}) satInfo {
	var best satInfo
	best.degree = -1
	for v := range uncolored {
		seen := map[int]struct{}{}
		for u := range adj[v] {
			if c, ok := colors[u]; ok {
				seen[c] = struct{}{}
			}
		}
		sat := len(seen)
		deg := len(adj[v])
		if sat > best.saturation || (sat == best.saturation && deg > best.degree) {
			best = satInfo{v: v, saturation: sat, degree: deg}
		}
	}
	return best
}
