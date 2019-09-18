// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package memo_test

import (
	memo "digest_gopl/ch9/memo5"
	"testing"

	"digest_gopl/ch9/memotest"
)

var httpGetBody = memotest.HTTPGetBody

func Test(t *testing.T) {
	m := memo.New(httpGetBody)
	defer m.Close()
	memotest.Sequential(t, m)
}

// NOTE: not concurrency-safe!  Test fails.
func TestConcurrent(t *testing.T) {
	m := memo.New(httpGetBody)
	defer m.Close()
	memotest.Concurrent(t, m)
}

// go test -v exercises-the_go_programming_language/ch9/memo5
// go test -run=TestConcurrent -race -v exercises-the_go_programming_language/ch9/memo5
