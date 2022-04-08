package views

import (
	"lazydocker/cells"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type View struct {
	containers *cells.Table
	images     *cells.Table
	keys       *widgets.Paragraph
	search     *cells.Input

	activeStyle ui.Style
	activeSort  []cells.Cell
	activeIndex int
}

func NewView() *View {
	return &View{
		containers:  cells.NewTable(),
		images:      cells.NewTable(),
		keys:        widgets.NewParagraph(),
		search:      cells.NewInput(),
		activeStyle: ui.NewStyle(51),
		activeSort:  make([]cells.Cell, 0),
	}
}

func (v *View) Init() *View {
	initKeys(v)
	initContainers(v)
	initSearch(v)
	initImages(v)
	v.activeSort = []cells.Cell{v.containers, v.images}
	v.containers.Active(v.activeStyle)
	return v
}

func (v *View) Render() {
	ui.Render(v.keys, v.containers, v.images)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Tab>":
			v.OnActive()
		case "<Resize>":
			v.OnResize(e.Payload.(ui.Resize))
		case "k", "<Up>":
			v.OnUp()
		case "j", "<Down>":
			v.OnDown()
		case "h", "<Left>":
		case "l", "<Right>":
		case "s":
			v.OnSwithStatus()
		case "/":
			v.OnSearch(uiEvents)
		}
	}
}

func (v *View) OnActive() {
	for i := range v.activeSort {
		v.activeSort[i].InActive()
	}
	v.activeIndex = (v.activeIndex + 1) % len(v.activeSort)
	v.activeSort[v.activeIndex].Active(v.activeStyle)
	v.ReRender()
}

func (v *View) OnResize(size ui.Resize) {
	cells.TerminalWidth = size.Width
	cells.TerminalHeight = size.Height
	v.containers.ResetSize(0, 0, 50, cells.TerminalHeight/2)
	v.images.ResetSize(0, cells.TerminalHeight/2, 50, cells.TerminalHeight-1)
	v.keys.SetRect(0, cells.TerminalHeight-1, cells.TerminalWidth, cells.TerminalHeight)
	v.ReRender()
}

func (v *View) OnUp() {
	v.activeSort[v.activeIndex].FocusUp()
	ui.Render(v.activeSort[v.activeIndex])
}

func (v *View) OnDown() {
	v.activeSort[v.activeIndex].FocusDown()
	ui.Render(v.activeSort[v.activeIndex])
}

func (v *View) OnSwithStatus() {
	if v.activeIndex == 0 {
		freshContainers(containerOption(), v.containers)
		v.ReRender()
	}
}

func (v *View) OnSearch(e <-chan ui.Event) {
	go v.search.Scan()
	v.search.ListenKeyboard(e)
	if v.activeIndex == 0 {
		freshContainers(string(v.search.Stdin), v.containers)
	} else {
		freshImages(string(v.search.Stdin), v.images)
	}
	v.ReRender()
}

func (v *View) ReRender() {
	ui.Clear()
	ui.Render(v.keys, v.containers, v.images)
}
