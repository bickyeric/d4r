package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	arr, err := cli.ImageList(context.Background(), types.ImageListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	for _, i := range arr {
		fmt.Printf("%s %v\n", i.ID, i)
	}
}
