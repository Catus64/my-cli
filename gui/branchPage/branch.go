package branchPage

import (
	"gocmd/gui"
	"gocmd/gui/sidebar"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Show(gui *gui.MyApp, pathName string, window fyne.Window, onHome func(), onSave func(), onModified func(), onIgnored func(), onSaveFile func(), onHistory func(), onHelp func()) {
	sidebar := sidebar.SideBar(gui, window, pathName, "branch", 
		func() { onHome() }, 
		func() { onSave() },
		func() { onModified() }, 
		func() { onIgnored() }, 
		func() { onSaveFile() },
		func() { onHistory() },
		func() { onHelp() },
	)

	fullContent := container.NewBorder(nil, nil, sidebar, nil, branchContent(pathName, window, gui.App))
	window.SetContent(fullContent)
}
