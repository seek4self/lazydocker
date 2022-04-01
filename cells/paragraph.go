package cells

import (
	"lazydocker/docker"

	"github.com/gizak/termui/v3/widgets"
)

type Paragraph struct {
	*widgets.Paragraph

	GetText func(name string) string
}

func NewParagraph() *Paragraph {
	return &Paragraph{
		Paragraph: widgets.NewParagraph(),
		GetText:   docker.Inspect,
	}
}
