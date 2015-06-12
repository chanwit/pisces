package main

import (
	"fmt"

	"github.com/codegangsta/cli"
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
	fmt.Println("build " + c.Args()[0])
}
