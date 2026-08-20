package rpn

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestHappyPath(t *testing.T) {
	// (3 + 4) * 2 = 14
	result, err := Evaluate("3 4 + 2 *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(result, 14) {
		t.Errorf("got %g, want 14", result)
	}
}

func TestPrecedenceViaRPN(t *testing.T) {
	// Standard infix: 5 + ((1 + 2) * 4) - 3 = 14, expressed in RPN
	result, err := Evaluate("5 1 2 + 4 * + 3 -")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(result, 14) {
		t.Errorf("got %g, want 14", result)
	}
}

func TestAllOperators(t *testing.T) {
	cases := []struct {
		name string
		expr string
		want float64
	}{
		{"add", "2 3 +", 5},
		{"subtract", "5 3 -", 2},
		{"multiply", "4 3 *", 12},
		{"divide", "10 4 /", 2.5},
		{"power", "2 10 ^", 1024},
		{"sqrt", "16 sqrt", 4},
		{"sin", "0 sin", 0},
		{"cos", "0 cos", 1},
		{"tan", "0 tan", 0},
		{"log", "1 log", 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Evaluate(tc.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(got, tc.want) {
				t.Errorf("Evaluate(%q) = %g, want %g", tc.expr, got, tc.want)
			}
		})
	}
}

func TestStackUnderflow(t *testing.T) {
	_, err := Evaluate("1 +")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDivisionByZero(t *testing.T) {
	_, err := Evaluate("5 0 /")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUnknownOperator(t *testing.T) {
	_, err := Evaluate("3 4 foo")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNonNumericToken(t *testing.T) {
	_, err := Evaluate("3 abc +")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestEmptyExpression(t *testing.T) {
	_, err := Evaluate("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTooManyOperandsLeft(t *testing.T) {
	_, err := Evaluate("1 2 3 +")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSqrtNegative(t *testing.T) {
	_, err := Evaluate("-4 sqrt")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLogNonPositive(t *testing.T) {
	_, err := Evaluate("0 log")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
