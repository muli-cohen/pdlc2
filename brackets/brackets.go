package brackets

import "errors"

type frame struct {
	opener byte
	idx    int
}

func scan(text string) (mismatchIdx int, hasBracket bool, balanced bool) {
	var stack []frame
	for i := 0; i < len(text); i++ {
		b := text[i]
		switch b {
		case '(', '[', '{':
			hasBracket = true
			stack = append(stack, frame{b, i})
		case ')', ']', '}':
			hasBracket = true
			var expected byte
			switch b {
			case ')':
				expected = '('
			case ']':
				expected = '['
			case '}':
				expected = '{'
			}
			if len(stack) == 0 || stack[len(stack)-1].opener != expected {
				return i, true, false
			}
			stack = stack[:len(stack)-1]
		}
	}
	if len(stack) > 0 {
		return stack[0].idx, true, false
	}
	return -1, hasBracket, true
}

func IsBalanced(text string) bool {
	_, _, balanced := scan(text)
	return balanced
}

func FirstMismatch(text string) (int, error) {
	mismatchIdx, hasBracket, balanced := scan(text)
	if !hasBracket {
		return 0, errors.New("no brackets found in input")
	}
	if !balanced {
		return mismatchIdx, nil
	}
	return -1, nil
}
