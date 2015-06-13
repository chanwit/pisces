package main

import (
	"fmt"
	"os"
	"path"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Version = "0.3.0"
	app.Usage = fmt.Sprintf("A Fig-clone that understands Docker Swarm")

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "project,p",
			Value: "",
			Usage: "project name, default is the base name of the current directory",
		},
	}

	app.Commands = []cli.Command{
		cmdBuild,
		cmdUp,
	}

	/*
		app.Action = func(c *cli.Context) {
			project := c.String("project")
			if project == "" {
				dir, _ := os.Getwd()
				project := path.Base(dir)
			}
		}
	*/

	app.Run(os.Args)
}
