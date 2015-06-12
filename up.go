package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var cmdUp = cli.Command{
	Name:   "up",
	Action: up,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "d",
			Usage: "run in background",
		},
	},
}

func up(c *cli.Context) {
	if c.Bool("d") {
		fmt.Println("up -d " + c.Args()[0])
	}
}
