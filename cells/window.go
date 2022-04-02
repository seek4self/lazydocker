package cells

import "github.com/olekukonko/ts"

var (
	TerminalWidth  int
	TerminalHeight int
)

func init() {
	size, _ := ts.GetSize()
	TerminalWidth = size.Col()
	TerminalHeight = size.Row()
}
