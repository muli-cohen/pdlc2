package roman

import (
	"fmt"
	"strings"
)

var table = []struct {
	val int
	sym string
}{
	{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
	{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
	{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
}

func ToRoman(n int) (string, error) {
	if n < 1 || n > 3999 {
		return "", fmt.Errorf("value %d out of range: roman numerals support 1-3999", n)
	}
	var sb strings.Builder
	for _, e := range table {
		for n >= e.val {
			sb.WriteString(e.sym)
			n -= e.val
		}
	}
	return sb.String(), nil
}

var charVal = map[byte]int{
	'I': 1, 'V': 5, 'X': 10, 'L': 50,
	'C': 100, 'D': 500, 'M': 1000,
}

func FromRoman(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("roman numeral cannot be empty")
	}
	upper := strings.ToUpper(s)
	for i := 0; i < len(upper); i++ {
		if _, ok := charVal[upper[i]]; !ok {
			return 0, fmt.Errorf("invalid roman numeral %q: unrecognised character %q", s, string(upper[i]))
		}
	}
	sum := 0
	for i := 0; i < len(upper); i++ {
		cur := charVal[upper[i]]
		if i+1 < len(upper) && cur < charVal[upper[i+1]] {
			sum -= cur
		} else {
			sum += cur
		}
	}
	canonical, err := ToRoman(sum)
	if err != nil || canonical != upper {
		return 0, fmt.Errorf("invalid roman numeral %q: not in canonical form", s)
	}
	return sum, nil
}
