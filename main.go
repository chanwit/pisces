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

	app.Commands = []cli.Command{
		cmdUp,
	}

	app.Run(os.Args)
}
