package jokes_test

import (
	"math/rand"
	"testing"

	"github.com/muli-cohen/pdlc2/jokes"
)

func TestAllLength(t *testing.T) {
	if got := len(jokes.All()); got != 15 {
		t.Errorf("All() length = %d; want 15", got)
	}
}

func TestAllNonEmpty(t *testing.T) {
	for i, j := range jokes.All() {
		if j == "" {
			t.Errorf("All()[%d] is empty", i)
		}
	}
}

func TestRandomMembership(t *testing.T) {
	r := rand.New(rand.NewSource(99))
	all := jokes.All()
	for range 50 {
		got := jokes.Random(r)
		found := false
		for _, j := range all {
			if j == got {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Random() returned %q which is not in All()", got)
		}
	}
}

func TestRandomDeterminism(t *testing.T) {
	r1 := rand.New(rand.NewSource(42))
	first := jokes.Random(r1)

	r2 := rand.New(rand.NewSource(42))
	second := jokes.Random(r2)

	if first != second {
		t.Errorf("same seed produced different jokes: %q vs %q", first, second)
	}
}
