package buildpacker_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBuildpacker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "buildpacker Client Suite")
}
