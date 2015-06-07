package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/chanwit/pisces"
	"github.com/codegangsta/cli"
)

func main() {
	app := pisces.NewApp("up")
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "d",
			Usage: "run in background",
		},
	}
	app.Action = func(c *cli.Context) {
		if c.Bool("help") {
			cli.ShowAppHelp(c)
			os.Exit(0)
		}
		action(c)
	}
	app.Run(os.Args)
}

func action(c *cli.Context) {
	if pisces.CheckDockerHostVar() == false {
		os.Exit(1)
	}

	conf, err := pisces.ReadConfig()
	if err != nil {
		fmt.Printf("Config error: %s", err)
		os.Exit(1)
	}

	dir, err := os.Getwd()
	project := path.Base(dir)
	// up -d front web db
	services := c.Args()

	// check if unit is defined in the conf
	for _, s := range services {
		if _, exist := conf[s]; !exist {
			fmt.Printf("'%s' is not defined in docker-compose.yml.\n", s)
			os.Exit(1)
		}
	}

	/*
			certPath := os.Getenv("DOCKER_CERT_PATH")
			ca := path.Join(certPath, "ca.pem")
			cert := path.Join(certPath, "cert.pem")
			key := path.Join(certPath, "key.pem")
			verify := os.Getenv("DOCKER_TLS_VERIFY") == "1"
			tlsConfig, err := pisces.LoadTLSConfig(ca, cert, key, verify)
			docker, err := dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), tlsConfig)

			images, err := docker.ListImages()
			for _, image := range images {
				fmt.Println(image)
			}

			containerConfig := &dockerclient.ContainerConfig{
		        Image: "kapook_web",
		        AttachStdin: false,
		        Tty:   false}
		    _, err = docker.CreateContainer(containerConfig, "")
		    if err != nil {
		        fmt.Println(err)
		    }
	*/

	for service, info := range pisces.FilterService(conf, services) {

		projectKey := fmt.Sprintf("%s_%s", project, service)
		namespace := fmt.Sprintf("pisces.%s.id", projectKey)

		// if image: is specify, just use it
		imageName := projectKey
		if info.Image != "" {
			imageName = info.Image
		}

		output, _ := exec.Command(
			"docker", "ps",
			"-a", "-q",
			"--filter", "label="+namespace).Output()

		nextId := len(strings.Split(string(output), "\n"))

		createArgs := []string{"create"}
		createArgs = append(createArgs, "-e", "affinity:image=="+imageName)
		for _, env := range info.Environment {
			createArgs = append(createArgs, "-e", env)
		}
		for _, port := range info.Ports {
			createArgs = append(createArgs, "-p", port)
		}

		createArgs = append(createArgs, "-l", "pisces.container=true")
		createArgs = append(createArgs, "-l", fmt.Sprintf("%s=%d", namespace, nextId))
		createArgs = append(createArgs, "-t", imageName)

		cmd := exec.Command("docker", createArgs...)
		cmd.Stderr = os.Stderr
		output, err = cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		containerId := strings.TrimSpace(string(output))
		startArgs := []string{"start"}
		if c.Bool("d") == false {
			startArgs = append(startArgs, "-a")
		}
		startArgs = append(startArgs, containerId)

		cmd = exec.Command("docker", startArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}

}
