package metrics

import (
	"strings"

	"github.com/adrg/strutil/internal/util"
)

// SorensenDice represents the Sorensen-Dice metric for measuring the
// similarity between sequences.
//   For more information see https://en.wikipedia.org/wiki/Sorensen-Dice_coefficient.
type SorensenDice struct {
	// CaseSensitive specifies if the string comparison is case sensitive.
	CaseSensitive bool

	// NgramSize represents the size (in characters) of the tokens generated
	// when comparing the input sequences.
	NgramSize int
}

// NewSorensenDice returns a new Sorensen-Dice string metric.
//
// Default options:
//   CaseSensitive: true
//   NGramSize: 2
func NewSorensenDice() *SorensenDice {
	return &SorensenDice{
		CaseSensitive: true,
		NgramSize:     2,
	}
}

// Compare returns the Sorensen-Dice similarity coefficient of a and b. The
// returned similarity is a number between 0 and 1. Larger similarity numbers
// indicate closer matches.
// An n-gram size of 2 is used if the provided size is less than or equal to 0.
func (m *SorensenDice) Compare(a, b string) float64 {
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

	// Calculate n-gram intersection and union.
	_, common, totalA, totalB := util.NgramIntersection(runesA, runesB, size)

	total := totalA + totalB
	if total == 0 {
		return 0
	}

	// Return similarity.
	return 2 * float64(common) / float64(total)
}
