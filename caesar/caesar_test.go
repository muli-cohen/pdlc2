package caesar

import "testing"

func TestEncryptLowercase(t *testing.T) {
	got, err := Encrypt("abc", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "def" {
		t.Errorf("got %q, want %q", got, "def")
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

func TestWrapAroundLower(t *testing.T) {
	got, err := Encrypt("xyz", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "yza" {
		t.Errorf("got %q, want %q", got, "yza")
	}
}

func TestWrapAroundUpper(t *testing.T) {
	got, err := Encrypt("XYZ", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "YZA" {
		t.Errorf("got %q, want %q", got, "YZA")
	}
}

func TestNonLetterPassThrough(t *testing.T) {
	got, err := Encrypt("hello, world! 123", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "khoor, zruog! 123" {
		t.Errorf("got %q, want %q", got, "khoor, zruog! 123")
	}
}

func TestShiftZero(t *testing.T) {
	got, err := Encrypt("abc", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}

func TestShift26(t *testing.T) {
	got, err := Encrypt("abc", 26)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}

func TestNegativeShift(t *testing.T) {
	got, err := Encrypt("abc", -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "zab" {
		t.Errorf("got %q, want %q", got, "zab")
	}
}

func TestEmptyString(t *testing.T) {
	got, err := Encrypt("", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}

func TestRoundTrip(t *testing.T) {
	input := "attack at dawn"
	encoded, err := Encrypt(input, 3)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	decoded, err := Decrypt(encoded, 3)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	if decoded != input {
		t.Errorf("round-trip failed: got %q, want %q", decoded, input)
	}
}

func TestEncryptReturnsNilErrorForASCII(t *testing.T) {
	_, err := Encrypt("hello", 5)
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDecryptReturnsNilErrorForASCII(t *testing.T) {
	_, err := Decrypt("hello", 5)
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNonASCIIReturnsError(t *testing.T) {
	_, err := Encrypt("caf\xc3\xa9", 3)
	if err == nil {
		t.Error("expected non-nil error for non-ASCII input, got nil")
	}
}
