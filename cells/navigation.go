package cells

import (
	"fmt"
	"image"

	ui "github.com/gizak/termui/v3"
)

type Navigation struct {
	*ui.Block
	Text      string
	TextStyle ui.Style
	WrapText  bool
	Rows      [][]ui.Cell

	Header         []string
	Target         string
	ContentHandler map[string]func(string) []byte
	Active         bool
	ActiveCol      int
	ActiveStyle    ui.Style

	input      string
	x, y       int
	offset     int
	TextHeight int
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
	n.drawHeader(buf)
	n.drawProgress(buf)

	for i := 0; i < n.visibleRows() && (i+n.offset) < n.totalRows(); i++ {
		y := i + n.y
		row := ui.TrimCells(n.Rows[i+n.offset], n.Inner.Dx())
		for _, cx := range ui.BuildCellWithXArray(row) {
			x, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(n.Inner.Min))
		}
	}
}

func (n *Navigation) parseText(text []byte) {
	cells := ui.ParseStyles(string(text), n.TextStyle)
	if n.WrapText {
		cells = ui.WrapCells(cells, uint(n.Inner.Dx()))
	}
	n.Rows = ui.SplitCells(cells, '\n')
}

func (n *Navigation) FreshContent(input string) {
	n.input = input
	n.getContent()
}

func (n *Navigation) visibleRows() int {
	return n.Inner.Dy() - 2
}

func (n *Navigation) totalRows() int {
	return len(n.Rows)
}

func (n *Navigation) PageUp() {
	if n.offset-n.visibleRows() <= 0 {
		n.offset = 0
		return
	}
	n.offset -= n.visibleRows() - 3
}

func (n *Navigation) PageDown() {
	if n.offset+n.visibleRows() > n.totalRows() {
		return
	}
	n.offset += n.visibleRows() - 3
}

func (n *Navigation) FocusRight() {
	if n.ActiveCol+1 == len(n.Header) {
		return
	}
	n.ActiveCol++
	n.getContent()
}

func (n *Navigation) getContent() {
	key := n.Header[n.ActiveCol]
	if getContent, ok := n.ContentHandler[key]; ok {
		n.parseText(getContent(n.input))
	}
	n.offset = 0
}

func (n *Navigation) FocusLeft() {
	if n.ActiveCol == 0 {
		return
	}
	n.ActiveCol--
	n.getContent()
}

func (n *Navigation) drawHeader(buf *ui.Buffer) {
	n.x = n.Inner.Min.X
	n.y = n.Inner.Min.Y
	for i := 0; i < len(n.Header); i++ {
		style := n.TextStyle
		if i == n.ActiveCol {
			style = n.ActiveStyle
		}
		col := ui.ParseStyles(n.Header[i], style)
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
	n.y++
	n.x = n.Inner.Min.X
	horizontalCell := ui.NewCell(ui.HORIZONTAL_LINE, n.Block.BorderStyle)
	buf.Fill(horizontalCell, image.Rect(n.Inner.Min.X, n.y, n.Inner.Max.X, n.y+1))
}

func (n *Navigation) progress() string {
	if n.offset+n.visibleRows() >= n.totalRows() {
		return "Bot"
	}
	if n.offset == 0 {
		return "Top"
	}
	progress := float64(n.offset+1) / float64(n.totalRows()) * 100
	if progress == 100 {
		return "Bot"
	}
	return fmt.Sprintf("%.0f%%", progress)
}

func (n *Navigation) drawProgress(buf *ui.Buffer) {
	col := ui.ParseStyles(n.progress(), n.TextStyle)
	for _, cx := range ui.BuildCellWithXArray(col) {
		buf.SetCell(cx.Cell, image.Pt(n.Max.X-5+cx.X, n.Inner.Max.Y))
	}
}
