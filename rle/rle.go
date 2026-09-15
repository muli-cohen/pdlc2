package rle

import (
	"fmt"
	"strconv"
	"strings"
)

// Encode replaces each run of identical bytes with the byte followed by its decimal count.
// Returns an error if text contains any ASCII digit.
func Encode(text string) (string, error) {
	if text == "" {
		return "", nil
	}
	for i := 0; i < len(text); i++ {
		if text[i] >= '0' && text[i] <= '9' {
			return "", fmt.Errorf("input contains digit(s): digits are not allowed in encode input")
		}
	}
	var b strings.Builder
	i := 0
	for i < len(text) {
		ch := text[i]
		j := i + 1
		for j < len(text) && text[j] == ch {
			j++
		}
		b.WriteByte(ch)
		b.WriteString(strconv.Itoa(j - i))
		i = j
	}
	return b.String(), nil
}

// Decode reverses Encode: each token is one non-digit byte followed by one or more digit bytes.
// Returns an error if the input is malformed (missing count, zero count, or leading digit).
func Decode(text string) (string, error) {
	if text == "" {
		return "", nil
	}
	var b strings.Builder
	i := 0
	for i < len(text) {
		if text[i] >= '0' && text[i] <= '9' {
			return "", fmt.Errorf("malformed encoded input: unexpected digit at position %d", i)
		}
		ch := text[i]
		i++
		j := i
		for j < len(text) && text[j] >= '0' && text[j] <= '9' {
			j++
		}
		if j == i {
			return "", fmt.Errorf("malformed encoded input: character %q has no following count", ch)
		}
		count, err := strconv.Atoi(text[i:j])
		if err != nil || count == 0 {
			return "", fmt.Errorf("malformed encoded input: character %q has a count of zero", ch)
		}
		for k := 0; k < count; k++ {
			b.WriteByte(ch)
		}
		i = j
	}
	return b.String(), nil
}
