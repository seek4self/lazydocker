package views

import (
	ui "github.com/gizak/termui/v3"

	"github.com/gizak/termui/v3/widgets"
)

func Keys() *widgets.Paragraph {
	header := widgets.NewParagraph()
	header.Text = "Press q to quit, Press j or k to switch tabs"
	header.SetRect(0, 0, 50, 1)
	header.Border = false
	header.TextStyle.Bg = ui.ColorBlue
	return header
}
