package views

import (
	"lazydocker/cells"
	"lazydocker/docker"
	"strings"

	ui "github.com/gizak/termui/v3"
)

func ContainerStatus(option string) *cells.Table {
	table := cells.NewTable()
	table.Header = []string{"Name", "Status", "Age"}
	// table.Rows = append(table.Rows, []string{"dms", "Exited", "--"})
	// table.Rows = append(table.Rows, []string{"ums", "Up", "3 days"})
	// table.Rows = append(table.Rows, []string{"gws", "Up", "3 days"})
	UpdateContainers(option, table)
	table.ColumnWidths = []int{20, 8, 12}
	table.ColumnAlignment = []ui.Alignment{ui.AlignLeft, ui.AlignCenter, ui.AlignLeft}
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.SetRect(0, 1, 40, 10)
	table.Title = "Containers"
	// table.ColumnSeparator = true
	return table
}

func UpdateContainers(option string, table *cells.Table) {
	table.Rows = make([][]string, 0)
	status := docker.PS(option)
	for _, s := range status {
		// fmt.Println(s)
		age := "--"
		if s.Status[0:2] == "Up" {
			age = s.Status[3:]
		}
		table.Rows = append(table.Rows, []string{s.Name, strings.Fields(s.Status)[0], age})
	}
	table.ActiveRowIndex = 0
	table.Page = 0
}
