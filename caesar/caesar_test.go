package caesar

import (
	"math"
	"testing"
)

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

func TestEncodeBasicShift(t *testing.T) {
	got, err := Encode("attack at dawn", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "dwwdfn dw gdzq" {
		t.Errorf("got %q, want %q", got, "dwwdfn dw gdzq")
	}
}

func TestEncodeWrapAroundShift3(t *testing.T) {
	got, err := Encode("XYZ", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ABC" {
		t.Errorf("got %q, want %q", got, "ABC")
	}
}

func TestEncodeShift27EqualShift1(t *testing.T) {
	got27, err := Encode("abc", 27)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got1, err := Encode("abc", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got27 != got1 {
		t.Errorf("shift 27 produced %q but shift 1 produced %q; expected equal", got27, got1)
	}
}

func TestEncodeNegativeShift3(t *testing.T) {
	got, err := Encode("dwwdfn dw gdzq", -3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "attack at dawn" {
		t.Errorf("got %q, want %q", got, "attack at dawn")
	}
}

func TestRoundTripMinInt(t *testing.T) {
	input := "abc"
	encoded, err := Encode(input, math.MinInt)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	decoded, err := Decode(encoded, math.MinInt)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if decoded != input {
		t.Errorf("round-trip failed for MinInt shift: got %q, want %q", decoded, input)
	}
}

func TestEncodeNonASCIIPassthrough(t *testing.T) {
	input := "caf\xc3\xa9"
	got, err := Encode(input, 3)
	if err != nil {
		t.Errorf("expected nil error for non-ASCII input, got: %v", err)
	}
	// letters c, a, f shift by 3; non-ASCII bytes \xc3\xa9 pass through unchanged
	if got != "fdi\xc3\xa9" {
		t.Errorf("got %q, want %q", got, "fdi\xc3\xa9")
	}
}
