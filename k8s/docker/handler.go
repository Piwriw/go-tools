package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"io"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var (
	DockerUserName     = "piwriw"
	DockerUserPassword = "Joohwan."
)

type ClientDocker struct {
	Client *client.Client
}

func NewDockerWithHost(host string) (*ClientDocker, error) {
	cli, err := client.NewClientWithOpts(client.WithHost(fmt.Sprintf("tcp://%s:2375", host)), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &ClientDocker{Client: cli}, nil
}

func NewDocker() (*ClientDocker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &ClientDocker{Client: cli}, nil
}

// DownloadFile 从指定的网址获取文件，并保存到 savePath 中
func DownloadFile(fileUrl string, savePath string) error {
	// 从url中读取下载文件的名称
	fileURL, err := url.Parse(fileUrl)
	if err != nil {
		fmt.Errorf("DownloadFile:出错 :%v", err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := strings.Split(segments[len(segments)-1], "?")[0]
	filePath := filepath.Join(savePath, fileName)
	//todo : config.AIModelPath = filePath
	// 判断该模型文件是否已经存在，如果存在，就删除重新下载
	tf := DirFileIsExist(path)
	if tf {
		if err = os.Remove(filePath); err != nil {
			return fmt.Errorf("DownloadFile:出错:%v", err)
		}
	}
	_ = MakeDirAll(savePath)
	// create blank file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("DownloadFile: 出错: %v", err)
	}
	c := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	defer file.Close()
	// Put content on file
	resp, err := c.Get(fileUrl)
	if err != nil {
		return fmt.Errorf("4 DownloadFile: 出错: %v", err)
	}
	defer resp.Body.Close() // 不使用 EnsureReaderClosed

	_, err = io.Copy(file, resp.Body)

	return nil
}

func (c *ClientDocker) GetContainerStatus(containerID string) (string, error) {
	inspect, err := c.Client.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return "", errors.Wrap(err, "inspect container failed")
	}
	return inspect.State.Status, nil
}

// removeRecordByRemovedStatus 写成协程的方式进行工作
func (c *ClientDocker) removeRecordByRemovedStatus(containerID string) error {
	ctx := context.Background()
	status, err := c.GetContainerStatus(containerID)
	if err != nil {
		return errors.Wrap(err, "get container status failed")
	}
	klog.Infof("this container status is %+v", status)
	klog.Infof("this container ID:%+v -> is stop, it will be removed", containerID)

	if err := c.Client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         true,
	}); err != nil {
		return errors.Wrap(err, "remove container failed")
	}

	// 删除 db 中的对应的任务数据
	//klog.Infof("delete this record by container ID ")
	//_, err = dao.DeleteNodeTasksCByContainerID(containerID)
	//if err != nil {
	//	return errors.Wrap(err, "delete record by container ID")
	//}
	return nil
}

// IsImageExist 判断镜像本地是否存在
func (c *ClientDocker) IsImageExist(imageAddr string) (bool, error) {
	ctx := context.Background()
	images, err := c.Client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return false, err
	}
	for _, image := range images {
		if imageAddr == image.RepoTags[0] {
			return true, nil
		}
	}
	return false, nil
}

// GetImageID 根据镜像地址获取镜像ID
func (c *ClientDocker) GetImageID(imageAddr string) (string, error) {
	imageID := ""
	if exist, _ := c.IsImageExist(imageAddr); !exist {
		return "", fmt.Errorf("镜像名为：%v 的镜像不存在", imageAddr)
	}
	ctx := context.Background()
	images, err := c.Client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return "", err
	}
	for _, image := range images {
		if imageAddr == image.RepoTags[0] {
			imageID = image.ID
			break
		}
	}
	return imageID, nil
}

// PullImage 从给定的镜像地址拉取镜像
func (c *ClientDocker) PullImage(imageAddr string) error {
	authConfig := registry.AuthConfig{
		Username: DockerUserName,
		Password: DockerUserPassword,
	}
	marshal, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(marshal)
	out, err := c.Client.ImagePull(context.Background(), imageAddr, types.ImagePullOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		return fmt.Errorf("PullImage:出错 :%v", err)
	}
	defer out.Close()
	klog.Infof("info:%v", out)
	return nil
}

// PauseContainer 暂停一个容器
func (c *ClientDocker) PauseContainer(containerID string) error {
	paused, err := c.IsContainerPaused(containerID)
	if err != nil {
		return err
	}
	if paused {
		return errors.Errorf("this container has paused, containerID:%s", containerID)
	}
	if err := c.Client.ContainerPause(context.Background(), containerID); err != nil {
		return err
	}
	return nil
}

// UnPauseContainer Unpause Container
func (c *ClientDocker) UnPauseContainer(containerID string) error {
	containerExist, err := c.IsContainerExist(containerID)
	if err != nil {
		return err
	}
	if !containerExist {
		return errors.Errorf("[%s]容器不存在", containerID)
	}
	exist, err := c.IsContainerPaused(containerID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.Errorf("[%s]容器没有pause", containerID)
	}
	if err := c.Client.ContainerUnpause(context.Background(), containerID); err != nil {
		klog.Errorf("Unpause[%s]容器失败，失败原因为: %v\n", containerID, err)
		return err
	}
	return nil
}

// IsContainerPaused 通过containerID判断某容器是否已经Pause
func (c *ClientDocker) IsContainerPaused(containerID string) (bool, error) {
	inspect, err := c.Client.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return false, errors.Wrap(err, "inspect container failed")
	}
	if inspect.State.Status == "paused" {
		return true, nil
	}

	return false, nil
}

// IsContainerExist 根据传来的 容器ID 来判断该容器是否已经存在
func (c *ClientDocker) IsContainerExist(containerID string) (bool, error) {
	containers, err := c.Client.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return false, err
	}

	for _, c := range containers {
		if strings.Contains(c.ID, containerID) {
			return true, nil
		}
	}
	return false, nil
}

// StopContainer 根据containerID 停止某一个容器
func (c *ClientDocker) StopContainer(containerID string) error {
	exist, err := c.IsContainerExist(containerID)
	if err != nil {
		return err
	}
	if !exist {
		klog.Infof("container [ %v ] does not exist\n", containerID)
		return nil
	}
	// delete container
	if err = c.Client.ContainerStop(context.Background(), containerID, container.StopOptions{}); err != nil {
		klog.Errorf("stop [ %v ] container error \n", containerID)
		return err
	}
	return nil
}

// DeleteContainer 根据containerID 删除某一个容器
func (c *ClientDocker) DeleteContainer(containerID string) error {
	if err := c.StopContainer(containerID); err != nil {
		return err
	}
	if err := c.Client.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}); err != nil {
		klog.Errorf("移除[%v]容器失败，失败原因为: %v\n", containerID, err)
		return err
	}
	return nil
}
