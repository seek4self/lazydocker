package docker

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

var (
	headers = []string{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "NAMES"}
	index   = make([]int, len(headers))
)

type ContainerStatus struct {
	ID     string
	Image  string
	CMD    string
	Create string
	Status string
	Ports  string
	Name   string
	index  []int
}

func PS() []ContainerStatus {
	cmd := exec.Command("docker", "ps", "-a")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("out info")
	status := make([]ContainerStatus, 0)
	for {
		line, err := stdout.ReadString('\n')
		if err == io.EOF {
			break
		}
		// fmt.Println(len(fields), fields)
		if line[0:9] == "CONTAINER" {
			parseHeader(line)
			continue
		}
		if len(line) == 0 {
			continue
		}
		cs := ContainerStatus{
			ID:     strings.TrimSpace(line[index[0]:index[1]]),
			Image:  strings.TrimSpace(line[index[1]:index[2]]),
			CMD:    strings.TrimSpace(line[index[2]:index[3]]),
			Create: strings.TrimSpace(line[index[3]:index[4]]),
			Status: strings.TrimSpace(line[index[4]:index[5]]),
			Ports:  strings.TrimSpace(line[index[5]:index[6]]),
			Name:   strings.TrimSpace(line[index[6]:]),
		}
		status = append(status, cs)
	}
	return status
	// fmt.Println(stdout.String())
	// fmt.Println("out error")
	// fmt.Println(stderr.String())
}

func parseHeader(line string) {
	for i, h := range headers {
		index[i] = strings.Index(line, h)
	}
}
