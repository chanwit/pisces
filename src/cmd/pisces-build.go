package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/chanwit/pisces"
	"github.com/codegangsta/cli"
)

func main() {
	app := pisces.NewApp("build")
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "no-cache",
			Usage: "disable cache during build",
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
		return
	}

	conf, err := pisces.ReadConfig()
	if err != nil {
		fmt.Printf("Config error: %s", err)
		return
	}

	dir, err := os.Getwd()
	project := path.Base(dir)
	nodes := pisces.Nodes()

	// build front web db
	units := c.Args()

	filteredService, order := pisces.FilterService(conf, units)
	for _, key := range order {

		info, exist := filteredService[key]
		if exist == false {
			continue
		}

		// skip if build: is not specified
		if info.Build == "" {
			continue
		}

		imageName := fmt.Sprintf("%s_%s", project, key)
		buildDir := path.Join(dir, info.Build)

		// TODO build on *every node* of the cluster
		for _, node := range nodes {
			args := append(pisces.MachineConfig(node), "build")
			if c.Bool("no-cache") {
				args = append(args, "--no-cache")
			}
			args = append(args, "-t", imageName, ".")
			cmd := exec.Command("docker", args...)
			cmd.Dir = buildDir
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr

			fmt.Printf("Building image '%s' on node: '%s'\n", imageName, node)
			err := cmd.Run()
			if err != nil {
				fmt.Printf("%s\n", err)
			}
			fmt.Println()
		}

	}

}
