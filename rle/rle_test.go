package rle

import (
	"strings"
	"testing"
)

func TestEncodeBasicRun(t *testing.T) {
	got, err := Encode("aaabbc")
	if err != nil || got != "a3b2c1" {
		t.Errorf("Encode(\"aaabbc\") = (%q, %v), want (\"a3b2c1\", nil)", got, err)
	}
}

func TestEncodeSingleCharRuns(t *testing.T) {
	got, err := Encode("abc")
	if err != nil || got != "a1b1c1" {
		t.Errorf("Encode(\"abc\") = (%q, %v), want (\"a1b1c1\", nil)", got, err)
	}
}

func TestEncodeLongRun(t *testing.T) {
	got, err := Encode("aaaaaaaaaaaa")
	if err != nil || got != "a12" {
		t.Errorf("Encode(12 a's) = (%q, %v), want (\"a12\", nil)", got, err)
	}
}

func TestEncodeMixedCase(t *testing.T) {
	got, err := Encode("aAaa")
	if err != nil || got != "a1A1a2" {
		t.Errorf("Encode(\"aAaa\") = (%q, %v), want (\"a1A1a2\", nil)", got, err)
	}
}

func TestEncodePunctuation(t *testing.T) {
	got, err := Encode("hello, world!")
	want := "h1e1l2o1,1 1w1o1r1l1d1!1"
	if err != nil || got != want {
		t.Errorf("Encode(\"hello, world!\") = (%q, %v), want (%q, nil)", got, err, want)
	}
}

func TestEncodeEmpty(t *testing.T) {
	got, err := Encode("")
	if err != nil || got != "" {
		t.Errorf("Encode(\"\") = (%q, %v), want (\"\", nil)", got, err)
	}
}

func TestEncodeDigitRejection(t *testing.T) {
	got, err := Encode("abc123")
	if err == nil {
		t.Errorf("Encode(\"abc123\") = (%q, nil), want non-nil error", got)
	}
	if got != "" {
		t.Errorf("Encode(\"abc123\") returned non-empty string %q on error", got)
	}
	if !strings.Contains(err.Error(), "digit") {
		t.Errorf("error message %q does not mention \"digit\"", err.Error())
	}
}

func TestDecodeMissingCount(t *testing.T) {
	got, err := Decode("a3b")
	if err == nil {
		t.Errorf("Decode(\"a3b\") = (%q, nil), want non-nil error", got)
	}
	if !strings.Contains(err.Error(), "no following count") {
		t.Errorf("error message %q does not contain \"no following count\"", err.Error())
	}
}

func TestDecodeZeroCount(t *testing.T) {
	got, err := Decode("a0")
	if err == nil {
		t.Errorf("Decode(\"a0\") = (%q, nil), want non-nil error", got)
	}
	if got != "" {
		t.Errorf("Decode(\"a0\") returned non-empty string %q on error", got)
	}
}

func TestDecodeEmpty(t *testing.T) {
	got, err := Decode("")
	if err != nil || got != "" {
		t.Errorf("Decode(\"\") = (%q, %v), want (\"\", nil)", got, err)
	}
}

func TestRoundTrip(t *testing.T) {
	input := "hello, world!"
	encoded, err := Encode(input)
	if err != nil {
		t.Fatalf("Encode(%q) error: %v", input, err)
	}
	got, err := Decode(encoded)
	if err != nil || got != input {
		t.Errorf("Decode(Encode(%q)) = (%q, %v), want (%q, nil)", input, got, err, input)
	}
}

func TestDecodeMultiDigit(t *testing.T) {
	got, err := Decode("a12")
	want := "aaaaaaaaaaaa"
	if err != nil || got != want {
		t.Errorf("Decode(\"a12\") = (%q, %v), want (%q, nil)", got, err, want)
	}
}
