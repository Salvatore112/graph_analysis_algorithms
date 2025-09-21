package multigraph_algo

// hopcroftKarp принимает adjacency: left -> []right и список всех left-вершин (sorted, возможно с нулевой степенью).
// Возвращает pairU (left->right) и pairV (right->left). Если вершина не в паре — значение "".
func hopcroftKarp(adj map[string][]string, leftVerts []string) (map[string]string, map[string]string) {
	pairU := make(map[string]string)
	pairV := make(map[string]string)
	dist := make(map[string]int)

	rightVertsSet := make(map[string]bool)
	for _, nbrs := range adj {
		for _, v := range nbrs {
			rightVertsSet[v] = true
		}
	}
	for _, u := range leftVerts {
		pairU[u] = ""
	}
	for v := range rightVertsSet {
		pairV[v] = ""
	}

	const INF = 1 << 30

	bfs := func() bool {
		queue := make([]string, 0)
		for _, u := range leftVerts {
			if pairU[u] == "" {
				dist[u] = 0
				queue = append(queue, u)
			} else {
				dist[u] = INF
			}
		}
		foundAug := false
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range adj[u] {
				pu := pairV[v]
				if pu == "" {
					foundAug = true
				} else {
					if dist[pu] == INF {
						dist[pu] = dist[u] + 1
						queue = append(queue, pu)
					}
				}
			}
		}
		return foundAug
	}

	var dfs func(u string) bool
	dfs = func(u string) bool {
		for _, v := range adj[u] {
			pu := pairV[v]
			if pu == "" || (dist[pu] == dist[u]+1 && dfs(pu)) {
				pairU[u] = v
				pairV[v] = u
				return true
			}
		}
		dist[u] = INF
		return false
	}

	for bfs() {
		for _, u := range leftVerts {
			if pairU[u] == "" {
				dfs(u)
			}
		}
	}

	return pairU, pairV
}
