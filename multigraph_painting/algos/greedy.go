package multigraph_algo

import (
	"sort"

	"github.com/Salvatore112/graph_analysis_algorithms/graphs"
)

// GreedyEdgeColoring возвращает (colorsUsed, edges, colorsMap)
func GreedyEdgeColoring(g *graphs.MultiGraph) (int, []Edge, map[int]int) {
	edges := EdgeList(g)

	inc := make(map[string][]int)
	for _, e := range edges {
		inc[e.U] = append(inc[e.U], e.ID)
		inc[e.V] = append(inc[e.V], e.ID)
	}

	type it struct {
		id     int
		weight int
	}
	items := make([]it, 0, len(edges))
	for _, e := range edges {
		w := Degree(g, e.U) + Degree(g, e.V)
		items = append(items, it{id: e.ID, weight: w})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].weight == items[j].weight {
			return items[i].id < items[j].id
		}
		return items[i].weight > items[j].weight
	})

	edgeByID := make(map[int]Edge, len(edges))
	for _, e := range edges {
		edgeByID[e.ID] = e
	}

	colors := make(map[int]int)
	maxColor := -1

	for _, it := range items {
		e := edgeByID[it.id]

		used := make(map[int]bool)
		for _, eid := range inc[e.U] {
			if c, ok := colors[eid]; ok {
				used[c] = true
			}
		}
		for _, eid := range inc[e.V] {
			if c, ok := colors[eid]; ok {
				used[c] = true
			}
		}

		c := 0
		for {
			if !used[c] {
				break
			}
			c++
		}
		colors[e.ID] = c
		if c > maxColor {
			maxColor = c
		}
	}

	return maxColor + 1, edges, colors
}
