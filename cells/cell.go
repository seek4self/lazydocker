package cells

import ui "github.com/gizak/termui/v3"

type Cell interface {
	ui.Drawable
	Active(ui.Style)
	InActive()
	FocusUp()
	FocusDown()
}
