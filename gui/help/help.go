package help

import (
	"gocmd/gui"
	"gocmd/gui/sidebar"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Show(gui *gui.MyApp, pathName string, window fyne.Window, onHome func(), onSave func(), onModified func(), onIgnored func(), onHistory func()) {
	sidebar := sidebar.SideBar(gui, window, pathName, "help", 
		func() { onHome() }, 
		func() { onSave() },
		func() { onModified() },
		func() { onIgnored() },
		func() { onHistory() },
		func() {}, // stay at default page
	)

	fullContent := container.NewBorder(nil, nil, sidebar, nil, HelpContent())
	window.SetContent(fullContent)
}
