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
		localbuildpath: localbuildpath,
	}
}

type BPacker struct {
	buildpack      string
	localbuildpath string
}

func (s *BPacker) Build(endpoint string, certpath string) {
	cert := fmt.Sprintf(certFileFormat, certpath)
	key := fmt.Sprintf(keyFileFormat, certpath)
	ca := fmt.Sprintf(caFileFormat, certpath)
	client, err := docker.NewTLSClient(endpoint, cert, key, ca)
	outputbuf, errbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	dockerfileString := s.CreateDockerFile()
	fmt.Println(dockerfileString)
	ioutil.WriteFile("./Dockerfile", []byte(dockerfileString), os.ModePerm)
	defer os.Remove("Dockerfile")
	reader := strings.NewReader("")
	inputbuf := ioutil.NopCloser(reader)
	endpoint = strings.TrimPrefix(endpoint, fmt.Sprintf("%s://", DefaultProto))
	dcli := dkr.NewDockerCli(inputbuf, outputbuf, errbuf, key, DefaultProto, endpoint, client.TLSConfig)
	err = dcli.CmdBuild("./")
	fmt.Println(err)
	fmt.Println(outputbuf)
}

func (s *BPacker) CreateDockerFile() (dockerFileString string) {
	var buffer bytes.Buffer
	buildpacks := fmt.Sprintf("%s/%s", BuildpackerRoot, BuildpackDir)
	builddir := fmt.Sprintf("%s/%s", BuildpackerRoot, BuildDir)
	buildpackZipPath := fmt.Sprintf("%s/%s/%s", BuildpackerRoot, BuildpackDir, BuildpackZip)
	buffer.WriteString(fmt.Sprintf("FROM %s\n", DefaultBox))
	buffer.WriteString("RUN rm /bin/sh && ln -s /bin/bash /bin/sh \n")
	buffer.WriteString("RUN apt-get install -y unzip curl ruby gcc \n")
	buffer.WriteString(fmt.Sprintf("RUN mkdir -p %s \n", BuildpackerRoot))
	buffer.WriteString(fmt.Sprintf("ADD %s %s \n", s.localbuildpath, builddir))
	buffer.WriteString(fmt.Sprintf("ADD %s %s \n", s.buildpack, buildpackZipPath))
	buffer.WriteString(fmt.Sprintf("RUN unzip %s -d %s/unpacked \n", buildpackZipPath, buildpacks))
	buffer.WriteString(fmt.Sprintf("RUN cd %s && if [ $(ls ./unpacked | wc -l) == 1 ]; then mv ./unpacked/$(ls ./unpacked) ./tmp && rm -fR ./unpacked && mv ./tmp ./unpacked; fi && rm -fR %s\n", buildpacks, buildpackZipPath))
	buffer.WriteString(fmt.Sprintf("RUN %s/unpacked/bin/detect %s/code\n", buildpacks, BuildpackerRoot))
	buffer.WriteString(fmt.Sprintf("RUN %s/unpacked/bin/compile %s/code /tmp\n", buildpacks, BuildpackerRoot))
	dockerFileString = buffer.String()
	return
}
