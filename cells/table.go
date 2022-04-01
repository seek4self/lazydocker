package cells

import (
	"image"
	"math"

	ui "github.com/gizak/termui/v3"
)

type Table struct {
	// *widgets.Table
	*ui.Block
	Header          []string
	Rows            [][]string
	ColumnWidths    []int
	TextStyle       ui.Style
	RowSeparator    bool
	ColumnSeparator bool
	ColumnAlignment []ui.Alignment
	RowStyles       map[int]ui.Style
	FillRow         bool

	// ColumnResizer is called on each Draw. Can be used for custom column sizing.
	ColumnResizer    func()
	ActiveRowIndex   int
	ActiveRowStyle   ui.Style
	InactiveRowStyle ui.Style
	RowTab           map[int]ui.Drawable

	Page int

	textAlignment ui.Alignment
	x, y          int // drawing coordinate
}

func NewTable() *Table {
	return &Table{
		Block:            ui.NewBlock(),
		TextStyle:        ui.Theme.Table.Text,
		RowStyles:        make(map[int]ui.Style),
		ColumnResizer:    func() {},
		ActiveRowIndex:   0,
		ActiveRowStyle:   ui.Theme.Tab.Active,
		InactiveRowStyle: ui.Theme.Tab.Inactive,
	}
}

func (t *Table) FocusDown() {
	if t.activeRow() < len(t.Rows)-1 &&
		t.ActiveRowIndex < t.rowsCount()-1 {
		t.ActiveRowIndex++
	}
}

func (t *Table) activeRow() int {
	return t.ActiveRowIndex + t.Page*t.rowsCount()
}

func (t *Table) FocusUp() {
	if t.ActiveRowIndex > 0 {
		t.ActiveRowIndex--
	}
}

func (t *Table) NextPage() {
	if t.Page < t.totalPage()-1 {
		t.Page++
		if t.activeRow() >= len(t.Rows) {
			t.ActiveRowIndex = 0
		}
	}
}

func (t *Table) PrePage() {
	if t.Page > 0 {
		t.Page--
	}
}

func (t *Table) Draw(buf *ui.Buffer) {
	t.Block.Draw(buf)
	t.drawTable(buf)
	t.drawTabPane(buf)
	t.drawTabPage()
}

func (t *Table) height() int {
	return t.Inner.Dy()
}

func (t *Table) width() int {
	return t.Inner.Dx()
}

func (t *Table) rowsCount() int {
	return t.Inner.Dy() - 2
}

func (t *Table) totalPage() int {
	return int(math.Ceil(float64(len(t.Rows)) / float64(t.rowsCount())))
}

// fix columns width
func (t *Table) fixWidths() {
	if len(t.ColumnWidths) > 0 {
		return
	}
	columnCount := len(t.Rows[0])
	t.ColumnWidths = make([]int, 0, columnCount)
	avgWidth := t.Inner.Dx() / columnCount
	for i := 0; i < columnCount; i++ {
		t.ColumnWidths = append(t.ColumnWidths, avgWidth)
	}
}

func (t *Table) drawHeader(buf *ui.Buffer) {
	t.drawRow(-1, buf)
	t.drawHorizontalSep(-1, buf)
}

func (t *Table) rowStyle(rowNum int, buf *ui.Buffer) ui.Style {
	rowStyle := t.TextStyle
	// get the row style if one exists
	if style, ok := t.RowStyles[rowNum]; ok {
		rowStyle = style
	}
	if t.FillRow {
		blankCell := ui.NewCell(' ', rowStyle)
		buf.Fill(blankCell, image.Rect(t.Inner.Min.X, t.y, t.Inner.Max.X, t.y+1))
	}
	return rowStyle
}

func (t *Table) columnAlignment(num int) {
	t.textAlignment = ui.AlignLeft
	if len(t.ColumnAlignment) > 0 && num < len(t.Header) {
		t.textAlignment = t.ColumnAlignment[num]
	}
}

