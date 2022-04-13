package app

import (
	"lazydocker/views"

	ui "github.com/gizak/termui/v3"
)

func initKeys(a *App) {
	a.keys.Text = "Press q to quit, Press j or k to switch container"
	a.keys.SetRect(0, views.TerminalHeight-1, views.TerminalWidth, views.TerminalHeight)
	a.keys.Border = false
	a.keys.TextStyle.Bg = ui.ColorBlue
}
