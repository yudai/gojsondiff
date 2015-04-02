package gojsondiff_test

import (
	. "github.com/yudai/gojsondiff"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yudai/gojsondiff/test"
)

var _ = Describe("Gojsondiff", func() {
	Describe("Diff", func() {
		var (
			a, b map[string]interface{}
		)

		Context("There are no difference between the two JSON strings", func() {
			It("Detects nothing", func() {
				a = LoadFixture("FIXTURES/base.json")
				b = LoadFixture("FIXTURES/base.json")

				diff := CompareObjects(a, b)

				Expect(diff.Modified()).To(BeFalse())

				s := diff.Structure()
				Expect(s["str"]).To(Equal(Same{OldValue: "abcde", NewValue: "abcde"}))
				// actually, they are not int values, they are float64.
				Expect(s["num_int"]).To(Equal(Same{OldValue: 13.0, NewValue: 13.0}))
				Expect(s["num_float"]).To(Equal(Same{OldValue: 39.39, NewValue: 39.39}))
				Expect(s["bool"]).To(Equal(Same{OldValue: true, NewValue: true}))

				arr := s["arr"].([]interface{})
				Expect(arr[0]).To(Equal(Same{OldValue: "arr0", NewValue: "arr0"}))
				Expect(arr[1]).To(Equal(Same{OldValue: 21.0, NewValue: 21.0}))
				arrobj := arr[2].(map[string]interface{})
				Expect(arrobj["str"]).To(Equal(Same{OldValue: "pek3f", NewValue: "pek3f"}))
				Expect(arrobj["num"]).To(Equal(Same{OldValue: 1.0, NewValue: 1.0}))
				arrarr := arr[3].([]interface{})
				Expect(arrarr[0]).To(Equal(Same{OldValue: 0.0, NewValue: 0.0}))
				Expect(arrarr[1]).To(Equal(Same{OldValue: "1", NewValue: "1"}))

				obj := s["obj"].(map[string]interface{})
				Expect(obj["str"]).To(Equal(Same{OldValue: "bcded", NewValue: "bcded"}))
				Expect(obj["num"]).To(Equal(Same{OldValue: 19.0, NewValue: 19.0}))
				objarr := obj["arr"].([]interface{})
				Expect(objarr[0]).To(Equal(Same{OldValue: 17.0, NewValue: 17.0}))
				Expect(objarr[1]).To(Equal(Same{OldValue: "str", NewValue: "str"}))
				objobj := obj["obj"].(map[string]interface{})
				Expect(objobj["str"]).To(Equal(Same{OldValue: "efj3", NewValue: "efj3"}))
				Expect(objobj["num"]).To(Equal(Same{OldValue: 14.0, NewValue: 14.0}))
			})
		})

		Context("There are some values modified", func() {
			It("Detects changes", func() {
				a = LoadFixture("FIXTURES/base.json")
				b = LoadFixture("FIXTURES/base_changed.json")

				diff := CompareObjects(a, b)

				Expect(diff.Modified()).To(BeTrue())

				s := diff.Structure()

				Expect(s["str"]).To(Equal(Same{OldValue: "abcde", NewValue: "abcde"}))
				// actually, they are not int values, they are float64.
				Expect(s["num_int"]).To(Equal(Same{OldValue: 13.0, NewValue: 13.0}))
				Expect(s["num_float"]).To(Equal(Same{OldValue: 39.39, NewValue: 39.39}))
				Expect(s["bool"]).To(Equal(Same{OldValue: true, NewValue: true}))

				arr := s["arr"].([]interface{})
				Expect(arr[0]).To(Equal(Same{OldValue: "arr0", NewValue: "arr0"}))
				Expect(arr[1]).To(Equal(Same{OldValue: 21.0, NewValue: 21.0}))
				arrobj := arr[2].(map[string]interface{})
				Expect(arrobj["str"]).To(Equal(Modified{OldValue: "pek3f", NewValue: "changed"}))
				Expect(arrobj["num"]).To(Equal(Same{OldValue: 1.0, NewValue: 1.0}))
				arrarr := arr[3].([]interface{})
				Expect(arrarr[0]).To(Equal(Same{OldValue: 0.0, NewValue: 0.0}))
				Expect(arrarr[1]).To(Equal(Modified{OldValue: "1", NewValue: "changed"}))

				obj := s["obj"].(map[string]interface{})
				Expect(obj["str"]).To(Equal(Same{OldValue: "bcded", NewValue: "bcded"}))
				Expect(obj["num"]).To(Equal(Deleted{OldValue: 19.0}))
				Expect(obj["new"]).To(Equal(Added{NewValue: "added"}))
				objarr := obj["arr"].([]interface{})
				Expect(objarr[0]).To(Equal(Same{OldValue: 17.0, NewValue: 17.0}))
				Expect(objarr[1]).To(Equal(Same{OldValue: "str", NewValue: "str"}))
				objobj := obj["obj"].(map[string]interface{})
				Expect(objobj["str"]).To(Equal(Modified{OldValue: "efj3", NewValue: "changed"}))
				Expect(objobj["num"]).To(Equal(Modified{OldValue: 14.0, NewValue: 9999.0}))
			})
		})

		Context("There are values only types are changed", func() {
			It("Detects changed types", func() {
				a := LoadFixture("FIXTURES/changed_types_from.json")
				b := LoadFixture("FIXTURES/changed_types_to.json")

				diff := CompareObjects(a, b)

				Expect(diff.Modified()).To(BeTrue())

				s := diff.Structure()
				Expect(s["str"]).To(Equal(Modified{OldValue: "true", NewValue: true}))
				Expect(s["num_int"]).To(Equal(Modified{OldValue: 3.0, NewValue: "3"}))
				Expect(s["num_float"]).To(Equal(Modified{OldValue: 99.99, NewValue: "99.99"}))
				Expect(s["bool"]).To(Equal(Modified{OldValue: true, NewValue: "true"}))
				Expect(s["obj"]).To(Equal(
					Modified{
						OldValue: map[string]interface{}{"0": "i0", "1": "i1"},
						NewValue: []interface{}{"i0", "i1"},
					}),
				)

				// By default, all JSON Numbers are float64
				Expect(s["num_int_float"]).To(Equal(Same{OldValue: 5.0, NewValue: 5.0}))
			})
		})

	})
})
