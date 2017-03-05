package gojsondiff

import (
	"container/list"
	"encoding/json"
	"reflect"

	dmp "github.com/sergi/go-diff/diffmatchpatch"
	"github.com/yudai/golcs"
)

// Compare finds changes between two JSON values and returnes as a Diff.
func (differ *Differ) Compare(
	left interface{},
	right interface{},
) Diff {
	_, delta := differ.compareValues(left, right)

	return &diff{delta: delta}
}

// CompareBytes do the same as Compare() but accepts JSON strings as []byte.
func (differ *Differ) CompareBytes(
	left []byte,
	right []byte,
) (Diff, error) {
	var leftValue, rightValue interface{}
	err := json.Unmarshal(left, &leftValue)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(right, &rightValue)
	if err != nil {
		return nil, err
	}

	_, delta := differ.compareValues(leftValue, rightValue)

	return &diff{delta: delta}, nil
}

func (differ *Differ) compareMaps(
	left map[string]interface{},
	right map[string]interface{},
) (deltas map[string]Delta) {
	deltas = make(map[string]Delta, 0)

	names := sortedKeys(left) // stabilize delta order
	for _, name := range names {
		if rightValue, ok := right[name]; ok {
			same, delta := differ.compareValues(left[name], rightValue)
			if !same {
				deltas[name] = delta
			}
		} else {
			deltas[name] = NewDeleted(left[name])
		}
	}

	// find added items
	names = sortedKeys(right) // stabilize delta order
	for _, name := range names {
		if _, ok := left[name]; !ok {
			deltas[name] = NewAdded(right[name])
		}
	}

	return deltas
}

// item in an array
type maybe struct {
	index    int
	lcsIndex int
	item     interface{}
}

func (differ *Differ) compareArrays(
	left []interface{},
	right []interface{},
) (preDeltas map[int]Delta, postDeltas map[int]Delta) {
	// for delete, move
	preDeltas = make(map[int]Delta, 0)
	// for add, modify
	postDeltas = make(map[int]Delta, 0)

	// LCS index pairs, which holds the indexes of matched items
	lcsPairs := lcs.New(left, right).IndexPairs()

	// list up items in the left not in the index pairs, they are maybe deleted
	maybeDeleted := list.New() // but could be moved or modified
	lcsI := 0
	for i, leftValue := range left {
		if lcsI < len(lcsPairs) && lcsPairs[lcsI].Left == i {
			lcsI++
		} else {
			maybeDeleted.PushBack(maybe{index: i, lcsIndex: lcsI, item: leftValue})
		}
	}

	// list up items not in the right LCS, they are maybe Added
	maybeAdded := list.New() // but could be moved or modified
	lcsI = 0
	for i, rightValue := range right {
		if lcsI < len(lcsPairs) && lcsPairs[lcsI].Right == i {
			lcsI++
		} else {
			maybeAdded.PushBack(maybe{index: i, lcsIndex: lcsI, item: rightValue})
		}
	}

	// find moved items by comparing maybeAdded and maybeDeleted using a nasted loop
	var delNext *list.Element // for prefetch to remove item in iteration
	for delCandidate := maybeDeleted.Front(); delCandidate != nil; delCandidate = delNext {
		delCan := delCandidate.Value.(maybe)
		delNext = delCandidate.Next()

		for addCandidate := maybeAdded.Front(); addCandidate != nil; addCandidate = addCandidate.Next() {
			addCan := addCandidate.Value.(maybe)
			if reflect.DeepEqual(delCan.item, addCan.item) {
				// found matched item in added and deleted, which means it's moved
				preDeltas[delCan.index] = NewMoved(addCan.index, delCan.item, nil)
				// remove from the lists
				maybeAdded.Remove(addCandidate)
				maybeDeleted.Remove(delCandidate)
				break
			}
		}
	}

	// find modified or add+del
	// we want to maximise the total smilality of items, so use LCS here again
	// in each block fenced by the lcs pairs.
	prevIndexDel := 0
	prevIndexAdd := 0
	delElement := maybeDeleted.Front()
	addElement := maybeAdded.Front()
	for i := 0; i <= len(lcsPairs); i++ { // not "< len(lcsPairs)"
		var lcsPair lcs.IndexPair
		var delSize, addSize int

		// count items in the deleted and added lists
		// between the previous LCS pair and the current LCS pair
		if i < len(lcsPairs) {
			lcsPair = lcsPairs[i]
			delSize = lcsPair.Left - prevIndexDel - 1
			addSize = lcsPair.Right - prevIndexAdd - 1
			prevIndexDel = lcsPair.Left
			prevIndexAdd = lcsPair.Right
		}

		// collect deleted items between two LCS pair
		var delSlice []maybe
		if delSize > 0 {
			delSlice = make([]maybe, 0, delSize)
		} else {
			delSlice = make([]maybe, 0, maybeDeleted.Len())
		}
		for ; delElement != nil; delElement = delElement.Next() {
			d := delElement.Value.(maybe)
			if d.lcsIndex != i {
				break
			}
			delSlice = append(delSlice, d)
		}

		// collect added items between two LCS pair
		var addSlice []maybe
		if addSize > 0 {
			addSlice = make([]maybe, 0, addSize)
		} else {
			addSlice = make([]maybe, 0, maybeAdded.Len())
		}
		for ; addElement != nil; addElement = addElement.Next() {
			a := addElement.Value.(maybe)
			if a.lcsIndex != i {
				break
			}
			addSlice = append(addSlice, a)
		}

		if len(delSlice) > 0 && len(addSlice) > 0 {
			var bestDeltas map[int]Delta
			// calculate LCS here again
			bestDeltas, delSlice, addSlice = differ.maximizeSimilarities(delSlice, addSlice)
			for key, delta := range bestDeltas {
				// modifed
				postDeltas[key] = delta
			}
		}

		for _, del := range delSlice {
			preDeltas[del.index] = NewDeleted(del.item)
		}
		for _, add := range addSlice {
			postDeltas[add.index] = NewAdded(add.item)
		}
	}

	return preDeltas, postDeltas
}

