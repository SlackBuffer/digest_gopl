package main

import (
	"fmt"
	"net/url"
)

/*
// net/url
// maps a string key to a list of values
type Values map[string][]string

// returns the first value associated with the given key,
// or "" if there are none
func (v Values) Get(key string) string {
	if vs := v[key]; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (v Values) Add(key, value string) {
	v[key] = append(v[key], value)
}
*/

func main() {
	m := url.Values{"lang": {"en"}} // direct construction
	m.Add("item", "1")
	m.Add("item", "2")

	fmt.Println(m.Get("lang"))
	fmt.Println(m.Get("q"))
	fmt.Println(m.Get("item")) // (first value)
	fmt.Println(m["item"])     // (direct map access)

	m = nil
	fmt.Println(m.Get("item"))
	m.Add("item", "3") // panic
}
