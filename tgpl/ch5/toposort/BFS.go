package main

import (
	"sort"
)

// breadth-first search
// queue
func topoSort2(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)

	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				order = append(order, item)

				for _, it := range m[item] {
					if !seen[it] {
						seen[it] = true
						order = append(order, it)
					}
				}
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