// draw row cell
func (t *Table) drawRowCell(col []ui.Cell, width int, buf *ui.Buffer) {
	var offset int
	if len(col) > width || t.textAlignment == ui.AlignLeft {
		offset = t.x
	} else if t.textAlignment == ui.AlignCenter {
		offset = t.x + (width-len(col))/2
	} else if t.textAlignment == ui.AlignRight {
		offset = ui.MinInt(t.x+width, t.Inner.Max.X) - len(col)
	}
	for _, cx := range ui.BuildCellWithXArray(col) {
		if cx.X == width || offset+cx.X == t.Inner.Max.X {
			cx.Cell.Rune = ui.ELLIPSES
			buf.SetCell(cx.Cell, image.Pt(offset+cx.X-1, t.y))
			break
		}
		buf.SetCell(cx.Cell, image.Pt(offset+cx.X, t.y))
	}
}

func (t *Table) drawRow(rowNum int, buf *ui.Buffer) {
	var row []string
	if rowNum == -1 {
		row = t.Header
	} else {
		row = t.Rows[rowNum]
	}
	t.x = t.Inner.Min.X

	rowStyle := t.rowStyle(rowNum, buf)
	// draw row cells
	for i := 0; i < len(row); i++ {
		col := ui.ParseStyles(row[i], rowStyle)
		t.columnAlignment(i)
		t.drawRowCell(col, t.ColumnWidths[i], buf)
		t.x += t.ColumnWidths[i] + 1
	}

	t.drawVerticalSep(rowStyle, buf)
	t.y++
	if t.RowSeparator {
		t.drawHorizontalSep(rowNum, buf)
	}
}

// drawVerticalSep draw vertical separators
func (t *Table) drawVerticalSep(rowStyle ui.Style, buf *ui.Buffer) {
	sepStyle := t.Block.BorderStyle
	sep := ' '
	if t.ColumnSeparator {
		sep = ui.VERTICAL_LINE
	}
	sepX := t.Inner.Min.X
	verticalCell := ui.NewCell(sep, sepStyle)
	for i, width := range t.ColumnWidths {
		verticalCell.Style.Bg = sepStyle.Bg
		if t.FillRow && i < len(t.ColumnWidths)-1 {
			verticalCell.Style.Bg = rowStyle.Bg
		}

		sepX += width
		buf.SetCell(verticalCell, image.Pt(sepX, t.y))
		sepX++
	}
}

// drawHorizontalSep draw horizontal separators
func (t *Table) drawHorizontalSep(rowNum int, buf *ui.Buffer) {
	horizontalCell := ui.NewCell(ui.HORIZONTAL_LINE, t.Block.BorderStyle)
	if rowNum == -1 || (t.y < t.Inner.Max.Y && rowNum != len(t.Rows)-1) {
		buf.Fill(horizontalCell, image.Rect(t.Inner.Min.X, t.y, t.Inner.Max.X, t.y+1))
		t.y++
	}
}

func (t *Table) drawTable(buf *ui.Buffer) {
	t.ColumnResizer()
	t.fixWidths()

	t.y = t.Inner.Min.Y
	t.drawHeader(buf)

	// draw rows
	for i := t.Page * t.rowsCount(); i < len(t.Rows) && t.y < t.Inner.Max.Y; i++ {
		// fmt.Println("                                                     ", i, t.rowsCount())
		t.drawRow(i, buf)
	}
}

func (t *Table) drawTabPage() {
	if item, ok := t.RowTab[t.ActiveRowIndex]; ok {
		ui.Render(item)
	}
}

func (t *Table) drawTabPane(buf *ui.Buffer) {
	yCoordinate := t.Inner.Min.Y + t.ActiveRowIndex + 2
	t.columnAlignment(0)
	text := t.Rows[t.activeRow()][0]
	width := t.ColumnWidths[0]
	offset := 0
	if len(text) >= width || t.textAlignment == ui.AlignLeft {

	} else if t.textAlignment == ui.AlignCenter {
		offset = (width - len(text)) / 2
	} else {
		offset = width - len(text)
	}
	buf.SetString(
		ui.TrimString(text, width),
		t.ActiveRowStyle,
		image.Pt(t.Inner.Min.X+offset, yCoordinate),
	)
}