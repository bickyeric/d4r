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
	pages  *tview.Pages

	logger *Debugger
	prompt *Prompt

	main *tview.Flex

	keybordHandlers map[string]func(*tcell.EventKey) *tcell.EventKey
}

func NewApp() *App {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	app := &App{
		Application:     tview.NewApplication(),
		docker:          cli,
		pages:           tview.NewPages(),
		logger:          NewDebugger(),
		keybordHandlers: make(map[string]func(*tcell.EventKey) *tcell.EventKey),
	}

	app.init()

	return app
}

func (a *App) init() {
	a.prompt = NewPrompt(a)
	a.prompt.init()

	a.main = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Header"), 10, 1, false).
		AddItem(a.prompt, 0, 0, false).
		AddItem(a.pages, 0, 1, false).
		AddItem(a.logger, 10, 1, false)

	a.main.SetFocusFunc(func() {
		a.logger.Println(a.pages.GetPageCount())
		a.resetKeys()
		if a.pages.GetPageCount() < 1 {
			content := NewContainer(a.docker)
			content.init()
			a.pages.AddPage("Container", content, true, true)
		}
		a.main.SetInputCapture(a.keyboardHandler)
		a.Application.SetFocus(a.pages)
	})

	a.Application.SetRoot(a.main, true)
}

func (a *App) bindKeys() {
	a.keybordHandlers[":"] = a.activatePrompt
}

func (a *App) clearKeys() {
	for k := range a.keybordHandlers {
		delete(a.keybordHandlers, k)
	}
}

func (a *App) resetKeys() {
	a.clearKeys()
	a.bindKeys()
}

func (a *App) activatePrompt(_ *tcell.EventKey) *tcell.EventKey {
	a.main.ResizeItem(a.prompt, 3, 1)
	a.Application.SetFocus(a.prompt)
	return nil
}

func (a *App) deactivatePrompt(_ *tcell.EventKey) *tcell.EventKey {
	a.main.ResizeItem(a.prompt, 0, 0)
	a.SetFocus(a.main)
	return nil
}

func (a *App) keyboardHandler(event *tcell.EventKey) *tcell.EventKey {
	a.logger.Println("from flex", event.Name())

	if fun, ok := a.keybordHandlers[string(event.Rune())]; ok {
		return fun(event)
	}

	if fun, ok := a.keybordHandlers[event.Name()]; ok {
		return fun(event)
	}

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
		modal := tview.NewModal()
		modal.SetText(fmt.Sprintf("command %s tidak ditemukan", cmd))
		a.pages.AddPage("Error", modal, true, true)
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
