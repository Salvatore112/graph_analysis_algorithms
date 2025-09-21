package multigraph_algo

import (
	"math/bits"
	"sort"

	"github.com/Salvatore112/graph_analysis_algorithms/graphs"
)

// ExactEdgeColoring находит ИМЕННО минимальное число цветов для рёберной раскраски.
// Возвращает (colorsUsed, edges, colorsMap).
// Внутри: итеративный поиск по k = Δ..UB, где UB берётся из жадной раскраски (если она есть)
// или из оценки Шеннона (⌊3Δ/2⌋) как запасной вариант.
func ExactEdgeColoring(g *graphs.MultiGraph) (int, []Edge, map[int]int) {
	edges := EdgeList(g)
	if len(edges) == 0 {
		return 0, edges, map[int]int{}
	}

	type it struct {
		idx int
		w   int
	}
	order := make([]it, 0, len(edges))
	for _, e := range edges {
		w := Degree(g, e.U) + Degree(g, e.V)
		order = append(order, it{idx: e.ID, w: w})
	}
	sort.Slice(order, func(i, j int) bool {
		if order[i].w == order[j].w {
			return order[i].idx < order[j].idx
		}
		return order[i].w > order[j].w
	})

	pos2eid := make([]int, len(order))
	eid2pos := make(map[int]int, len(order))
	for i, it := range order {
		pos2eid[i] = it.idx
		eid2pos[it.idx] = i
	}

	edgeByID := make(map[int]Edge, len(edges))
	for _, e := range edges {
		edgeByID[e.ID] = e
	}

	delta := MaxDegree(g)
	lb := delta

	ubGreedy, _, greedyColors := GreedyEdgeColoring(g)
	ub := ubGreedy
	if ub <= 0 {
		if delta == 0 {
			ub = 0
		} else {
			ub = (3 * delta) / 2
		}
	}

	for k := lb; k <= ub; k++ {
		if ok, sol := colorWithK(edgeByID, pos2eid, k); ok {
			return k, edges, sol
		}
	}

	return ubGreedy, edges, greedyColors
}

// ------------------------------
// ВНУТРЕННОСТИ ПОИСКА С ОГРАНИЧЕНИЕМ K
// ------------------------------

type solverK interface {
	tryAssign(u, v string, c int) bool
	unassign(u, v string, c int)
	availMask(u, v string) uint64 // только если k<=64
	availCount(u, v string, k int) int
}

// bitsetSolver — быстрый путь (k <= 64): у каждой вершины битовая маска занятых цветов.
// used[v] — битовая маска, где 1 означает "цвет занят у вершины v".
type bitsetSolver struct {
	kMask    uint64
	used     map[string]uint64
	colors   map[int]int
	edgeByID map[int]Edge
}

func newBitsetSolver(k int, edgeByID map[int]Edge) *bitsetSolver {
	return &bitsetSolver{
		kMask:    (uint64(1) << uint(k)) - 1,
		used:     make(map[string]uint64),
		colors:   make(map[int]int),
		edgeByID: edgeByID,
	}
}

func (s *bitsetSolver) tryAssign(u, v string, c int) bool {
	b := uint64(1) << uint(c)
	if (s.used[u]&b) != 0 || (s.used[v]&b) != 0 {
		return false
	}
	s.used[u] |= b
	s.used[v] |= b
	return true
}

func (s *bitsetSolver) unassign(u, v string, c int) {
	b := ^(uint64(1) << uint(c))
	s.used[u] &= b
	s.used[v] &= b
}

func (s *bitsetSolver) availMask(u, v string) uint64 {
	return ^(s.used[u] | s.used[v]) & s.kMask
}

func (s *bitsetSolver) availCount(u, v string, _ int) int {
	return bits.OnesCount64(s.availMask(u, v))
}

// sliceSolver — общий путь (k > 64): у каждой вершины массив занятости цветов.
type sliceSolver struct {
	k        int
	used     map[string][]bool
	colors   map[int]int
	edgeByID map[int]Edge
}

