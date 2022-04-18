package app

import (
	"lazydocker/views"
)

func initKeys(a *App) {
	a.keys.Text = "q: quit, j/↓ or k/↑, Tab: switch containers/images, h/← ok l/→ : switch navigation, s: change containers status"
	a.keys.SetRect(0, views.TerminalHeight-1, views.TerminalWidth, views.TerminalHeight)
	a.keys.Border = false
	a.keys.TextStyle = a.activeStyle
}
