package pisces

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
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

type ContainerConfig struct {
	Project   string
	Service   string
	Namespace string
	Info
}

type Config map[string]Info

func ReadConfig() (Config, error) {
	// FIXME traverse back to parent until finding .yml
	data, err := ioutil.ReadFile("./docker-compose.yml")
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	config := make(Config)
	err = yaml.Unmarshal(data, &config)
	return config, err
}

func Nodes() []string {
	output, err := exec.Command("docker", "info").Output()
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}
	return getNodes(string(output))
}

func getNodes(str string) []string {
	lines := strings.Split(str, "\n")
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

func CheckDockerHostVar() bool {
	// detect DOCKER_HOST
	if os.Getenv("DOCKER_HOST") == "" {
		fmt.Println(`Environment variable "DOCKER_HOST" is required`)
		fmt.Println(`You can set it by calling: eval "$(docker-machine env --swarm <NAME>)"`)
		fmt.Println()

		// detect if SWARM (master) available to connect
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
	}

	return true
}

func FilterService(config Config, services []string) Config {
	var filteredConf Config
	if len(services) > 0 {
		// filter only matched includes
		filteredConf = make(Config)
		for _, s := range services {
			filteredConf[s] = config[s]
		}
	} else {
		filteredConf = config
	}

	return filteredConf
}

func MachineConfig(node string) []string {
	config, err := exec.Command("docker-machine", "config", node).Output()
	if err != nil {
		return nil
	}

	// FIXME: need to properly parse config, rather than splitting using " "
	args := strings.Split(string(config), " ")
	return args
}

// Load the TLS certificates/keys and, if verify is true, the CA.
func LoadTLSConfig(ca, cert, key string, verify bool) (*tls.Config, error) {
	c, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("Couldn't load X509 key pair (%s, %s): %s. Key encrypted?",
			cert, key, err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{c},
		MinVersion:   tls.VersionTLS10,
	}

	if verify {
		certPool := x509.NewCertPool()
		file, err := ioutil.ReadFile(ca)
		if err != nil {
			return nil, fmt.Errorf("Couldn't read CA certificate: %s", err)
		}
		certPool.AppendCertsFromPEM(file)
		config.RootCAs = certPool
		config.ClientAuth = tls.RequireAndVerifyClientCert
		config.ClientCAs = certPool
	} else {
		// If --tlsverify is not supplied, disable CA validation.
		config.InsecureSkipVerify = true
	}

	return config, nil
}

func CountContainers(namespace string) int {
	output, _ := exec.Command(
		"docker", "ps",
		"-a", "-q",
		"--filter", "label="+namespace).Output()

	return len(strings.Split(string(output), "\n")) - 1
}

func findContainerByNamespace(namespace string) (string, error) {
	output, _ := exec.Command(
		"docker", "ps",
		"-a", "-q", "-n", "1",
		"--filter", "label="+namespace).Output()
	containerId := strings.TrimSpace(string(output))
	if containerId == "" {
		return "", fmt.Errorf("No container found")
	}
	return containerId, nil
}

func CreateContainer(cc *ContainerConfig, i int) string {
	projectKey := fmt.Sprintf("%s_%s", cc.Project, cc.Service)
	// if image: is specify, just use it
	imageName := projectKey
	if cc.Info.Image != "" {
		imageName = cc.Info.Image
	}

	createArgs := []string{"create"}
	createArgs = append(createArgs, "-e", "affinity:image=="+imageName)
	for _, env := range cc.Info.Environment {
		createArgs = append(createArgs, "-e", env)
	}
	for _, port := range cc.Info.Ports {
		createArgs = append(createArgs, "-p", port)
	}

	// if there's a link
	// just use container ID
	//
	// Swarm's dependency filter will do the job
	for _, link := range cc.Info.Links {
		if strings.Contains(link, ":") == false {
			service := link // e.g. db
			// web_1 will link to db_1 for example
			namespaceToLink := fmt.Sprintf("pisces.%s_%s.id=%d", cc.Project, service, i)
			id, err := findContainerByNamespace(namespaceToLink)
			if err != nil {
				// nothing to link, so relax a bit
				// web_1 will link to db_?
				namespaceToLink = fmt.Sprintf("pisces.%s_%s.id", cc.Project, service)
				id, err = findContainerByNamespace(namespaceToLink)
			}
			link = id + ":" + link
		}
		createArgs = append(createArgs, "--link", link)
	}

	createArgs = append(createArgs, "-l", "pisces.container=true")
	createArgs = append(createArgs, "-l", fmt.Sprintf("%s=%d", cc.Namespace, i))
	createArgs = append(createArgs, "-t", imageName)

	cmd := exec.Command("docker", createArgs...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	containerId := strings.TrimSpace(string(output))
	return containerId
}

func StartContainer(containerId string, daemon bool) {
	startArgs := []string{"start"}
	if daemon == false {
		startArgs = append(startArgs, "-a")
	}
	startArgs = append(startArgs, containerId)

	cmd := exec.Command("docker", startArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func RemoveContainer(namespace string, id int) error {
	output, _ := exec.Command(
		"docker", "ps",
		"-a", "-q",
		"--filter", "label="+fmt.Sprintf("%s=%d", namespace, id)).Output()

	containerId := strings.TrimSpace(string(output))
	cmd := exec.Command("docker", "rm", "-f", containerId)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
