package views

import (
	"fmt"
	"lazydocker/cells"
	"lazydocker/docker"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
)

const (
	B  = 1
	KB = B << 10
	MB = KB << 10
	GB = MB << 10
	TB = GB << 10
)

func freshImages(option string, table *cells.Table) {
	table.Rows = make([][]string, 0)
	images := docker.Images(option)
	for _, i := range images {
		table.Rows = append(table.Rows, []string{
			i.ID[7:19],
			time.Unix(i.Created, 0).Format("2006-01-02 15:04:05"),
			parseSize(i.Size),
		})
	}
	table.ResetSize(0, cells.TerminalHeight/2, 50, cells.TerminalHeight-1)
	table.Title = strings.Join([]string{"Images ", option}, ": ")
}

func initImages(v *View) {
	v.images.Header = []string{"ImageID", "Created", "Size"}
	v.images.Rows = make([][]string, 0)
	freshImages("", v.images)
	v.images.ColumnWidths = []int{16, 24, 8}
	v.images.ColumnAlignment = []ui.Alignment{ui.AlignLeft, ui.AlignLeft, ui.AlignRight}
	v.images.TabTitle = "Image info"
	v.images.TabContent = docker.ImageInspect
}

func parseSize(size int64) string {
	if size < KB {
		return fmt.Sprintf("%.1fB", float64(size)/B)
	}
	if size < MB {
		return fmt.Sprintf("%.1fK", float64(size)/KB)
	}
	if size < GB {
		return fmt.Sprintf("%.1fM", float64(size)/MB)
	}
	if size < TB {
		return fmt.Sprintf("%.1fG", float64(size)/GB)
	}
	return fmt.Sprintf("%.1fT", float64(size)/TB)
}
