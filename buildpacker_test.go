package buildpacker_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/buildpacker"
)

var _ = Describe("Buildpacker", func() {
	Context("test poc build", func() {
		It("should do something?", func() {
			Build(os.Getenv("DOCKER_HOST"), os.Getenv("DOCKER_CERT_PATH"))
			Ω(true).Should(Equal(false))
		})
	})
})
