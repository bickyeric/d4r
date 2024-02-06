package view

import (
	"strings"

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
	p.TextArea.SetText(":", true)
	p.bindKeys()
}

func (p *Prompt) bindKeys() {
	p.app.clearKeys()

	p.app.keybordHandlers["Enter"] = NewAction("Execute", p.execute)
	p.app.keybordHandlers["Esc"] = NewAction("Deactivate Prompt", p.app.deactivatePrompt)
}

func (p *Prompt) execute() error {
	text := p.GetText()
	p.app.deactivatePrompt()

	if text == ":q" {
		p.app.Stop()
	}

	cmd, _ := strings.CutPrefix(text, ":")
	p.app.execCommand(cmd)
	return nil
}
