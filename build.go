package main

import (
	"fmt"
	"os"

	"github.com/chanwit/pisces/conf"
	"github.com/chanwit/pisces/swarm"
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

	dir, _ := os.Getwd()
	project := util.ProjectName(c)

	config, err := conf.ReadConfig()
	if err != nil {
		fmt.Printf("Config error: %s", err)
		os.Exit(2)
	}

	// filter and reorder according to DAG
	// build does not require the service order
	filteredConfig, _ := config.FilterServices(c.Args())
	for service, info := range filteredConfig.Services {
		for name, addr := range swarm.Nodes() {
			spec := swarm.BuildSpec{
				Info:       info,
				NodeName:   name,
				NodeAddr:   addr,
				ProjectDir: dir,
				Project:    project,
				Service:    service,
				NoCache:    c.Bool("no-cache"),
			}
			imageId := swarm.Build(spec)
			if imageId == "" {
				os.Exit(5)
			}
			fmt.Println(imageId)
		}
	}

}
