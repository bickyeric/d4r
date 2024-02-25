package view

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/rivo/tview"
)

type Container struct {
	*tview.Table
	app     *App
	content []types.Container
	viewAll bool
}

func NewContainer(app *App) *Container {
	return &Container{
		Table: tview.NewTable(),
		app:   app,
	}
}

func (c *Container) init() {
	c.reloadData()
	c.SetSelectable(true, false)
	c.SetBorder(true)

	c.bindKeys()
}

func (c *Container) reloadData() {
	c.Clear()

	containers, err := c.app.docker.ContainerList(context.Background(), types.ContainerListOptions{
		All: c.viewAll,
	})
	if err != nil {
		panic(err)
	}

	c.content = containers

	for i, name := range []string{"CONTAINER ID", "NAMES", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS"} {
		c.SetCell(0, i, tview.NewTableCell(name).SetSelectable(false))
	}

	for index, container := range c.content {
		c.SetCell(index+1, 0, tview.NewTableCell(container.ID[7:19]))
		c.SetCell(index+1, 1, tview.NewTableCell(container.Names[0][1:]))
		c.SetCell(index+1, 2, tview.NewTableCell(container.Image))

		if len(container.Command) > 19 {
			c.SetCell(index+1, 3, tview.NewTableCell(strconv.Quote(fmt.Sprintf("%s…", container.Command[:19]))))
		} else {
			c.SetCell(index+1, 3, tview.NewTableCell(strconv.Quote(container.Command)))
		}

		createdAt := time.Unix(container.Created, 0)
		createdDuration := time.Since(createdAt).Round(time.Second)
		c.SetCell(index+1, 4, tview.NewTableCell(createdDuration.String()+" ago"))

		c.SetCell(index+1, 5, tview.NewTableCell(container.Status))

		var ports []string
		for _, port := range container.Ports {
			ports = append(ports, c.formatPort(port))
		}
		c.SetCell(index+1, 6, tview.NewTableCell(strings.Join(ports, ", ")))
	}

	c.SetTitle(fmt.Sprintf(" %s [%d] ", "Container", c.GetRowCount()))
}

func (c *Container) bindKeys() {
	c.app.resetKeys()

	c.app.keybordHandlers["a"] = NewAction("View All", c.toggleViewAll)
	c.app.keybordHandlers["d"] = NewAction("Delete Container", c.delete)
	c.app.keybordHandlers["k"] = NewAction("Kill Container", c.kill)
	c.app.keybordHandlers["s"] = NewAction("Start Container", c.start)
}

func (c *Container) toggleViewAll() error {
	c.viewAll = !c.viewAll
	c.reloadData()
	return nil
}

func (c *Container) delete() error {
	selectedRow, _ := c.Table.GetSelection()
	selectedContainer := c.content[selectedRow-1]
	err := c.app.docker.ContainerRemove(context.Background(), selectedContainer.ID, types.ContainerRemoveOptions{})
	if err != nil {
		c.app.logger.Println(err)
	}
	c.reloadData()
	return nil
}

func (c *Container) kill() error {
	selectedRow, _ := c.Table.GetSelection()
	selectedContainer := c.content[selectedRow-1]

	go func(container types.Container) {
		err := c.app.docker.ContainerKill(context.Background(), container.ID, "")
		if err != nil {
			c.app.logger.Println(err)
		}
	}(selectedContainer)

	return nil
}

func (c *Container) start() error {
	selectedRow, _ := c.Table.GetSelection()
	selectedContainer := c.content[selectedRow-1]

	go func(container types.Container) {
		err := c.app.docker.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})
		if err != nil {
			c.app.logger.Println(err)
		}
	}(selectedContainer)

	return nil
}

func (c *Container) formatPort(p types.Port) string {
	var s string
	if p.IP != "" {
		s = s + p.IP + ":"
	}

	if p.PublicPort > 0 {
		s = s + fmt.Sprint(p.PrivatePort, "->", p.PublicPort)
	} else {
		s = s + fmt.Sprint(p.PrivatePort)
	}
	return s + "/" + p.Type
}
