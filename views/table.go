package views

import (
	"lazydocker/cells"
	"lazydocker/docker"
	"strings"

	ui "github.com/gizak/termui/v3"
)

func ContainerStatus() *cells.Table {
	table := cells.NewTable()
	status := docker.PS()
	table.Rows = append(table.Rows, []string{"Name", "Status", "Age"})
	// table.Rows = append(table.Rows, []string{"dms", "Exited", "--"})
	// table.Rows = append(table.Rows, []string{"ums", "Up", "3 days"})
	// table.Rows = append(table.Rows, []string{"gws", "Up", "3 days"})
	for _, s := range status {
		// fmt.Println(s)
		age := "--"
		if s.Status[0:2] == "Up" {
			age = s.Status[3:]
		}
		table.Rows = append(table.Rows, []string{s.Name, strings.Fields(s.Status)[0], age})
	}
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.SetRect(0, 1, 40, len(status)+3+1)
	table.Title = "Containers"
	// table.TextAlignment = ui.AlignCenter
	table.RowSeparator = false
	return table
}
