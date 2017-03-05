package formatter

import (
	"testing"

	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/tests"
)

func TestAsciiFormatter(t *testing.T) {
	differ := diff.New()

	cases := []struct {
		left   string
		right  string
		output string
	}{
		{
			"../FIXTURES/base.json", "../FIXTURES/base_changed.json",
			` {
   "arr": [
     "arr0",
     21,
     {
       "num": 1,
-      "str": "pek3f"
+      "str": "changed"
     },
     [
       0,
-      "1"
+      "changed"
     ]
   ],
   "bool": true,
-  "null": null,
   "num_float": 39.39,
   "num_int": 13,
   "obj": {
     "arr": [
       17,
       "str",
       {
-        "str": "eafeb"
+        "str": "changed"
       }
     ],
-    "num": 19,
     "obj": {
-      "num": 14,
+      "num": 9999,
-      "str": "efj3"
+      "str": "changed"
     },
     "str": "bcded"
+    "new": "added"
   },
   "str": "abcde"
 }
`,
		},
		{
			"../FIXTURES/add_delete_from.json", "../FIXTURES/add_delete_to.json",
			` {
-  "delete": {
-    "l0a": [
-      "abcd",
-      [
-        "efcg"
-      ]
-    ],
-    "l0o": {
-      "l1o": {
-        "l2s": "efed"
-      },
-      "l1s": "abcd"
-    }
-  }
+  "add": {
+    "l0a": [
+      "abcd",
+      [
+        "efcg"
+      ]
+    ],
+    "l0o": {
+      "l1o": {
+        "l2s": "efed"
+      },
+      "l1s": "abcd"
+    }
+  }
 }
`,
		},
	}

	for i, c := range cases {
		left := tests.LoadFixture(t, c.left)
		right := tests.LoadFixture(t, c.right)
		diff := differ.Compare(left, right)

		f := NewAsciiFormatter(left, AsciiFormatterDefaultConfig)
		output, err := f.Format(diff)
		if err != nil {
			t.Errorf("unexpected format error at test case `%d`: left `%s`, right `%s` ", i, c.left, c.right)
		}
		if output != c.output {
			t.Errorf("unexpected output at test case `%d` (left `%s`, right `%s`): output: %s, expected: %s ", i, c.left, c.right, output, c.output)
		}
	}
}
