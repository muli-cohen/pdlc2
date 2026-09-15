package brackets

import "testing"

func TestSimpleBalancedPair(t *testing.T) {
	if !IsBalanced("()") {
		t.Errorf("expected IsBalanced(\"()\") to be true")
	}
}

func TestNestedBrackets(t *testing.T) {
	if !IsBalanced("{[()]}") {
		t.Errorf("expected IsBalanced(\"{[()]}\") to be true")
	}
}

func TestMixedBalanced(t *testing.T) {
	if !IsBalanced("(a[b]{c})") {
		t.Errorf("expected IsBalanced(\"(a[b]{c})\") to be true")
	}
}

func TestWrongOrder(t *testing.T) {
	if IsBalanced("([)]") {
		t.Errorf("expected IsBalanced(\"([)]\") to be false")
	}
}

func TestTrailingUnmatchedOpener(t *testing.T) {
	if IsBalanced("(a[b") {
		t.Errorf("expected IsBalanced(\"(a[b\") to be false")
	}
}

func TestLeadingUnmatchedCloser(t *testing.T) {
	if IsBalanced("]hello") {
		t.Errorf("expected IsBalanced(\"]hello\") to be false")
	}
}

func TestNoBrackets(t *testing.T) {
	if !IsBalanced("no brackets here") {
		t.Errorf("expected IsBalanced(\"no brackets here\") to be true")
	}
}

func TestBracketsWithProse(t *testing.T) {
	if !IsBalanced("hello (world [foo] bar)") {
		t.Errorf("expected IsBalanced(\"hello (world [foo] bar)\") to be true")
	}
}

func TestFirstMismatchIndex(t *testing.T) {
	idx, err := FirstMismatch("(a[b)]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx != 4 {
		t.Errorf("got index %d, want 4", idx)
	}
}

func TestFirstMismatchNoBrackets(t *testing.T) {
	_, err := FirstMismatch("no brackets here")
	if err == nil {
		t.Errorf("expected non-nil error for no-bracket input")
	}
}
