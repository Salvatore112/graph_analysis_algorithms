package algos

type Adj map[string]map[string]struct{}

func cloneAdj(m map[string][]string) Adj {
	res := make(Adj, len(m))
	for v, lst := range m {
		if res[v] == nil {
			res[v] = map[string]struct{}{}
		}
		for _, u := range lst {
			if res[u] == nil {
				res[u] = map[string]struct{}{}
			}
			res[v][u] = struct{}{}
			res[u][v] = struct{}{}
		}
	}
	for v := range res {
		delete(res[v], v)
	}
	return res
}
