package util

import (
	"fmt"
	"os"
	"path"

	"github.com/codegangsta/cli"
)

func CheckDockerHostVar() bool {
	// detect DOCKER_HOST
	if os.Getenv("DOCKER_HOST") == "" {
		fmt.Println(`Environment variable "DOCKER_HOST" is required.`)
		fmt.Println(`You can set it by calling: eval "$(docker-machine env --swarm <NAME>)".`)
		fmt.Println()

		return false
	}

	return true
}

func ProjectName(c *cli.Context) string {
	project := c.String("project")
	if project == "" {
		dir, _ := os.Getwd()
		project = path.Base(dir)
	}

	return project
}
