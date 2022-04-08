package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func Images(option string) []types.ImageSummary {
	images, err := cli.ImageList(context.Background(), parseImageFileter(option))
	if err != nil {
		panic(err)
	}
	return images
}

func parseImageFileter(option string) (options types.ImageListOptions) {
	switch option {
	case OptAll:
		options.All = true
	case "":

	default:
		options.Filters = filters.NewArgs(filters.Arg("reference", option))
	}
	return
}
