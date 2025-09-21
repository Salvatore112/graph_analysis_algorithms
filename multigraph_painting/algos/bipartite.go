package multigraph_algo

import (
	"errors"
	"fmt"
	"sort"

	"github.com/Salvatore112/graph_analysis_algorithms/graphs"
)

const pairSep = "\x00"

// IsBipartite проверяет двудольность графа и возвращает map vertex->partition (0 или 1).
// Если не двудольный — возвращает (nil, false).
func IsBipartite(g *graphs.MultiGraph) (map[string]int, bool) {
	part := make(map[string]int)
	visited := make(map[string]bool)

	for start := range g.Vertices {
		if visited[start] {
			continue
		}
		// BFS
		queue := []string{start}
		part[start] = 0
		visited[start] = true
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for v := range g.Vertices[u] {
				if u == v {
					return nil, false
				}
				if !visited[v] {
					part[v] = 1 - part[u]
					visited[v] = true
					queue = append(queue, v)
				} else {
					if part[v] == part[u] {
						return nil, false
					}
				}
			}
		}
	}
	return part, true
}

// keyOf возвращает упорядоченную пару (left, right) при известном разбиении part.
// Если один из аргументов в части 0, другой в части 1 — возвращает (left,right).
// В иных случаях возвращает (a,b) как есть.
func keyOf(a, b string, part map[string]int) (string, string) {
	if part[a] == 0 && part[b] == 1 {
		return a, b
	}
	if part[b] == 0 && part[a] == 1 {
		return b, a
	}
	return a, b
}

// BipartiteEdgeColoring раскрашивает рёбра двудольного мультиграфа в Δ цветов.
// Алгоритм: повторно находит максимальное паросочетание на остаче рёбер Δ раз (Δ = MaxDegree).
// Для каждого найденного паросочетания все выбранные рёбра получают текущий цвет.
// Возвращает (colorsUsed, edges, colorsMap, error).
func BipartiteEdgeColoring(g *graphs.MultiGraph) (int, []Edge, map[int]int, error) {
	part, ok := IsBipartite(g)
	if !ok {
		return 0, nil, nil, errors.New("graph is not bipartite")
	}

	edges := EdgeList(g)
	pairEdges := make(map[string][]int)
	for _, e := range edges {
		u, v := keyOf(e.U, e.V, part)
		k := u + pairSep + v
		pairEdges[k] = append(pairEdges[k], e.ID)
	}

	leftVerts := make([]string, 0)
	leftSet := make(map[string]bool)
	rightSet := make(map[string]bool)
	for v, p := range part {
		if p == 0 {
			leftVerts = append(leftVerts, v)
			leftSet[v] = true
		} else {
			rightSet[v] = true
		}
	}
	sort.Strings(leftVerts)

	colors := make(map[int]int)
	Delta := MaxDegree(g)

	for color := 0; color < Delta; color++ {
		adj := make(map[string][]string)
		for k, list := range pairEdges {
			if len(list) == 0 {
				continue
			}
			parts := make([]string, 2)
			for i := 0; i < len(k); i++ {
				if k[i] == pairSep[0] {
					parts[0] = k[:i]
					parts[1] = k[i+1:]
					break
				}
			}
			if parts[0] == "" {
				continue
			}
			u := parts[0]
			v := parts[1]
			if leftSet[u] && rightSet[v] {
				adj[u] = append(adj[u], v)
			} else if leftSet[v] && rightSet[u] {
				adj[v] = append(adj[v], u)
			}
		}

		pairU, _ := hopcroftKarp(adj, leftVerts)

		for u, v := range pairU {
			if v == "" {
				continue
			}
			k := u + pairSep + v
			list := pairEdges[k]
			if len(list) == 0 {
				continue
			}
			eid := list[len(list)-1]
			pairEdges[k] = list[:len(list)-1]
			colors[eid] = color
		}
	}

	for _, e := range edges {
		if _, ok := colors[e.ID]; !ok {
			return 0, edges, colors, fmt.Errorf("failed to color all edges; remaining edge id %d", e.ID)
		}
	}

	return Delta, edges, colors, nil
}
