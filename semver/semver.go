package semver

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease []string
	Build      string
}

func Parse(s string) (Version, error) {
	rest := s
	var build string

	if idx := strings.Index(rest, "+"); idx >= 0 {
		build = rest[idx+1:]
		if build == "" {
			return Version{}, fmt.Errorf("invalid version %q: empty build metadata", s)
		}
		rest = rest[:idx]
	}

	var preRelease []string
	if idx := strings.Index(rest, "-"); idx >= 0 {
		pre := rest[idx+1:]
		if pre == "" {
			return Version{}, fmt.Errorf("invalid version %q: empty pre-release identifier", s)
		}
		parts := strings.Split(pre, ".")
		for _, p := range parts {
			if p == "" {
				return Version{}, fmt.Errorf("invalid version %q: empty pre-release identifier", s)
			}
		}
		preRelease = parts
		rest = rest[:idx]
	}

	core := strings.Split(rest, ".")
	if len(core) != 3 {
		return Version{}, fmt.Errorf("invalid version %q: expected MAJOR.MINOR.PATCH", s)
	}

	nums := make([]int, 3)
	for i, part := range core {
		if part == "" {
			return Version{}, fmt.Errorf("invalid version %q: expected MAJOR.MINOR.PATCH", s)
		}
		for _, c := range part {
			if c < '0' || c > '9' {
				return Version{}, fmt.Errorf("invalid version %q: non-numeric component", s)
			}
		}
		if len(part) > 1 && part[0] == '0' {
			return Version{}, fmt.Errorf("invalid version %q: leading zero in numeric identifier", s)
		}
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 {
			return Version{}, fmt.Errorf("invalid version %q: invalid numeric component", s)
		}
		nums[i] = n
	}

	return Version{
		Major:      nums[0],
		Minor:      nums[1],
		Patch:      nums[2],
		PreRelease: preRelease,
		Build:      build,
	}, nil
}

func Compare(a, b Version) int {
	if d := cmpInt(a.Major, b.Major); d != 0 {
		return d
	}
	if d := cmpInt(a.Minor, b.Minor); d != 0 {
		return d
	}
	if d := cmpInt(a.Patch, b.Patch); d != 0 {
		return d
	}

	aPre := len(a.PreRelease) > 0
	bPre := len(b.PreRelease) > 0
	if !aPre && !bPre {
		return 0
	}
	if aPre && !bPre {
		return -1
	}
	if !aPre && bPre {
		return 1
	}

	n := len(a.PreRelease)
	if len(b.PreRelease) < n {
		n = len(b.PreRelease)
	}
	for i := 0; i < n; i++ {
		d := cmpIdentifier(a.PreRelease[i], b.PreRelease[i])
		if d != 0 {
			return d
		}
	}

	return cmpInt(len(a.PreRelease), len(b.PreRelease))
}

func cmpInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func cmpIdentifier(a, b string) int {
	an, aerr := strconv.Atoi(a)
	bn, berr := strconv.Atoi(b)

	aNum := aerr == nil
	bNum := berr == nil

	switch {
	case aNum && bNum:
		return cmpInt(an, bn)
	case aNum && !bNum:
		return -1
	case !aNum && bNum:
		return 1
	default:
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
}

func Sort(versions []Version) {
	sort.SliceStable(versions, func(i, j int) bool {
		return Compare(versions[i], versions[j]) < 0
	})
}

func (v Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if len(v.PreRelease) > 0 {
		s += "-" + strings.Join(v.PreRelease, ".")
	}
	if v.Build != "" {
		s += "+" + v.Build
	}
	return s
}
