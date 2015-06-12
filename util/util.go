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
		// detect if SWARM (master) available to connect
		/*
		if output, err := exec.Command("docker-machine", "ls").Output(); err == nil {
			masterCount := 0
			lines := strings.Split(string(output), "\n")
			fmt.Println(lines[0]) // header
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasSuffix(line, "(master)") {
					masterCount++
					fmt.Println(line)
				}
			}
			if masterCount == 0 {
				fmt.Println("No Swarm master available, exiting")
				return false
			}
		} else {
			fmt.Println("Docker Machine List failed")
			return false
		}

		return false
	} */

	return true
}