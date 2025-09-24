package blossom

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadGraph(filename string) (map[int][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ReadGraphFromReader(file)
}

func ReadGraphFromReader(r io.Reader) (map[int][]int, error) {
	scanner := bufio.NewScanner(r)
	graph := make(map[int][]int)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		u, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("некорректный идентификатор вершины %q: %v", fields[0], err)
		}
		for _, fs := range fields[1:] {
			v, err := strconv.Atoi(fs)
			if err != nil {
				return nil, fmt.Errorf("некорректный идентификатор соседа %q: %v", fs, err)
			}
			graph[u] = append(graph[u], v)
		}
		if _, exists := graph[u]; !exists {
			graph[u] = []int{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return graph, nil
}
