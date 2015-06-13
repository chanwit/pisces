package swarm

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// map[name]address
func Nodes() map[string]string {
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

	result := make(map[string]string)
	for i := 0; i < num; i++ {
		line := lines[found+1+(i*5)]
		parts := strings.SplitN(line, ":", 2)
		name := strings.TrimSpace(parts[0])
		addr := strings.TrimSpace(parts[1])
		result[name] = addr
	}

	return result
}

func Build(onNode string, service string, noCache bool) string {
	return ""
}
