package printer

import (
	"fmt"
	"github.com/k0kubun/pp"
	diff "github.com/yudai/gojsondiff"
	"reflect"
	"regexp"
	"sort"
	"strconv"
)

func NewAsciiPrinter() *AsciiPrinter {
	return &AsciiPrinter{
		path:    []string{},
		size:    []int{},
		inArray: []bool{},
	}
}

type AsciiPrinter struct {
	buffer  string
	path    []string
	size    []int
	inArray []bool
}

func (i *AsciiPrinter) SortObjectFields() bool {
	return true
}

func (i *AsciiPrinter) EnterRoot(size int) bool {
	i.printIndent(same)
	i.println("{")
	i.push("ROOT", size, false)
	return false
}

func (i *AsciiPrinter) ExitRoot() bool {
	i.pop()
	i.printIndent(same)
	i.println("}")
	return false
}

func (i *AsciiPrinter) EnterObject(name string, size int) bool {
	i.printKeyWithIndent(name, same)
	i.println("{")
	i.push("ROOT", size, false)
	return false
}
func (i *AsciiPrinter) ExitObject(name string) bool {
	i.pop()
	i.printIndent(same)
	i.print("}")
	i.printComma()
	return false
}

func (i *AsciiPrinter) EnterArray(name string, size int) bool {
	i.printKeyWithIndent(name, same)
	i.println("[")
	i.push(name, size, true)
	return false
}

func (i *AsciiPrinter) ExitArray(name string) bool {
	i.pop()
	i.printIndent(same)
	i.print("]")
	i.printComma()
	return false
}

func (i *AsciiPrinter) VisitSame(name string, d diff.Same) bool {
	i.printRecursive(name, d.OldValue, same)
	return false
}

func (i *AsciiPrinter) VisitAdded(name string, d diff.Added) bool {
	i.printRecursive(name, d.NewValue, added)
	return false
}

func (i *AsciiPrinter) VisitModified(name string, d diff.Modified) bool {
	savedSize := i.size[len(i.size)-1]
	i.printRecursive(name, d.OldValue, deleted)
	i.size[len(i.size)-1] = savedSize
	i.printRecursive(name, d.NewValue, added)
	return false
}

func (i *AsciiPrinter) VisitDeleted(name string, d diff.Deleted) bool {
	i.printRecursive(name, d.OldValue, deleted)
	return false
}

func (i *AsciiPrinter) Result() string {
	return i.buffer
}

func (i *AsciiPrinter) ResultWithoutColor() string {
	colorFilter, _ := regexp.Compile("\\x1b\\[[0-9;]*m")
	return colorFilter.ReplaceAllString(i.buffer, "")
}

func (i *AsciiPrinter) push(name string, size int, array bool) {
	i.path = append(i.path, name)
	i.size = append(i.size, size)
	i.inArray = append(i.inArray, array)
}

func (i *AsciiPrinter) pop() {
	i.path = i.path[0 : len(i.path)-1]
	i.size = i.size[0 : len(i.size)-1]
	i.inArray = i.inArray[0 : len(i.inArray)-1]
}

func (i *AsciiPrinter) printRecursive(name string, value interface{}, marker string) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Map {
		i.printKeyWithIndent(name, marker)
		i.println("{")

		m := value.(map[string]interface{})
		size := len(m)
		i.push(name, size, false)

		keys := []string{}
		for key, _ := range m {
			keys = append(keys, key)
		}
		sort.StringSlice(keys).Sort()

		for _, key := range keys {
			i.printRecursive(key, m[key], marker)
		}
		i.pop()

		i.printIndent(marker)
		i.print("}")
		i.printComma()
	} else if v.Kind() == reflect.Slice {
		i.printKeyWithIndent(name, marker)
		i.println("[")

		s := value.([]interface{})
		size := len(s)
		i.push(name, size, true)
		for key, item := range s {
			keyStr := strconv.Itoa(key)
			i.printRecursive(keyStr, item, marker)
		}
		i.pop()

		i.printIndent(marker)
		i.print("]")
		i.printComma()
	} else {
		i.printKeyWithIndent(name, marker)
		i.printValue(value)
		i.printComma()
	}
}

const (
	same    = " "
	added   = "+"
	deleted = "-"
)

func (i *AsciiPrinter) printIndent(marker string) {
	i.print(marker)
	for n := 0; n < len(i.path); n++ {
		i.print("  ")
	}
}

func (i *AsciiPrinter) printKeyWithIndent(name string, marker string) {
	i.printIndent(marker)
	if i.inArray[len(i.inArray)-1] {
		i.printf(`%s: `, name)
	} else {
		i.printf(`"%s": `, name)
	}
}

func (i *AsciiPrinter) printComma() {
	i.size[len(i.size)-1]--
	if i.size[len(i.size)-1] > 0 {
		i.println(",")
	} else {
		i.println()
	}
}

func (i *AsciiPrinter) printValue(v interface{}) {
	i.buffer += pp.Sprint(v)
}

func (i *AsciiPrinter) print(a ...interface{}) {
	i.buffer += fmt.Sprint(a...)
}

func (i *AsciiPrinter) printf(format string, a ...interface{}) {
	i.buffer += fmt.Sprintf(format, a...)
}

func (i *AsciiPrinter) println(a ...interface{}) {
	i.buffer += fmt.Sprintln(a...)
}
