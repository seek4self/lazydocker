package docker

import (
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/daemon/network"
	"github.com/docker/go-connections/nat"
)

type State struct {
	Status string
}

type HostConfig struct {
	// Applicable to all platforms
	Binds           []string                // List of volume bindings for this container
	ContainerIDFile string                  // File (path) where the containerId is written
	NetworkMode     container.NetworkMode   // Network mode to use for the container
	PortBindings    nat.PortMap             // Port mapping between the exposed port (container) and the host
	RestartPolicy   container.RestartPolicy // Restart policy to be used for the container
	AutoRemove      bool                    // Automatically remove container when it exits
	VolumeDriver    string                  // Name of the volume driver used to mount volumes
	VolumesFrom     []string                // List of volumes to take from other container
}

type NetworkSettings struct {
	Bridge   string      // Bridge is the Bridge name the network uses(e.g. `docker0`)
	Ports    nat.PortMap // Ports is a collection of PortBinding indexed by Port
	Networks map[string]*network.EndpointSettings
}

type ContainerInfo struct {
	ID      string `json:"Id"`
	Created string
	Path    string
	Args    []string
	State   *types.ContainerState
	// Image        string
	Name         string
	RestartCount int
	Driver       string
	Platform     string
	// MountLabel      string
	// ProcessLabel    string
	AppArmorProfile string
	// ExecIDs         []string
	HostConfig *HostConfig
	SizeRw     *int64 `json:",omitempty"`
	SizeRootFs *int64 `json:",omitempty"`
	// Mounts          []types.MountPoint
	Config          *container.Config
	NetworkSettings *NetworkSettings
}

func ContainerInspect(name string) []byte {
	_, raw, err := cli.ContainerInspectWithRaw(context.Background(), name, false)
	if err != nil {
		panic(err)
	}
	var info ContainerInfo
	_ = json.Unmarshal(raw, &info)
	buf, _ := json.MarshalIndent(info, "", "    ")
	return buf

	// cmd := exec.Command("docker", "inspect", name)
	// var stdout bytes.Buffer
	// var stderr bytes.Buffer
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	// err := cmd.Run()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// return stdout.String()
}

type ImageInfo struct {
	ID            string `json:"Id"`
	RepoTags      []string
	RepoDigests   []string
	Parent        string
	Comment       string
	Created       string
	Container     string
	DockerVersion string
	Author        string
	Config        *container.Config
	Architecture  string
	Variant       string `json:",omitempty"`
	Os            string
	OsVersion     string `json:",omitempty"`
	Size          int64
	VirtualSize   int64
}

func ImageInspect(id string) []byte {
	_, raw, err := cli.ImageInspectWithRaw(context.Background(), id)
	if err != nil {
		panic(err)
	}
	var info ImageInfo
	_ = json.Unmarshal(raw, &info)
	buf, _ := json.MarshalIndent(info, "", "    ")
	return buf
}
