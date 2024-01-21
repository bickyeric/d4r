package view

import (
	"context"
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
	pages  *tview.Pages

	logger *Debugger
	prompt *Prompt

	main *tview.Flex
}

func NewApp() *App {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	app := &App{
		Application: tview.NewApplication(),
		docker:      cli,
		pages:       tview.NewPages(),
		logger:      NewDebugger(),
	}

	app.init()

	return app
}

func (a *App) init() {
	content := NewContainer(a.docker)
	content.init()

	a.pages.AddPage("Container", content, true, true)

	a.prompt = NewPrompt(a)
	a.prompt.init()

	a.main = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Header"), 10, 1, false).
		AddItem(a.prompt, 0, 0, false).
		AddItem(a.pages, 0, 1, false).
		AddItem(a.logger, 10, 1, false)

	a.main.SetFocusFunc(func() {
		a.main.SetInputCapture(a.mainFunc)
		a.Application.SetFocus(a.pages)
	})

	a.Application.SetRoot(a.main, true)
}

func (a *App) mainFunc(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == ':' {
		a.Application.SetFocus(a.prompt)
		a.main.SetInputCapture(nil)
	}

	a.logger.Println("from flex", event.Name())
	return event
}

func (a *App) execCommand(cmd string) {
	// TODO
	var primitive tview.Primitive

	if cmd == "image" {
		img := NewImage(a.docker)
		img.init()
		primitive = img
		a.pages.AddPage("Image", primitive, true, true)
	} else if cmd == "ps" || cmd == "container" {
		ps := NewContainer(a.docker)
		ps.init()
		primitive = ps
		a.pages.AddPage("Container", primitive, true, true)
	} else {
		panic("apanih")
	}
}

func (a *App) Logger() *log.Logger {
	return a.logger.Logger
}

func (a *App) Run(ctx context.Context) error {
	return a.Application.Run()
}

func (a *App) Stop() {
	a.docker.Close()
	a.Application.Stop()
}
