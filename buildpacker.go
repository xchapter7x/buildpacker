package buildpacker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	dkr "github.com/docker/docker/api/client"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	BuildpackerRoot = "/var/buildpacker"
	BuildDir        = "code"
	BuildpackDir    = "buildpack"
	BuildpackZip    = "bp.zip"
	DefaultBox      = "cloudfoundry/cflinuxfs2"
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

type (
	DockerFileBucket struct {
		DefaultBox       string
		BuildpackerRoot  string
		LocalBuildPath   string
		BuildDir         string
		Buildpack        string
		BuildpackZipPath string
		BuildpackDir     string
	}
	BPacker struct {
		buildpack      string
		localbuildpath string
	}
)

func (s *BPacker) Build(endpoint string, certpath string, imagename string) {
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
	err = dcli.CmdBuild("--force-rm=true", "--rm=true", fmt.Sprintf("--tag=\"%s\"", imagename), "./")
	fmt.Println(err)
	fmt.Println(outputbuf)
}

func (s *BPacker) CreateDockerFile() string {
	dfBucket := DockerFileBucket{
		DefaultBox:       DefaultBox,
		BuildpackerRoot:  BuildpackerRoot,
		LocalBuildPath:   s.localbuildpath,
		BuildDir:         fmt.Sprintf("%s/%s", BuildpackerRoot, BuildDir),
		Buildpack:        s.buildpack,
		BuildpackZipPath: fmt.Sprintf("%s/%s/%s", BuildpackerRoot, BuildpackDir, BuildpackZip),
		BuildpackDir:     fmt.Sprintf("%s/%s", BuildpackerRoot, BuildpackDir),
	}
	return dfBucket.Dockerfile()
}

func (s *DockerFileBucket) Dockerfile() string {
	var buffer bytes.Buffer
	dfTemplateString := s.dockerfileTemplateString()

	if tmpl, err := template.New("Dockerfile").Parse(dfTemplateString); err == nil {

		if err = tmpl.Execute(&buffer, s); err != nil {
			panic(err)
		}
	}
	return buffer.String()
}

func (s *DockerFileBucket) dockerfileTemplateString() (dockerfileTemplate string) {
	dockerfileTemplate = `FROM {{.DefaultBox}}
RUN rm /bin/sh && ln -s /bin/bash /bin/sh 
RUN apt-get install -y unzip curl ruby 
RUN mkdir -p {{.BuildpackerRoot}}
ADD {{.LocalBuildPath}} {{.BuildDir}}
ADD {{.Buildpack}} {{.BuildpackZipPath}}
RUN unzip {{.BuildpackZipPath}} -d {{.BuildpackDir}}/unpacked
RUN cd {{.BuildpackDir}} && if [ $(ls ./unpacked | wc -l) == 1 ]; then mv ./unpacked/$(ls ./unpacked) ./tmp && rm -fR ./unpacked && mv ./tmp ./unpacked; fi && rm -fR {{.BuildpackZipPath}}
RUN {{.BuildpackDir}}/unpacked/bin/detect {{.BuildpackerRoot}}/code
RUN {{.BuildpackDir}}/unpacked/bin/compile {{.BuildpackerRoot}}/code /tmp
`
	return
}
