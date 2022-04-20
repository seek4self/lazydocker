package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// containers option
const (
	OptRuning = "runing"
	OptAll    = "all"
	OptExited = "exited"
)

var duration = map[string]string{
	"second":  "s",
	"seconds": "s",
	"minute":  "m",
	"minutes": "m",
	"hour":    "h",
	"hours":   "h",
	"days":    "d",
	"weeks":   "w",
	"months":  "M",
	"years":   "y",
}

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
	containers, err := cli.ContainerList(context.Background(), parseFileter(option))
	if err != nil {
		panic(err)
	}

	status := make([]ContainerStatus, 0)
	for _, v := range containers {
		status = append(status, ContainerStatus{
			ID:     v.ID[:12],
			Image:  v.Image,
			CMD:    v.Command,
			Age:    parseAge(v.Status),
			Status: v.State,
			Name:   v.Names[0][1:],
		})
	}
	return status
}

func parseAge(status string) string {
	if status[0:2] != "Up" {
		return "--"
	}
	fields := strings.Fields(status)
	subfix := duration[fields[len(fields)-1]]
	age := fields[len(fields)-2]
	if age[0] == 'a' {
		age = "1"
	}
	return strings.Join([]string{age, subfix}, "")
}

func parseFileter(option string) (options types.ContainerListOptions) {
	switch option {
	case OptAll:
		options.All = true
	case OptExited:
		options.Filters = filters.NewArgs(filters.Arg("status", "exited"))
	case OptRuning:

	default:
		options.Filters = filters.NewArgs(filters.Arg("name", option))
	}
	return
}
