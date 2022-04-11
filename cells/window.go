package cells

import (
	ui "github.com/gizak/termui/v3"
)

var (
	TerminalWidth  int
	TerminalHeight int
)

func InitTerminal() {
	TerminalWidth, TerminalHeight = ui.TerminalDimensions()
}
