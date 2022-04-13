package main

import (
	"lazydocker/app"
	"log"

	ui "github.com/gizak/termui/v3"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	app.NewApp().Init().Render()
}
