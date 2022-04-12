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
		ShowStdout: true,
		Tail:       "500",
	})
	if err != nil {
		panic(err)
	}

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}
	n, err := stdcopy.StdCopy(&stdout, &stderr, body)
	if n == 0 {
		return nil
	}
	if err != nil {
		panic(err)
	}
	if stderr.Len() > 0 {
		return stderr.Bytes()
	}
	return stdout.Bytes()
}
