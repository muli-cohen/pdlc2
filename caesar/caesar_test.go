package caesar

import "testing"

func TestEncodeLowercase(t *testing.T) {
	got, err := Encrypt("abc", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "def" {
		t.Errorf("got %q, want %q", got, "def")
	}
}

func TestDecodeLowercase(t *testing.T) {
	got, err := Decrypt("def", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}

func TestEncodeUppercase(t *testing.T) {
	got, err := Encrypt("ABC", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "DEF" {
		t.Errorf("got %q, want %q", got, "DEF")
	}
}

func TestDecodeUppercase(t *testing.T) {
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

func TestEncryptReturnsNilError(t *testing.T) {
	_, err := Encrypt("hello", 5)
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDecryptReturnsNilError(t *testing.T) {
	_, err := Decrypt("hello", 5)
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

// FR16: basic shifting - "attack at dawn" with shift 3
func TestEncryptBasicShift(t *testing.T) {
	got, err := Encrypt("attack at dawn", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "dwwdfn dw gdzq" {
		t.Errorf("got %q, want %q", got, "dwwdfn dw gdzq")
	}
}

// FR17: wraparound - "XYZ" with shift 3 -> "ABC"
func TestEncryptWrapAroundShift3(t *testing.T) {
	got, err := Encrypt("XYZ", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ABC" {
		t.Errorf("got %q, want %q", got, "ABC")
	}
}

// FR20: shift 27 behaves identically to shift 1
func TestEncryptShift27EqualShift1(t *testing.T) {
	got27, err := Encrypt("abc", 27)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got1, err := Encrypt("abc", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got27 != got1 {
		t.Errorf("shift 27 produced %q but shift 1 produced %q; expected equal", got27, got1)
	}
}

// FR21: shift -3 on ciphertext returns plaintext
func TestEncryptNegativeShift3(t *testing.T) {
	got, err := Encrypt("dwwdfn dw gdzq", -3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "attack at dawn" {
		t.Errorf("got %q, want %q", got, "attack at dawn")
	}
}

// FR22: non-ASCII input causes Encrypt to return a non-nil error
func TestEncryptNonASCII(t *testing.T) {
	got, err := Encrypt("caf\xc3\xa9", 3)
	if err == nil {
		t.Errorf("expected non-nil error for non-ASCII input, got nil; result: %q", got)
	}
	if got != "" {
		t.Errorf("expected empty string on error, got %q", got)
	}
}

// AC: Encrypt("Hello, World!", 3) returns ("Khoor, Zruog!", nil) - mixed case with punctuation
func TestEncryptMixedCaseNonLetters(t *testing.T) {
	got, err := Encrypt("Hello, World!", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Khoor, Zruog!" {
		t.Errorf("got %q, want %q", got, "Khoor, Zruog!")
	}
}

// FR22: non-ASCII input causes Decrypt to return a non-nil error
func TestDecryptNonASCII(t *testing.T) {
	got, err := Decrypt("caf\xc3\xa9", 3)
	if err == nil {
		t.Errorf("expected non-nil error for non-ASCII input, got nil; result: %q", got)
	}
	if got != "" {
		t.Errorf("expected empty string on error, got %q", got)
	}
}
