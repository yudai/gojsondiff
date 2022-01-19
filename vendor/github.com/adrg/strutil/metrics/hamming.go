package metrics

import (
	"strings"
)

// Hamming represents the Hamming metric for measuring the similarity
// between sequences.
//   For more information see https://en.wikipedia.org/wiki/Hamming_distance.
type Hamming struct {
	// CaseSensitive specifies if the string comparison is case sensitive.
	CaseSensitive bool
}

// NewHamming returns a new Hamming string metric.
//
// Default options:
//   CaseSensitive: true
func NewHamming() *Hamming {
	return &Hamming{
		CaseSensitive: true,
	}
}

// Compare returns the Hamming similarity of a and b. The returned
// similarity is a number between 0 and 1. Larger similarity numbers indicate
// closer matches.
func (m *Hamming) Compare(a, b string) float64 {
	distance, maxLen := m.distance(a, b)
	return 1 - float64(distance)/float64(maxLen)
}

// Distance returns the Hamming distance between a and b. Lower distances
// indicate closer matches. A distance of 0 means the strings are identical.
func (m *Hamming) Distance(a, b string) int {
	distance, _ := m.distance(a, b)
	return distance
}

func (m *Hamming) distance(a, b string) (int, int) {
	// Lower terms if case insensitive comparison is specified.
	if !m.CaseSensitive {
		a = strings.ToLower(a)
		b = strings.ToLower(b)
	}
	runesA, runesB := []rune(a), []rune(b)

	// Check if both terms are empty.
	lenA, lenB := len(runesA), len(runesB)
	if lenA == 0 && lenB == 0 {
		return 0, 0
	}

	// If the lengths of the sequences are not equal, the distance is
	// initialized to their absolute difference. Otherwise, it is set to 0.
	if lenA > lenB {
		lenA, lenB = lenB, lenA
	}
	distance := lenB - lenA

	// Calculate Hamming distance.
	for i := 0; i < lenA; i++ {
		if runesA[i] != runesB[i] {
			distance++
		}
	}

	return distance, lenB
}
