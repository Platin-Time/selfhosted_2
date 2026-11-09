// Package calc содержит чистую логику вычислений без привязки к HTTP.
package calc

import (
	"errors"
	"math"
)

var (
	ErrDivByZero = errors.New("division by zero")
	ErrUnknownOp = errors.New("unknown operation, allowed: +, -, *, /, **")
	ErrOverflow  = errors.New("result is out of range")
)

// Calculate выполняет операцию op над a и b.
func Calculate(a, b int, op string) (float64, error) {
	var res float64
	x, y := float64(a), float64(b)

	switch op {
	case "+":
		res = x + y
	case "-":
		res = x - y
	case "*":
		res = x * y
	case "/":
		if b == 0 {
			return 0, ErrDivByZero
		}
		res = x / y
	case "**":
		res = math.Pow(x, y)
	default:
		return 0, ErrUnknownOp
	}

	if math.IsInf(res, 0) || math.IsNaN(res) {
		return 0, ErrOverflow
	}
	return res, nil
}
