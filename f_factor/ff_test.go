package ff

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func graphFromEdges(n int, edges [][2]int) *Graph {
	g := NewGraph(n)
	for _, e := range edges {
		g.AddEdge(e[0], e[1])
	}
	return g
}

func TestSimpleFFactor(t *testing.T) {
	g := graphFromEdges(4, [][2]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 0},
	})
	f := []int{1, 1, 1, 1}
	ok, _ := FindFFactor(FFactorInput{Original: g, F: f})
	if !ok {
		t.Errorf("Expected f-factor to exist, but got false")
	}
}

func TestNoFFactor(t *testing.T) {
	g := graphFromEdges(3, [][2]int{
		{0, 1}, {1, 2},
	})
	f := []int{1, 1, 1}
	ok, _ := FindFFactor(FFactorInput{Original: g, F: f})
	if ok {
		t.Errorf("Expected f-factor to not exist, but got true")
	}
}

func TestFFactor(t *testing.T) {
	g := graphFromEdges(8, [][2]int{
		{0, 1}, {0, 4}, {0, 6}, {1, 4}, {1, 2},
		{2, 5}, {2, 3}, {3, 3}, {3, 7}, {4, 5},
		{4, 5}, {4, 6}, {5, 7}, {6, 7},
	})

	f := []int{3, 3, 2, 3, 4, 2, 2, 1}
	ok, matching := FindFFactor(FFactorInput{Original: g, F: f})
	fmt.Println("Matching:", matching)
	fmt.Println("Result:", ok)
	if !ok {
		t.Errorf("Expected f-factor to not exist, but got true")
	}
}

func generateRandomGraph(n, m int) *Graph {
	g := NewGraph(n)
	edges := make(map[[2]int]bool)
	rand.Seed(time.Now().UnixNano())

	for len(edges) < m {
		u := rand.Intn(n)
		v := rand.Intn(n)
		if u != v && !edges[[2]int{u, v}] && !edges[[2]int{v, u}] {
			g.AddEdge(u, v)
			edges[[2]int{u, v}] = true
		}
	}
	return g
}

func BenchmarkFFactor(b *testing.B) {
	benchmarks := []struct {
		n     int  // число вершин
		dense bool // плотность: плотный или разреженный граф
	}{
		{50, false},
		{50, true},
		{100, false},
		{100, true},
		{500, false},
		{500, true},
	}

	for _, bm := range benchmarks {
		name := fmt.Sprintf("N=%d_Dense=%v", bm.n, bm.dense)
		b.Run(name, func(b *testing.B) {
			var m int
			if bm.dense {
				m = bm.n * (bm.n - 1) / 4
			} else {
				m = bm.n * 2
			}

			for i := 0; i < b.N; i++ {
				g := generateRandomGraph(bm.n, m)
				f := make([]int, bm.n)
				for i := range f {
					f[i] = rand.Intn(3)
				}
				FindFFactor(FFactorInput{Original: g, F: f})
			}
		})
	}
}
