package buildpacker

import (
	"archive/tar"
	"bytes"
	"fmt"
	"log"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

func Build(endpoint string, certpath string) {
	cert := fmt.Sprintf("%s/cert.pem", certpath)
	key := fmt.Sprintf("%s/key.pem", certpath)
	ca := fmt.Sprintf("%s/ca.pem", certpath)

	client, err := docker.NewTLSClient(endpoint, cert, key, ca)
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputbuf)
	tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: 10, ModTime: t, AccessTime: t, ChangeTime: t})
	tr.Write([]byte("FROM ubuntu\n"))
	tr.Close()
	opts := docker.BuildImageOptions{
		Name:         "test",
		InputStream:  inputbuf,
		OutputStream: outputbuf,
	}
	if err := client.BuildImage(opts); err != nil {
		log.Fatal(err)
	}
}
