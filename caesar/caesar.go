package caesar

import "fmt"

func Encrypt(text string, shift int) (string, error) {
	for i := 0; i < len(text); i++ {
		if text[i] > 127 {
			return "", fmt.Errorf("input contains non-ASCII character")
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

func Decrypt(text string, shift int) (string, error) {
	return Encrypt(text, -shift)
}
