package dashboard

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"gocmd/gui"
)

func DashBoard(screen *gui.MyApp) fyne.CanvasObject {
	dashboardTitle := canvas.NewText("DashBoard", color.White)
	dashboardTitle.TextSize = 36
	dashboardTitle.TextStyle = fyne.TextStyle{Bold: true}

	dashboardSubTitle := canvas.NewText("Manage your repository configuration", color.Gray{Y: 150})

	actionTitle := canvas.NewText("Repository Actions", color.White)
	actionTitle.TextSize = 20
	actionTitle.TextStyle = fyne.TextStyle{Bold: true}

	ScreenMargin := canvas.NewRectangle(color.Transparent)
	ScreenMargin.SetMinSize(fyne.NewSize(50, 0)) // Adjust 50 to make the gap

	ScreenHeightMargin := canvas.NewRectangle(color.Transparent)
	ScreenHeightMargin.SetMinSize(fyne.NewSize(0, 30)) // Adjust 50 to make the gap

	content := container.NewVBox(
		ScreenHeightMargin,
		dashboardTitle,
		dashboardSubTitle,
		layout.NewSpacer(),
		actionTitle,
		Action(screen),
		layout.NewSpacer(),
	)

	// NewBorder Syntax: Top, Bottom, Left, Right, Center
	return container.NewBorder(nil, nil, ScreenMargin, ScreenMargin, container.NewPadded(content))
}