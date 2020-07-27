package orchestrator

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/google/gnxi/gnxi_tester/config"
	"github.com/moby/moby/client"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

var cli *client.Client

// InitContainers will check if the containers are running and run them if not.
func InitContainers(names []string) error {
	var err error
	if cli == nil {
		cli, err = client.NewEnvClient()
		if err != nil {
			return err
		}
	}
	build := viper.GetString("docker.build")
	if err = pullImage(build); err != nil {
		return err
	}
	runtime := viper.GetString("docker.runtime")
	if err = pullImage(runtime); err != nil {
		return err
	}

	tests := config.GetTests()
	if len(names) == 0 {
		names = make([]string, len(tests))
		i := 0
		for name := range tests {
			names[i] = name
			i++
		}
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}
	for _, c := range containers {
		for _, name := range c.Names {
			for i, testName := range names {
				if name == testName {
					if c.Status != "running" {
						if err := cli.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{}); err != nil {
							return err
						}
					}
					copy(names[i:], names[i+1:])
					names[len(names)-1] = ""
					names = names[:len(names)-1]
				}
			}
		}
	}
	for _, name := range names {
		if err := createContainer(name); err != nil {
			return err
		}
	}
	return nil
}

func createContainer(name string) error {
	cli.ContainerCreate(context.Background(), &container.Config{}, &container.HostConfig{}, &network.NetworkingConfig{}, name)
	return nil
}

func pullImage(name string) error {
	results, err := cli.ImageSearch(context.Background(), name, types.ImageSearchOptions{})
	if err != nil {
		return err
	}
	if len(results) == 0 {
		closer, err := cli.ImagePull(context.Background(), name, types.ImagePullOptions{})
		if err != nil {
			return err
		}
		closer.Close()
	}
	return nil
}

// RunContainer runs an executable in a docker conatainer.
var RunContainer = func(name, args string) (out string, code int, err error) {
	return "", 0, nil
}
