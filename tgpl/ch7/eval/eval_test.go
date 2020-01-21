package eval

import (
	"fmt"
	"math"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}, "1729"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
		{"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
		{"5 / 9 * (F - 32)", Env{"F": 32}, "0"},
		{"5 / 9 * (F - 32)", Env{"F": 212}, "100"},
	}
	var prevExpr string
	for _, test := range tests {
		// print expr only when it changes
		if test.expr != prevExpr {
			fmt.Printf("\n%s\n", test.expr)
			prevExpr = test.expr
		}

		// Parse handles the recursive part
		expr, err := Parse(test.expr)
		if err != nil {
			t.Error(err) // parse error
			continue
		}
		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		fmt.Printf("\t%v => %s\n", test.env, got)
		if got != test.want {
			// reports an error
			t.Errorf("%s.Eval() in %v = %q, want %q\n", test.expr, test.env, got, test.want)
		}
	}
}

// go test -v digest_gopl/ch7/eval

/*
=== RUN   TestCoverage
--- PASS: TestCoverage (0.00s)
=== RUN   TestEval

sqrt(A / pi)
        map[A:87616 pi:3.141592653589793] => 167

pow(x, 3) + pow(y, 3)
        map[x:12 y:1] => 1729
        map[x:9 y:10] => 1729

5 / 9 * (F - 32)
        map[F:-40] => -40
        map[F:32] => 0
        map[F:212] => 100
--- PASS: TestEval (0.00s)
PASS
ok      digest_gopl/ch7/eval    0.008s
*/
