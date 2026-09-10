package password

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits    = "0123456789"
	symbols   = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

// Options selects which character classes participate in password generation.
type Options struct {
	Lowercase bool
	Uppercase bool
	Digits    bool
	Symbols   bool
}

// Generate returns a cryptographically secure random password of the given length.
// Every enabled character class is guaranteed to contribute at least one character.
func Generate(length int, opts Options) (string, error) {
	if length < 8 || length > 128 {
		return "", fmt.Errorf("length must be between 8 and 128")
	}

	classes := []string{}
	if opts.Lowercase {
		classes = append(classes, lowercase)
	}
	if opts.Uppercase {
		classes = append(classes, uppercase)
	}
	if opts.Digits {
		classes = append(classes, digits)
	}
	if opts.Symbols {
		classes = append(classes, symbols)
	}

	if len(classes) == 0 {
		return "", fmt.Errorf("at least one character class must be selected")
	}

	alphabet := ""
	for _, c := range classes {
		alphabet += c
	}

	buf := make([]byte, length)

	// Reserve one character per selected class to guarantee coverage.
	for i, class := range classes {
		idx, err := randIndex(len(class))
		if err != nil {
			return "", err
		}
		buf[i] = class[idx]
	}

	// Fill remaining positions from the full alphabet.
	for i := len(classes); i < length; i++ {
		idx, err := randIndex(len(alphabet))
		if err != nil {
			return "", err
		}
		buf[i] = alphabet[idx]
	}

	// Fisher-Yates shuffle using crypto/rand.
	for i := length - 1; i > 0; i-- {
		j, err := randIndex(i + 1)
		if err != nil {
			return "", err
		}
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf), nil
}

func randIndex(n int) (int, error) {
	max := big.NewInt(int64(n))
	v, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return int(v.Int64()), nil
}
