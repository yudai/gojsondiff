package metrics

import (
	"strings"

	"github.com/adrg/strutil/internal/util"
)

// OverlapCoefficient represents the overlap coefficient for measuring the
// similarity between sequences. The metric is also know as the
// Szymkiewicz-Simpson coefficient.
//   For more information see https://en.wikipedia.org/wiki/Overlap_coefficient.
type OverlapCoefficient struct {
	// CaseSensitive specifies if the string comparison is case sensitive.
	CaseSensitive bool

	// NgramSize represents the size (in characters) of the tokens generated
	// when comparing the input sequences.
	NgramSize int
}

// NewOverlapCoefficient returns a new overlap coefficient string metric.
//
// Default options:
//   CaseSensitive: true
//   NGramSize: 2
func NewOverlapCoefficient() *OverlapCoefficient {
	return &OverlapCoefficient{
		CaseSensitive: true,
		NgramSize:     2,
	}
}

// Compare returns the OverlapCoefficient similarity coefficient of a and b.
// The returned similarity is a number between 0 and 1. Larger similarity
// numbers indicate closer matches.
// An n-gram size of 2 is used if the provided size is less than or equal to 0.
func (m *OverlapCoefficient) Compare(a, b string) float64 {
	// Lower terms if case insensitive comparison is specified.
	if !m.CaseSensitive {
		a = strings.ToLower(a)
		b = strings.ToLower(b)
	}

	// Check if both terms are empty.
	runesA, runesB := []rune(a), []rune(b)
	if len(runesA) == 0 && len(runesB) == 0 {
		return 1
	}

	size := m.NgramSize
	if size <= 0 {
		size = 2
	}

	// Calculate n-gram intersection and minimum subset.
	_, common, totalA, totalB := util.NgramIntersection(runesA, runesB, size)

	min := util.Min(totalA, totalB)
	if min == 0 {
		return 0
	}

	// Return similarity.
	return float64(common) / float64(min)
}
