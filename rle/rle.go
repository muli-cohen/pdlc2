package rle

import (
	"fmt"
	"strconv"
	"strings"
)

// Encode replaces each contiguous run of identical bytes with the byte followed
// by its decimal count. Returns an error if the input contains any digit, since
// digits in the input would make the encoded output ambiguous.
func Encode(text string) (string, error) {
	for i := 0; i < len(text); i++ {
		if text[i] >= '0' && text[i] <= '9' {
			return "", fmt.Errorf("input must not contain digits")
		}
	}
	if len(text) == 0 {
		return "", nil
	}
	var b strings.Builder
	count := 1
	for i := 1; i < len(text); i++ {
		if text[i] == text[i-1] {
			count++
		} else {
			fmt.Fprintf(&b, "%c%d", text[i-1], count)
			count = 1
		}
	}
	fmt.Fprintf(&b, "%c%d", text[len(text)-1], count)
	return b.String(), nil
}

// Decode reverses Encode. Each token is one non-digit byte followed by one or
// more digit bytes parsed greedily as a decimal count. Returns an error for
// malformed input: a leading digit, a character with no following count, or a
// count of zero.
func Decode(text string) (string, error) {
	if len(text) == 0 {
		return "", nil
	}
	if text[0] >= '0' && text[0] <= '9' {
		return "", fmt.Errorf("malformed input: encoded string must not begin with a digit")
	}
	var b strings.Builder
	i := 0
	for i < len(text) {
		ch := text[i]
		i++
		j := i
		for j < len(text) && text[j] >= '0' && text[j] <= '9' {
			j++
		}
		if j == i {
			return "", fmt.Errorf("malformed input: character '%c' has no following count", ch)
		}
		count, err := strconv.Atoi(text[i:j])
		if err != nil {
			return "", fmt.Errorf("malformed input: invalid count")
		}
		if count == 0 {
			return "", fmt.Errorf("malformed input: count of zero is not valid")
		}
		b.WriteString(strings.Repeat(string(ch), count))
		i = j
	}
	return b.String(), nil
}
