// Package eval is an evaluator for simple arithmetic expressions
package eval

import (
	"fmt"
	"math"
	"strings"
)

// An Expr is an arithmetic expresson
type Expr interface {
	// Eval returns the value of this Expr in the environment env
	Eval(env Env) float64

	// Check reports errors in this Expr and adds its Vars to the set
	Check(vars map[Var]bool) error
}

// Env is an environment that maps variable name to values
type Env map[Var]float64

// 5 concrete types that represent particular kinds of expression: Var, literal, unary, binary, call

// We'll also need each kind of expresion to define an Eval method that returns the expression's value in a given environment.
// (Since every expression must provide Eval, we add it to the Expr interface)

// A Var represents a reference to a variable
type Var string

// Eval is an environment that maps variable name to values
func (v Var) Eval(env Env) float64 {
	return env[v]
}
func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

// A literal is a floating-point constant
type literal float64

func (l literal) Eval(_ Env) float64 {
	return float64(l)
}
func (literal) Check(vars map[Var]bool) error {
	return nil
}

// A unary represents a unary operator expression, e.g., -x
type unary struct {
	op rune // one of '+', '-'
	x  Expr
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}
func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

// A binary represents a binary operator expression, e.g., x+y
type binary struct {
	op   rune // one of '+', '-', '*', '/'
	x, y Expr
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", b.op))
}
func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/", b.op) {
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

// A call represents a function call expression, e.g., sin(x)
type call struct {
	fn   string // one of "pow", "sin", "sqrt"
	args []Expr
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", c.fn))
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d", c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return nil
		}
	}
	return nil
}
