package views

import (
	"image"
	"time"

	ui "github.com/gizak/termui/v3"
)

type Input struct {
	*ui.Block
	Text      string
	TextStyle ui.Style
	WrapText  bool
	Stdin     []rune
	Length    int
	cursor    bool
	close     chan struct{}
}

func NewInput() *Input {
	return &Input{
		Block:     ui.NewBlock(),
		TextStyle: ui.Theme.Paragraph.Text,
		WrapText:  true,
		Stdin:     make([]rune, 0),
		Length:    20,
		cursor:    true,
		close:     make(chan struct{}),
	}
}

func (i *Input) ListenKeyboard(event <-chan ui.Event) {
	i.Stdin = make([]rune, 0, i.Length+1)
	for {
		e := <-event
		if e.ID == "<Enter>" {
			i.Close()
			return
		}
		if e.ID == "<Backspace>" {
			if len(i.Stdin) == 0 {
				continue
			}
			i.Stdin = i.Stdin[:len(i.Stdin)-1]
			i.Text = i.appendCursor()
			ui.Render(i)
			continue
		}
		if e.ID[0] == '<' || len(i.Stdin) > i.Length {
			continue
		}
		i.Stdin = append(i.Stdin, []rune(e.ID)[0])
		i.Text = i.appendCursor()
		ui.Render(i)
	}
}

func (i *Input) Draw(buf *ui.Buffer) {
	i.Block.Draw(buf)

	cells := ui.ParseStyles(i.Text, i.TextStyle)
	if i.WrapText {
		cells = ui.WrapCells(cells, uint(i.Inner.Dx()))
	}

	rows := ui.SplitCells(cells, '\n')

	for y, row := range rows {
		if y+i.Inner.Min.Y >= i.Inner.Max.Y {
			break
		}
		row = ui.TrimCells(row, i.Inner.Dx())
		for _, cx := range ui.BuildCellWithXArray(row) {
			x, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(i.Inner.Min))
		}
	}
}

func (i *Input) appendCursor() string {
	if i.cursor {
		return string(append(i.Stdin, '_'))
	}
	return string(append(i.Stdin, ' '))
}

func (i *Input) Scan() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	i.cursor = true
	for {
		select {
		case <-ticker.C:
			i.Text = i.appendCursor()
			i.cursor = !i.cursor
			ui.Render(i)
		case <-i.close:
			return
		}
	}
}

func (i *Input) Close() {
	i.close <- struct{}{}
}
