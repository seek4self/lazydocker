package views

import (
	"lazydocker/cells"
	"lazydocker/docker"
	"strings"

	ui "github.com/gizak/termui/v3"
)

// containers status
const (
	Running = iota
	Exited
	All
	Total
)

var options = []string{docker.OptRuning, docker.OptExited, docker.OptAll}

var currentStatus = Running

func containerOption() string {
	currentStatus = (currentStatus + 1) % Total
	return options[currentStatus]
}

func freshContainers(option string, table *cells.Table) {
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
	table.ResetSize(0, 0, 50, cells.TerminalHeight/2)
	table.Title = strings.Join([]string{"Containers ", option}, ": ")
}

func initContainers(v *View) {
	v.containers.Header = []string{"Name", "Status", "Age"}
	// table.Rows = append(table.Rows, []string{"dms", "Exited", "--"})
	// table.Rows = append(table.Rows, []string{"ums", "Up", "3 days"})
	// table.Rows = append(table.Rows, []string{"gws", "Up", "3 days"})
	freshContainers(docker.OptRuning, v.containers)
	v.containers.ColumnWidths = []int{35, 8, 5}
	v.containers.ColumnAlignment = []ui.Alignment{ui.AlignLeft, ui.AlignLeft, ui.AlignRight}
	v.containers.TextStyle = ui.NewStyle(ui.ColorWhite)
	v.containers.TabTitle = "Container info"
	v.containers.TabContent = docker.Inspect
}
