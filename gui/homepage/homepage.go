package homepage

import (
	"gocmd/gui"
	"gocmd/gui/sidebar"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Show(gui *gui.MyApp, pathName string, window fyne.Window, onSave func(), onModified func(), onIgnored func(), onHistory func(), onHelp func()) {
	sidebar := sidebar.SideBar(gui, window, pathName, "home", 
		func() {}, // stay at default page
		func() { onSave() },
		func() { onModified() },
		func() { onIgnored() },
		func() { onHistory() },
		func() { onHelp() },
	)

	fullContent := container.NewBorder(nil, nil, sidebar, nil, HomePageContent())
	window.SetContent(fullContent)
}
