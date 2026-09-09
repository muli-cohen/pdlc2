package calc

import "testing"

func TestAdd(t *testing.T) {
	got := Add(1, 2)
	if got != 3 {
		t.Errorf("Add(1, 2) = %g; want 3", got)
	}
}

func TestSubtract(t *testing.T) {
	got := Subtract(5, 3)
	if got != 2 {
		t.Errorf("Subtract(5, 3) = %g; want 2", got)
	}
}

func TestMultiply(t *testing.T) {
	got := Multiply(3, 4)
	if got != 12 {
		t.Errorf("Multiply(3, 4) = %g; want 12", got)
	}
}

func TestDivide(t *testing.T) {
	got, err := Divide(6, 3)
	if err != nil {
		t.Fatalf("Divide(6, 3) returned unexpected error: %v", err)
	}
	if got != 2 {
		t.Errorf("Divide(6, 3) = %g; want 2", got)
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := Divide(1, 0)
	if err == nil {
		t.Fatal("Divide(1, 0) expected non-nil error, got nil")
	}
}
