package metrics

import (
	"strings"

	"github.com/adrg/strutil/internal/util"
)

// SmithWatermanGotoh represents the Smith-Waterman-Gotoh metric for measuring
// the similarity between sequences.
//   For more information see https://en.wikipedia.org/wiki/Smith-Waterman_algorithm.
type SmithWatermanGotoh struct {
	// CaseSensitive specifies if the string comparison is case sensitive.
	CaseSensitive bool

	// GapPenalty defines a score penalty for character insertions or deletions.
	// For relevant results, the gap penalty should be a non-positive number.
	GapPenalty float64

	// Substitution represents a substitution function which is used to
	// calculate a score for character substitutions.
	Substitution Substitution
}

// NewSmithWatermanGotoh returns a new Smith-Waterman-Gotoh string metric.
//
// Default options:
//   CaseSensitive: true
//   GapPenalty: -0.5
//   Substitution: MatchMismatch{
//   	Match:    1,
//   	Mismatch: -2,
//   },
func NewSmithWatermanGotoh() *SmithWatermanGotoh {
	return &SmithWatermanGotoh{
		CaseSensitive: true,
		GapPenalty:    -0.5,
		Substitution: MatchMismatch{
			Match:    1,
			Mismatch: -2,
		},
	}
}

// Compare returns the Smith-Waterman-Gotoh similarity of a and b. The returned
// similarity is a number between 0 and 1. Larger similarity numbers indicate
// closer matches.
func (m *SmithWatermanGotoh) Compare(a, b string) float64 {
	gap := m.GapPenalty

	// Lower terms if case insensitive comparison is specified.
	if !m.CaseSensitive {
		a = strings.ToLower(a)
		b = strings.ToLower(b)
	}
	runesA, runesB := []rune(a), []rune(b)

	// Check if both terms are empty.
	lenA, lenB := len(runesA), len(runesB)
	if lenA == 0 && lenB == 0 {
		return 1
	}

	// Check if one of the terms is empty.
	if lenA == 0 || lenB == 0 {
		return 0
	}

	// Use default substitution, if none is specified.
	subst := m.Substitution
	if subst == nil {
		subst = MatchMismatch{
			Match:    1,
			Mismatch: -2,
		}
	}

	// Calculate max distance.
	maxDistance := util.Minf(float64(lenA), float64(lenB)) * util.Maxf(subst.Max(), gap)

	// Calculate distance.
	v0 := make([]float64, lenB)
	v1 := make([]float64, lenB)

	distance := util.Maxf(0, gap, subst.Compare(runesA, 0, runesB, 0))
	v0[0] = distance

	for i := 1; i < lenB; i++ {
		v0[i] = util.Maxf(0, v0[i-1]+gap, subst.Compare(runesA, 0, runesB, i))
		distance = util.Maxf(distance, v0[i])
	}

	for i := 1; i < lenA; i++ {
		v1[0] = util.Maxf(0, v0[0]+gap, subst.Compare(runesA, i, runesB, 0))
		distance = util.Maxf(distance, v1[0])

		for j := 1; j < lenB; j++ {
			v1[j] = util.Maxf(0, v0[j]+gap, v1[j-1]+gap, v0[j-1]+subst.Compare(runesA, i, runesB, j))
			distance = util.Maxf(distance, v1[j])
		}

		for j := 0; j < lenB; j++ {
			v0[j] = v1[j]
		}
	}

	// Return similarity.
	return distance / maxDistance
}
