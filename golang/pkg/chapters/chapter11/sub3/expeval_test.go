package sub3_test

import (
	"fmt"
	evaluator "golang/pkg/chapters/chapter7"
	"math"
	"testing"
)

func TestCoverage(t *testing.T) {
	const (
		errTemplate              = "%s: got %q, want %q"
		unexpectedResultTemplate = "%s: %v => %s, want %s"
	)

	var (
		tests = []struct {
			input string
			env   evaluator.Environment
			// Expected error from Parse/Check or result from Eval
			want string
		}{
			{"x % 2", nil, "unexpected '%'"},
			{"!true", nil, "unexpected '!'"},
			{"log(10)", nil, `unknown function "log"`},
			{"sqrt(1, 2)", nil, "call to sqrt has 2 args, want 1"},
			{"sqrt(A / pi)", evaluator.Environment{"A": 87616, "pi": math.Pi}, "167"},
			{"pow(x, 3) + pow(y, 3)", evaluator.Environment{"x": 9, "y": 10}, "1729"},
			{"5 / 9 * (F - 32)", evaluator.Environment{"F": -40}, "-40.000000"},
		}

		curExpr evaluator.Expression
		curErr  error
	)

	// ops
	for _, test := range tests {
		curExpr, curErr = evaluator.Parse(test.input)

		if curErr == nil {
			curErr = curExpr.Check(make(map[evaluator.Variable]evaluator.Empty))
		}

		if curErr != nil && curErr.Error() != test.want {
			t.Errorf(errTemplate, test.input, curErr, test.want)
		} else {
			continue
		}

		if got := fmt.Sprintf("%.6g", curExpr.Eval(test.env)); got != test.want {
			t.Errorf(unexpectedResultTemplate, test.input, test.env, got, test.want)
		}

	}
}
