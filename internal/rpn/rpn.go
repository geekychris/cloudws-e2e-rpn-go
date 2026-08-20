// Package rpn implements a scientific calculator that evaluates
// space-separated Reverse Polish Notation (RPN) expressions.
package rpn

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// binaryOps maps a token to a function of two operands (a, b) where a was
// pushed before b, i.e. the expression "a b OP" computes a OP b.
var binaryOps = map[string]func(a, b float64) (float64, error){
	"+": func(a, b float64) (float64, error) { return a + b, nil },
	"-": func(a, b float64) (float64, error) { return a - b, nil },
	"*": func(a, b float64) (float64, error) { return a * b, nil },
	"/": func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	},
	"^": func(a, b float64) (float64, error) { return math.Pow(a, b), nil },
}

// unaryOps maps a token to a function of a single operand.
var unaryOps = map[string]func(a float64) (float64, error){
	"sqrt": func(a float64) (float64, error) {
		if a < 0 {
			return 0, fmt.Errorf("sqrt of negative number %g", a)
		}
		return math.Sqrt(a), nil
	},
	"sin": func(a float64) (float64, error) { return math.Sin(a), nil },
	"cos": func(a float64) (float64, error) { return math.Cos(a), nil },
	"tan": func(a float64) (float64, error) { return math.Tan(a), nil },
	"log": func(a float64) (float64, error) {
		if a <= 0 {
			return 0, fmt.Errorf("log of non-positive number %g", a)
		}
		return math.Log(a), nil
	},
}

// Evaluate parses and evaluates a space-separated RPN expression, returning
// the numeric result or a descriptive error.
func Evaluate(expression string) (float64, error) {
	tokens := strings.Fields(expression)
	if len(tokens) == 0 {
		return 0, fmt.Errorf("empty expression")
	}

	var stack []float64

	for _, tok := range tokens {
		if fn, ok := binaryOps[tok]; ok {
			if len(stack) < 2 {
				return 0, fmt.Errorf("stack underflow: operator %q requires 2 operands, got %d", tok, len(stack))
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			res, err := fn(a, b)
			if err != nil {
				return 0, err
			}
			stack = append(stack, res)
			continue
		}

		if fn, ok := unaryOps[tok]; ok {
			if len(stack) < 1 {
				return 0, fmt.Errorf("stack underflow: operator %q requires 1 operand, got %d", tok, len(stack))
			}
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			res, err := fn(a)
			if err != nil {
				return 0, err
			}
			stack = append(stack, res)
			continue
		}

		val, err := strconv.ParseFloat(tok, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid token %q: not a number or known operator", tok)
		}
		stack = append(stack, val)
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression: %d values left on stack, expected 1", len(stack))
	}

	return stack[0], nil
}
