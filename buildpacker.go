package buildpacker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	dkr "github.com/docker/docker/api/client"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	BuildpackerRoot = "/var/buildpacker"
	BuildDir        = "code"
	BuildpackDir    = "buildpack"
	BuildpackZip    = "bp.zip"
	DefaultBox      = "ubuntu"
	certFileFormat  = "%s/cert.pem"
	keyFileFormat   = "%s/key.pem"
	caFileFormat    = "%s/ca.pem"
	DefaultProto    = "tcp"
)

func New(buildpack string, localbuildpath string) *BPacker {
	return &BPacker{
		buildpack:      buildpack,
		codepath:       codepath,
		localbuildpath: localbuildpath,
	}
}

type BPacker struct {
	buildpack      string
	codepath       string
	localbuildpath string
}

func (s *BPacker) Build(endpoint string, certpath string) {
	cert := fmt.Sprintf(certFileFormat, certpath)
	key := fmt.Sprintf(keyFileFormat, certpath)
	ca := fmt.Sprintf(caFileFormat, certpath)
	client, err := docker.NewTLSClient(endpoint, cert, key, ca)
	outputbuf, errbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	ioutil.WriteFile("./Dockerfile", []byte(s.createDockerFile()), os.ModePerm)

	reader := strings.NewReader("")
	inputbuf := ioutil.NopCloser(reader)
	endpoint = strings.TrimPrefix(endpoint, fmt.Sprintf("%s://", DefaultProto))
	dcli := dkr.NewDockerCli(inputbuf, outputbuf, errbuf, key, DefaultProto, endpoint, client.TLSConfig)
	err = dcli.CmdBuild("./")
	fmt.Println(err)
	fmt.Println(outputbuf)
	os.Remove("Dockerfile")
}

func (s *BPacker) createDockerFile() (dockerFileString string) {
	var buffer bytes.Buffer
	buildpacks := fmt.Sprintf("%s/%s", BuildpackerRoot, BuildpackDir)
	builddir := fmt.Sprintf("%s/%s", BuildpackerRoot, BuildDir)
	buildpackZipPath := fmt.Sprintf("%s/%s/%s", BuildpackerRoot, BuildpackDir, BuildpackZip)
	buffer.WriteString(fmt.Sprintf("FROM %s\n", DefaultBox))
	buffer.WriteString("RUN apt-get install unzip \n")
	buffer.WriteString(fmt.Sprintf("RUN mkdir -p %s \n", BuildpackerRoot))
	buffer.WriteString(fmt.Sprintf("ADD %s %s \n", s.localbuildpath, builddir))
	buffer.WriteString(fmt.Sprintf("ADD %s %s \n", s.buildpack, buildpackZipPath))
	buffer.WriteString(fmt.Sprintf("RUN unzip %s -d %s \n", buildpackZipPath, buildpacks))
	buffer.WriteString(fmt.Sprintf("RUN export \"out=$(ls %s | grep -v .zip)\" && cd %s && mv $out/* ./ && rm -fR $out %s\n", buildpacks, buildpacks, buildpackZipPath))
	buffer.WriteString(fmt.Sprintf("RUN %s/bin/compile %s/code /tmp", buildpacks, BuildpackerRoot))
	dockerFileString = buffer.String()
	fmt.Println(dockerFileString)
	return
}
