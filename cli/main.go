package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/xchapter7x/buildpacker"
)

func main() {
	app := cli.NewApp()
	app.Name = "buildpacker"
	app.Usage = "bootstrap a docker image from a buildpack... yeah i said it"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "buildpack, bp",
			Usage:  "url or relative path to buildpack .zip file",
			EnvVar: "BP_BUILDPACK",
		},
		cli.StringFlag{
			Name:   "codepath, c",
			Usage:  "the relative path to your codebase",
			EnvVar: "BP_CODEPATH",
		},
		cli.StringFlag{
			Name:   "dockerhost, dh",
			Usage:  "the host of your docker server (ex. tcp://192.168.59.103:2376) ",
			EnvVar: "DOCKER_HOST",
		},
		cli.StringFlag{
			Name:   "dockercert, dc",
			Usage:  "the path of your docker certs ",
			EnvVar: "DOCKER_CERT_PATH",
		},
		cli.StringFlag{
			Name:   "imagename, in",
			Usage:  "the name of the image your will create",
			EnvVar: "BP_IMAGENAME",
		},
	}
	app.Action = func(c *cli.Context) {

		if c.String("buildpack") != "" && c.String("codepath") != "" && c.String("dockerhost") != "" && c.String("dockercert") != "" && c.String("imagename") != "" {
			bdPkr := buildpacker.New(c.String("buildpack"), c.String("codepath"))
			bdPkr.Build(c.String("dockerhost"), c.String("dockercert"), c.String("imagename"))

		} else {
			fmt.Println(c.String("buildpack"), c.String("codepath"), c.String("dockerhost"), c.String("dockercert"), c.String("imagename"))
			//cli.ShowAppHelp(c)
		}
	}
	app.Run(os.Args)
}
