package main

import (
	"lazydocker/cells"
	"lazydocker/views"
	"log"

	ui "github.com/gizak/termui/v3"
)

var (
	options     = []string{"up", "all", "exit", ""}
	statusTimes = 0
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	keys := views.Keys()
	container := views.ContainerStatus("up")
	ui.Render(keys, container)
	input := cells.NewInput()
	input.SetRect(40, 10, 80, 13)
	input.Title = "find container with:"

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Resize>":
			size := e.Payload.(ui.Resize)
			cells.TerminalWidth = size.Width
			cells.TerminalHeight = size.Height
			ui.Clear()
			ui.Render(keys, container)
		case "k", "<Up>":
			container.FocusUp()
			ui.Clear()
			ui.Render(keys, container)
		case "j", "<Down>":
			container.FocusDown()
			ui.Clear()
			ui.Render(keys, container)
		case "h", "<Left>":
			// container.PrePage()
			// ui.Clear()
			// ui.Render(keys, container)
		case "l", "<Right>":
			// container.NextPage()
			// ui.Clear()
			// ui.Render(keys, container)
		case "s":
			statusTimes = (statusTimes + 1) % len(options)
			if options[statusTimes] == "" {
				statusTimes = (statusTimes + 1) % len(options)
			}
			views.UpdateContainers(options[statusTimes], container)
			ui.Clear()
			ui.Render(keys, container)
		case "/":
			go input.Scan()
			input.ListenKeyboard(uiEvents)
			views.UpdateContainers(string(input.Stdin), container)
			ui.Clear()
			ui.Render(keys, container)
		}
	}
}
