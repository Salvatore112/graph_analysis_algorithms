package blossom

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func equalMatchings(a, b [][2]int) bool {
	if len(a) != len(b) {
		return false
	}
	norm := func(m [][2]int) [][2]int {
		res := make([][2]int, len(m))
		for i, p := range m {
			u, v := p[0], p[1]
			if u > v {
				u, v = v, u
			}
			res[i] = [2]int{u, v}
		}
		sort.Slice(res, func(i, j int) bool {
			if res[i][0] != res[j][0] {
				return res[i][0] < res[j][0]
			}
			return res[i][1] < res[j][1]
		})
		return res
	}
	na := norm(a)
	nb := norm(b)
	return reflect.DeepEqual(na, nb)
}

func TestEmptyGraph(t *testing.T) {
	graph := map[int][]int{}
	match := MaxMatching(graph)
	if len(match) != 0 {
		t.Errorf("ожидалось пустое паросочетание, получили %v", match)
	}
}

func TestSingleEdge(t *testing.T) {
	graph := map[int][]int{
		1: {2},
		2: {1},
	}
	expected := [][2]int{{1, 2}}
	match := MaxMatching(graph)
	if !equalMatchings(match, expected) {
		t.Errorf("ожидалось %v, получили %v", expected, match)
	}
}

func TestOddCycleTriangle(t *testing.T) {
	graph := map[int][]int{
		0: {1, 2},
		1: {0, 2},
		2: {0, 1},
	}
	match := MaxMatching(graph)
	if len(match) != 1 {
		t.Errorf("в треугольнике ожидается 1 ребро в паросочетании, получили %d: %v", len(match), match)
	}
}

func TestEvenCycleSquare(t *testing.T) {
	graph := map[int][]int{
		0: {1, 3},
		1: {0, 2},
		2: {1, 3},
		3: {0, 2},
	}
	match := MaxMatching(graph)
	if len(match) != 2 {
		t.Errorf("в квадрате ожидается 2 ребра в паросочетании, получили %d: %v", len(match), match)
	}
}

func TestChainGraph(t *testing.T) {
	graph := map[int][]int{
		0: {1},
		1: {0, 2},
		2: {1, 3},
		3: {2},
	}
	match := MaxMatching(graph)
	if len(match) != 2 {
		t.Errorf("в цепочке ожидается 2 ребра в паросочетании, получили %d: %v", len(match), match)
	}
}

func TestDisconnected(t *testing.T) {
	graph := map[int][]int{
		0: {1},
		1: {0},
		2: {3},
		3: {2},
	}
	expected := [][2]int{{0, 1}, {2, 3}}
	match := MaxMatching(graph)
	if !equalMatchings(match, expected) {
		t.Errorf("ожидалось %v, получили %v", expected, match)
	}
}

func TestMixedGraph(t *testing.T) {
	graph := map[int][]int{
		0: {1, 2},
		1: {0, 2},
		2: {0, 1},
		3: {4},
		4: {3, 5},
		5: {4},
	}
	match := MaxMatching(graph)
	if len(match) != 2 {
		t.Errorf("ожидается 2 ребра в смешанном графе, получили %d: %v", len(match), match)
	}
}

func TestReadGraph(t *testing.T) {
	data := `
# пример графа
0 1 2
1 0 2
2 0 1 3
3 2
`
	reader := strings.NewReader(data)
	graph, err := ReadGraphFromReader(reader)
	if err != nil {
		t.Fatalf("ошибка при чтении графа: %v", err)
	}
	expected := map[int][]int{
		0: {1, 2},
		1: {0, 2},
		2: {0, 1, 3},
		3: {2},
	}
	for u, nbrs := range expected {
		got := graph[u]
		sort.Ints(got)
		sort.Ints(nbrs)
		if !reflect.DeepEqual(got, nbrs) {
			t.Errorf("для вершины %d ожидались соседи %v, получили %v", u, nbrs, got)
		}
	}
}
