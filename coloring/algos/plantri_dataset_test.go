package algos_test

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/graph6"

	algo "github.com/Salvatore112/graph_analysis_algorithms/coloring/algos"
)

func g6ToAdj(g graph6.Graph) [][]int {
	nodes := graph.NodesOf(g.Nodes())
	n := len(nodes)
	adj := make([][]int, n)
	for _, v := range nodes {
		vi := int(v.ID())
		for _, u := range graph.NodesOf(g.From(v.ID())) {
			ui := int(u.ID())
			if ui != vi {
				adj[vi] = append(adj[vi], ui)
			}
		}
	}
	return adj
}

func TestPlantriDataset(t *testing.T) {
	dsDir := filepath.Join("..", "dataset")
	info, err := os.Stat(dsDir)
	if err != nil || !info.IsDir() {
		t.Skipf("нет каталога с датасетом (%s). Сначала запустите: coloring/generate_plantri_dataset.sh", dsDir)
	}

	matches, err := filepath.Glob(filepath.Join(dsDir, "*.g6"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Skipf("в %s нет файлов *.g6. Сначала сгенерируйте датасет", dsDir)
	}

	total := 0
	for _, file := range matches {
		f, err := os.Open(file)
		if err != nil {
			t.Fatalf("open %s: %v", file, err)
		}
		sc := bufio.NewScanner(f)

		lineNo := 0
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" {
				continue
			}

			g := graph6.Graph(line)
			if !graph6.IsValid(g) {
				t.Fatalf("invalid graph6 in %s at line %d", file, lineNo+1)
			}
			adj := g6ToAdj(g)

			cols := algo.FourColor(adj)
			if cols == nil {
				t.Fatalf("no 4-coloring found for %s at line %d", file, lineNo+1)
			}
			for v := range adj {
				if cols[v] < 0 || cols[v] > 3 {
					t.Fatalf("bad color for v=%d in %s line %d", v, file, lineNo+1)
				}
				for _, u := range adj[v] {
					if cols[u] == cols[v] {
						t.Fatalf("conflict %d-%d in %s line %d", v, u, file, lineNo+1)
					}
				}
			}
			total++
			lineNo++
		}
		if err := sc.Err(); err != nil {
			t.Fatalf("scan %s: %v", file, err)
		}
		_ = f.Close()
	}

	if total == 0 {
		t.Skip("dataset пуст. Сгенерируйте графы plantri скриптом.")
	}
}
