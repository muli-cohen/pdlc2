package caesar

import "testing"

// FR15: basic shifting
func TestEncryptBasic(t *testing.T) {
	got, err := Encrypt("attack at dawn", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "dwwdfn dw gdzq" {
		t.Errorf("got %q, want %q", got, "dwwdfn dw gdzq")
	}
}

// FR16: wraparound at end of alphabet
func TestEncryptWrapAround(t *testing.T) {
	got, err := Encrypt("XYZ", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ABC" {
		t.Errorf("got %q, want %q", got, "ABC")
	}
}

func TestEncryptWrapAroundLower(t *testing.T) {
	got, err := Encrypt("xyz", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "yza" {
		t.Errorf("got %q, want %q", got, "yza")
	}
}

// FR17: mixed case treated independently
func TestEncryptMixedCase(t *testing.T) {
	got, err := Encrypt("Hello", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Khoor" {
		t.Errorf("got %q, want %q", got, "Khoor")
	}
}

func TestEncryptUppercase(t *testing.T) {
	got, err := Encrypt("ABC", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "DEF" {
		t.Errorf("got %q, want %q", got, "DEF")
	}
}

func TestDecryptUppercase(t *testing.T) {
	got, err := Decrypt("DEF", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ABC" {
		t.Errorf("got %q, want %q", got, "ABC")
	}
}

// FR18: non-letter characters left untouched
func TestEncryptNonLetterPassthrough(t *testing.T) {
	got, err := Encrypt("hello, world! 123", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "khoor, zruog! 123" {
		t.Errorf("got %q, want %q", got, "khoor, zruog! 123")
	}
}

// FR19: shift normalization including negative shifts
func TestEncryptShiftZero(t *testing.T) {
	got, err := Encrypt("abc", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}

func TestEncryptShift26(t *testing.T) {
	got, err := Encrypt("abc", 26)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}

func TestEncryptShift27(t *testing.T) {
	got, err := Encrypt("abc", 27)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "bcd" {
		t.Errorf("got %q, want %q", got, "bcd")
	}
}

func TestEncryptNegativeShift(t *testing.T) {
	got, err := Encrypt("dwwdfn dw gdzq", -3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "attack at dawn" {
		t.Errorf("got %q, want %q", got, "attack at dawn")
	}
}

func TestEncryptNegativeShiftMinus1(t *testing.T) {
	got, err := Encrypt("abc", -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "zab" {
		t.Errorf("got %q, want %q", got, "zab")
	}
}

// FR20: non-ASCII input rejected
func TestEncryptNonASCII(t *testing.T) {
	got, err := Encrypt("caf\xc3\xa9", 3)
	if err == nil {
		t.Errorf("expected error for non-ASCII input, got nil; result: %q", got)
	}
	if got != "" {
		t.Errorf("expected empty string on error, got %q", got)
	}
}

func TestDecryptNonASCII(t *testing.T) {
	got, err := Decrypt("caf\xc3\xa9", 3)
	if err == nil {
		t.Errorf("expected error for non-ASCII input, got nil; result: %q", got)
	}
	if got != "" {
		t.Errorf("expected empty string on error, got %q", got)
	}
}

// FR21: round-trip
func TestRoundTrip(t *testing.T) {
	input := "attack at dawn"
	encoded, err := Encrypt(input, 13)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	decoded, err := Decrypt(encoded, 13)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	if decoded != input {
		t.Errorf("round-trip failed: got %q, want %q", decoded, input)
	}
}

func TestRoundTripDecryptEncrypt(t *testing.T) {
	input := "dwwdfn dw gdzq"
	decoded, err := Decrypt(input, 3)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	encoded, err := Encrypt(decoded, 3)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	if encoded != input {
		t.Errorf("round-trip failed: got %q, want %q", encoded, input)
	}
}

func TestEncryptEmpty(t *testing.T) {
	got, err := Encrypt("", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}

func TestDecryptLowercase(t *testing.T) {
	got, err := Decrypt("def", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}
