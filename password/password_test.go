package password

import (
	"strings"
	"testing"
	"unicode"
)

func TestDefaultOptions(t *testing.T) {
	result, err := Generate(16, Options{Lowercase: true, Uppercase: true, Digits: true, Symbols: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 16 {
		t.Fatalf("expected length 16, got %d", len(result))
	}
	assertContainsClass(t, result, lowercase, "lowercase")
	assertContainsClass(t, result, uppercase, "uppercase")
	assertContainsClass(t, result, digits, "digits")
	assertContainsClass(t, result, symbols, "symbols")
}

func TestExcludeUppercase(t *testing.T) {
	result, err := Generate(16, Options{Lowercase: true, Uppercase: false, Digits: true, Symbols: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, ch := range result {
		if unicode.IsUpper(ch) {
			t.Fatalf("expected no uppercase letters, got %q in %q", ch, result)
		}
	}
}

func TestExcludeDigits(t *testing.T) {
	result, err := Generate(16, Options{Lowercase: true, Uppercase: true, Digits: false, Symbols: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, ch := range result {
		if unicode.IsDigit(ch) {
			t.Fatalf("expected no digits, got %q in %q", ch, result)
		}
	}
}

func TestExcludeSymbols(t *testing.T) {
	result, err := Generate(16, Options{Lowercase: true, Uppercase: true, Digits: true, Symbols: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.ContainsAny(result, symbols) {
		t.Fatalf("expected no symbols, got %q", result)
	}
}

func TestClassCoverageGuarantee(t *testing.T) {
	result, err := Generate(8, Options{Lowercase: true, Uppercase: true, Digits: true, Symbols: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 8 {
		t.Fatalf("expected length 8, got %d", len(result))
	}
	assertContainsClass(t, result, lowercase, "lowercase")
	assertContainsClass(t, result, uppercase, "uppercase")
	assertContainsClass(t, result, digits, "digits")
	assertContainsClass(t, result, symbols, "symbols")
}

func TestLengthBelowMinimum(t *testing.T) {
	_, err := Generate(4, Options{Lowercase: true, Uppercase: true, Digits: true, Symbols: true})
	if err == nil {
		t.Fatal("expected error for length 4, got nil")
	}
	if !strings.Contains(err.Error(), "8") {
		t.Fatalf("error message should mention minimum length 8, got: %q", err.Error())
	}
}

func TestLengthAboveMaximum(t *testing.T) {
	_, err := Generate(200, Options{Lowercase: true, Uppercase: true, Digits: true, Symbols: true})
	if err == nil {
		t.Fatal("expected error for length 200, got nil")
	}
	if !strings.Contains(err.Error(), "128") {
		t.Fatalf("error message should mention maximum length 128, got: %q", err.Error())
	}
}

func TestNoClassSelected(t *testing.T) {
	_, err := Generate(16, Options{Lowercase: false, Uppercase: false, Digits: false, Symbols: false})
	if err == nil {
		t.Fatal("expected error when no class selected, got nil")
	}
	if !strings.Contains(err.Error(), "at least one") {
		t.Fatalf("error message should mention 'at least one', got: %q", err.Error())
	}
}

func assertContainsClass(t *testing.T, result, class, name string) {
	t.Helper()
	if !strings.ContainsAny(result, class) {
		t.Fatalf("expected at least one %s character in %q", name, result)
	}
}
