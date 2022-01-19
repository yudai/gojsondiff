package metrics

import (
	"strings"
	"unicode/utf8"

	"github.com/adrg/strutil/internal/util"
)

// Jaro represents the Jaro metric for measuring the similarity
// between sequences.
//   For more information see https://en.wikipedia.org/wiki/Jaro-Winkler_distance.
type Jaro struct {
	// CaseSensitive specifies if the string comparison is case sensitive.
	CaseSensitive bool
}

// NewJaro returns a new Jaro string metric.
//
// Default options:
//   CaseSensitive: true
func NewJaro() *Jaro {
	return &Jaro{
		CaseSensitive: true,
	}
}

// Compare returns the Jaro similarity of a and b. The returned similarity is
// a number between 0 and 1. Larger similarity numbers indicate closer matches.
func (m *Jaro) Compare(a, b string) float64 {
	// Check if both terms are empty.
	lenA, lenB := utf8.RuneCountInString(a), utf8.RuneCountInString(b)
	if lenA == 0 && lenB == 0 {
		return 1
	}

	// Check if one of the terms is empty.
	if lenA == 0 || lenB == 0 {
		return 0
	}

	// Lower terms if case insensitive comparison is specified.
	if !m.CaseSensitive {
		a = strings.ToLower(a)
		b = strings.ToLower(b)
	}

	// Get matching runes.
	halfLen := util.Max(0, util.Max(lenA, lenB)/2)
	mrA := matchingRunes(a, b, halfLen)
	mrB := matchingRunes(b, a, halfLen)

	fmLen, smLen := len(mrA), len(mrB)
	if fmLen == 0 || smLen == 0 {
		return 0.0
	}

	// Return similarity.
	return (float64(fmLen)/float64(lenA) +
		float64(smLen)/float64(lenB) +
		float64(fmLen-transpositions(mrA, mrB)/2)/float64(fmLen)) / 3.0
}

func matchingRunes(a, b string, limit int) []rune {
	common := []rune{}
	runesB := []rune(b)
	lenB := len(runesB)

	for i, r := range a {
		end := util.Min(i+limit+1, lenB)
		for j := util.Max(0, i-limit); j < end; j++ {
			if r == runesB[j] && runesB[j] != -1 {
				common = append(common, runesB[j])
				runesB[j] = -1
				break
			}
		}
	}

	return common
}

func transpositions(a, b []rune) int {
	var count int

	minLen := util.Min(len(a), len(b))
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			count++
		}
	}

	return count
}
