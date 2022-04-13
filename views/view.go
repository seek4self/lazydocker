package views

import ui "github.com/gizak/termui/v3"

type View interface {
	ui.Drawable
	Active(ui.Style)
	InActive()
	FocusUp()
	FocusDown()
	ActiveText() string
}
