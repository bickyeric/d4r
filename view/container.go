package view

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/rivo/tview"
)

type Container struct {
	client  *client.Client
	content []types.Container
}

func NewContainer(cli *client.Client) *Container {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}
	return &Container{
		client:  cli,
		content: containers,
	}
}

func (c *Container) Title() (name string, count int) {
	return "Container", len(c.content)
}

func (c *Container) Headers() []string {
	return []string{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "PORTS", "NAMES"}
}

func (c *Container) SetTable(t *tview.Table) {
	for index, container := range c.content {
		t.SetCell(index+1, 0, tview.NewTableCell(container.ID[7:19]))
		t.SetCell(index+1, 1, tview.NewTableCell(container.Image))
		t.SetCell(index+1, 2, tview.NewTableCell(container.Command))
		t.SetCell(index+1, 3, tview.NewTableCell(fmt.Sprint(time.Unix(container.Created, 0))))
		t.SetCell(index+1, 4, tview.NewTableCell(fmt.Sprint(container.Ports)))
		t.SetCell(index+1, 5, tview.NewTableCell(fmt.Sprint(container.Names[0][1:])))
	}
}
