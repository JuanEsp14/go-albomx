package example_test

import (
	. "github.com/onsi/ginkgo"
)

const baseHelloWorldUrl = "/hello"

var _ = Describe("AlbomxComicsController", func() {
	Context("Happy path", func() {
		helloSpecs()
	})
})

func helloSpecs() {
	When("Valid Hello Request without name param", func() {
		It("Should return Hello World", func() {

		})
	})
}