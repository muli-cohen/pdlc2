package rle

import (
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"single char run", "c", "c1", false},
		{"basic run", "aaabbc", "a3b2c1", false},
		{"run longer than 9", strings.Repeat("a", 12), "a12", false},
		{"mixed case", "aAaa", "a1A1a2", false},
		{"whitespace and punctuation", "hello, world!", "h1e1l2o1,1 1w1o1r1l1d1!1", false},
		{"empty input", "", "", false},
		{"digit in input", "abc1", "", true},
		{"digit only", "5", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Encode(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Encode(%q): expected error, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("Encode(%q): unexpected error: %v", tc.input, err)
				return
			}
			if got != tc.want {
				t.Errorf("Encode(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"basic decode", "a3b2c1", "aaabbc", false},
		{"count greater than 9", "a12", strings.Repeat("a", 12), false},
		{"single char", "c1", "c", false},
		{"empty input", "", "", false},
		{"missing count at end", "a3b", "", true},
		{"zero count", "a0b3", "", true},
		{"leading digit", "3a", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Decode(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Decode(%q): expected error, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("Decode(%q): unexpected error: %v", tc.input, err)
				return
			}
			if got != tc.want {
				t.Errorf("Decode(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	inputs := []string{
		"aaabbc",
		"hello, world!",
		"aAaa",
		"c",
		strings.Repeat("a", 12),
		"",
	}
	for _, text := range inputs {
		encoded, err := Encode(text)
		if err != nil {
			t.Errorf("Encode(%q): unexpected error: %v", text, err)
			continue
		}
		decoded, err := Decode(encoded)
		if err != nil {
			t.Errorf("Decode(Encode(%q)): unexpected error: %v", text, err)
			continue
		}
		if decoded != text {
			t.Errorf("round-trip(%q): got %q, want %q", text, decoded, text)
		}
	}
}
