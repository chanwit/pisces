package conf

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	PodSpec  map[string]int
	Services map[string]Info
}

type Info struct {
	Build       string
	Environment []string
	Links       []string
	Ports       []string
	Image       string
}

func (config *Config) FilterServices(services []string) (*Config, []string) {
	g := make(graph)
	for service, info := range config.Services {
		value := []string{}
		for _, link := range info.Links {
			parts := strings.SplitN(link, ":", 2)
			value = append(value, parts[len(parts)-1])
		}
		g[service] = value
	}
	order, _ := topoSortDFS(g)

	// nothing to filter, return the whole
	if len(services) == 0 {
		return config, order
	}

	filteredConf := &Config{config.PodSpec, make(map[string]Info)}
	for _, s := range services {
		filteredConf.Services[s] = config.Services[s]
	}
	return filteredConf, order
}

func ReadConfig() (*Config, error) {
	// FIXME traverse back to parent until finding .yml
	data, err := ioutil.ReadFile("docker-compose.yml")
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}

	return parseConfig(data)
}

func parseConfig(content []byte) (*Config, error) {
	services := make(map[string]Info)
	err := yaml.Unmarshal(content, &services)

	// check if first line defines pod
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if strings.HasPrefix(lines[0], "# pod:") {
		podYaml := strings.TrimLeft(lines[0], "#")
		podSpec := make(map[string]map[string]int)
		err = yaml.Unmarshal([]byte(podYaml), &podSpec)
		if err != nil {
			return nil, err
		}

		return &Config{podSpec["pod"], services}, nil
	}

	return &Config{nil, services}, err
}

type graph map[string][]string

func topoSortDFS(g graph) (order, cyclic []string) {
	L := make([]string, len(g))
	i := 0 // len(L)
	temp := map[string]bool{}
	perm := map[string]bool{}
	var cycleFound bool
	var cycleStart string
	var visit func(string)
	visit = func(n string) {
		switch {
		case temp[n]:
			cycleFound = true
			cycleStart = n
			return
		case perm[n]:
			return
		}
		temp[n] = true
		for _, m := range g[n] {
			visit(m)
			if cycleFound {
				if cycleStart > "" {
					cyclic = append(cyclic, n)
					if n == cycleStart {
						cycleStart = ""
					}
				}
				return
			}
		}
		delete(temp, n)
		perm[n] = true
		L[i] = n
		i++
	}
	for n := range g {
		if perm[n] {
			continue
		}
		visit(n)
		if cycleFound {
			return nil, cyclic
		}
	}
	return L, nil
}