package gedecomp

import (
	"reflect"
	"testing"
)

func graphFromEdges(n int, edges [][2]int) *Graph {
	g := NewGraph(n)
	for _, edge := range edges {
		g.AddEdge(edge[0], edge[1])
	}
	return g
}

// TestCycleGraph проверяет граф-цикл из 6 вершин.
// Ожидается: размер паросочетания = 3, разложение: D = {}, A = {}, C = {0,1,2,3,4,5}.
func TestCycleGraph(t *testing.T) {
	edges := [][2]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 0},
	}
	g := graphFromEdges(6, edges)

	size := edmondsMaximumMatchingSize(g)
	if size != 3 {
		t.Errorf("TestCycleGraph: Expected matching size 3, got %d", size)
	}

	// Ожидаемые множества
	expectedD := map[int]bool{}
	expectedA := map[int]bool{}
	expectedC := map[int]bool{
		0: true, 1: true, 2: true, 3: true, 4: true, 5: true,
	}

	D, A, C := gallaiEdmondsDecomposition(g)
	if !reflect.DeepEqual(D, expectedD) {
		t.Errorf("TestCycleGraph: Expected D = %v, got %v", expectedD, D)
	}
	if !reflect.DeepEqual(A, expectedA) {
		t.Errorf("TestCycleGraph: Expected A = %v, got %v", expectedA, A)
	}
	if !reflect.DeepEqual(C, expectedC) {
		t.Errorf("TestCycleGraph: Expected C = %v, got %v", expectedC, C)
	}
}

// TestGraph7 проверяет граф из 7 вершин.
// Ожидается: размер паросочетания = 3, разложение: D = {0,2,4,6}, A = {1,3,5}, C = {}.
func TestGraph7(t *testing.T) {
	edges := [][2]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6},
	}
	g := graphFromEdges(7, edges)

	size := edmondsMaximumMatchingSize(g)
	if size != 3 {
		t.Errorf("TestGraph7: Expected matching size 3, got %d", size)
	}

	expectedD := map[int]bool{0: true, 2: true, 4: true, 6: true}
	expectedA := map[int]bool{1: true, 3: true, 5: true}
	expectedC := map[int]bool{}

	D, A, C := gallaiEdmondsDecomposition(g)
	if !reflect.DeepEqual(D, expectedD) {
		t.Errorf("TestGraph7: Expected D = %v, got %v", expectedD, D)
	}
	if !reflect.DeepEqual(A, expectedA) {
		t.Errorf("TestGraph7: Expected A = %v, got %v", expectedA, A)
	}
	if !reflect.DeepEqual(C, expectedC) {
		t.Errorf("TestGraph7: Expected C = %v, got %v", expectedC, C)
	}
}

// TestComplexGraph проверяет сложный граф с 23 вершинами.
// Ожидается: размер паросочетания = 11 (покрываются 22 вершины),
// разложение: D = {0,1,...,11}, A = {12,13,14}, C = {15,...,22}.
func TestComplexGraph(t *testing.T) {
	edges := [][2]int{
		{0, 1}, {0, 2}, {0, 14},
		{1, 2}, {1, 12},
		{2, 12},
		{3, 4}, {3, 5}, {3, 12},
		{4, 5}, {4, 13},
		{5, 13},
		{6, 7}, {6, 9}, {6, 10},
		{7, 8}, {7, 12},
		{8, 9}, {8, 10}, {8, 14},
		{9, 10}, {9, 14},
		{11, 13}, {11, 14},
		{12, 16}, {12, 19},
		{13, 14}, {13, 18},
		{14, 21}, {14, 22},
		{15, 16}, {15, 18}, {15, 20},
		{16, 17},
		{17, 18}, {17, 19},
		{19, 20},
		{21, 22},
	}
	g := graphFromEdges(23, edges)

	size := edmondsMaximumMatchingSize(g)
	if size != 11 {
		t.Errorf("TestComplexGraph: Expected matching size 11, got %d", size)
	}

	expectedD := map[int]bool{}
	for i := 0; i <= 11; i++ {
		expectedD[i] = true
	}
	expectedA := map[int]bool{12: true, 13: true, 14: true}
	expectedC := map[int]bool{}
	for i := 15; i <= 22; i++ {
		expectedC[i] = true
	}

	D, A, C := gallaiEdmondsDecomposition(g)
	if !reflect.DeepEqual(D, expectedD) {
		t.Errorf("TestComplexGraph: Expected D = %v, got %v", expectedD, D)
	}
	if !reflect.DeepEqual(A, expectedA) {
		t.Errorf("TestComplexGraph: Expected A = %v, got %v", expectedA, A)
	}
	if !reflect.DeepEqual(C, expectedC) {
		t.Errorf("TestComplexGraph: Expected C = %v, got %v", expectedC, C)
	}
}
