package formatter

import (
	"reflect"
	"testing"

	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/tests"
)

func TestDeltaFormatter(t *testing.T) {
	differ := diff.New()

	cases := []struct {
		left   string
		right  string
		output interface{}
	}{
		{
			"../FIXTURES/base.json", "../FIXTURES/base_changed.json",
			map[string]interface{}{
				"arr": map[string]interface{}{
					"_t": "a",
					"2": map[string]interface{}{
						"str": []interface{}{
							"pek3f", "changed",
						},
					},
					"3": map[string]interface{}{
						"_t": "a",
						"1": []interface{}{
							"1", "changed",
						},
					},
				},
				"null": []interface{}{nil, 0, 0},
				"obj": map[string]interface{}{
					"arr": map[string]interface{}{
						"_t": "a",
						"2": map[string]interface{}{
							"str": []interface{}{
								"eafeb", "changed",
							},
						},
					},
					"num": []interface{}{
						float64(19),
						0,
						0,
					},
					"obj": map[string]interface{}{
						"num": []interface{}{
							float64(14),
							float64(9999),
						},
						"str": []interface{}{
							"efj3",
							"changed",
						},
					},
					"new": []interface{}{
						"added",
					},
				},
			},
		},
		{
			"../FIXTURES/long_text_from.json", "../FIXTURES/long_text_to.json",
			map[string]interface{}{
				"str": []interface{}{
					"@@ -27,14 +27,15 @@\n 40fj\n-q048hf\n+nafefea\n bgvz\n@@ -55,25 +55,20 @@\n m48q\n+a3\n 9p8qfh\n-qn4gbqqq4qp\n+nafe\n 94hq\n",
					0,
					2,
				},
			},
		},
	}

	for i, c := range cases {
		left := tests.LoadFixture(t, c.left)
		right := tests.LoadFixture(t, c.right)
		diff := differ.Compare(left, right)

		f := NewDeltaFormatter()
		deltaJson, err := f.FormatAsJson(diff)
		if err != nil {
			t.Errorf("unexpected format error at test case `%d`: left `%s`, right `%s` ", i, c.left, c.right)
		}
		if !reflect.DeepEqual(deltaJson, c.output) {
			t.Errorf("unexpected output at test case `%d` (left `%s`, right `%s`): actual %#v, expected %#v", i, c.left, c.right, deltaJson, c.output)
		}
	}
}
