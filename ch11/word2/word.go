// Package word provides utilities for word games
package word

import "unicode"

// IsPalindrome reports whether s reads the same forward and backward.
// Letter case is ignored, as are non-letters
func IsPalindrome(s string) bool {
	// var letters []rune
	// pre-allocate a sufficiently large array for use by letters, rather than expand it by successive calls to append

	// Consolidating the allocations in a single call to make eliminated 75% of the allocations and halved the quantity of allocated memory
	letters := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToLower(r))
		}
	}

	// use rune sequences
	n := len(letters) / 2
	for i := 0; i < n; i++ {
		if letters[i] != letters[len(letters)-1-i] {
			return false
		}
	}
	return true
}
