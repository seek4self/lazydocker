package app

import (
	"lazydocker/views"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type App struct {
	containers *views.Table
	images     *views.Table
	keys       *widgets.Paragraph
	navigation *views.Navigation
	search     *views.Input

	activeStyle ui.Style
	activeSort  []views.View
	activeIndex int
}

func NewApp() *App {
	views.InitTerminal()

	return &App{
		containers:  views.NewTable(),
		images:      views.NewTable(),
		keys:        widgets.NewParagraph(),
		navigation:  views.NewNavigation(),
		search:      views.NewInput(),
		activeStyle: ui.NewStyle(51),
		activeSort:  make([]views.View, 0),
	}
}

func (a *App) Init() *App {
	initKeys(a)
	initContainers(a)
	initSearch(a)
	initImages(a)
	initNavigation(a)
	a.activeSort = []views.View{a.containers, a.images}
	a.containers.Active(a.activeStyle)
	return a
}

func (a *App) Render() {
	ui.Render(a.keys, a.containers, a.images, a.navigation)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Tab>":
			a.OnSwitch()
		case "<Resize>":
			a.OnResize(e.Payload.(ui.Resize))
		case "k", "<Up>":
			a.OnUp()
		case "j", "<Down>":
			a.OnDown()
		case "h", "<Left>":
			a.navigation.FocusLeft()
			ui.Render(a.navigation)
		case "l", "<Right>":
			a.navigation.FocusRight()
			ui.Render(a.navigation)
		case "<PageUp>":
			a.navigation.PageUp()
			ui.Render(a.navigation)
		case "<PageDown>":
			a.navigation.PageDown()
			ui.Render(a.navigation)
		case "s":
			a.OnSwithStatus()
		case "/":
			a.OnSearch(uiEvents)
		}
	}
}

func (a *App) OnSwitch() {
	for i := range a.activeSort {
		a.activeSort[i].InActive()
	}
	a.activeIndex = (a.activeIndex + 1) % len(a.activeSort)
	a.activeSort[a.activeIndex].Active(a.activeStyle)
	a.navigation.Fresh(a.activeSort[a.activeIndex].ActiveText(), a.activeIndex)
	a.ReRender()
}

func (a *App) OnResize(size ui.Resize) {
	views.TerminalWidth = size.Width
	views.TerminalHeight = size.Height
	a.containers.ResetSize(0, 0, 50, views.TerminalHeight/2)
	a.images.ResetSize(0, views.TerminalHeight/2, 50, views.TerminalHeight-1)
	a.keys.SetRect(0, views.TerminalHeight-1, views.TerminalWidth, views.TerminalHeight)
	a.navigation.SetRect(a.containers.Inner.Max.X+1, 0, views.TerminalWidth, views.TerminalHeight-1)
	// a.navigation.FreshContent(a.containers.ActiveText(), a.activeIndex)
	a.ReRender()
}

func (a *App) OnUp() {
	a.activeSort[a.activeIndex].FocusUp()
	a.navigation.Fresh(a.activeSort[a.activeIndex].ActiveText(), a.activeIndex)
	ui.Render(a.activeSort[a.activeIndex], a.navigation)
}

func (a *App) OnDown() {
	a.activeSort[a.activeIndex].FocusDown()
	a.navigation.Fresh(a.activeSort[a.activeIndex].ActiveText(), a.activeIndex)
	ui.Render(a.activeSort[a.activeIndex], a.navigation)
}

func (a *App) OnSwithStatus() {
	if a.activeIndex == 0 {
		freshContainers(containerOption(), a.containers)
		a.navigation.Fresh(a.containers.ActiveText(), a.activeIndex)
		a.ReRender()
	}
}

func (a *App) OnSearch(e <-chan ui.Event) {
	go a.search.Scan()
	a.search.ListenKeyboard(e)
	if a.activeIndex == 0 {
		freshContainers(string(a.search.Stdin), a.containers)
	} else {
		freshImages(string(a.search.Stdin), a.images)
	}
	a.ReRender()
}

func (a *App) ReRender() {
	ui.Clear()
	ui.Render(a.keys, a.containers, a.images, a.navigation)
}
