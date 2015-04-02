package gojsondiff

import (
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
)

type Diff interface {
	Modified() bool
	Structure() map[string]interface{}
	Iterate(i Iterator) bool
}

type Difference interface{}

func (diff *concreteDiff) Modified() bool {
	if diff.state == unknown {
		if diff.Iterate(&modificationDetector{}) {
			diff.state = modified
		} else {
			diff.state = same
		}
	}
	return diff.state == modified
}

func (diff *concreteDiff) Iterate(i Iterator) bool {
	return !(!i.EnterRoot(len(diff.structure)) &&
		!iterateMap(diff.structure, i) &&
		!i.ExitRoot())
}

func (diff *concreteDiff) Structure() map[string]interface{} {
	return diff.structure
}

type Same concreteDifference
type Added concreteDifference
type Modified concreteDifference
type Deleted concreteDifference

func Compare(
	a []byte,
	b []byte,
) (Diff, error) {
	var aObj, bObj map[string]interface{}
	err := json.Unmarshal(a, &aObj)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &bObj)
	if err != nil {
		return nil, err
	}

	return CompareObjects(aObj, bObj), nil
}

func CompareObjects(
	a map[string]interface{},
	b map[string]interface{},
) Diff {
	structure := CompareMaps(a, b)
	return &concreteDiff{
		structure: structure,
		state:     unknown,
	}
}

func CompareMaps(
	a map[string]interface{},
	b map[string]interface{},
) map[string]interface{} {
	structure := map[string]interface{}{}

	for key, aValue := range a {
		if bValue, ok := b[key]; ok {
			structure[key] = CompareValues(aValue, bValue)
		} else {
			structure[key] = Deleted{OldValue: aValue}
		}
	}
	for key, bValue := range b {
		if _, ok := a[key]; !ok {
			structure[key] = Added{NewValue: bValue}
		}
	}
	return structure
}

func CompareArrays(
	a []interface{},
	b []interface{},
) []interface{} {
	structure := []interface{}{}
	for i, aValue := range a {
		if i < len(b) {
			structure = append(structure, CompareValues(aValue, b[i]))
		} else {
			structure = append(structure, Deleted{OldValue: aValue})
		}
	}
	for i := len(a); i < len(b); i++ {
		structure = append(structure, Added{NewValue: b[i]})
	}
	return structure
}

func CompareValues(
	aValue interface{},
	bValue interface{},
) Difference {
	av := reflect.ValueOf(aValue)
	bv := reflect.ValueOf(bValue)

	if av.Type() != bv.Type() {
		return Modified{OldValue: aValue, NewValue: bValue}
	}

	if av.Kind() == reflect.Map {
		return CompareMaps(aValue.(map[string]interface{}), bValue.(map[string]interface{}))
	} else if av.Kind() == reflect.Slice {
		return CompareArrays(aValue.([]interface{}), bValue.([]interface{}))
	} else if !reflect.DeepEqual(aValue, bValue) {
		return Modified{OldValue: aValue, NewValue: bValue}
	} else {
		return Same{OldValue: aValue, NewValue: bValue}
	}
}

type concreteDiff struct {
	structure map[string]interface{}
	state     state
}

type state int

const (
	unknown  state = -1
	same     state = 0
	modified state = 1
)

type concreteDifference struct {
	OldValue interface{} `json:",omitempty"`
	NewValue interface{} `json:",omitempty"`
}

type modificationDetector struct {
	NullIterator
}

func (i *modificationDetector) VisitAdded(name string, d Added) bool       { return true }
func (i *modificationDetector) VisitModified(name string, d Modified) bool { return true }
func (i *modificationDetector) VisitDeleted(name string, d Deleted) bool   { return true }

func iterateMap(m map[string]interface{}, i Iterator) bool {
	if i.SortObjectFields() {
		keys := []string{}
		for key, _ := range m {
			keys = append(keys, key)
		}
		sort.StringSlice(keys).Sort()
		for _, key := range keys {
			abort := iterateValue(key, m[key], i)
			if abort {
				return true
			}
		}
	} else {
		for key, value := range m {
			abort := iterateValue(key, value, i)
			if abort {
				return true
			}
		}
	}
	return false
}

func iterateSlice(s []interface{}, i Iterator) bool {
	for key, value := range s {
		abort := iterateValue(strconv.Itoa(key), value, i)
		if abort {
			return true
		}
	}
	return false
}

func iterateValue(key string, value interface{}, i Iterator) bool {
	kind := reflect.ValueOf(value).Kind()
	if kind == reflect.Map {
		m := value.(map[string]interface{})
		return !(!i.EnterObject(key, len(m)) &&
			!iterateMap(m, i) &&
			!i.ExitObject(key))
	} else if kind == reflect.Slice {
		s := value.([]interface{})
		return !(!i.EnterArray(key, len(s)) &&
			!iterateSlice(s, i) &&
			!i.ExitArray(key))
	} else {
		switch value.(type) {
		case Same:
			return i.VisitSame(key, value.(Same))
		case Added:
			return i.VisitAdded(key, value.(Added))
		case Modified:
			return i.VisitModified(key, value.(Modified))
		case Deleted:
			return i.VisitDeleted(key, value.(Deleted))
		default:
			panic("Unkwon type found")
		}
	}
}
