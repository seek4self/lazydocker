package docker

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Inspect(name string) string {
	cmd := exec.Command("docker", "inspect", name)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	return stdout.String()
}
