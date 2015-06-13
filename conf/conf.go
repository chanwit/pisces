package conf

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Info struct {
	Build       string
	Environment []string
	Links       []string
	Ports       []string
	Image       string
}

type Config struct {
	PodSpec  map[string]int
	Services map[string]Info
}

func ReadConfig() (Config, error) {
	// FIXME traverse back to parent until finding .yml
	data, err := ioutil.ReadFile("docker-compose.yml")
	if err != nil {
		fmt.Printf("%s\n", err)
		return Config{}, err
	}

	return parseConfig(data)
}

func parseConfig(content []byte) (Config, error) {
	services := make(map[string]Info)
	err := yaml.Unmarshal(content, &services)

	// check if first line defines pod
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if strings.HasPrefix(lines[0], "# pod:") {
		podYaml := strings.TrimLeft(lines[0], "#")
		podSpec := make(map[string]map[string]int)
		err = yaml.Unmarshal([]byte(podYaml), &podSpec)
		if err != nil {
			return Config{}, err
		}

		return Config{podSpec["pod"], services}, nil
	}

	return Config{nil, services}, err
}
