package gui

import (
	"lazydocker/docker"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func StatusView() *widgets.Table {
	table := widgets.NewTable()
	status := docker.PS()
	table.Rows = append(table.Rows, []string{"Name", "Status", "Age"})
	for _, s := range status {
		// fmt.Println(s)
		age := "--"
		if s.Status[0:2] == "Up" {
			age = s.Status[3:]
		}
		table.Rows = append(table.Rows, []string{s.Name, strings.Fields(s.Status)[0], age})
	}
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.SetRect(0, 0, 40, len(status)+3)
	table.Title = "Containers"
	table.TextAlignment = ui.AlignCenter
	table.RowSeparator = false
	return table
}
