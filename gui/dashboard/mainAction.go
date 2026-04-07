package dashboard

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"gocmd/gui"
    "gocmd/gui/createRepo"
    "gocmd/gui/openRepo"
)

func Action(screen *gui.MyApp) fyne.CanvasObject {
	CreateButton := widget.NewButton("Create Repository", func() {
		screen.Window.Hide()
		createRepo.Show(screen)
	})
	OpenButton := widget.NewButton("Open Repostiory", func() {
		screen.Window.Hide()
		openRepo.Show(screen)
	})
	ExitButton := widget.NewButton("Exit", func() { screen.Window.Close() })

	// Put the button into the action box
	btnSpacer := canvas.NewRectangle(color.Transparent)
	btnSpacer.SetMinSize(fyne.NewSize(0, 15))

	LRMargin := canvas.NewRectangle(color.Transparent)
	LRMargin.SetMinSize(fyne.NewSize(30, 0))

	TDMargin := canvas.NewRectangle(color.Transparent)
	TDMargin.SetMinSize(fyne.NewSize(0, 20))

	// Wrap the buttons in a Border layout to apply the margins to the sides
	buttonBox := container.NewVBox(CreateButton, btnSpacer, OpenButton, btnSpacer, ExitButton)
	squeezedButtons := container.NewBorder(TDMargin, TDMargin, LRMargin, LRMargin, buttonBox)

	actionBox := canvas.NewRectangle(color.RGBA{R: 2, G: 35, B: 62, A: 255})
	actionBox.StrokeColor = color.RGBA{R: 204, G: 200, B: 200, A: 255}
	actionBox.StrokeWidth = 2
	actionBox.CornerRadius = 8

	return container.NewStack(actionBox, container.NewPadded(squeezedButtons))
}