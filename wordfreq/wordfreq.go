package wordfreq

import (
	"bufio"
	"io"
	"sort"
	"strings"
	"unicode"
)

// Entry holds a word and the number of times it appeared in the input.
type Entry struct {
	Word  string
	Count int
}

// Count reads all text from r, tokenises it into runs of Unicode letters and
// decimal digits, lowercases each token, and returns a map of word to count.
// It returns an empty (non-nil) map and a nil error for empty or
// whitespace-only input, and a non-nil error only if reading from r fails.
func Count(r io.Reader) (map[string]int, error) {
	counts := make(map[string]int)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.FieldsFunc(line, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		})
		for _, t := range tokens {
			counts[strings.ToLower(t)]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return counts, nil
}

// Top returns a slice of at most n entries from counts, sorted by descending
// Count with ties broken alphabetically by Word. If n <= 0 or counts is
// nil/empty, it returns an empty non-nil slice.
func Top(counts map[string]int, n int) []Entry {
	if n <= 0 || len(counts) == 0 {
		return []Entry{}
	}
	entries := make([]Entry, 0, len(counts))
	for w, c := range counts {
		entries = append(entries, Entry{Word: w, Count: c})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Word < entries[j].Word
	})
	if n > len(entries) {
		n = len(entries)
	}
	return entries[:n]
}
