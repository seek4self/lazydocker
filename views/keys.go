package views

import (
	"lazydocker/cells"

	ui "github.com/gizak/termui/v3"
)

func initKeys(v *View) {
	v.keys.Text = "Press q to quit, Press j or k to switch container"
	v.keys.SetRect(0, cells.TerminalHeight-1, cells.TerminalWidth, cells.TerminalHeight)
	v.keys.Border = false
	v.keys.TextStyle.Bg = ui.ColorBlue
}
