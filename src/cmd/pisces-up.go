package main

import (
	"fmt"
	"os"
	"path"

	"github.com/chanwit/pisces"
	"github.com/codegangsta/cli"
)

func main() {
	app := pisces.NewApp("up")
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "d",
			Usage: "run in background",
		},
	}
	app.Action = func(c *cli.Context) {
		if c.Bool("help") {
			cli.ShowAppHelp(c)
			os.Exit(0)
		}
		action(c)
	}
	app.Run(os.Args)
}

func action(c *cli.Context) {
	if pisces.CheckDockerHostVar() == false {
		os.Exit(1)
	}

	conf, err := pisces.ReadConfig()
	if err != nil {
		fmt.Printf("Config error: %s", err)
		os.Exit(1)
	}

	dir, err := os.Getwd()
	project := path.Base(dir)
	// up -d front web db
	services := c.Args()

	if err := pisces.CheckServices(conf, services); err != nil {
		os.Exit(1)
	}

	filteredService, order := pisces.FilterService(conf, services)
	for _, service := range order {

		info, exist := filteredService.Services[service]
		if exist == false {
			continue
		}

		projectKey := fmt.Sprintf("%s_%s", project, service)
		namespace := fmt.Sprintf("pisces.%s.id", projectKey)

		containerConfig := &pisces.ContainerConfig{
			project,
			service,
			namespace,
			conf.PodSpec,
			info,
		}

		nextId := pisces.CountContainers(namespace) + 1
		id := pisces.CreateContainer(containerConfig, nextId)
		pisces.StartContainer(id, c.Bool("d"))

	}

}
