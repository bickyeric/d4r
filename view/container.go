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
	*tview.Table
	client  *client.Client
	content []types.Container
}

func NewContainer(cli *client.Client) *Container {
	return &Container{
		Table:  tview.NewTable(),
		client: cli,
	}
}

func (c *Container) init() {
	containers, err := c.client.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	c.content = containers

	for i, name := range []string{"CONTAINER ID", "NAMES", "IMAGE", "COMMAND", "CREATED", "PORTS"} {
		c.SetCell(0, i, tview.NewTableCell(name).SetSelectable(false))
	}

	for index, container := range c.content {
		c.SetCell(index+1, 0, tview.NewTableCell(container.ID[7:19]))
		c.SetCell(index+1, 1, tview.NewTableCell(fmt.Sprint(container.Names[0][1:])))
		c.SetCell(index+1, 2, tview.NewTableCell(container.Image))
		c.SetCell(index+1, 3, tview.NewTableCell(container.Command))
		c.SetCell(index+1, 4, tview.NewTableCell(fmt.Sprint(time.Unix(container.Created, 0))))
		c.SetCell(index+1, 5, tview.NewTableCell(fmt.Sprint(container.Ports)))
	}
	c.SetSelectable(true, false)
	c.SetBorder(true).SetTitle(fmt.Sprintf(" %s [%d] ", "Container", c.GetRowCount()))
}
