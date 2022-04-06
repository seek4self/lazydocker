package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

type ImageInfo struct {
}

func Images() []types.ImageSummary {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	return images
}
