package view

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/gdamore/tcell/v2"
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
	c.SetBorder(true).SetTitle(fmt.Sprintf(" %s [%d] ", "Container", c.GetRowCount()))

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
}

func (c *Container) bindKeys() {
	c.app.resetKeys()

	c.app.keybordHandlers["a"] = c.toggleViewAll
	c.app.keybordHandlers["d"] = c.delete
}

func (c *Container) toggleViewAll(ek *tcell.EventKey) *tcell.EventKey {
	c.viewAll = !c.viewAll
	c.reloadData()
	return nil
}

func (c *Container) delete(ek *tcell.EventKey) *tcell.EventKey {
	selectedRow, _ := c.Table.GetSelection()
	selectedContainer := c.content[selectedRow-1]
	err := c.app.docker.ContainerRemove(context.Background(), selectedContainer.ID, types.ContainerRemoveOptions{})
	if err != nil {
		c.app.logger.Println(err)
	}
	c.reloadData()
	return nil
}
