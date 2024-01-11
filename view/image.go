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
	client  *client.Client
	content []types.ImageSummary
}

func NewImage(cli *client.Client) *Image {
	summaries, err := cli.ImageList(context.Background(), types.ImageListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}
	return &Image{
		client:  cli,
		content: summaries,
	}
}

func (i *Image) Title() (name string, count int) {
	return "Image", len(i.content)
}

func (i *Image) Headers() []string {
	return []string{"REPOSITORY", "TAG", "IMAGE ID", "CREATED", "SIZE"}
}

func (i *Image) SetTable(t *tview.Table) {
	for index, img := range i.content {
		tags := strings.Split(img.RepoTags[0], ":")
		t.SetCell(index+1, 0, tview.NewTableCell(tags[0]))
		t.SetCell(index+1, 1, tview.NewTableCell(tags[1]))
		t.SetCell(index+1, 2, tview.NewTableCell(img.ID[7:19]))
		t.SetCell(index+1, 3, tview.NewTableCell(fmt.Sprint(time.Unix(img.Created, 0))))
		t.SetCell(index+1, 4, tview.NewTableCell(fmt.Sprint(img.Size)))
	}
}
