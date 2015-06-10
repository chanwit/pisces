package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/chanwit/pisces"
	"github.com/codegangsta/cli"
)

func main() {
	app := pisces.NewApp("scale")
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
	// scale web=2 db=3
	serviceNums := c.Args()
	serviceMap := make(map[string]int)
	services := []string{}
	for _, serviceNum := range serviceNums {
		parts := strings.SplitN(serviceNum, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("Error: %s not in the correct format", serviceNum)
			os.Exit(1)
		}
		service := parts[0]
		if _, exist := conf[service]; !exist {
			fmt.Printf("Service '%s' not defined in docker-compose.yml", service)
			os.Exit(1)
		}
		num, _ := strconv.Atoi(parts[1])

		services = append(services, service)
		serviceMap[service] = num
	}

	filteredService, order := pisces.FilterService(conf, services)
	for _, service := range order {

		info, exist := filteredService[service]
		if exist == false {
			continue
		}

		projectKey := fmt.Sprintf("%s_%s", project, service)
		namespace := fmt.Sprintf("pisces.%s.id", projectKey)

		containerConfig := &pisces.ContainerConfig{
			project,
			service,
			namespace,
			info,
		}

		count := pisces.CountContainers(namespace)
		nextId := count + 1
		num := serviceMap[service]
		// count = nextId - 1 (1-1 = 0)
		// delta = num - count
		if num == count {
			// do nothing
		} else if num > count {
			// scale up
			conIds := []string{}
			for i := nextId; i <= num; i++ {
				id := pisces.CreateContainer(containerConfig, i)
				conIds = append(conIds, id)
			}
			for _, id := range conIds {
				pisces.StartContainer(id, true)
			}
		} else if count > num {
			// scale down
			for i := count; i > num; i-- {
				pisces.RemoveContainer(namespace, i)
			}
		}
	}
}
