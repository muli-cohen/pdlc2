package caesar

import "fmt"

// Encrypt applies a Caesar shift to every ASCII letter in text, preserving case.
// Non-letter bytes are copied unchanged. shift is normalised modulo 26.
// Returns an error if text contains any byte with value > 127.
func Encrypt(text string, shift int) (string, error) {
	for i := 0; i < len(text); i++ {
		if text[i] > 127 {
			return "", fmt.Errorf("input contains non-ASCII bytes")
		}
	}
	n := ((shift % 26) + 26) % 26
	buf := make([]byte, len(text))
	for i := 0; i < len(text); i++ {
		b := text[i]
		switch {
		case b >= 'a' && b <= 'z':
			buf[i] = 'a' + (b-'a'+byte(n))%26
		case b >= 'A' && b <= 'Z':
			buf[i] = 'A' + (b-'A'+byte(n))%26
		default:
			buf[i] = b
		}
	}
	return string(buf), nil
}

// Decrypt reverses Encrypt by applying the inverse shift.
func Decrypt(text string, shift int) (string, error) {
	return Encrypt(text, -shift)
}
