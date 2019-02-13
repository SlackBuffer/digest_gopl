package main

import "fmt"

// remove adjacent duplicates in a []string slice
func remove(s []string) []string {
	result := s[:0]
	for _, v := range s {
		if len(result) == 0 {
			result = append(result, v)
		} else {
			if result[len(result)-1] != v {
				result = append(result, v)
			}
		}
	}
	return result
}

func main() {
	s := []string{"a", "b", "b", "c", "c", "c"}
	fmt.Println(remove(s))
}
