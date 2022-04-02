package views

import (
	"lazydocker/cells"

	ui "github.com/gizak/termui/v3"

	"github.com/gizak/termui/v3/widgets"
)

func Keys() *widgets.Paragraph {
	header := widgets.NewParagraph()
	header.Text = "Press q to quit, Press j or k to switch container"
	header.SetRect(0, cells.TerminalHeight-1, cells.TerminalWidth, cells.TerminalHeight)
	header.Border = false
	header.TextStyle.Bg = ui.ColorBlue
	return header
}
