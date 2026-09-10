package caesar

import "testing"

func TestEncodeLowercase(t *testing.T) {
	got, err := Encode("abc", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "def" {
		t.Errorf("got %q, want %q", got, "def")
	}
}

func TestDecodeLowercase(t *testing.T) {
	got, err := Decode("def", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}

func TestEncodeUppercase(t *testing.T) {
	got, err := Encode("ABC", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "DEF" {
		t.Errorf("got %q, want %q", got, "DEF")
	}
}

func TestDecodeUppercase(t *testing.T) {
	got, err := Decode("DEF", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ABC" {
		t.Errorf("got %q, want %q", got, "ABC")
	}
}

func TestWrapAroundLower(t *testing.T) {
	got, err := Encode("xyz", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "yza" {
		t.Errorf("got %q, want %q", got, "yza")
	}
}

func TestWrapAroundUpper(t *testing.T) {
	got, err := Encode("XYZ", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "YZA" {
		t.Errorf("got %q, want %q", got, "YZA")
	}
}

func TestNonLetterPassThrough(t *testing.T) {
	got, err := Encode("hello, world! 123", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "khoor, zruog! 123" {
		t.Errorf("got %q, want %q", got, "khoor, zruog! 123")
	}
}

func TestShiftZero(t *testing.T) {
	got, err := Encode("abc", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}

func TestShift26(t *testing.T) {
	got, err := Encode("abc", 26)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}

func TestNegativeShift(t *testing.T) {
	got, err := Encode("abc", -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "zab" {
		t.Errorf("got %q, want %q", got, "zab")
	}
}

func TestEmptyString(t *testing.T) {
	got, err := Encode("", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}

func TestRoundTrip(t *testing.T) {
	input := "attack at dawn"
	encoded, err := Encode(input, 3)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	decoded, err := Decode(encoded, 3)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if decoded != input {
		t.Errorf("round-trip failed: got %q, want %q", decoded, input)
	}
}

func TestEncodeReturnsNilError(t *testing.T) {
	_, err := Encode("hello", 5)
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDecodeReturnsNilError(t *testing.T) {
	_, err := Decode("hello", 5)
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}
