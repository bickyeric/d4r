package view

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/client"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Page interface {
	Title() (name string, count int)
	Headers() []string
	SetTable(t *tview.Table)
}

type App struct {
	*tview.Application
	docker *client.Client
	page   Page
}

func NewApp() *App {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	return &App{
		Application: tview.NewApplication(),
		docker:      cli,
		page:        NewContainer(cli),
	}
}

func (a *App) Run(ctx context.Context) error {
	table := tview.NewTable()

	name, count := a.page.Title()
	table.SetBorder(true).SetTitle(fmt.Sprintf(" %s [%d] ", name, count))

	for i, name := range a.page.Headers() {
		table.SetCell(0, i, tview.NewTableCell(name).SetSelectable(false))
	}

	a.page.SetTable(table)
	table.SetSelectable(true, false)

	footer := tview.NewTextView().ScrollToEnd()

	l := log.Default()
	l.SetOutput(footer)

	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			a.Stop()
		}
	})

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	grid := tview.NewGrid().
		SetRows(3, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(footer, 2, 0, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(table, 1, 0, 1, 3, 0, 0, false)
	a.Application.SetRoot(grid, true).SetFocus(table)
	return a.Application.Run()
}

func (a *App) Stop() {
	a.docker.Close()
	a.Application.Stop()
}