func newSliceSolver(k int, edgeByID map[int]Edge) *sliceSolver {
	return &sliceSolver{
		k:        k,
		used:     make(map[string][]bool),
		colors:   make(map[int]int),
		edgeByID: edgeByID,
	}
}

func (s *sliceSolver) ensure(v string) {
	if s.used[v] == nil {
		s.used[v] = make([]bool, s.k)
	}
}

func (s *sliceSolver) tryAssign(u, v string, c int) bool {
	s.ensure(u)
	s.ensure(v)
	if s.used[u][c] || s.used[v][c] {
		return false
	}
	s.used[u][c] = true
	s.used[v][c] = true
	return true
}

func (s *sliceSolver) unassign(u, v string, c int) {
	s.used[u][c] = false
	s.used[v][c] = false
}

func (s *sliceSolver) availMask(_ string, _ string) uint64 {
	return 0
}

func (s *sliceSolver) availCount(u, v string, k int) int {
	s.ensure(u)
	s.ensure(v)
	cnt := 0
	for c := 0; c < k; c++ {
		if !s.used[u][c] && !s.used[v][c] {
			cnt++
		}
	}
	return cnt
}

// colorWithK пробует раскрасить все рёбра в k цветов.
// Использует backtracking + MRV (выбираем нераскрашенное ребро с минимальным числом доступных цветов).
func colorWithK(edgeByID map[int]Edge, pos2eid []int, k int) (bool, map[int]int) {

	var sol map[int]int
	var sv solverK
	if k <= 64 {
		sv = newBitsetSolver(k, edgeByID)
		sol = sv.(*bitsetSolver).colors
	} else {
		sv = newSliceSolver(k, edgeByID)
		sol = sv.(*sliceSolver).colors
	}

	N := len(pos2eid)

	nextEdge := func() (pos int, avail []int, ok bool) {
		bestPos := -1
		bestCnt := 1 << 30
		var bestAvail []int

		for p := 0; p < N; p++ {
			eid := pos2eid[p]
			if _, colored := sol[eid]; colored {
				continue
			}
			e := edgeByID[eid]

			if k <= 64 {
				mask := sv.availMask(e.U, e.V)
				cnt := bits.OnesCount64(mask)
				if cnt == 0 {
					return -1, nil, false
				}
				if cnt < bestCnt {
					bestCnt = cnt
					bestPos = p
					bestAvail = bestAvail[:0]
					for c := 0; c < k; c++ {
						if (mask>>uint(c))&1 == 1 {
							bestAvail = append(bestAvail, c)
						}
					}
					if bestCnt == 1 {
						break
					}
				}
			} else {
				cnt := sv.availCount(e.U, e.V, k)
				if cnt == 0 {
					return -1, nil, false
				}
				if cnt < bestCnt {
					bestCnt = cnt
					bestPos = p
					tmp := make([]int, 0, cnt)
					for c := 0; c < k; c++ {
						if s, ok := sv.(*sliceSolver); ok {
							s.ensure(e.U)
							s.ensure(e.V)
							if !s.used[e.U][c] && !s.used[e.V][c] {
								tmp = append(tmp, c)
							}
						}
					}
					bestAvail = tmp
					if bestCnt == 1 {
						break
					}
				}
			}
		}
		if bestPos == -1 {
			return -1, nil, false
		}
		return bestPos, bestAvail, true
	}

	var dfs func(colored int) bool
	dfs = func(colored int) bool {
		if colored == N {
			return true
		}
		pos, avail, ok := nextEdge()
		if !ok {
			return false
		}
		eid := pos2eid[pos]
		e := edgeByID[eid]

		for _, c := range avail {
			if !sv.tryAssign(e.U, e.V, c) {
				continue
			}
			sol[eid] = c

			if dfs(colored + 1) {
				return true
			}

			delete(sol, eid)
			sv.unassign(e.U, e.V, c)
		}
		return false
	}

	ok := dfs(0)
	return ok, sol
}
