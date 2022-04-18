package docker

import (
	"bytes"
	"sync"

	"github.com/docker/docker/client"
)

var cli *client.Client

var bufPool = &sync.Pool{New: func() interface{} { return bytes.NewBuffer(nil) }}

func init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
}

func resetBuf(buf *bytes.Buffer) {
	buf.Reset()
	bufPool.Put(buf)
}
