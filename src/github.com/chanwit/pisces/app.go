package pisces

import (
    "os"
	"path"
	"fmt"
	"github.com/codegangsta/cli"
)

func NewApp(subcmd string) *cli.App {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Version = "0.3.0"
	app.Usage = fmt.Sprintf("(%s) A Fig-clone that understands Docker Swarm", subcmd)
	return app
}