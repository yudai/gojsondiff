package printer_test

import (
	. "github.com/yudai/gojsondiff/printer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yudai/gojsondiff/test"

	"regexp"

	diff "github.com/yudai/gojsondiff"
)

var _ = Describe("Ascii", func() {
	Describe("AsciiPrinter", func() {
		var (
			a, b           map[string]interface{}
			colorFilter, _ = regexp.Compile("\\x1b\\[[0-9;]*m")
		)

		It("Prints the given diff", func() {
			a = LoadFixture("../FIXTURES/base.json")
			b = LoadFixture("../FIXTURES/base_changed.json")

			diff := diff.CompareObjects(a, b)
			Expect(diff.Modified()).To(BeTrue())

			printer := NewAsciiPrinter()
			diff.Iterate(printer)
			Expect(colorFilter.ReplaceAllString(printer.Result(), "")).To(
				Equal(
					` {
   "arr": [
     0: "arr0",
     1: 21,
     2: {
       "num": 1,
-      "str": "pek3f"
+      "str": "changed"
     },
     3: [
       0: 0,
-      1: "1"
+      1: "changed"
     ]
   ],
   "bool": true,
   "num_float": 39.39,
   "num_int": 13,
   "obj": {
     "arr": [
       0: 17,
       1: "str",
       2: {
-        "str": "eafeb"
+        "str": "changed"
       }
     ],
+    "new": "added",
-    "num": 19,
     "obj": {
-      "num": 14,
+      "num": 9999,
-      "str": "efj3"
+      "str": "changed"
     },
     "str": "bcded"
   },
   "str": "abcde"
 }
`,
				),
			)
		})

		It("Prints the given diff", func() {
			a = LoadFixture("../FIXTURES/add_delete_from.json")
			b = LoadFixture("../FIXTURES/add_delete_to.json")

			diff := diff.CompareObjects(a, b)
			Expect(diff.Modified()).To(BeTrue())

			printer := NewAsciiPrinter()
			diff.Iterate(printer)
			Expect(printer.ResultWithoutColor()).To(
				Equal(
					` {
+  "add": {
+    "l0a": [
+      0: "abcd",
+      1: [
+        0: "efcg"
+      ]
+    ],
+    "l0o": {
+      "l1o": {
+        "l2s": "efed"
+      },
+      "l1s": "abcd"
+    }
+  },
-  "delete": {
-    "l0a": [
-      0: "abcd",
-      1: [
-        0: "efcg"
-      ]
-    ],
-    "l0o": {
-      "l1o": {
-        "l2s": "efed"
-      },
-      "l1s": "abcd"
-    }
-  }
 }
`,
				),
			)
		})
	})

})
