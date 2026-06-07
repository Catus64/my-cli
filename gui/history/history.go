package history

import (
	"gocmd/gui"
	"gocmd/gui/sidebar"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Show(g *gui.MyApp, pathName string, window fyne.Window, onHome func(), onSave func(), onModified func(), onIgnored func(), onSaveFile func(), onHistory func(), onHelp func()) {
	sidebar := sidebar.SideBar(g, window, pathName, "history",
		func() { onHome() },
		func() { onSave() },
		func() { onModified() },
		func() { onIgnored() },
		func() { onSaveFile() },
		func() { onHistory() },
		func() { onHelp() },
	)

	fullContent := container.NewBorder(nil, nil, sidebar, nil, HistoryPageContent(pathName, g.App, window))
	window.SetContent(fullContent)
}
