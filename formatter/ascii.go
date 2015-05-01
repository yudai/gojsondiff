package formatter

import (
	"errors"
	"fmt"
	"sort"

	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/pp"
)

func NewAsciiFormatter(left interface{}) *AsciiFormatter {
	return &AsciiFormatter{
		left:           left,
		ShowArrayIndex: false,
	}
}

type AsciiFormatter struct {
	left           interface{}
	ShowArrayIndex bool
	buffer         string
	path           []diff.Position
	size           []int
	preShifts      []int
	postShifts     []int
	inArray        []bool
}

func (f *AsciiFormatter) Format(d diff.Diff) (result string, err error) {
	f.buffer = ""
	f.path = []diff.Position{}
	f.size = []int{1}
	f.preShifts = []int{0}
	f.postShifts = []int{0}
	f.inArray = []bool{false}

	f.processItem(f.left, []diff.Delta{d.Delta()}, diff.Root)

	return f.buffer, nil
}

func (f *AsciiFormatter) processArray(array []interface{}, deltas []diff.Delta) error {
	for index, value := range array {
		f.processItem(value, deltas, diff.Index(index))
	}

	// additional Added
	for _, delta := range deltas {
		switch delta.(type) {
		case *diff.Added:
			d := delta.(*diff.Added)
			// skip items already processed
			if int(d.Position.(diff.Index)) < len(array) {
				continue
			}
			f.printRecursive(d.Position, d.Value, AsciiAdded)
		}
	}

	return nil
}

func (f *AsciiFormatter) processObject(object map[string]interface{}, deltas []diff.Delta) error {
	names := sortedKeys(object)
	for _, name := range names {
		value := object[name]
		f.processItem(value, deltas, diff.Name(name))
	}

	// Added
	for _, delta := range deltas {
		switch delta.(type) {
		case *diff.Added:
			d := delta.(*diff.Added)
			f.printRecursive(d.Position, d.Value, AsciiAdded)
		}
	}

	return nil
}

func (f *AsciiFormatter) processItem(value interface{}, deltas []diff.Delta, position diff.Position) error {
	matchedDeltas := f.searchDeltas(deltas, position)
	pp.Println(position)
	pp.Println(*f.preShift())
	pp.Println(*f.postShift())
	pp.Println(value)
	pp.Println(matchedDeltas)
	if len(matchedDeltas) > 0 {
		for _, matchedDelta := range matchedDeltas {

			switch matchedDelta.(type) {
			case *diff.Object:
				d := matchedDelta.(*diff.Object)
				switch value.(type) {
				case map[string]interface{}:
					//ok
				default:
					return errors.New("Type mismatch")
				}
				o := value.(map[string]interface{})

				f.printKeyWithIndent(position, AsciiSame)
				f.println("{")
				f.push(position, len(o), false)
				f.processObject(o, d.Deltas)
				f.pop()
				f.printIndent(AsciiSame)
				f.print("}")
				f.printComma()

			case *diff.Array:
				d := matchedDelta.(*diff.Array)
				switch value.(type) {
				case []interface{}:
					//ok
				default:
					return errors.New("Type mismatch")
				}
				a := value.([]interface{})

				f.printKeyWithIndent(position, AsciiSame)
				f.println("[")
				f.push(position, len(a), true)
				f.processArray(a, d.Deltas)
				f.pop()
				f.printIndent(AsciiSame)
				f.print("]")
				f.printComma()

			case *diff.Added:
				d := matchedDelta.(*diff.Added)
				f.printRecursive(position, d.Value, AsciiAdded)
				*f.postShift()++
				f.printRecursive(position, value, AsciiSame)

			case *diff.Moved:
				d := matchedDelta.(*diff.Moved)
				if position == d.PrePosition() {
					f.printRecursive(position, d.Value, AsciiDeleted)
					*f.preShift()--
				} else {
					f.printRecursive(position, d.Value, AsciiAdded)
					*f.postShift()++
					f.printRecursive(position, value, AsciiSame)
				}

			case *diff.Modified:
				d := matchedDelta.(*diff.Modified)
				savedSize := f.size[len(f.size)-1]
				f.printRecursive(position, d.OldValue, AsciiDeleted)
				f.size[len(f.size)-1] = savedSize
				f.printRecursive(position, d.NewValue, AsciiAdded)

			case *diff.TextDiff:
				savedSize := f.size[len(f.size)-1]
				d := matchedDelta.(*diff.TextDiff)
				f.printRecursive(position, d.OldValue, AsciiDeleted)
				f.size[len(f.size)-1] = savedSize
				f.printRecursive(position, d.NewValue, AsciiAdded)

			case *diff.Deleted:
				d := matchedDelta.(*diff.Deleted)
				f.printRecursive(position, d.Value, AsciiDeleted)
				*f.preShift()--

			default:
				return errors.New("Unknown Delta type detected")
			}

		}
	} else {
		f.printRecursive(position, value, AsciiSame)
	}

	return nil
}

