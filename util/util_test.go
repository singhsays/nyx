package util_test

import (
	. "nyx/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {
	Describe("ToAmount", func() {
		Context("given string input", func() {
			It("Parses amounts with decimal correctly.", func() {
				Expect(ToAmount("12.456695")).To(Equal(12.456695))
			})
			It("Parses amounts without decimal correctly.", func() {
				Expect(ToAmount("34")).To(Equal(34.0))
			})
			It("Parses amounts with - prefix correctly.", func() {
				Expect(ToAmount("-12.45")).To(Equal(-12.45))
			})
			It("Parses amounts with enclosing () correctly.", func() {
				Expect(ToAmount("(12.45)")).To(Equal(-12.45))
			})
			It("Parses amounts with commas correctly.", func() {
				Expect(ToAmount("(12,456,213.95)")).To(Equal(-12456213.95))
			})
			It("Parses amounts with $ prefix correctly.", func() {
				Expect(ToAmount("$12.45")).To(Equal(12.45))
				Expect(ToAmount("($12.45)")).To(Equal(-12.45))
				Expect(ToAmount("-$12.45")).To(Equal(-12.45))
				Expect(ToAmount("$1,235,123.45")).To(Equal(1235123.45))
			})
		})
	})
})
