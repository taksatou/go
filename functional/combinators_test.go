package functional_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/taksatou/go/functional"
)

var (
	_ = fmt.Println
)
var _ = Describe("Combinators", func() {
	It("returns itself", func() {
		Ω(I(1)).Should(Equal(1))
		Ω(I(nil)).Should(BeNil())
		Ω(I("")).Should(Equal(""))
	})
})
