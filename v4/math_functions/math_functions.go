// math_functions.go

package math_functions

import (
	"fmt"
	"sort"
)

// Topological_sort performs a deterministic topological sort using Kahn's algorithm.
// Nodes with the same precedence are sorted alphabetically for consistent output.
func Topological_sort(graph map[string][]string) ([]string, error) {
	in_degree := make(map[string]int)
	for node := range graph {
		in_degree[node] = 0
	}
	for _, deps := range graph {
		for _, dep := range deps {
			in_degree[dep]++
		}
	}

	// Collect nodes with in-degree 0 and sort for deterministic order
	var zero_in_degree []string
	for node, degree := range in_degree {
		if degree == 0 {
			zero_in_degree = append(zero_in_degree, node)
		}
	}
	sort.Strings(zero_in_degree) // sort alphabetically
	queue := append([]string(nil), zero_in_degree...)

	var sorted []string
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		sorted = append(sorted, current)

		for _, neighbor := range graph[current] {
			in_degree[neighbor]--
		}

		// Add new zero-in-degree nodes, sort, and append to queue
		var newly_zero []string
		for _, neighbor := range graph[current] {
			if in_degree[neighbor] == 0 {
				newly_zero = append(newly_zero, neighbor)
			}
		}
		sort.Strings(newly_zero)
		queue = append(queue, newly_zero...)
	}

	if len(sorted) != len(graph) {
		return nil, fmt.Errorf("cycle detected: only sorted %d of %d nodes", len(sorted), len(graph))
	}

	return sorted, nil
}

// Reverse_topological_sort performs a deterministic topological sort and returns the reversed order.
// Useful for teardown operations or viewing leaf-to-root dependencies.
func Reverse_topological_sort(graph map[string][]string) ([]string, error) {
	sorted, err := Topological_sort(graph)
	if err != nil {
		return nil, err
	}

	// Reverse the sorted slice in-place
	for i, j := 0, len(sorted)-1; i < j; i, j = i+1, j-1 {
		sorted[i], sorted[j] = sorted[j], sorted[i]
	}

	return sorted, nil
}