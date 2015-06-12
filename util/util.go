package util

import (
	"os"
	"fmt"
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