package semver

import (
	"strings"
	"testing"
)

func mustParse(t *testing.T, s string) Version {
	t.Helper()
	v, err := Parse(s)
	if err != nil {
		t.Fatalf("Parse(%q) unexpected error: %v", s, err)
	}
	return v
}

func TestParseBasic(t *testing.T) {
	v := mustParse(t, "1.2.3")
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Errorf("got %+v, want Major=1 Minor=2 Patch=3", v)
	}
	if v.PreRelease != nil {
		t.Errorf("expected nil PreRelease, got %v", v.PreRelease)
	}
	if v.Build != "" {
		t.Errorf("expected empty Build, got %q", v.Build)
	}
}

func TestParsePreRelease(t *testing.T) {
	v := mustParse(t, "1.0.0-alpha")
	if len(v.PreRelease) != 1 || v.PreRelease[0] != "alpha" {
		t.Errorf("got PreRelease=%v, want [alpha]", v.PreRelease)
	}
}

func TestParsePreReleaseMulti(t *testing.T) {
	v := mustParse(t, "1.0.0-alpha.1")
	if len(v.PreRelease) != 2 || v.PreRelease[0] != "alpha" || v.PreRelease[1] != "1" {
		t.Errorf("got PreRelease=%v, want [alpha 1]", v.PreRelease)
	}
}

func TestParseBuild(t *testing.T) {
	v := mustParse(t, "1.0.0+build.5")
	if v.Build != "build.5" {
		t.Errorf("got Build=%q, want %q", v.Build, "build.5")
	}
	if v.PreRelease != nil {
		t.Errorf("expected nil PreRelease, got %v", v.PreRelease)
	}
}

func TestParsePreReleaseAndBuild(t *testing.T) {
	v := mustParse(t, "1.0.0-alpha+001")
	if len(v.PreRelease) != 1 || v.PreRelease[0] != "alpha" {
		t.Errorf("got PreRelease=%v, want [alpha]", v.PreRelease)
	}
	if v.Build != "001" {
		t.Errorf("got Build=%q, want %q", v.Build, "001")
	}
}

func TestParseRejectMissingComponent(t *testing.T) {
	_, err := Parse("1.2")
	if err == nil {
		t.Error("expected error for missing component, got nil")
	}
}

func TestParseRejectNonNumeric(t *testing.T) {
	_, err := Parse("1.x.0")
	if err == nil {
		t.Error("expected error for non-numeric component, got nil")
	}
}

func TestParseRejectLeadingZero(t *testing.T) {
	_, err := Parse("1.01.0")
	if err == nil {
		t.Error("expected error for leading zero, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "1.01.0") {
		t.Errorf("error message must contain offending version, got: %v", err)
	}
}

func TestParseRejectEmptyPreRelease(t *testing.T) {
	cases := []string{"1.0.0-", "1.0.0-.1", "1.0.0-1..2"}
	for _, s := range cases {
		_, err := Parse(s)
		if err == nil {
			t.Errorf("expected error for %q, got nil", s)
		}
	}
}

func TestParseRejectEmptyBuild(t *testing.T) {
	_, err := Parse("1.0.0+")
	if err == nil {
		t.Error("expected error for empty build metadata, got nil")
	}
}

func TestParseErrorContainsOffendingString(t *testing.T) {
	bad := "1.0.01"
	_, err := Parse(bad)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), bad) {
		t.Errorf("error %q does not contain offending version %q", err.Error(), bad)
	}
}

func TestParseRejectNegativeComponent(t *testing.T) {
	_, err := Parse("-1.0.0")
	if err == nil {
		t.Error("expected error for negative-looking version, got nil")
	}
}

func TestCompareNumericOrdering(t *testing.T) {
	a := mustParse(t, "1.10.0")
	b := mustParse(t, "1.2.3")
	if got := Compare(a, b); got != 1 {
		t.Errorf("Compare(1.10.0, 1.2.3) = %d, want 1", got)
	}
}

