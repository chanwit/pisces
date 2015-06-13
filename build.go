package main

import (
	"fmt"
	"os"

	"github.com/chanwit/pisces/conf"
	"github.com/chanwit/pisces/util"
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
	if util.CheckDockerHostVar() == false {
		os.Exit(1)
	}

	project := util.ProjectName(c)

	config, err := conf.ReadConfig()
	if err != nil {
		fmt.Printf("Config error: %s", err)
		return
	}

	// filter and reorder according to DAG
	services := config.FilterServices(c.Args())
	for _, service := range services {

		imageId := swarm.Build(service, c.Bool("no-cache"))

	}

	fmt.Println("build " + project + "_" + c.Args()[0])
}
