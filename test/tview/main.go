package main

import (
	"lazydocker/docker"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	// box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	app := tview.NewApplication()
	table := tview.NewTable()
	form := tview.NewForm().SetHorizontal(true)
	form.SetBorder(false).SetTitle("Control")

	form.AddDropDown("Docker", []string{"container", "image"}, 0, nil).
		AddDropDown("Status", []string{docker.OptRuning, docker.OptExited, docker.OptAll}, 0, nil).
		AddInputField("Filter", "", 10, nil, nil).
		AddButton("Search", func() {
			_, opt := form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
			if opt == "container" {
				_, opt := form.GetFormItem(1).(*tview.DropDown).GetCurrentOption()
				freshTable(table, opt)
				app.SetFocus(table)
			}
		})

	table.SetFixed(1, 1)
	freshTable(table, docker.OptRuning)
	table.SetSelectable(true, false)
	table.Select(1, 0)
	table.SetDoneFunc(func(key tcell.Key) {
		app.SetFocus(form)
	})

	table.SetBorder(true).SetTitle("Containers").SetRect(0, 0, 50, 10)
	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(form, 3, 0, true).
			AddItem(table, 20, 0, false),
			60, 0, true)
	app.SetRoot(flex, false)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func freshTable(table *tview.Table, opt string) {
	cs := docker.PS(opt)
	table.SetCell(0, 0, tview.NewTableCell("Name").SetMaxWidth(20).SetAttributes(tcell.AttrBold))
	table.SetCell(0, 1, tview.NewTableCell("Status").SetAttributes(tcell.AttrBold).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Age").SetAttributes(tcell.AttrBold))
	table.SetCell(0, 3, tview.NewTableCell("ID").SetAttributes(tcell.AttrBold))
	table.SetCell(0, 4, tview.NewTableCell("Ports").SetAttributes(tcell.AttrBold))
	table.SetCell(0, 5, tview.NewTableCell("CMD").SetAttributes(tcell.AttrBold))
	for i, v := range cs {
		table.SetCell(i+1, 0, tview.NewTableCell(v.Name).SetMaxWidth(20).SetExpansion(1))
		table.SetCell(i+1, 1, tview.NewTableCell(v.Status).SetExpansion(1))
		table.SetCell(i+1, 2, tview.NewTableCell(v.Age).SetAlign(tview.AlignRight).SetExpansion(1))
		table.SetCell(i+1, 3, tview.NewTableCell(v.ID).SetExpansion(1))
		table.SetCell(i+1, 4, tview.NewTableCell(v.Ports).SetExpansion(1))
		table.SetCell(i+1, 5, tview.NewTableCell(v.CMD).SetExpansion(1))
	}
}
