package main

import (
	"lazydocker/views"
	"log"

	ui "github.com/gizak/termui/v3"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	keys := views.Keys()
	container := views.ContainerStatus()
	ui.Render(keys, container)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "k", "<Up>":
			container.FocusUp()
			ui.Clear()
			ui.Render(keys, container)
		case "j", "<Down>":
			container.FocusDown()
			ui.Clear()
			ui.Render(keys, container)
		case "h", "<Left>":
			container.PrePage()
			ui.Clear()
			ui.Render(keys, container)
		case "l", "<Right>":
			container.NextPage()
			ui.Clear()
			ui.Render(keys, container)
		}
	}
}
