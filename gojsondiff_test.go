package gojsondiff

import (
	"reflect"
	"testing"

	"github.com/yudai/gojsondiff/tests"
)

func TestWithSameJSONs(t *testing.T) {
	differ := New()
	left := tests.LoadFixture(t, "FIXTURES/base.json")
	right := tests.LoadFixture(t, "FIXTURES/base.json")

	diff := differ.Compare(left, right)
	if diff.Modified() {
		t.Fail()
	}
}

func TestWithChanges(t *testing.T) {
	differ := New()

	cases := []struct {
		left  string
		right string
	}{
		{"FIXTURES/base.json", "FIXTURES/base_changed.json"},
		{"FIXTURES/changed_types_from.json", "FIXTURES/changed_types_to.json"},
		{"FIXTURES/move_from.json", "FIXTURES/move_to.json"},
		{"FIXTURES/long_text_from.json", "FIXTURES/long_text_to.json"},
		{"FIXTURES/array.json", "FIXTURES/array_changed.json"},
		{"FIXTURES/string.json", "FIXTURES/string_changed.json"},
	}

	for i, c := range cases {
		left := tests.LoadFixture(t, c.left)
		right := tests.LoadFixture(t, c.right)

		diff := differ.Compare(left, right)
		if !diff.Modified() {
			t.Errorf("Unexpected unmodified result at test case `%d`: left `%s`, right `%s` ", i, c.left, c.right)
		}
		patched, err := differ.ApplyPatch(left, diff)
		if err != nil {
			t.Errorf("unexpected patch error at test case `%d`: left `%s`, right `%s` ", i, c.left, c.right)
		}
		if !reflect.DeepEqual(tests.LoadFixture(t, c.right), patched) {
			t.Errorf("patched JSON doesn't match to expected JSON at test case `%d`: left `%s`, right `%s` ", i, c.left, c.right)
		}
	}
}
