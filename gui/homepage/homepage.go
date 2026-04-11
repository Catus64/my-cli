package homepage

import (
	"gocmd/gui"
	"gocmd/gui/sidebar"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Show(gui *gui.MyApp, pathName string) {
	myWindow := gui.App.NewWindow(gui.Window.Title())
	sidebar := sidebar.SideBar(gui, myWindow, pathName, "home")

	fullContent := container.NewBorder(nil, nil, sidebar, nil, HomePageContent())
	myWindow.SetContent(fullContent)
	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.Show()
}
