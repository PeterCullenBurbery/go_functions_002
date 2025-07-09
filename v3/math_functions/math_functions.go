// math_functions.go

package math_functions

import (
	"fmt"
)

// Topological_sort performs a topological sort on a DAG using Kahn's algorithm.
// Returns the ordered list of tasks, or an error if a cycle is detected.
func Topological_sort(graph map[string][]string) ([]string, error) {
	in_degree := make(map[string]int)

	// Count in-degrees
	for node := range graph {
		in_degree[node] = 0
	}
	for _, deps := range graph {
		for _, dep := range deps {
			in_degree[dep]++
		}
	}

	// Start with all 0 in-degree nodes
	var queue []string
	for node, degree := range in_degree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	var sorted []string
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		sorted = append(sorted, current)

		for _, neighbor := range graph[current] {
			in_degree[neighbor]--
			if in_degree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if len(sorted) != len(graph) {
		return nil, fmt.Errorf("cycle detected: only sorted %d of %d nodes", len(sorted), len(graph))
	}

	return sorted, nil
}

// Reverse_topological_sort performs a topological sort and returns the reversed order.
// Useful for teardown operations or viewing leaf-to-root dependencies.
func Reverse_topological_sort(graph map[string][]string) ([]string, error) {
	sorted, err := Topological_sort(graph)
	if err != nil {
		return nil, err
	}

	// Reverse the sorted slice
	for i, j := 0, len(sorted)-1; i < j; i, j = i+1, j-1 {
		sorted[i], sorted[j] = sorted[j], sorted[i]
	}

	return sorted, nil
}