package multigraph_algo

import (
	"sort"

	"github.com/Salvatore112/graph_analysis_algorithms/graphs"
)

// Edge — физический экземпляр ребра, генерируемый из graphs.MultiGraph
type Edge struct {
	ID int
	U  string
	V  string
}

// EdgeList — разворачивает кратности из graphs.MultiGraph в список физических рёбер.
// Генерирует детерминированный порядок: сортировка вершин, затем для каждой пары (u<=v)
// создаёт cnt экземпляров.
func EdgeList(g *graphs.MultiGraph) []Edge {
	edges := make([]Edge, 0)
	verts := make([]string, 0, len(g.Vertices))
	for v := range g.Vertices {
		verts = append(verts, v)
	}
	sort.Strings(verts)

	id := 0
	for _, u := range verts {
		neigh := g.Vertices[u]
		if neigh == nil {
			continue
		}
		for v, cnt := range neigh {
			if u > v {
				continue
			}
			for i := 0; i < cnt; i++ {
				edges = append(edges, Edge{ID: id, U: u, V: v})
				id++
			}
		}
	}
	return edges
}

// Degree считает степень (с учётом кратностей) вершины v.
func Degree(g *graphs.MultiGraph, v string) int {
	if g.Vertices[v] == nil {
		return 0
	}
	sum := 0
	for _, c := range g.Vertices[v] {
		sum += c
	}
	return sum
}

func MaxDegree(g *graphs.MultiGraph) int {
	_max := 0
	for v := range g.Vertices {
		if d := Degree(g, v); d > _max {
			_max = d
		}
	}
	return _max
}

func Mu(g *graphs.MultiGraph) int {
	m := 0
	seen := make(map[string]map[string]bool)
	for u, neigh := range g.Vertices {
		if seen[u] == nil {
			seen[u] = make(map[string]bool)
		}
		for v, cnt := range neigh {
			if seen[v] != nil && seen[v][u] {
				continue
			}
			if cnt > m {
				m = cnt
			}
			seen[u][v] = true
		}
	}
	return m
}
