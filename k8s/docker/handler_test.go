package docker

import (
	"testing"
)

var DockerCli *ClientDocker

func init() {
	cli, err := NewDockerWithHost("10.10.102.96")
	if err != nil {
		panic(err)
	}
	DockerCli = cli
}
func TestPullImage(t *testing.T) {
	err := DockerCli.PullImage("alpine:3")
	if err != nil {
		t.Error(err)
		return
	}
}
func TestGetContainerStatus(t *testing.T) {
	status, err := DockerCli.GetContainerStatus("d1c8abb7bfe4")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(status)
}
func TestIsImageExist(t *testing.T) {
	exist, err := DockerCli.IsImageExist("alpine:3")
	if err != nil {
		t.Error(err)
		return
	}
	if !exist {
		t.Log("image 不存在")
	}
}

func TestGetImageID(t *testing.T) {
	imageID, err := DockerCli.GetImageID("alpine:3")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(imageID)
}

func TestRemoveContainer(t *testing.T) {
	err := DockerCli.DeleteContainer("50d3ea4e5a25")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestPauseContainer(t *testing.T) {
	err := DockerCli.PauseContainer("d1c8abb7bfe4")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestUnPauseContainer(t *testing.T) {
	err := DockerCli.UnPauseContainer("d1c8abb7bfe4")
	if err != nil {
		t.Error(err)
		return
	}
}
func TestStopContainer(t *testing.T) {
	err := DockerCli.StopContainer("d1c8abb7bfe4")
	if err != nil {
		t.Error(err)
		return
	}
}