func TestCompareBuildMetadataIgnored(t *testing.T) {
	a := mustParse(t, "1.0.0")
	b := mustParse(t, "1.0.0+build.5")
	if got := Compare(a, b); got != 0 {
		t.Errorf("Compare(1.0.0, 1.0.0+build.5) = %d, want 0", got)
	}
}

func TestComparePreReleaseBelowRelease(t *testing.T) {
	a := mustParse(t, "1.0.0-alpha")
	b := mustParse(t, "1.0.0")
	if got := Compare(a, b); got != -1 {
		t.Errorf("Compare(1.0.0-alpha, 1.0.0) = %d, want -1", got)
	}
}

func TestCompareSemVerSpecChain(t *testing.T) {
	// 1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta < 1.0.0-beta < 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0
	chain := []string{
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-alpha.beta",
		"1.0.0-beta",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0",
	}
	for i := 0; i < len(chain)-1; i++ {
		a := mustParse(t, chain[i])
		b := mustParse(t, chain[i+1])
		if got := Compare(a, b); got != -1 {
			t.Errorf("Compare(%s, %s) = %d, want -1", chain[i], chain[i+1], got)
		}
		if got := Compare(b, a); got != 1 {
			t.Errorf("Compare(%s, %s) = %d, want 1", chain[i+1], chain[i], got)
		}
	}
}

func TestCompareNumericVsAlpha(t *testing.T) {
	// numeric identifier "1" ranks below alphanumeric "alpha" at same position
	a := mustParse(t, "1.0.0-1")
	b := mustParse(t, "1.0.0-alpha")
	if got := Compare(a, b); got != -1 {
		t.Errorf("Compare(1.0.0-1, 1.0.0-alpha) = %d, want -1", got)
	}
}

func TestComparePreReleasePrefix(t *testing.T) {
	a := mustParse(t, "1.0.0-alpha")
	b := mustParse(t, "1.0.0-alpha.1")
	if got := Compare(a, b); got != -1 {
		t.Errorf("Compare(1.0.0-alpha, 1.0.0-alpha.1) = %d, want -1", got)
	}
}

func TestCompareNumericBeta(t *testing.T) {
	// beta.2 < beta.11 because 2 < 11 numerically
	a := mustParse(t, "1.0.0-beta.2")
	b := mustParse(t, "1.0.0-beta.11")
	if got := Compare(a, b); got != -1 {
		t.Errorf("Compare(1.0.0-beta.2, 1.0.0-beta.11) = %d, want -1", got)
	}
}

func TestSortMixedList(t *testing.T) {
	inputs := []string{"1.10.0", "1.2.3", "1.0.0-rc.1", "1.0.0"}
	versions := make([]Version, len(inputs))
	for i, s := range inputs {
		versions[i] = mustParse(t, s)
	}
	Sort(versions)
	want := []string{"1.0.0-rc.1", "1.0.0", "1.2.3", "1.10.0"}
	for i, v := range versions {
		if v.String() != want[i] {
			t.Errorf("sorted[%d] = %s, want %s", i, v.String(), want[i])
		}
	}
}

func TestSortStable(t *testing.T) {
	// Two versions equal in M.M.P with different Build must preserve original order.
	a := mustParse(t, "1.0.0+aaa")
	b := mustParse(t, "1.0.0+zzz")
	versions := []Version{a, b}
	Sort(versions)
	if versions[0].Build != "aaa" || versions[1].Build != "zzz" {
		t.Errorf("stable sort violated: got builds %q %q, want aaa zzz", versions[0].Build, versions[1].Build)
	}
}

func TestCompareEqual(t *testing.T) {
	a := mustParse(t, "2.3.4")
	b := mustParse(t, "2.3.4")
	if got := Compare(a, b); got != 0 {
		t.Errorf("Compare(2.3.4, 2.3.4) = %d, want 0", got)
	}
}

func TestStringRoundTrip(t *testing.T) {
	cases := []string{"1.2.3", "1.0.0-alpha.1", "1.0.0+build.5", "1.0.0-rc.1+meta"}
	for _, s := range cases {
		v := mustParse(t, s)
		if got := v.String(); got != s {
			t.Errorf("String() = %q, want %q", got, s)
		}
	}
}
