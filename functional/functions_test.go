package functional_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/taksatou/go/functional"
	"unicode"
)

var _ = Describe("Functional", func() {
	var (
		_ = fmt.Println
	)

	Context("basic functions", func() {
		It("curries add function", func() {
			var curriedAdd func(int) func(int) int
			add := func(a, b int) int {
				return a + b
			}
			e := Curry(add, &curriedAdd)
			Ω(e).Should(BeNil())
			Ω(curriedAdd(1)(2)).Should(Equal(3))
		})

		It("flips arguments", func() {
			var flipped func(float64, float64) float64
			div := func(a, b float64) float64 {
				return a / b
			}
			e := Flip(div, &flipped)
			Ω(e).Should(BeNil())
			Ω(flipped(2, 3)).Should(Equal(div(3, 2)))
		})

		It("composes two functions", func() {
			var composed func(int) int
			var e error
			inc := func(a int) int { return a + 1 }
			double := func(a int) int { return a * 2 }
			e = Compose(inc, double, &composed)
			Ω(e).Should(BeNil())
			Ω(composed(1)).Should(Equal(3))

			e = Compose(double, inc, &composed)
			Ω(e).Should(BeNil())
			Ω(composed(1)).Should(Equal(4))
		})

		It("fails if two functions have wrong types", func() {
			var composed func(int) int
			inc := func(a int) int { return a + 1 }
			double := func(a float64) float64 { return a * 2.0 }
			e := Compose(inc, double, &composed)
			Ω(e).Should(HaveOccured())
		})
	})

	Context("Fold", func() {
		It("folds int slice into a int value", func() {
			lis := []int{1, 2, 3}
			res, e := Fold(func(a, b int) int {
				return a + b
			}, lis, nil)
			Ω(e).Should(BeNil())
			Ω(res).Should(Equal(6))
		})

		// // not implemented yet
		// It("folds int slice into a float64 value", func() {
		// 	lis := []int{1, 2, 3}
		// 	res, e := Fold(func(a float64, b int) float64 {
		// 		return a + float64(b)
		// 	}, lis, 0.1)
		// 	Ω(e).Should(BeNil())
		// 	Ω(res).Should(Equal(6.1))
		// })

		It("folds string slice into a value", func() {
			lis := []string{"a", "b", "c"}
			res, e := Fold(func(a, b string) string {
				return a + b
			}, lis, "")
			Ω(e).Should(BeNil())
			Ω(res).Should(Equal("abc"))
		})

		It("should fail for wrong type", func() {
			lis := []string{"a", "b", "c"}
			_, e := Fold(func(a, b int) int {
				return a + b
			}, lis, nil)
			Ω(e).Should(HaveOccured())
		})
	})

	Context("Map", func() {
		It("doubles each elements", func() {
			lis := []int{1, 2, 3}
			res, e := Map(lis, func(a int) int {
				return a * 2
			})
			Ω(e).Should(BeNil())
			Ω(res).Should(Equal([]int{2, 4, 6}))
		})

		It("convert type", func() {
			lis := []int{1, 2, 3}
			res, e := Map(lis, func(a int) int64 {
				return int64(a)
			})
			Ω(e).Should(BeNil())
			Ω(res).Should(Equal([]int64{1, 2, 3}))
		})

		It("should fail if it passed function has wrong type", func() {
			lis := []int{1, 2, 3}
			_, e := Map(lis, func(a int) {})
			Ω(e).Should(HaveOccured())
		})

		It("transforms each characters into upcase", func() {
			res, e := Map("abc", func(a rune) rune {
				return unicode.ToUpper(a)
			})
			Ω(e).Should(BeNil())
			Ω(res).Should(Equal("ABC"))
		})

		It("transforms map elements", func() {
			lis := map[string]int{
				"abc": 1,
				"def": 2,
				"ghi": 3,
			}
			res, e := Map(lis, func(k string, v int) (string, int) {
				return string(bytes.ToUpper([]byte(k))), v * 2
			})
			Ω(e).Should(BeNil())
			Ω(res).Should(Equal(map[string]int{
				"ABC": 2,
				"DEF": 4,
				"GHI": 6,
			}))
		})
	})

	Context("Each", func() {
		It("calls specified function with each elements", func() {
			lis := [...]int{1, 2, 3}
			res := []int{}
			Each(lis, func(a int) {
				res = append(res, a)
			})
			Ω(res).Should(Equal(lis[:]))
		})

		It("calls specified function with each elements", func() {
			type T struct {
				a int
				b string
			}
			lis := []T{
				{1, "a"},
				{2, "b"},
				{3, "c"},
			}
			res := []T{}
			Each(lis, func(a T) {
				res = append(res, a)
			})
			Ω(res).Should(Equal(lis))
		})

		It("calls specified function with each character", func() {
			s := "asdf"
			res := []byte{}
			Each(s, func(a byte) {
				res = append(res, a)
			})
			Ω(string(res)).Should(Equal(s))

		})

		It("calls specified function with each key and value", func() {
			lis := map[string]int{
				"a": 1,
				"b": 10,
				"c": 100,
			}
			res := map[string]int{}
			Each(lis, func(k string, v int) {
				res[k] = v
			})
			Ω(res).Should(Equal(lis))
		})
	})

	Context("Filter", func() {
		It("filters elements by specified function", func() {
			lis := []int{1, 2, 3}
			res, e := Filter(lis, func(a int) bool {
				return a%2 == 0
			})
			Ω(e).Should(BeNil())
			Ω(res).Should(Equal([]int{2}))
		})

		It("filters elements in interface slice by type", func() {
			lis := []interface{}{1, 2.3, "abc"}
			res, e := Filter(lis, func(a interface{}) bool {
				_, ok := a.(string)
				return ok
			})
			Ω(e).Should(BeNil())
			Ω(res).Should(Equal([]interface{}{"abc"}))

		})

		It("filters string by character", func() {
			s := "aBcdEFg"
			res, e := Filter(s, func(c rune) bool {
				return unicode.IsLower(c)
			})
			Ω(e).Should(BeNil())
			Ω(res).Should(Equal("acdg"))
		})

		It("filters map by key and value", func() {
			lis := map[string]int{
				"a": 1,
				"b": 10,
				"c": 100,
			}
			res, e := Filter(lis, func(k string, v int) bool {
				return k == "b" && v == 10
			})
			Ω(e).Should(BeNil())
			Ω(res).Should(Equal(map[string]int{"b": 10}))
		})
	})
})
