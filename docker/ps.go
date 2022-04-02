package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	units "github.com/docker/go-units"
)

var (
	headers = []string{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "NAMES"}
	index   = make([]int, len(headers))
	args    = []string{"ps", "-a"}
)

type ContainerStatus struct {
	ID     string
	Image  string
	CMD    string
	Age    string
	Status string
	Ports  string
	Name   string
}

func PS(option string) []ContainerStatus {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(context.Background(), parseFileter(option))
	if err != nil {
		panic(err)
	}

	status := make([]ContainerStatus, 0)
	for _, v := range containers {
		age := "--"
		if v.Status[0:2] == "Up" {
			age = v.Status[3:]
		}
		status = append(status, ContainerStatus{
			ID:     v.ID,
			Image:  v.Image,
			CMD:    v.Command,
			Age:    age,
			Status: v.State,
			Name:   v.Names[0][1:],
		})
	}

	// cmd := exec.Command("docker", getArgs(option)...)
	// var stdout bytes.Buffer
	// var stderr bytes.Buffer
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	// err := cmd.Run()
	// if err != nil {
	// 	fmt.Println(stderr.String())
	// 	fmt.Println(cmd.String(), err)
	// }
	// // fmt.Println("out info")
	// for {
	// 	line, err := stdout.ReadString('\n')
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	// fmt.Println(len(fields), fields)
	// 	if line[0:9] == "CONTAINER" {
	// 		parseHeader(line)
	// 		continue
	// 	}
	// 	if len(line) == 0 {
	// 		continue
	// 	}
	// 	cs := ContainerStatus{
	// 		ID:     strings.TrimSpace(line[index[0]:index[1]]),
	// 		Image:  strings.TrimSpace(line[index[1]:index[2]]),
	// 		CMD:    strings.TrimSpace(line[index[2]:index[3]]),
	// 		Create: strings.TrimSpace(line[index[3]:index[4]]),
	// 		Status: strings.TrimSpace(line[index[4]:index[5]]),
	// 		Ports:  strings.TrimSpace(line[index[5]:index[6]]),
	// 		Name:   strings.TrimSpace(line[index[6]:]),
	// 	}
	// 	status = append(status, cs)
	// }
	return status
	// fmt.Println(stdout.String())
	// fmt.Println("out error")
	// fmt.Println(stderr.String())
}

func parseAge(created int64) string {
	age := time.Since(time.Unix(created, 0))
	hours := age.Hours()
	if hours < 24 {
		return age.String()
	}
	units.HumanDuration(age)
	day := hours / 24
	subfix := "days"
	if day > 30 {
		day = day / 30
		subfix = "months"
	}
	return fmt.Sprintf("%.2f%s", day, subfix)
}

func parseHeader(line string) {
	for i, h := range headers {
		index[i] = strings.Index(line, h)
	}
}

func parseFileter(option string) (options types.ContainerListOptions) {
	switch option {
	case "all":
		options.All = true
	case "exit":
		options.Filters = filters.NewArgs(filters.Arg("status", "exited"))
	case "up":

	default:
		options.Filters = filters.NewArgs(filters.Arg("name", option))
	}
	return
}

func getArgs(option string) []string {
	switch option {
	case "all":
		return []string{"ps", "-a"}
	case "exit":
		return []string{"ps", "-f", "status=exited"}
	case "up":
		return []string{"ps"}
	default:
		return []string{"ps", "-f", fmt.Sprintf("name=%s", option)}
	}
}
