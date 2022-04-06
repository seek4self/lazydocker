package views

import (
	"lazydocker/cells"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var options = []string{"up", "all", "exit", ""}

type View struct {
	containers      *cells.Table
	images          *cells.Table
	keys            *widgets.Paragraph
	search          *cells.Input
	containerStatus int
}

func NewView() *View {
	return &View{
		containers: cells.NewTable(),
		images:     cells.NewTable(),
		keys:       widgets.NewParagraph(),
		search:     cells.NewInput(),
	}
}

func (v *View) Init() *View {
	initKeys(v)
	initContainers(v)
	initSearch(v)
	return v
}

func (v *View) Render() {
	ui.Render(v.keys, v.containers)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
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

func (v *View) OnResize(size ui.Resize) {
	cells.TerminalWidth = size.Width
	cells.TerminalHeight = size.Height
	v.containers.ResetSize(0, 0, 40, cells.TerminalHeight/2)
	v.keys.SetRect(0, cells.TerminalHeight-1, cells.TerminalWidth, cells.TerminalHeight)
	v.ReRender()
}

func (v *View) OnUp() {
	v.containers.FocusUp()
	v.ReRender()
}

func (v *View) OnDown() {
	v.containers.FocusDown()
	v.ReRender()
}

func (v *View) OnSwithStatus() {
	v.containerStatus = (v.containerStatus + 1) % len(options)
	if options[v.containerStatus] == "" {
		v.containerStatus = (v.containerStatus + 1) % len(options)
	}
	setTableRows(options[v.containerStatus], v.containers)
	v.ReRender()
}

func (v *View) OnSearch(e <-chan ui.Event) {
	go v.search.Scan()
	v.search.ListenKeyboard(e)
	setTableRows(string(v.search.Stdin), v.containers)
	v.ReRender()
}

func (v *View) ReRender() {
	ui.Clear()
	ui.Render(v.keys, v.containers)
}
