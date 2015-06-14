package main

import (
	log "github.com/Sirupsen/logrus"
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
		log.Infof("up -d %s", c.Args()[0])
	}
}
