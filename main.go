package main

import (
	"log/slog"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"gocmd/gui"
	"gocmd/gui/dashboard"
	logger "gocmd/testfiles/Helper"
)

func main() {
	logger.Init("ezgit.log", slog.LevelDebug)

	ezgit := app.New()
	ezgit.Settings().SetTheme(&gui.MyDarkTheme{})

	gui := &gui.MyApp{
		App: ezgit,
		Window: ezgit.NewWindow("Ezgit (Single-User Version Control System)"),
	}

	fullcontent := container.NewBorder(gui.Header(), nil, gui.SideRect(), nil, dashboard.DashBoard(gui))

	gui.Window.SetContent(fullcontent)
	gui.Window.Resize(fyne.NewSize(1000, 600))
	gui.Window.ShowAndRun()
}
