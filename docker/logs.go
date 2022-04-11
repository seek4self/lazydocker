package docker

import (
	"bytes"
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
)

func Logs(container string) []byte {
	body, err := cli.ContainerLogs(context.Background(), container, types.ContainerLogsOptions{
		ShowStderr: true,
		Tail:       "500",
	})
	if err != nil {
		panic(err)
	}

	buf := bytes.Buffer{}
	_, err = stdcopy.StdCopy(nil, &buf, body)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
