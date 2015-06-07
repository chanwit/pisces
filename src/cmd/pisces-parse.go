package main

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type Info struct {
	Build       string
	Environment []string
	Links       []string
	Ports       []string
	Image       string
}

func main() {
	data, err := ioutil.ReadFile("./docker-compose.yml")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	// fmt.Printf("%s\n----------\n", string(data))
	config := make(map[string]Info)
	err = yaml.Unmarshal(data, &config)
	if err == nil {
		fmt.Printf("%s", config)
	}
}
