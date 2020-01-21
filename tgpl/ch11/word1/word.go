// Package word provides utilities for word games
package word

// IsPalindrome reports whether s reads the same forward and backward.
func IsPalindrome(s string) bool {
	// 1. use bytes sequences here, not rune sequences
	// 2. not ignoring spaces here
	for i := range s {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}
