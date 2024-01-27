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
	p.TextArea.SetText(":", true)
	p.bindKeys()
}

func (p *Prompt) bindKeys() {
	p.app.clearKeys()

	p.app.keybordHandlers["Enter"] = p.execute
	p.app.keybordHandlers["Esc"] = p.app.deactivatePrompt
	p.app.keybordHandlers["Backspace2"] = func(ek *tcell.EventKey) *tcell.EventKey {
		if p.GetText() == ":" {
			return p.app.deactivatePrompt(ek)
		}
		return ek
	}
}

func (p *Prompt) execute(_ *tcell.EventKey) *tcell.EventKey {
	text := p.GetText()
	p.app.deactivatePrompt(nil)

	if text == ":q" {
		p.app.Stop()
	}

	cmd, _ := strings.CutPrefix(text, ":")
	p.app.execCommand(cmd)
	return nil
}
