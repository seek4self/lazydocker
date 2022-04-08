package cells

import (
	"image"

	ui "github.com/gizak/termui/v3"
)

type Navigation struct {
	*ui.Block
	Text      string
	TextStyle ui.Style
	WrapText  bool

	Header         []string
	Target         string
	ContentHandler map[string]func(string) []byte
	Active         bool
	ActiveCol      int
	ActiveStyle    ui.Style

	x, y int
}

func NewNavigation() *Navigation {
	return &Navigation{
		Block:          ui.NewBlock(),
		TextStyle:      ui.Theme.Paragraph.Text,
		WrapText:       true,
		ContentHandler: make(map[string]func(string) []byte),

		ActiveStyle: ui.NewStyle(51),
	}
}

func (n *Navigation) Draw(buf *ui.Buffer) {
	n.Block.Draw(buf)

	cells := ui.ParseStyles(n.Text, n.TextStyle)
	if n.WrapText {
		cells = ui.WrapCells(cells, uint(n.Inner.Dx()))
	}

	rows := ui.SplitCells(cells, '\n')

	for y, row := range rows {
		if y+n.Inner.Min.Y >= n.Inner.Max.Y {
			break
		}
		row = ui.TrimCells(row, n.Inner.Dx())
		for _, cx := range ui.BuildCellWithXArray(row) {
			x, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(n.Inner.Min))
		}
	}
}

func (n *Navigation) drawHeader(buf *ui.Buffer) {
	n.x = n.Inner.Min.X
	n.y = n.Inner.Min.Y
	for i := 0; i < len(n.Header); i++ {
		col := ui.ParseStyles(n.Header[i], n.TextStyle)
		for _, cx := range ui.BuildCellWithXArray(col) {
			if cx.X == n.Inner.Dx() || n.x+cx.X == n.Inner.Max.X {
				cx.Cell.Rune = ui.ELLIPSES
				buf.SetCell(cx.Cell, image.Pt(n.x+cx.X-1, n.y))
				break
			}
			buf.SetCell(cx.Cell, image.Pt(n.x+cx.X, n.y))
		}
		n.x += len(n.Header[i]) + 5
	}
}
