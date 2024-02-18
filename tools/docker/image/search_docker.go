package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

func main() {
	_, err := searchImage("nginx")
	if err != nil {
		fmt.Println(err)
	}

}

func searchImage(imageName string) ([]registry.SearchResult, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	images, err := cli.ImageSearch(ctx, imageName, types.ImageSearchOptions{})
	if err != nil {
		return nil, err
	}

	return images, nil
}
