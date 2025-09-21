package algos

import (
	"strconv"

	"github.com/Salvatore112/graph_analysis_algorithms/graphs"
)

func FourColor(adj [][]int) []int {
	n := len(adj)

	g := graphs.NewBasicGraph()
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = strconv.Itoa(i)
	}

	for i := 0; i < n; i++ {
		vi := names[i]
		for _, j := range adj[i] {
			if j < 0 || j >= n || i == j {
				continue
			}
			if i < j {
				g.AddEdge(vi, names[j])
			}
		}
	}

	colsMap, err := FourColorPlanar(g)
	if err != nil {
		return nil
	}

	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = colsMap[names[i]]
	}
	return out
}
