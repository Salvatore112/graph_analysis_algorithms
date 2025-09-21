package algos

import (
	"errors"
	"slices"

	"github.com/Salvatore112/graph_analysis_algorithms/graphs"
)

func FiveColorPlanar(g *graphs.BasicGraph) (map[string]int, error) {
	adj := cloneAdj(g.Vertices)
	order := make([]removalRecord, 0, len(adj))

	for len(adj) > 0 {
		var pick string
		found := false
		for v := range adj {
			if len(adj[v]) <= 5 {
				pick = v
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("не нашли вершину степени ≤5; граф может быть не планарным/непростым")
		}
		neis := keysOf(adj[pick])
		order = append(order, removalRecord{v: pick, neighbors: slices.Clone(neis)})

		for u := range adj[pick] {
			delete(adj[u], pick)
		}
		delete(adj, pick)
	}

	colors := make(map[string]int, len(order))
	adj = Adj{}
	for i := len(order) - 1; i >= 0; i-- {
		rec := order[i]
		v := rec.v
		if adj[v] == nil {
			adj[v] = map[string]struct{}{}
		}
		for _, u := range rec.neighbors {
			if adj[u] == nil {
				adj[u] = map[string]struct{}{}
			}
			adj[v][u] = struct{}{}
			adj[u][v] = struct{}{}
		}

		used := usedNeighborColors(v, adj, colors)
		if len(used) < 5 {
			for c := 0; c < 5; c++ {
				if _, ok := used[c]; !ok {
					colors[v] = c
					goto colored
				}
			}
		} else {
			color2nei := make(map[int]string, 5)
			for u := range adj[v] {
				if cu, ok := colors[u]; ok {
					color2nei[cu] = u
				}
			}
			pairs := [][2]int{
				{0, 1}, {0, 2}, {0, 3}, {0, 4},
				{1, 2}, {1, 3}, {1, 4},
				{2, 3}, {2, 4},
				{3, 4},
			}
			swapped := false
			for _, ab := range pairs {
				a, b := ab[0], ab[1]
				uA, okA := color2nei[a]
				uB, okB := color2nei[b]
				if !okA || !okB {
					if !okA {
						colors[v] = a
					} else {
						colors[v] = b
					}
					swapped = true
					break
				}
				if sameKempeComponent(uA, uB, a, b, adj, colors) {
					continue
				}
				swapKempe(uA, a, b, adj, colors)
				colors[v] = a
				swapped = true
				break
			}
			if !swapped {
				return nil, errors.New("не удалось выполнить перестановку Кемпе; проверьте планарность графа")
			}
		}
	colored:
	}
	return colors, nil
}

type removalRecord struct {
	v         string
	neighbors []string
}

func keysOf(m map[string]struct{}) []string {
	res := make([]string, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

func usedNeighborColors(v string, adj Adj, colors map[string]int) map[int]struct{} {
	used := map[int]struct{}{}
	for u := range adj[v] {
		if c, ok := colors[u]; ok {
			used[c] = struct{}{}
		}
	}
	return used
}

func sameKempeComponent(s, t string, a, b int, adj Adj, colors map[string]int) bool {
	if s == t {
		return true
	}
	seen := map[string]struct{}{s: {}}
	q := []string{s}
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		for y := range adj[x] {
			cy, ok := colors[y]
			if !ok || (cy != a && cy != b) {
				continue
			}
			if _, vis := seen[y]; vis {
				continue
			}
			seen[y] = struct{}{}
			if y == t {
				return true
			}
			q = append(q, y)
		}
	}
	return false
}

func swapKempe(start string, a, b int, adj Adj, colors map[string]int) {
	target := map[string]struct{}{start: {}}
	q := []string{start}
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		for y := range adj[x] {
			cy, ok := colors[y]
			if !ok || (cy != a && cy != b) {
				continue
			}
			if _, vis := target[y]; vis {
				continue
			}
			target[y] = struct{}{}
			q = append(q, y)
		}
	}
	for v := range target {
		if colors[v] == a {
			colors[v] = b
		} else {
			colors[v] = a
		}
	}
}
