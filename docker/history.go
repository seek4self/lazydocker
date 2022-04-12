package docker

import (
	"context"
	"encoding/json"
)

func History(imageID string) []byte {
	history, err := cli.ImageHistory(context.Background(), imageID)
	if err != nil {
		panic(err)
	}
	buf, _ := json.MarshalIndent(history, "", "    ")
	return buf
}
