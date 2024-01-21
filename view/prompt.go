package view

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Prompt struct {
	app *App
	*tview.TextArea
}

func NewPrompt(app *App) *Prompt {
	return &Prompt{
		app:      app,
		TextArea: tview.NewTextArea(),
	}
}

func (p *Prompt) init() {
	p.TextArea.SetBorder(true)
	p.TextArea.SetFocusFunc(p.start)
}

func (p *Prompt) start() {
	p.TextArea.SetText("", false)
	p.app.main.ResizeItem(p, 3, 1)
	p.SetInputCapture(p.promptFunc)
}

func (p *Prompt) promptFunc(event *tcell.EventKey) *tcell.EventKey {
	text := p.GetText()

	if event.Key() == tcell.KeyBackspace2 && text == ":" {
		p.backToMain()
	}
	if event.Key() == tcell.KeyEsc {
		p.backToMain()
	} else if event.Key() == tcell.KeyEnter {
		p.app.main.ResizeItem(p, 0, 0)
		p.TextArea.SetInputCapture(nil)
		p.app.SetFocus(p.app.main)

		if text == ":q" {
			p.app.Stop()
		}

		cmd, _ := strings.CutPrefix(text, ":")
		p.app.execCommand(cmd)
	}
	return event
}

func (p *Prompt) backToMain() {
	p.app.main.ResizeItem(p, 0, 0)
	p.TextArea.SetInputCapture(nil)
	p.app.SetFocus(p.app.main)
}
