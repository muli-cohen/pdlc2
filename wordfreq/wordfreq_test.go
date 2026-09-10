package wordfreq_test

import (
	"strings"
	"testing"

	"github.com/muli-cohen/pdlc2/wordfreq"
)

func TestCountBasic(t *testing.T) {
	counts, err := wordfreq.Count(strings.NewReader("the cat the hat"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]int{"the": 2, "cat": 1, "hat": 1}
	for w, c := range want {
		if counts[w] != c {
			t.Errorf("counts[%q] = %d, want %d", w, counts[w], c)
		}
	}
	if len(counts) != len(want) {
		t.Errorf("len(counts) = %d, want %d", len(counts), len(want))
	}
}

func TestCountCaseFolding(t *testing.T) {
	counts, err := wordfreq.Count(strings.NewReader("The the THE"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if counts["the"] != 3 {
		t.Errorf("counts[\"the\"] = %d, want 3", counts["the"])
	}
	if len(counts) != 1 {
		t.Errorf("expected 1 distinct word, got %d", len(counts))
	}
}

func TestCountPunctuation(t *testing.T) {
	counts, err := wordfreq.Count(strings.NewReader("hello, world!"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]int{"hello": 1, "world": 1}
	for w, c := range want {
		if counts[w] != c {
			t.Errorf("counts[%q] = %d, want %d", w, counts[w], c)
		}
	}
	if len(counts) != len(want) {
		t.Errorf("len(counts) = %d, want %d", len(counts), len(want))
	}
}

func TestCountEmpty(t *testing.T) {
	counts, err := wordfreq.Count(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if counts == nil {
		t.Fatal("expected non-nil map for empty input")
	}
	if len(counts) != 0 {
		t.Errorf("expected empty map, got %v", counts)
	}
}

func TestTopTieBreak(t *testing.T) {
	counts := map[string]int{"cat": 1, "bat": 1, "ant": 1}
	entries := wordfreq.Top(counts, 3)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	want := []string{"ant", "bat", "cat"}
	for i, e := range entries {
		if e.Word != want[i] {
			t.Errorf("entries[%d].Word = %q, want %q", i, e.Word, want[i])
		}
	}
}

func TestTopNLargerThanWords(t *testing.T) {
	counts := map[string]int{"a": 1, "b": 2, "c": 3}
	entries := wordfreq.Top(counts, 10)
	if len(entries) != 3 {
		t.Errorf("expected 3 entries (all words), got %d", len(entries))
	}
}

func TestTopZeroN(t *testing.T) {
	counts := map[string]int{"a": 1}
	entries := wordfreq.Top(counts, 0)
	if entries == nil {
		t.Fatal("expected non-nil slice for n=0")
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice for n=0, got %v", entries)
	}
}

func TestTopNilMap(t *testing.T) {
	entries := wordfreq.Top(nil, 5)
	if entries == nil {
		t.Fatal("expected non-nil slice for nil map")
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice for nil map, got %v", entries)
	}
}
