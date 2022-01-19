package metrics

// MatchMismatch represents a substitution function which returns the match or
// mismatch value depeding on the equality of the compared characters. The
// match value must be greater than the mismatch value.
type MatchMismatch struct {
	// Match represents the score of equal character substitutions.
	Match float64

	// Mismatch represents the score of unequal character substitutions.
	Mismatch float64
}

// Compare returns the match value if a[idxA] is equal to b[idxB] or the
// mismatch value otherwise.
func (m MatchMismatch) Compare(a []rune, idxA int, b []rune, idxB int) float64 {
	if a[idxA] == b[idxB] {
		return m.Match
	}

	return m.Mismatch
}

// Max returns the match value.
func (m MatchMismatch) Max() float64 {
	return m.Match
}

// Min returns the mismatch value.
func (m MatchMismatch) Min() float64 {
	return m.Mismatch
}
