package views

import (
	"lazydocker/cells"
	"lazydocker/docker"

	ui "github.com/gizak/termui/v3"
)

func setTableRows(option string, table *cells.Table) {
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

func initContainers(v *View) {
	v.containers.Header = []string{"Name", "Status", "Age"}
	// table.Rows = append(table.Rows, []string{"dms", "Exited", "--"})
	// table.Rows = append(table.Rows, []string{"ums", "Up", "3 days"})
	// table.Rows = append(table.Rows, []string{"gws", "Up", "3 days"})
	setTableRows("up", v.containers)
	v.containers.ColumnWidths = []int{20, 8, 12}
	v.containers.ColumnAlignment = []ui.Alignment{ui.AlignLeft, ui.AlignCenter, ui.AlignLeft}
	v.containers.TextStyle = ui.NewStyle(ui.ColorWhite)
	v.containers.Title = "Containers"
	v.containers.TabTitle = "Container info"
	v.containers.TabContent = docker.Inspect
}
