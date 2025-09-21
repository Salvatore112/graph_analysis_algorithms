package multigraph_algo

import (
	"fmt"
)

// VerifyEdgeColoring проверяет, что на каждой вершине инцидентные рёбра имеют разные цвета.
// edges и colors берутся из EdgeList/GreedyEdgeColoring/BipartiteEdgeColoring.
func VerifyEdgeColoring(edges []Edge, colors map[int]int) error {
	inc := make(map[string][]int)
	for _, e := range edges {
		inc[e.U] = append(inc[e.U], e.ID)
		inc[e.V] = append(inc[e.V], e.ID)
	}
	for v, eids := range inc {
		used := make(map[int]int)
		for _, eid := range eids {
			c, ok := colors[eid]
			if !ok {
				return fmt.Errorf("edge %d incident to vertex %s is not colored", eid, v)
			}
			if prev, exists := used[c]; exists {
				return fmt.Errorf("conflict at vertex %s: edges %d and %d have same color %d", v, prev, eid, c)
			}
			used[c] = eid
		}
	}
	return nil
}
