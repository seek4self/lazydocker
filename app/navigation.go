package app

import (
	"lazydocker/docker"
	"lazydocker/views"
)

const (
	containerInfo  = "Container Info"
	containerLog   = "Log"
	containerStats = "Stats"
	imageInfo      = "Image Info"
	imageHistory   = "History"
)

var (
	headerContainer = []views.Header{
		{Name: containerInfo, WrapText: true, Detail: docker.ContainerInspect},
		{Name: containerLog, WrapText: true, BotUp: true, Detail: docker.Logs},
		{Name: containerStats, Fresh: true, Detail: docker.Stats},
	}
	headerImage = []views.Header{
		{Name: imageInfo, WrapText: true, Detail: docker.ImageInspect},
		{Name: imageHistory, WrapText: true, Detail: docker.History},
	}
)

func initNavigation(a *App) {
	a.navigation.Header = [][]views.Header{headerContainer, headerImage}
	a.navigation.Title = "Navigation"
	a.navigation.SetRect(a.containers.Inner.Max.X+1, 0, views.TerminalWidth, views.TerminalHeight-1)
	a.navigation.Fresh(a.containers.ActiveText(), a.activeIndex)
}
