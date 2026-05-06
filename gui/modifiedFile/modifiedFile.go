package modifiedFile

import (
	"gocmd/gui"
	"gocmd/gui/sidebar"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Show(gui *gui.MyApp, pathName string, window fyne.Window, onHome func(), onSave func(), onModified func(), onIgnored func(), onHistory func(), onHelp func()) {
	sidebar := sidebar.SideBar(gui, window, pathName, "file-directory", 
		func() { onHome() }, 
		func() { onSave() },
		func() { onModified() }, 
		func() { onIgnored() }, 
		func() { onHistory() },
		func() { onHelp() },
	)

	fullContent := container.NewBorder(nil, nil, sidebar, nil, FolderDirectory(pathName))
	window.SetContent(fullContent)
}
