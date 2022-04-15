package views

import (
	"log"

	ui "github.com/gizak/termui/v3"
)

var (
	TerminalWidth  int
	TerminalHeight int
)

func InitTerminal() {
	TerminalWidth, TerminalHeight = ui.TerminalDimensions()
	if TerminalHeight < 10 || TerminalWidth < 50 {
		log.Panicln("no space to render")
	}
}
