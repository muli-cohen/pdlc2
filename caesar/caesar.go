package caesar

import "errors"

// Encrypt applies a Caesar shift to every ASCII letter in text, preserving case.
// Non-letter ASCII bytes are copied unchanged. shift is normalised modulo 26.
// Returns an error if any byte in text has value > 127 (non-ASCII).
func Encrypt(text string, shift int) (string, error) {
	n := ((shift % 26) + 26) % 26
	buf := make([]byte, len(text))
	for i := 0; i < len(text); i++ {
		b := text[i]
		if b > 127 {
			return "", errors.New("input contains non-ASCII characters")
		}
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

// Decrypt reverses Encrypt by applying the negated shift.
func Decrypt(text string, shift int) (string, error) {
	return Encrypt(text, -shift)
}
