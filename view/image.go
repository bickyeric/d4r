package view

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/rivo/tview"
)

type Image struct {
	*tview.Table
	app     *App
	content []types.ImageSummary
	viewAll bool
}

func NewImage(app *App) *Image {
	return &Image{
		Table: tview.NewTable(),
		app:   app,
	}
}

func (c *Image) init() {
	c.reloadData()
	c.SetSelectable(true, false)
	c.SetBorder(true)

	c.Table.SetFocusFunc(func() {
		c.bindKeys()
	})
}

func (c *Image) reloadData() {
	c.Clear()

	images, err := c.app.docker.ImageList(context.Background(), types.ImageListOptions{
		All: c.viewAll,
	})
	if err != nil {
		panic(err)
	}

	c.content = images

	for i, name := range []string{"REPOSITORY", "TAG", "IMAGE ID", "CREATED", "SIZE"} {
		c.SetCell(0, i, tview.NewTableCell(name).SetSelectable(false))
	}

	for index, image := range c.content {
		repoName := strings.Split(image.RepoTags[0], ":")
		c.SetCell(index+1, 0, tview.NewTableCell(repoName[0]))

		c.SetCell(index+1, 1, tview.NewTableCell(repoName[1]))
		fullID := image.ID
		longID := strings.Split(fullID, ":")[1]
		c.SetCell(index+1, 2, tview.NewTableCell(longID[0:12]))

		c.SetCell(index+1, 3, tview.NewTableCell(c.formatAge(image.Created)))

		c.SetCell(index+1, 4, tview.NewTableCell(fmt.Sprint(c.formatSize(image.Size))))
	}

	c.SetTitle(fmt.Sprintf(" %s [%d] ", "Image", c.GetRowCount()))
}

func (c *Image) bindKeys() {
	c.app.resetKeys()

	c.app.keybordHandlers["a"] = NewAction("View All", c.toggleViewAll)
	c.app.keybordHandlers["D"] = NewAction("Delete Container", c.delete)
}

func (c *Image) toggleViewAll() error {
	c.viewAll = !c.viewAll
	c.reloadData()
	return nil
}

func (c *Image) delete() error {
	selectedRow, _ := c.Table.GetSelection()
	selectedContainer := c.content[selectedRow-1]
	err := c.app.docker.ContainerRemove(context.Background(), selectedContainer.ID, types.ContainerRemoveOptions{})
	if err != nil {
		c.app.logger.Println(err)
	}
	c.reloadData()
	return nil
}

func (c *Image) formatSize(sizeInBytes int64) string {
	kilo := float64(1024)
	mega := kilo * 1024
	giga := mega * 1024

	if sizeInBytes > int64(giga) {
		return fmt.Sprintf("%2.2fGB", float64(sizeInBytes)/giga)
	}

	if sizeInBytes > int64(mega) {
		return fmt.Sprintf("%2.2fMB", float64(sizeInBytes)/mega)
	}

	if sizeInBytes > int64(kilo) {
		return fmt.Sprintf("%2.2fKB", float64(sizeInBytes)/kilo)
	}

	return fmt.Sprint(sizeInBytes, "B")
}

func (c *Image) formatAge(createdAt int64) string {
	yearInHours := float64(12 * 24 * 30)
	monthInHours := float64(24 * 30)
	weekInHours := float64(24 * 7)
	dayInHours := float64(24)

	createdTime := time.Unix(createdAt, 0)
	createdDuration := time.Since(createdTime)

	if createdDuration.Hours() > yearInHours {
		return fmt.Sprintf("%2.f years ago", createdDuration.Hours()/yearInHours)
	} else if createdDuration.Hours() > monthInHours {
		return fmt.Sprintf("%2.f months ago", createdDuration.Hours()/monthInHours)
	} else if createdDuration.Hours() > weekInHours {
		return fmt.Sprintf("%2.f weeks ago", createdDuration.Hours()/weekInHours)
	} else if createdDuration.Hours() > dayInHours {
		return fmt.Sprintf("%2.f days ago", createdDuration.Hours()/dayInHours)
	} else if createdDuration.Hours() > 0 {
		return fmt.Sprintf("%2.f hours ago", createdDuration.Hours())
	}

	return createdDuration.String() + " ago"
}
