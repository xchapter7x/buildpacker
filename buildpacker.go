package buildpacker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	dkr "github.com/docker/docker/api/client"
	docker "github.com/fsouza/go-dockerclient"
)

func Build(endpoint string, certpath string) {
	cert := fmt.Sprintf("%s/cert.pem", certpath)
	key := fmt.Sprintf("%s/key.pem", certpath)
	ca := fmt.Sprintf("%s/ca.pem", certpath)
	client, err := docker.NewTLSClient(endpoint, cert, key, ca)
	outputbuf, errbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	reader := strings.NewReader("FROM redis")
	inputbuf := ioutil.NopCloser(reader)
	endpoint = strings.TrimPrefix(endpoint, "tcp://")
	dcli := dkr.NewDockerCli(inputbuf, outputbuf, errbuf, key, "tcp", endpoint, client.TLSConfig)
	err = dcli.CmdBuild("-")
	fmt.Println(err)
	fmt.Println(outputbuf)
}
