package cells

import (
	"image"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Table struct {
	*widgets.Table
	ActiveRowIndex   int
	ActiveRowStyle   ui.Style
	InactiveRowStyle ui.Style
	RowTab           map[int]ui.Drawable
}

func NewTable() *Table {
	return &Table{
		Table:            widgets.NewTable(),
		ActiveRowIndex:   1,
		ActiveRowStyle:   ui.Theme.Tab.Active,
		InactiveRowStyle: ui.Theme.Tab.Inactive,
	}
}

func (t *Table) FocusDown() {
	if t.ActiveRowIndex < len(t.Rows)-1 {
		t.ActiveRowIndex++
	}
}

func (t *Table) FocusUp() {
	if t.ActiveRowIndex > 1 {
		t.ActiveRowIndex--
	}
}

func (t *Table) Draw(buf *ui.Buffer) {
	// t.Block.Draw(buf)
	t.Table.Draw(buf)
	t.renderTab()

	yCoordinate := t.Inner.Min.Y + 1
	for i, row := range t.Rows {
		if i == 0 {
			continue
		}
		ColorPair := t.InactiveRowStyle
		if i == t.ActiveRowIndex {
			ColorPair = t.ActiveRowStyle
		}
		// TODO fix TextAlignment offset
		buf.SetString(
			row[0],
			ColorPair,
			image.Pt(t.Inner.Min.X, yCoordinate),
		)

		yCoordinate++
	}
}

func (t *Table) renderTab() {
	if item, ok := t.RowTab[t.ActiveRowIndex]; ok {
		ui.Render(item)
	}
}
