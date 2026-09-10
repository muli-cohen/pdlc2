package roman

import "testing"

func TestToRoman_subtractiveForm(t *testing.T) {
	cases := []struct {
		n    int
		want string
	}{
		{4, "IV"}, {9, "IX"}, {40, "XL"}, {90, "XC"},
		{400, "CD"}, {900, "CM"}, {1994, "MCMXCIV"},
	}
	for _, c := range cases {
		got, err := ToRoman(c.n)
		if err != nil {
			t.Errorf("ToRoman(%d) unexpected error: %v", c.n, err)
			continue
		}
		if got != c.want {
			t.Errorf("ToRoman(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

func TestToRoman_boundaries(t *testing.T) {
	cases := []struct {
		n    int
		want string
	}{
		{1, "I"}, {3999, "MMMCMXCIX"},
	}
	for _, c := range cases {
		got, err := ToRoman(c.n)
		if err != nil {
			t.Errorf("ToRoman(%d) unexpected error: %v", c.n, err)
			continue
		}
		if got != c.want {
			t.Errorf("ToRoman(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

func TestToRoman_outOfRange(t *testing.T) {
	for _, n := range []int{0, 4000, -1} {
		s, err := ToRoman(n)
		if err == nil {
			t.Errorf("ToRoman(%d) expected error, got %q", n, s)
		}
		if s != "" {
			t.Errorf("ToRoman(%d) expected empty string on error, got %q", n, s)
		}
	}
}

func TestFromRoman_roundTrip(t *testing.T) {
	for n := 1; n <= 3999; n++ {
		s, err := ToRoman(n)
		if err != nil {
			t.Fatalf("ToRoman(%d) unexpected error: %v", n, err)
		}
		got, err := FromRoman(s)
		if err != nil {
			t.Errorf("FromRoman(%q) unexpected error: %v", s, err)
			continue
		}
		if got != n {
			t.Errorf("FromRoman(ToRoman(%d)) = %d, want %d", n, got, n)
		}
	}
}

func TestFromRoman_caseInsensitive(t *testing.T) {
	got, err := FromRoman("mcmxciv")
	if err != nil {
		t.Fatalf("FromRoman(\"mcmxciv\") unexpected error: %v", err)
	}
	if got != 1994 {
		t.Errorf("FromRoman(\"mcmxciv\") = %d, want 1994", got)
	}
}

func TestFromRoman_rejectEmpty(t *testing.T) {
	_, err := FromRoman("")
	if err == nil {
		t.Error("FromRoman(\"\") expected error, got nil")
	}
}

func TestFromRoman_rejectUnknownChar(t *testing.T) {
	_, err := FromRoman("ABC")
	if err == nil {
		t.Error("FromRoman(\"ABC\") expected error, got nil")
	}
}

func TestFromRoman_rejectNonCanonical(t *testing.T) {
	for _, s := range []string{"IIII", "VX", "IC"} {
		_, err := FromRoman(s)
		if err == nil {
			t.Errorf("FromRoman(%q) expected error, got nil", s)
		}
	}
}
