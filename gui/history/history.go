package history

import (
	"gocmd/gui"
	"gocmd/gui/sidebar"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Show(gui *gui.MyApp, pathName string, window fyne.Window, onHome func(), onModified func(), onHelp func()) {
	sidebar := sidebar.SideBar(gui, window, pathName, "history", 
		func() { onHome() }, 
		func() { onModified() },
		func() {}, // stay at default page
		func() { onHelp() },
	)

	fullContent := container.NewBorder(nil, nil, sidebar, nil, nil)
	window.SetContent(fullContent)
}