func (differ *Differ) compareValues(
	left interface{},
	right interface{},
) (same bool, delta Delta) {
	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return false, NewModified(left, right)
	}

	switch l := left.(type) {
	case map[string]interface{}:
		childDeltas := differ.compareMaps(l, right.(map[string]interface{}))
		if len(childDeltas) > 0 {
			return false, NewObject(childDeltas)
		}

	case []interface{}:
		preChildDeltas, postChildDeltas := differ.compareArrays(l, right.([]interface{}))

		if len(preChildDeltas) > 0 || len(postChildDeltas) > 0 {
			return false, NewArray(preChildDeltas, postChildDeltas)
		}

	default:
		if !reflect.DeepEqual(left, right) {
			if reflect.ValueOf(left).Kind() == reflect.String &&
				reflect.ValueOf(right).Kind() == reflect.String &&
				differ.textDiffMinimumLength <= len(left.(string)) {
				textDiff := dmp.New()
				patchs := textDiff.PatchMake(left.(string), right.(string))
				return false, NewTextDiff(patchs, left, right)

			} else {
				x := NewModified(left, right)
				return false, x
			}
		}
	}

	return true, nil
}

// array LCS
func (differ *Differ) maximizeSimilarities(left []maybe, right []maybe) (resultDeltas map[int]Delta, freeLeft, freeRight []maybe) {
	deltaTable := make([][]Delta, len(left))
	for i := 0; i < len(left); i++ {
		deltaTable[i] = make([]Delta, len(right))
	}

	for i, leftValue := range left {
		for j, rightValue := range right {
			same, delta := differ.compareValues(leftValue.item, rightValue.item)
			if same {
				panic("unexpected samve value in LCS")
			}
			deltaTable[i][j] = delta
		}
	}

	sizeX := len(left) + 1 // margins for both sides
	sizeY := len(right) + 1

	// fill out with similarities
	dpTable := make([][]float64, sizeX)
	for i := 0; i < sizeX; i++ {
		dpTable[i] = make([]float64, sizeY)
	}

	for x := sizeX - 2; x >= 0; x-- {
		for y := sizeY - 2; y >= 0; y-- {
			prevX := dpTable[x+1][y]
			prevY := dpTable[x][y+1]
			score := deltaTable[x][y].Similarity() + dpTable[x+1][y+1]

			dpTable[x][y] = max(prevX, prevY, score)
		}
	}

	minLength := len(left)
	if minLength > len(right) {
		minLength = len(right)
	}
	maxInvalidLength := minLength - 1

	freeLeft = make([]maybe, 0, len(left)-minLength)
	freeRight = make([]maybe, 0, len(right)-minLength)
	resultDeltas = make(map[int]Delta, minLength)
	var x, y int
	for x, y = 0, 0; x <= sizeX-2 && y <= sizeY-2; {
		current := dpTable[x][y]
		nextX := dpTable[x+1][y]
		nextY := dpTable[x][y+1]

		xValidLength := len(left) - maxInvalidLength + y
		yValidLength := len(right) - maxInvalidLength + x

		if x+1 < xValidLength && current == nextX {
			freeLeft = append(freeLeft, left[x])
			x++
		} else if y+1 < yValidLength && current == nextY {
			freeRight = append(freeRight, right[y])
			y++
		} else {
			resultDeltas[right[y].index] = deltaTable[x][y]
			x++
			y++
		}
	}
	for ; x < sizeX-1; x++ {
		freeLeft = append(freeLeft, left[x-1])
	}
	for ; y < sizeY-1; y++ {
		freeRight = append(freeRight, right[y-1])
	}

	return resultDeltas, freeLeft, freeRight
}
