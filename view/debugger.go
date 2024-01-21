package view

import (
	"log"

	"github.com/rivo/tview"
)

type Debugger struct {
	*log.Logger
	*tview.TextView
}

func NewDebugger() *Debugger {
	tv := tview.NewTextView().ScrollToEnd()
	log := log.Default()
	log.SetOutput(tv)

	return &Debugger{Logger: log, TextView: tv}
}
