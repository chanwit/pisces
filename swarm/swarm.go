package swarm

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func Nodes() []string {
	output, err := exec.Command("docker", "info").Output()
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}

	lines := strings.Split(string(output), "\n")
	found := 0
	num := 0
	for i, line := range lines {
		line := strings.TrimSpace(line)
		if strings.HasPrefix(line, "\bNodes") {
			num, _ = strconv.Atoi(strings.TrimSpace(strings.SplitN(line, ":", 2)[1]))
			found = i
			break
		}
	}

	result := []string{}
	for i := 0; i < num; i++ {
		line := lines[found+1+(i*5)]
		name := strings.TrimSpace(strings.SplitN(line, ":", 2)[0])
		result = append(result, name)
	}
	return result
}
