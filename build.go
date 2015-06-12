package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/chanwit/pisces/util"
)

var cmdBuild = cli.Command{
	Name:   "build",
	Action: build,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "no-cache",
			Usage: "build refreshly without cache",
		},
	},
}

func build(c *cli.Context) {
	if util.CheckDockerHostVar() == false {
		os.Exit(1)
	}

	fmt.Println("build " + c.Args()[0])
}
