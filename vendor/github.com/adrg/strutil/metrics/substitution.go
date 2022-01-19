package metrics

// Substitution represents a substitution function which is used to
// calculate a score for character substitutions.
type Substitution interface {
	// Compare returns the substitution score of characters a[idxA] and b[idxB].
	Compare(a []rune, idxA int, b []rune, idxB int) float64

	// Returns the maximum score of a character substitution operation.
	Max() float64

	// Returns the minimum score of a character substitution operation.
	Min() float64
}