func (f *AsciiFormatter) searchDeltas(deltas []diff.Delta, position diff.Position) (results []diff.Delta) {
	results = make([]diff.Delta, 0)

	for _, delta := range deltas {
		switch delta.(type) {
		case diff.PreDelta:
			if delta.(diff.PreDelta).PrePosition() == position {
				results = append(results, delta)
			}
		}
		switch delta.(type) {
		case diff.PostDelta:
			switch position.(type) {
			case diff.Index:
				if int(delta.(diff.PostDelta).PostPosition().(diff.Index)) ==
					int(position.(diff.Index))+*f.preShift()+*f.postShift() {
					results = append(results, delta)
				}
			default:
				if delta.(diff.PostDelta).PostPosition() == position {
					results = append(results, delta)
				}
			}
		}
	}
	return
}

const (
	AsciiSame    = " "
	AsciiAdded   = "+"
	AsciiDeleted = "-"
)

func (f *AsciiFormatter) push(position diff.Position, size int, array bool) {
	f.path = append(f.path, position)
	f.size = append(f.size, size)
	f.preShifts = append(f.preShifts, 0)
	f.postShifts = append(f.postShifts, 0)
	f.inArray = append(f.inArray, array)
}

func (f *AsciiFormatter) pop() {
	f.path = f.path[0 : len(f.path)-1]
	f.size = f.size[0 : len(f.size)-1]
	f.preShifts = f.preShifts[0 : len(f.preShifts)-1]
	f.postShifts = f.postShifts[0 : len(f.postShifts)-1]
	f.inArray = f.inArray[0 : len(f.inArray)-1]
}

func (f *AsciiFormatter) preShift() *int {
	return &f.preShifts[len(f.preShifts)-1]
}

func (f *AsciiFormatter) postShift() *int {
	return &f.postShifts[len(f.postShifts)-1]
}

func (f *AsciiFormatter) printIndent(marker string) {
	f.print(marker)
	for n := 0; n < len(f.path); n++ {
		f.print("  ")
	}
}

func (f *AsciiFormatter) printKeyWithIndent(position diff.Position, marker string) {
	f.printIndent(marker)
	if position == diff.Root {
		// nothing to do
	} else if !f.inArray[len(f.inArray)-1] {
		f.printf(`"%s": `, position.String())
	} else if f.ShowArrayIndex {
		f.printf(`%d: `, (int(position.(diff.Index)) + *f.preShift() + *f.postShift()))
	}
}

func (f *AsciiFormatter) printComma() {
	f.size[len(f.size)-1]--
	if f.size[len(f.size)-1] + +*f.preShift() + *f.postShift() > 0 {
		f.println(",")
	} else {
		f.println()
	}
}

func (f *AsciiFormatter) printValue(value interface{}) {
	switch value.(type) {
	case string:
		f.buffer += fmt.Sprintf(`"%s"`, value)
	default:
		f.buffer += fmt.Sprintf(`%#v`, value)
	}
}

func (f *AsciiFormatter) print(a ...interface{}) {
	f.buffer += fmt.Sprint(a...)
}

func (f *AsciiFormatter) printf(format string, a ...interface{}) {
	f.buffer += fmt.Sprintf(format, a...)
}

func (f *AsciiFormatter) println(a ...interface{}) {
	f.buffer += fmt.Sprintln(a...)
}

func (f *AsciiFormatter) printRecursive(position diff.Position, value interface{}, marker string) {
	switch value.(type) {
	case map[string]interface{}:
		f.printKeyWithIndent(position, marker)
		f.println("{")

		m := value.(map[string]interface{})
		size := len(m)
		f.push(position, size, false)

		keys := sortedKeys(m)
		for _, key := range keys {
			f.printRecursive(diff.Name(key), m[key], marker)
		}
		f.pop()

		f.printIndent(marker)
		f.print("}")
		f.printComma()
	case []interface{}:
		f.printKeyWithIndent(position, marker)
		f.println("[")

		s := value.([]interface{})
		size := len(s)
		f.push(position, size, true)
		for i, item := range s {
			f.printRecursive(diff.Index(i), item, marker)
		}
		f.pop()

		f.printIndent(marker)
		f.print("]")
		f.printComma()
	default:
		f.printKeyWithIndent(position, marker)
		f.printValue(value)
		f.printComma()
	}
}

func sortedKeys(m map[string]interface{}) (keys []string) {
	keys = make([]string, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}
