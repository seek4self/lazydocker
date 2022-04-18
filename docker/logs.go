package docker

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"

	"github.com/docker/docker/api/types"
)

func Logs(container string) []byte {
	body, err := cli.ContainerLogs(context.Background(), container, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Tail:       "500",
	})
	if err != nil {
		return []byte(err.Error())
	}
	// r := bufio.NewReader(body)
	buf := bytes.Buffer{}
	header := make([]byte, 8)
	for {
		n, err := io.ReadFull(body, header)
		if n < 8 || err == io.EOF {
			break
		}
		if header[1]|header[2]|header[3] != 0 {
			b, err := io.ReadAll(body)
			if err != nil {
				buf.Write([]byte(err.Error()))
				break
			}
			buf.Write(header)
			buf.Write(b)
			break
		}
		size := binary.BigEndian.Uint32(header[4:])
		frame := make([]byte, size)
		_, err = io.ReadFull(body, frame)
		if err == io.EOF {
			break
		}
		buf.Write(frame)
	}
	// buf, _ := io.ReadAll(body)
	return buf.Bytes()
}
