package formatter

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	gdiff "github.com/yudai/gojsondiff"
)

func NewAsciiFormatter(left interface{}, config AsciiFormatterConfig) *AsciiFormatter {
	return &AsciiFormatter{
		left:   left,
		config: config,
	}
}

type AsciiFormatter struct {
	left   interface{}
	config AsciiFormatterConfig
	buffer *bytes.Buffer
	line   *AsciiLine

	path  []string
	lasts []bool
}

type AsciiFormatterConfig struct {
	ShowArrayIndex bool
	Coloring       bool
}

var AsciiFormatterDefaultConfig = AsciiFormatterConfig{}

type AsciiLine struct {
	marker string
	indent int
	buffer *bytes.Buffer
}

func (f *AsciiFormatter) Format(diff gdiff.Diff) (result string, err error) {
	f.buffer = bytes.NewBuffer([]byte{})
	f.path = []string{}
	f.lasts = []bool{}

	f.processItem(f.left, diff.Delta())

	return f.buffer.String(), nil
}

func (f *AsciiFormatter) processArray(array []interface{}, delta *gdiff.Array) error {
	preDeltas, postDeltas := delta.WithoutMoved()

	keys := make([]int, 0, len(preDeltas)+len(postDeltas))
	for key, _ := range preDeltas {
		keys = append(keys, key)
	}
	for key, _ := range postDeltas {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	prev := -1
	for _, key := range keys {
		if key == prev {
			continue
		}
		prev = key

		preDelta, ok := preDeltas[key]
		if ok {
			f.push(fmt.Sprintf("%d", key), key == len(array)-1)
			err := f.processItem(array[key], preDelta)
			if err != nil {
				return err
			}
			f.pop()

		}

		postDelta, ok := postDeltas[key]
		if ok {
			f.push(fmt.Sprintf("%d", key), key == len(array)-1)
			err := f.processItem(array[key], postDelta)
			if err != nil {
				return err
			}
			f.pop()
		}
	}

	return nil
}

func (f *AsciiFormatter) processObject(object map[string]interface{}, delta *gdiff.Object) error {
	keys := make([]string, 0, len(delta.Deltas))
	for key, _ := range delta.Deltas {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for i, key := range keys {
		f.push(fmt.Sprintf(`"%s"`, key), i == len(keys)-1)
		err := f.processItem(object[key], delta.Deltas[key])
		if err != nil {
			return err
		}
		f.pop()
	}

	return nil
}

func (f *AsciiFormatter) processItem(value interface{}, delta gdiff.Delta) error {
	switch d := delta.(type) {
	case *gdiff.Object:
		o, ok := value.(map[string]interface{})
		if !ok {
			return errors.New("Type mismatch")
		}

		f.newLine(AsciiSame)
		f.printKey()
		f.print("{")
		f.closeLine()
		err := f.processObject(o, d)
		if err != nil {
			return err
		}
		f.newLine(AsciiSame)
		f.print("}")
		f.closeLine()
		return nil

	case *gdiff.Array:
		a, ok := value.([]interface{})
		if !ok {
			return errors.New("Type mismatch")
		}

		f.newLine(AsciiSame)
		f.printKey()
		f.print("[")
		f.closeLine()
		err := f.processArray(a, d)
		if err != nil {
			return err
		}
		f.newLine(AsciiSame)
		f.print("]")
		f.closeLine()

		return nil
	case *gdiff.Added:
		f.printRecursive(d.Value, AsciiAdded)
		return nil

	case *gdiff.Modified:
		f.printRecursive(d.OldValue, AsciiDeleted)
		f.printRecursive(d.NewValue, AsciiAdded)
		return nil

	case *gdiff.TextDiff:
		f.printRecursive(d.OldValue, AsciiDeleted)
		f.printRecursive(d.NewValue, AsciiAdded)
		return nil

	case *gdiff.Deleted:
		f.printRecursive(d.Value, AsciiDeleted)
		return nil

	default:
		panic("Unknown Delta type detected")
	}
}

const (
	// Space for lines not changed
	AsciiSame = " "
	// Mak for added lines
	AsciiAdded = "+"
	// Mak for deleted lines
	AsciiDeleted = "-"
)

var AsciiStyles = map[string]string{
	// Green color for added lines
	AsciiAdded: "30;42",
	// Red color for deleted lines
	AsciiDeleted: "30;41",
}

func (f *AsciiFormatter) push(name string, last bool) {
	f.path = append(f.path, name)
	f.lasts = append(f.lasts, last)
}

func (f *AsciiFormatter) pop() {
	f.path = f.path[0 : len(f.path)-1]
	f.lasts = f.lasts[0 : len(f.lasts)-1]
}

func (f *AsciiFormatter) addLineWith(marker string, value string) {
	f.line = &AsciiLine{
		marker: marker,
		indent: len(f.path),
		buffer: bytes.NewBufferString(value),
	}
	f.closeLine()
}

func (f *AsciiFormatter) newLine(marker string) {
	f.line = &AsciiLine{
		marker: marker,
		indent: len(f.path),
		buffer: bytes.NewBuffer([]byte{}),
	}
}

func (f *AsciiFormatter) closeLine() {
	style, ok := AsciiStyles[f.line.marker]
	if f.config.Coloring && ok {
		f.buffer.WriteString("\x1b[" + style + "m")
	}

	f.buffer.WriteString(f.line.marker)
	for n := 0; n < f.line.indent; n++ {
		f.buffer.WriteString("  ")
	}
	f.buffer.Write(f.line.buffer.Bytes())

	if f.config.Coloring && ok {
		f.buffer.WriteString("\x1b[0m")
	}

	f.buffer.WriteRune('\n')
}

func (f *AsciiFormatter) printKey() {
	if len(f.path) > 0 {
		fmt.Fprintf(f.line.buffer, `%s: `, f.path[len(f.path)-1])
	}
}

func (f *AsciiFormatter) printComma() {
	if len(f.lasts) == 0 {
		return
	}
	if !f.lasts[len(f.lasts)-1] {
		f.line.buffer.WriteRune(',')
	}
}

func (f *AsciiFormatter) printValue(value interface{}) {
	switch value.(type) {
	case string:
		fmt.Fprintf(f.line.buffer, `"%s"`, value)
	case nil:
		f.line.buffer.WriteString("null")
	default:
		fmt.Fprintf(f.line.buffer, `%#v`, value)
	}
}

func (f *AsciiFormatter) print(a string) {
	f.line.buffer.WriteString(a)
}

func (f *AsciiFormatter) printRecursive(value interface{}, marker string) {
	switch v := value.(type) {
	case map[string]interface{}:
		f.newLine(marker)
		f.printKey()
		f.print("{")
		f.closeLine()

		keys := sortedKeys(v)
		for i, key := range keys {
			f.push(fmt.Sprintf(`"%s"`, key), i == len(keys)-1)
			f.printRecursive(v[key], marker)
			f.pop()
		}

		f.newLine(marker)
		f.print("}")
		f.printComma()
		f.closeLine()

	case []interface{}:
		f.newLine(marker)
		f.printKey()
		f.print("[")
		f.closeLine()

		for i, item := range v {
			f.push(fmt.Sprintf("%d", i), i == len(v)-1)
			f.printRecursive(item, marker)
			f.pop()
		}

		f.newLine(marker)
		f.print("]")
		f.printComma()
		f.closeLine()

	default:
		f.newLine(marker)
		f.printKey()
		f.printValue(value)
		f.printComma()
		f.closeLine()
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
