package formatter_test

import (
	. "github.com/yanyan300300/gojsondiff/formatter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yanyan300300/gojsondiff/test"

	diff "github.com/yanyan300300/gojsondiff"
)

var _ = Describe("Ascii", func() {
	Describe("AsciiPrinter", func() {
		var (
			a, b map[string]interface{}
		)

		It("Prints the given diff", func() {
			a = LoadFixture("../FIXTURES/base.json")
			b = LoadFixture("../FIXTURES/base_changed.json")

			diff := diff.New().CompareObjects(a, b)
			Expect(diff.Modified()).To(BeTrue())

			f := NewAsciiFormatter(a)
			deltaJson, err := f.Format(diff)
			Expect(err).To(BeNil())
			Expect(deltaJson).To(Equal(
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
			),
			)
		})

		It("Prints the given diff", func() {
			a = LoadFixture("../FIXTURES/add_delete_from.json")
			b = LoadFixture("../FIXTURES/add_delete_to.json")

			diff := diff.New().CompareObjects(a, b)
			Expect(diff.Modified()).To(BeTrue())

			f := NewAsciiFormatter(a)
			deltaJson, err := f.Format(diff)
			Expect(err).To(BeNil())
			Expect(deltaJson).To(Equal(
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
			),
			)
		})
	})

})
