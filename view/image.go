package view

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/rivo/tview"
)

type Image struct {
	*tview.Table
	client  *client.Client
	content []types.ImageSummary
}

func NewImage(cli *client.Client) *Image {
	return &Image{
		Table:  tview.NewTable(),
		client: cli,
	}
}

func (i *Image) init() {
	summaries, err := i.client.ImageList(context.Background(), types.ImageListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	i.content = summaries

	for index, name := range []string{"REPOSITORY", "TAG", "IMAGE ID", "CREATED", "SIZE"} {
		i.SetCell(0, index, tview.NewTableCell(name).SetSelectable(false))
	}

	for index, img := range i.content {
		tags := strings.Split(img.RepoTags[0], ":")
		i.SetCell(index+1, 0, tview.NewTableCell(tags[0]))
		i.SetCell(index+1, 1, tview.NewTableCell(tags[1]))
		i.SetCell(index+1, 2, tview.NewTableCell(img.ID[7:19]))
		i.SetCell(index+1, 3, tview.NewTableCell(fmt.Sprint(time.Unix(img.Created, 0))))
		i.SetCell(index+1, 4, tview.NewTableCell(fmt.Sprint(img.Size)))
	}
	i.SetSelectable(true, false)
	i.SetBorder(true).SetTitle(fmt.Sprintf(" %s [%d] ", "Image", i.GetRowCount()))
}
