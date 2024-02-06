package view

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/client"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application
	docker *client.Client
	pages  *tview.Pages

	logger *Debugger
	prompt *Prompt

	main *tview.Flex

	keybordHandlers map[string]Action
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
		keybordHandlers: make(map[string]Action),
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
		if a.pages.GetPageCount() < 1 {
			content := NewContainer(a)
			content.init()
			a.pages.AddPage("Container", content, true, true)
		}
		a.main.SetInputCapture(a.keyboardHandler)
		a.Application.SetFocus(a.pages)
	})

	a.Application.SetRoot(a.main, true)
}

func (a *App) bindKeys() {
	a.keybordHandlers[":"] = NewAction("Active Prompt", a.activatePrompt)
	a.keybordHandlers["?"] = NewAction("Show Help", a.showKeyHelp)
	a.keybordHandlers["Esc"] = NewAction("Back", func() error {
		if a.pages.GetPageCount() > 1 {
			name, _ := a.pages.GetFrontPage()
			a.pages.RemovePage(name)
		}
		return nil
	})
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

func (a *App) activatePrompt() error {
	a.main.ResizeItem(a.prompt, 3, 1)
	a.Application.SetFocus(a.prompt)
	return nil
}

func (a *App) deactivatePrompt() error {
	a.main.ResizeItem(a.prompt, 0, 0)
	a.SetFocus(a.pages)
	return nil
}

func (a *App) showKeyHelp() error {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		flex := tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
		return flex
	}

	textArea := tview.NewTextView()
	textArea.SetBorder(true).
		SetTitle(" Help ")

	w := textArea.BatchWriter()
	for k, a := range a.keybordHandlers {
		fmt.Fprintln(w, k, a.name)
	}
	w.Close()

	a.pages.AddPage("modal", modal(textArea, 40, 10), true, true)
	return nil
}

func (a *App) keyboardHandler(event *tcell.EventKey) *tcell.EventKey {
	a.logger.Println("from flex", event.Name())

	if action, ok := a.keybordHandlers[string(event.Rune())]; ok {
		action.handler()
		return nil
	}

	if action, ok := a.keybordHandlers[event.Name()]; ok {
		action.handler()
		return nil
	}

	return event
}

func (a *App) execCommand(cmd string) {
	var primitive tview.Primitive
	if cmd == "ps" || cmd == "container" {
		ps := NewContainer(a)
		ps.init()
		primitive = ps
		a.pages.AddPage("Container", primitive, true, true)
	} else {
		a.Logger().Printf("command %s tidak ditemukan\n", cmd)
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
