package swarm

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/chanwit/pisces/conf"
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

type BuildSpec struct {
	Info       conf.Info
	NodeName   string
	NodeAddr   string
	ProjectDir string
	Project    string
	Service    string
	NoCache    bool
}

func Build(spec BuildSpec) string {
	home := os.Getenv("HOME")
	imageName := spec.Project + "_" + spec.Service
	args := []string{"build"}
	if spec.NoCache {
		args = append(args, "--no-cache")
	}
	args = append(args, "-t", imageName, ".")
	cmd := exec.Command("docker", args...)
	cmd.Env = append(cmd.Env, "DOCKER_HOST="+spec.NodeAddr)
	cmd.Dir = path.Join(spec.ProjectDir, spec.Info.Build)
	// cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if strings.HasSuffix(spec.NodeAddr, ":2376") ||
		strings.HasSuffix(spec.NodeAddr, ":3376") {
		// assume that it's listed in Docker-Machine
		certPath := path.Join(home, ".docker/machine/machines", spec.NodeName)
		cmd.Env = append(cmd.Env, "DOCKER_TLS_VERIFY=1")
		cmd.Env = append(cmd.Env, "DOCKER_CERT_PATH="+certPath)
	}

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	lastLine := lines[len(lines)-1]
	if strings.Contains(lastLine, "Successfully built ") {
		result := strings.SplitN(lastLine, "Successfully built ", 2)[1]
		return result
	}
	return ""

}
