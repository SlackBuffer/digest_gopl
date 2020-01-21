// topological sorting. the prerequisite information forms a directed graph with a node for each course and edges
// from each course to courses that it depends on. the graph is acyclic: there's no path from a course that leads
// back to itself
package main

import (
	"fmt"
	"sort"
)

// maps computer science courses to their prerequisites.
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

// depth-first search
// stack
func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)

	var visitAll func(items []string)
	visitAll = func(items []string) {
		fmt.Printf("%#v\n", items)
		for _, item := range items {
			fmt.Printf("Push stack ")
			if !seen[item] {
				seen[item] = true
				fmt.Println("$$$$$$$$$$$$$$$$", "saw this time:", item)
				visitAll(m[item])
				order = append(order, item)
				fmt.Println("Pop stack********************", "append:", item)
			} else {
				fmt.Println("Pop stack----------------", "seen already:", item)
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	// sort keys alphabetically
	sort.Strings(keys) // order matters
	// fmt.Printf("%#v\n\n", keys)

	visitAll(keys)
	return order
}

func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}

	fmt.Println()
	for i, course := range topoSort2(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

// go run *.go
