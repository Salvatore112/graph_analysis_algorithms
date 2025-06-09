package ff

import (
	"container/list"
)

type Graph struct {
	n   int
	adj [][]int
}

func NewGraph(n int) *Graph {
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		adj[i] = []int{}
	}
	return &Graph{
		n:   n,
		adj: adj,
	}
}

func (g *Graph) AddEdge(u, v int) {
	g.adj[u] = append(g.adj[u], v)
	g.adj[v] = append(g.adj[v], u)
}

type Blossom struct {
	n       int
	graph   [][]int
	match   []int
	parent  []int
	base    []int
	used    []bool
	blossom []bool
}

func NewBlossom(n int, graph [][]int) *Blossom {
	match := make([]int, n)
	parent := make([]int, n)
	base := make([]int, n)
	used := make([]bool, n)
	blossom := make([]bool, n)
	for i := 0; i < n; i++ {
		match[i] = -1
		parent[i] = -1
		base[i] = i
		used[i] = false
		blossom[i] = false
	}
	return &Blossom{
		n:       n,
		graph:   graph,
		match:   match,
		parent:  parent,
		base:    base,
		used:    used,
		blossom: blossom,
	}
}

func (b *Blossom) lca(a, b1 int) int {
	used := make([]bool, b.n)
	for {
		a = b.base[a]
		used[a] = true
		if b.match[a] == -1 {
			break
		}
		a = b.parent[b.match[a]]
	}
	for {
		b1 = b.base[b1]
		if used[b1] {
			return b1
		}
		b1 = b.parent[b.match[b1]]
	}
}

func (b *Blossom) markPath(v, bbase, x int) {
	for b.base[v] != bbase {
		b.blossom[b.base[v]] = true
		b.blossom[b.base[b.match[v]]] = true
		b.parent[v] = x
		x = b.match[v]
		v = b.parent[b.match[v]]
	}
}

func (b *Blossom) findPath(root int) int {
	for i := 0; i < b.n; i++ {
		b.used[i] = false
		b.parent[i] = -1
		b.base[i] = i
	}
	queue := list.New()
	queue.PushBack(root)
	b.used[root] = true

	for queue.Len() > 0 {
		v := queue.Remove(queue.Front()).(int)
		for _, u := range b.graph[v] {
			if b.base[v] == b.base[u] || b.match[v] == u {
				continue
			}
			if u == root || (b.match[u] != -1 && b.parent[b.match[u]] != -1) {
				curBase := b.lca(v, u)
				for i := 0; i < b.n; i++ {
					b.blossom[i] = false
				}
				b.markPath(v, curBase, u)
				b.markPath(u, curBase, v)
				for i := 0; i < b.n; i++ {
					if b.blossom[b.base[i]] {
						b.base[i] = curBase
						if !b.used[i] {
							b.used[i] = true
							queue.PushBack(i)
						}
					}
				}
			} else if b.parent[u] == -1 {
				b.parent[u] = v
				if b.match[u] == -1 {
					return u
				}
				if !b.used[b.match[u]] {
					b.used[b.match[u]] = true
					queue.PushBack(b.match[u])
				}
			}
		}
	}
	return -1
}

func (b *Blossom) augmentPath(start int) {
	v := start
	for v != -1 {
		pv := b.parent[v]
		w := -1
		if pv != -1 {
			w = b.match[pv]
		}
		b.match[v] = pv
		b.match[pv] = v
		v = w
	}
}

func (b *Blossom) Solve() []int {
	for i := 0; i < b.n; i++ {
		if b.match[i] == -1 {
			if endpoint := b.findPath(i); endpoint != -1 {
				b.augmentPath(endpoint)
			}
		}
	}
	return b.match
}

func edmondsBlossomMatching(g *Graph) []int {
	solver := NewBlossom(g.n, g.adj)
	return solver.Solve()
}

func maximumMatching(g *Graph) []int {
	return edmondsBlossomMatching(g)
}

// func edmondsMaximumMatchingSize(g *Graph) int {
// 	m := maximumMatching(g)
// 	count := 0
// 	for _, v := range m {
// 		if v != -1 {
// 			count++
// 		}
// 	}
// 	return count / 2
// }

type FFactorInput struct {
	Original *Graph
	F        []int
}

func FindFFactor(input FFactorInput) (bool, [][2]int) {
	G := input.Original
	f := input.F

	n := G.n
	degree := make([]int, n)
	for u := 0; u < n; u++ {
		degree[u] = len(G.adj[u])
		if f[u] > degree[u] {
			return false, nil
		}
	}

	sMap := make([][]int, n)
	vertexCounter := 0
	edgeMap := make(map[int][2]int)

	for u := 0; u < n; u++ {
		sMap[u] = make([]int, len(G.adj[u]))
		for i, w := range G.adj[u] {
			sMap[u][i] = vertexCounter
			edgeMap[vertexCounter] = [2]int{u, w}
			vertexCounter++
		}
	}

	GStar := NewGraph(0)
	GStar.adj = make([][]int, vertexCounter)
	GStar.n = vertexCounter

	used := make([][]bool, n)
	for u := range used {
		used[u] = make([]bool, len(G.adj[u]))
	}

	for u := 0; u < n; u++ {
		for i, w := range G.adj[u] {
			if used[u][i] {
				continue
			}
			// найдём индекс обратного ребра
			var j int
			for k, v := range G.adj[w] {
				if v == u {
					j = k
					break
				}
			}
			su := sMap[u][i]
			sw := sMap[w][j]
			GStar.AddEdge(su, sw)
			used[u][i] = true
			used[w][j] = true
		}
	}

	// T(v)
	for v := 0; v < n; v++ {
		delta := degree[v] - f[v]
		for i := 0; i < delta; i++ {
			tv := vertexCounter
			vertexCounter++
			GStar.adj = append(GStar.adj, []int{})
			GStar.n++
			for _, sv := range sMap[v] {
				GStar.AddEdge(tv, sv)
			}
		}
	}

	matching := maximumMatching(GStar)

	matchedCount := 0
	for _, m := range matching {
		if m != -1 {
			matchedCount++
		}
	}
	if matchedCount != GStar.n {
		return false, nil
	}

	fFactorEdges := make([][2]int, 0)
	seen := make(map[[2]int]bool)

	for v, u := range matching {
		if u == -1 || v > u {
			continue
		}
		vu, ok1 := edgeMap[v]
		uv, ok2 := edgeMap[u]
		if ok1 && ok2 {
			// Убедимся, что это одно и то же ребро
			if vu[0] == uv[1] && vu[1] == uv[0] {
				e := [2]int{vu[0], vu[1]}
				if !seen[e] && !seen[[2]int{e[1], e[0]}] {
					fFactorEdges = append(fFactorEdges, e)
					seen[e] = true
				}
			}
		}
	}

	return true, fFactorEdges
}
