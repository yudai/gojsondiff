package metrics

import (
	"strings"
	"unicode/utf8"

	"github.com/adrg/strutil/internal/util"
)

// JaroWinkler represents the Jaro-Winkler metric for measuring the similarity
// between sequences.
//   For more information see https://en.wikipedia.org/wiki/Jaro-Winkler_distance.
type JaroWinkler struct {
	// CaseSensitive specifies if the string comparison is case sensitive.
	CaseSensitive bool
}

// NewJaroWinkler returns a new Jaro-Winkler string metric.
//
// Default options:
//   CaseSensitive: true
func NewJaroWinkler() *JaroWinkler {
	return &JaroWinkler{
		CaseSensitive: true,
	}
}

// Compare returns the Jaro-Winkler similarity of a and b. The returned
// similarity is a number between 0 and 1. Larger similarity numbers indicate
// closer matches.
func (m *JaroWinkler) Compare(a, b string) float64 {
	// Lower terms if case insensitive comparison is specified.
	if !m.CaseSensitive {
		a = strings.ToLower(a)
		b = strings.ToLower(b)
	}

	// Calculate common prefix.
	lenPrefix := utf8.RuneCountInString(util.CommonPrefix(a, b))
	if lenPrefix > 4 {
		lenPrefix = 4
	}

	jaro := NewJaro()
	jaro.CaseSensitive = m.CaseSensitive

	// Return similarity.
	similarity := jaro.Compare(a, b)
	return similarity + (0.1 * float64(lenPrefix) * (1.0 - similarity))
}
