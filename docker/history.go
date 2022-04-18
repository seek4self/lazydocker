package docker

import (
	"context"
	"encoding/json"
)

func History(imageID string) []byte {
	history, err := cli.ImageHistory(context.Background(), imageID)
	if err != nil {
		return []byte(err.Error())
	}
	buf, _ := json.MarshalIndent(history, "", "    ")
	return buf
}
