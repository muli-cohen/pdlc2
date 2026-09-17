package caesar

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

func Decode(text string, shift int) (string, error) {
	inverse := (26 - ((shift%26+26)%26)) % 26
	return Encode(text, inverse)
}
