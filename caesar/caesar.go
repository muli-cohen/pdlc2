package caesar

// Encode applies a Caesar shift to every ASCII letter in text, preserving case.
// Non-letter bytes are copied unchanged. shift is normalised modulo 26.
// The error return is always nil for any int shift.
func Encode(text string, shift int) (string, error) {
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

// Decode reverses Encode by applying the inverse shift.
// Normalises shift before negating to avoid overflow at math.MinInt.
// The error return is always nil for any int shift.
func Decode(text string, shift int) (string, error) {
	n := ((shift % 26) + 26) % 26
	return Encode(text, 26-n)
}
