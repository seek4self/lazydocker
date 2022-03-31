package main

import (
	"lazydocker/docker"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := widgets.NewTable()
	status := docker.PS()
	table.Rows = append(table.Rows, []string{"Name", "Status"})
	for _, s := range status {
		// fmt.Println(s)
		table.Rows = append(table.Rows, []string{s.Name, s.Status})
	}
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.SetRect(0, 0, 60, len(status)+3)
	table.TextAlignment = ui.AlignCenter
	table.RowSeparator = false
	ui.Render(table)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
