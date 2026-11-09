package calc

import (
	"errors"
	"testing"
)

func TestCalculate(t *testing.T) {
	if got, err := Calculate(3, 4, "**"); err != nil || got != 81 {
		t.Fatalf("3**4 = %v, %v; want 81", got, err)
	}
	if _, err := Calculate(1, 0, "/"); !errors.Is(err, ErrDivByZero) {
		t.Fatalf("want ErrDivByZero, got %v", err)
	}
	if _, err := Calculate(1, 1, "?"); !errors.Is(err, ErrUnknownOp) {
		t.Fatalf("want ErrUnknownOp, got %v", err)
	}
}
