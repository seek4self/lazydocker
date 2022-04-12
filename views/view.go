package views

import (
	"lazydocker/cells"
	"lazydocker/docker"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	containerInfo = "Container Info"
	containerLog  = "Log"
	imageInfo     = "Image Info"
	imageHistory  = "History"
)

var (
	navigationContainer = []string{containerInfo, containerLog, "memory"}
	navigationImage     = []string{imageInfo, imageHistory, "memory"}
)

type View struct {
	containers *cells.Table
	images     *cells.Table
	keys       *widgets.Paragraph
	navigation *cells.Navigation
	search     *cells.Input

	activeStyle ui.Style
	activeSort  []cells.Cell
	activeIndex int
}

func NewView() *View {
	cells.InitTerminal()
	if cells.TerminalHeight < 10 || cells.TerminalWidth < 50 {
		log.Panicln("no space to render")
	}
	return &View{
		containers:  cells.NewTable(),
		images:      cells.NewTable(),
		keys:        widgets.NewParagraph(),
		navigation:  cells.NewNavigation(),
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
	v.navigation.Header = navigationContainer
	v.navigation.ContentHandler = map[string]func(string) []byte{
		containerInfo: docker.ContainerInspect,
		containerLog:  docker.Logs,
		imageInfo:     docker.ImageInspect,
		imageHistory:  docker.History,
	}
	v.navigation.SetRect(v.containers.Inner.Max.X+1, 0, cells.TerminalWidth, cells.TerminalHeight-1)
	v.navigation.FreshContent(v.containers.ActiveText())
	v.activeSort = []cells.Cell{v.containers, v.images}
	v.containers.Active(v.activeStyle)
	return v
}

func (v *View) Render() {
	ui.Render(v.keys, v.containers, v.images, v.navigation)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Tab>":
			v.OnSwitch()
		case "<Resize>":
			v.OnResize(e.Payload.(ui.Resize))
		case "k", "<Up>":
			v.OnUp()
		case "j", "<Down>":
			v.OnDown()
		case "h", "<Left>":
			v.navigation.FocusLeft()
			ui.Render(v.navigation)
		case "l", "<Right>":
			v.navigation.FocusRight()
			ui.Render(v.navigation)
		case "<PageUp>":
			v.navigation.PageUp()
			ui.Render(v.navigation)
		case "<PageDown>":
			v.navigation.PageDown()
			ui.Render(v.navigation)
		case "s":
			v.OnSwithStatus()
		case "/":
			v.OnSearch(uiEvents)
		}
	}
}

func (v *View) OnSwitch() {
	for i := range v.activeSort {
		v.activeSort[i].InActive()
	}
	v.activeIndex = (v.activeIndex + 1) % len(v.activeSort)
	v.activeSort[v.activeIndex].Active(v.activeStyle)
	if v.activeIndex == 0 {
		v.navigation.Header = navigationContainer
	} else {
		v.navigation.Header = navigationImage
	}
	v.ReRender()
}

func (v *View) OnResize(size ui.Resize) {
	cells.TerminalWidth = size.Width
	cells.TerminalHeight = size.Height
	v.containers.ResetSize(0, 0, 50, cells.TerminalHeight/2)
	v.images.ResetSize(0, cells.TerminalHeight/2, 50, cells.TerminalHeight-1)
	v.keys.SetRect(0, cells.TerminalHeight-1, cells.TerminalWidth, cells.TerminalHeight)
	v.navigation.SetRect(v.containers.Inner.Max.X+1, 0, cells.TerminalWidth, cells.TerminalHeight-1)
	v.navigation.FreshContent(v.containers.ActiveText())
	v.ReRender()
}

func (v *View) OnUp() {
	v.activeSort[v.activeIndex].FocusUp()
	val := v.activeSort[v.activeIndex].ActiveText()
	v.navigation.FreshContent(val)
	ui.Render(v.activeSort[v.activeIndex], v.navigation)
}

func (v *View) OnDown() {
	v.activeSort[v.activeIndex].FocusDown()
	val := v.activeSort[v.activeIndex].ActiveText()
	v.navigation.FreshContent(val)
	ui.Render(v.activeSort[v.activeIndex], v.navigation)
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
	ui.Render(v.keys, v.containers, v.images, v.navigation)
}
