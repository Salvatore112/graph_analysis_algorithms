package blossom

func MaxMatching(graph map[int][]int) [][2]int {
	var nodes []int
	nodeIndex := make(map[int]int)
	for u := range graph {
		nodes = append(nodes, u)
	}

	for _, nbrs := range graph {
		nodes = append(nodes, nbrs...)
	}

	uniq := make(map[int]bool)
	for _, u := range nodes {
		uniq[u] = true
	}
	nodes = nodes[:0]
	for u := range uniq {
		nodes = append(nodes, u)
	}

	for i, u := range nodes {
		nodeIndex[u] = i
	}
	n := len(nodes)

	adj := make([][]int, n)
	for u, nbrs := range graph {
		ui := nodeIndex[u]
		for _, v := range nbrs {
			if idx, ok := nodeIndex[v]; ok {
				adj[ui] = append(adj[ui], idx)
			}
		}
	}

	for u := 0; u < n; u++ {
		for _, v := range adj[u] {
			found := false
			for _, w := range adj[v] {
				if w == u {
					found = true
					break
				}
			}
			if !found {
				adj[v] = append(adj[v], u)
			}
		}
	}

	spouse := make([]int, n)
	next := make([]int, n)
	belong := make([]int, n)
	mark := make([]int, n)
	visited := make([]int, n)
	queue := make([]int, n)
	lcaTimer := 0

	for i := 0; i < n; i++ {
		spouse[i] = -1
	}

	var findb func(int) int
	findb = func(a int) int {
		if belong[a] == a {
			return a
		}
		belong[a] = findb(belong[a])
		return belong[a]
	}

	union := func(a, b int) {
		ra := findb(a)
		rb := findb(b)
		if ra != rb {
			belong[ra] = rb
		}
	}

	findLCA := func(a, b int) int {
		lcaTimer++
		x, y := a, b
		for {
			if x != -1 {
				x = findb(x)
				if visited[x] == lcaTimer {
					return x
				}
				visited[x] = lcaTimer
				if spouse[x] != -1 {
					x = next[spouse[x]]
				} else {
					x = -1
				}
			}
			x, y = y, x
		}
	}

	goup := func(a, p int, queue *[]int, tail *int) {
		for findb(a) != p {
			b := spouse[a]
			c := next[b]
			if findb(c) != p {
				next[c] = b
			}
			if mark[b] == 2 {
				mark[b] = 1
				(*queue)[*tail] = b
				(*tail)++
			}
			if mark[c] == 2 {
				mark[c] = 1
				(*queue)[*tail] = c
				(*tail)++
			}
			union(a, b)
			union(b, c)
			a = c
		}
	}

	for start := 0; start < n; start++ {
		if spouse[start] != -1 {
			continue
		}
		for i := 0; i < n; i++ {
			belong[i] = i
			next[i] = -1
			mark[i] = 0
			visited[i] = 0
		}
		head, tail := 0, 0
		queue[tail] = start
		tail++
		mark[start] = 1

		augmented := false
		for head < tail && !augmented {
			x := queue[head]
			head++
			for _, y := range adj[x] {
				if spouse[x] == y || findb(x) == findb(y) || mark[y] == 2 {
					continue
				}
				if mark[y] == 1 {
					p := findLCA(x, y)
					if findb(x) != p {
						next[x] = y
					}
					if findb(y) != p {
						next[y] = x
					}
					goup(x, p, &queue, &tail)
					goup(y, p, &queue, &tail)
				} else if spouse[y] == -1 {
					next[y] = x
					v := y
					for v != -1 {
						u := next[v]
						w := spouse[u]
						spouse[v] = u
						spouse[u] = v
						v = w
					}
					augmented = true
					break
				} else {
					next[y] = x
					mark[spouse[y]] = 1
					queue[tail] = spouse[y]
					tail++
					mark[y] = 2
				}
			}
		}
	}

	used := make([]bool, n)
	var matching [][2]int
	for i := 0; i < n; i++ {
		j := spouse[i]
		if j != -1 && !used[i] && !used[j] {
			matching = append(matching, [2]int{nodes[i], nodes[j]})
			used[i] = true
			used[j] = true
		}
	}
	return matching
}
