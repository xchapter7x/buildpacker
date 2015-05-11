package buildpacker_test

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/buildpacker"
)

var _ = Describe("Buildpacker", func() {
	Describe("DockerFileBucket", func() {
		Context("calling Dockerfile() w/ valid initialized struct", func() {
			var dfBucket DockerFileBucket

			BeforeEach(func() {
				dfBucket = DockerFileBucket{
					DefaultBox:       DefaultBox,
					BuildpackerRoot:  BuildpackerRoot,
					LocalBuildPath:   "./",
					BuildDir:         fmt.Sprintf("%s/%s", BuildpackerRoot, BuildDir),
					Buildpack:        "https://github.com/ryandotsmith/null-buildpack/archive/master.zip",
					BuildpackZipPath: fmt.Sprintf("%s/%s/%s", BuildpackerRoot, BuildpackDir, BuildpackZip),
					BuildpackDir:     fmt.Sprintf("%s/%s", BuildpackerRoot, BuildpackDir),
				}
			})

			It("should return the dockerfile we want in string format", func() {
				dockerfileString := dfBucket.Dockerfile()
				controlBytes, _ := ioutil.ReadFile("fixtures/Dockerfile-valid")
				Ω(len(dockerfileString)).Should(BeEquivalentTo(len(string(controlBytes[:]))))
			})
		})
	})

	XContext("test poc build", func() {
		It("should do something?", func() {
			bdPkr := New("https://github.com/ryandotsmith/null-buildpack/archive/master.zip", "./")
			//bdPkr := New("https://github.com/cloudfoundry/go-buildpack/releases/download/v1.3.1/go_buildpack-cached-v1.3.1.zip", "./")
			//bdPkr := New("go_buildpack-cached-v1.3.1.zip", "./")
			bdPkr.Build(os.Getenv("DOCKER_HOST"), os.Getenv("DOCKER_CERT_PATH"), "testimage")
			Ω(true).Should(Equal(true))
		})
	})
})
