package gojsondiff

import (
	"errors"
	"reflect"

	dmp "github.com/sergi/go-diff/diffmatchpatch"
	"github.com/yudai/golcs"
)

// A Delta represents an atomic difference between two JSON values.
type Delta interface {
	// Similarity calculates the similarity of the two values.
	// The returned value is normalized from 0 to 1,
	// 0 means that the values are completely different and 1 means the oppsite.
	Similarity() (similarity float64)
}

// To cache the calculated similarity,
// concrete Deltas can use similariter and similarityCache
type similariter interface {
	similarity() (similarity float64)
}

type similarityCache struct {
	similariter
	value float64
}

func newSimilarityCache(sim similariter) similarityCache {
	cache := similarityCache{similariter: sim, value: -1}
	return cache
}

func (cache similarityCache) Similarity() (similarity float64) {
	if cache.value < 0 {
		cache.value = cache.similariter.similarity()
	}
	return cache.value
}

// An Object is a Delta that represents changes in a JSON object
type Object struct {
	similarityCache

	// Deltas holds internal changes with property names
	Deltas map[string]Delta
}

// NewObject returns an Object
func NewObject(deltas map[string]Delta) *Object {
	d := Object{Deltas: deltas}
	d.similarityCache = newSimilarityCache(&d)
	return &d
}

func (d *Object) similarity() (similarity float64) {
	for _, delta := range d.Deltas {
		similarity += delta.Similarity()
	}
	return similarity / float64(len(d.Deltas))
}

// An Array is a Delta that represents changes in a JSON array
type Array struct {
	similarityCache

	// PreDeltas hold internal changes with index values which
	// indicate the positions in the original JSON data.
	// The map may have Delted and Moved Deltas.
	PreDeltas map[int]Delta

	// PostDeltas hold internal changes with index values which
	// indicate the positions in the changed JSON data.
	// The map may have Added and Modified Deltas.
	PostDeltas map[int]Delta
}

// NewArray returns an Array
func NewArray(preDeltas map[int]Delta, postDeltas map[int]Delta) *Array {
	d := Array{PreDeltas: preDeltas, PostDeltas: postDeltas}
	d.similarityCache = newSimilarityCache(&d)
	return &d
}

func (d *Array) similarity() (similarity float64) {
	for _, delta := range d.PreDeltas {
		similarity += delta.Similarity()
	}
	for _, delta := range d.PostDeltas {
		similarity += delta.Similarity()
	}

	return similarity / float64(len(d.PreDeltas)+len(d.PostDeltas))
}

// WithoutMoves returns internal changes without Moved Deltas.
// A Moved is decomposed to an Add and a Deleted
func (d *Array) WithoutMoved() (preDeltas, postDeltas map[int]Delta) {
	preDeltas = make(map[int]Delta, len(d.PreDeltas))
	postDeltas = make(map[int]Delta, len(d.PostDeltas))

	for key, delta := range d.PreDeltas {
		switch dl := delta.(type) {
		case *Moved:
			preDeltas[key] = NewDeleted(dl.Value)
			postDeltas[dl.NewPosition] = NewAdded(dl.Value)
		default:
			preDeltas[key] = delta
		}
	}

	for key, delta := range d.PostDeltas {
		postDeltas[key] = delta
	}

	return preDeltas, postDeltas
}

// An Added represents a new added field of an object or an array
type Added struct {
	similarityCache

	// Value holds the added value
	Value interface{}
}

// NewAdded returns a new Added
func NewAdded(value interface{}) *Added {
	d := Added{Value: value}
	return &d
}

func (d *Added) similarity() (similarity float64) {
	return 0
}

// A Modified represents a change of a value
type Modified struct {
	similarityCache

	// The value before modification
	OldValue interface{}

	// The value after modification
	NewValue interface{}
}

// NewModified returns a Modified
func NewModified(oldValue, newValue interface{}) *Modified {
	d := Modified{
		OldValue: oldValue,
		NewValue: newValue,
	}
	d.similarityCache = newSimilarityCache(&d)
	return &d

}

func (d *Modified) similarity() (similarity float64) {
	similarity += 0.3 // at least, they are at the same position
	if reflect.TypeOf(d.OldValue) == reflect.TypeOf(d.NewValue) {
		similarity += 0.3 // types are same

		switch d.OldValue.(type) {
		case string:
			similarity += 0.4 * stringSimilarity(d.OldValue.(string), d.NewValue.(string))
		case float64:
			ratio := d.OldValue.(float64) / d.NewValue.(float64)
			if ratio > 1 {
				ratio = 1 / ratio
			}
			similarity += 0.4 * ratio
		}
	}
	return
}

func stringSimilarity(left, right string) (similarity float64) {
	matchingLength := float64(
		lcs.New(
			stringToInterfaceSlice(left),
			stringToInterfaceSlice(right),
		).Length(),
	)
	similarity =
		(matchingLength / float64(len(left))) * (matchingLength / float64(len(right)))
	return
}

func stringToInterfaceSlice(str string) []interface{} {
	s := make([]interface{}, len(str))
	for i, v := range str {
		s[i] = v
	}
	return s
}

// A TextDiff represents a Modified with TextDiff between the old and the new values.
type TextDiff struct {
	Modified

	// Diff string
	Diff []dmp.Patch
}

// NewTextDiff returns
func NewTextDiff(diff []dmp.Patch, oldValue, newValue interface{}) *TextDiff {
	d := TextDiff{
		Modified: *NewModified(oldValue, newValue),
		Diff:     diff,
	}
	return &d
}

func (d *TextDiff) patch() error {
	if d.OldValue == nil {
		return errors.New("Old Value is not set")
	}
	patcher := dmp.New()
	patched, successes := patcher.PatchApply(d.Diff, d.OldValue.(string))
	for _, success := range successes {
		if !success {
			return errors.New("Failed to apply a patch")
		}
	}
	d.NewValue = patched
	return nil
}

func (d *TextDiff) DiffString() string {
	dmp := dmp.New()
	return dmp.PatchToText(d.Diff)
}

// A Delted represents deleted field or index of an Object or an Array.
type Deleted struct {
	// The value deleted
	Value interface{}
}

// NewDeleted returns a Deleted
func NewDeleted(value interface{}) *Deleted {
	d := Deleted{
		Value: value,
	}
	return &d

}

func (d Deleted) Similarity() (similarity float64) {
	return 0
}

// A Moved represents field that is moved, which means the index or name is
// changed. Note that, in this library, assigning a Moved and a Modified to
// a single position is not allowed. For the compatibility with jsondiffpatch,
// the Moved in this library can hold the old and new value in it.
type Moved struct {
	// The indexx moved to
	NewPosition int

	// The value before moving
	Value interface{}
	// The delta applied after moving (for compatibility)
	Delta interface{}

	similarityCache
}

func NewMoved(newPosition int, value interface{}, delta Delta) *Moved {
	d := Moved{
		NewPosition: newPosition,
		Value:       value,
		Delta:       delta,
	}
	d.similarityCache = newSimilarityCache(&d)
	return &d
}

func (d *Moved) similarity() (similarity float64) {
	similarity = 0.6 // as type and contens are same
	/*
		ratio := float64(d.PrePosition().(Index)) / float64(d.PostPosition().(Index))
		if ratio > 1 {
			ratio = 1 / ratio
		}
		similarity += 0.4 * ratio
	*/
	return
}
