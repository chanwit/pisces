package main

import (
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/chanwit/pisces"
	"github.com/codegangsta/cli"
)

func main() {
	app := pisces.NewApp("stop")
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

	// check if unit is defined in the conf
	for _, s := range services {
		if _, exist := conf[s]; !exist {
			fmt.Printf("'%s' is not defined in docker-compose.yml.\n", s)
			os.Exit(1)
		}
	}

	filteredService, order := pisces.FilterService(conf, services)
	// stop in reverse topology
	sort.Sort(sort.Reverse(sort.StringSlice(order)))
	for _, service := range order {

		_, exist := filteredService[service]
		if exist == false {
			continue
		}

		projectKey := fmt.Sprintf("%s_%s", project, service)
		namespace := fmt.Sprintf("pisces.%s.id", projectKey)

		pisces.StopContainers(namespace)

	}

}
