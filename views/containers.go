package views

import (
	"lazydocker/cells"
	"lazydocker/docker"

	ui "github.com/gizak/termui/v3"
)

func ContainerStatus() *cells.Table {
	table := cells.NewTable()
	table.Header = []string{"Name", "Status", "Age"}
	// table.Rows = append(table.Rows, []string{"dms", "Exited", "--"})
	// table.Rows = append(table.Rows, []string{"ums", "Up", "3 days"})
	// table.Rows = append(table.Rows, []string{"gws", "Up", "3 days"})
	UpdateContainers("up", table)
	table.ColumnWidths = []int{20, 8, 12}
	table.ColumnAlignment = []ui.Alignment{ui.AlignLeft, ui.AlignCenter, ui.AlignLeft}
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.Title = "Containers"
	table.TabTitle = "Contain info"
	table.TabContent = docker.Inspect
	// table.ColumnSeparator = true
	return table
}

func UpdateContainers(option string, table *cells.Table) {
	table.Rows = make([][]string, 0)
	status := docker.PS(option)
	for _, s := range status {
		// fmt.Println(s)
		// age := "--"
		// if s.Status[0:2] == "Up" {
		// 	age = s.Status[3:]
		// }
		table.Rows = append(table.Rows, []string{s.Name, s.Status, s.Age})
	}
	table.ResetSize(0, 0, 40, cells.TerminalHeight/2)
}
