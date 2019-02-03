package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func anagram(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	if sortStr(s1) == sortStr(s2) {
		return true
	}

	return false
}

func sortStr(s string) string {
	var sli []string
	for _, v := range []byte(s) {
		sli = append(sli, string(v))
	}
	sort.Strings(sli)
	return strings.Join(sli, "")
}

func main() {
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		s := strings.Split(string(strings.TrimSpace(input.Text())), " ")
		fmt.Println(anagram(s[0], s[1]))
	}
}
