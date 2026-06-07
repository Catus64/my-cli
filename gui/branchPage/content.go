package branchPage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

func branchContent() fyne.CanvasObject {
	title := canvas.NewText("Branch", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Manage your branches", color.Gray{Y: 150})
	subtitle.TextSize = 15

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	branchContent := container.NewVBox(heightMargin, title, subtitle, heightMargin)

	return container.NewBorder(nil, nil, widthMargin, widthMargin, container.NewPadded(branchContent))
}