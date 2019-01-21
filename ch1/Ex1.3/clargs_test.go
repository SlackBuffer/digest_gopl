package clargs

import "testing"

var args = [100000]string{}

func BenchmarkPrintArgs1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printArgs1(args[:])
	}
}
func BenchmarkPrintArgs2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printArgs2(args[:])
	}
}

// go test -bench=.
