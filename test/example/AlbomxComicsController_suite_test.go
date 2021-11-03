package example_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)


var _ = BeforeSuite(func() {

})

func TestAlbomxComicsController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AlbomxComicsController Suite")
}

var _ = AfterSuite(func() {

})