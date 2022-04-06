package views

import (
	"fmt"
	"lazydocker/cells"
	"lazydocker/docker"

	ui "github.com/gizak/termui/v3"
)

const (
	B  = 1
	KB = B << 10
	MB = KB << 10
	GB = MB << 10
	TB = GB << 10
)

func initImages(v *View) {
	v.images.Header = []string{"Name", "Size"}
	v.images.Rows = make([][]string, 0)
	images := docker.Images()
	for _, i := range images {
		name := i.ID
		if len(i.RepoTags) > 0 {
			name = i.RepoTags[0]
		}
		v.images.Rows = append(v.images.Rows, []string{name, parseSize(i.Size)})
	}
	v.images.ResetSize(0, cells.TerminalHeight/2, 40, cells.TerminalHeight-1)
	v.images.ColumnWidths = []int{31, 9}
	v.images.ColumnAlignment = []ui.Alignment{ui.AlignLeft, ui.AlignRight}
	v.images.Title = "Images"
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
