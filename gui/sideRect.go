package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"ezgit/img"
)

func (screen *MyApp) SideRect() fyne.CanvasObject {
	sideRect := canvas.NewRectangle(color.RGBA{R: 2, G: 35, B: 62, A: 255})
	sideRect.SetMinSize(fyne.NewSize(150, 0))
	sideRect.StrokeColor = color.RGBA{R: 204, G: 200, B: 200, A: 255}
	sideRect.StrokeWidth = 1

	data, _ := img.Assets.ReadFile("ezgitText.png")
    res := fyne.NewStaticResource("ezgitText.png", data)
    sideImg := canvas.NewImageFromResource(res)
	sideImg.FillMode = canvas.ImageFillContain

	sideBox := container.NewStack(sideRect, container.NewPadded(sideImg))

	Margin := canvas.NewRectangle(color.Transparent)
	Margin.SetMinSize(fyne.NewSize(0, 10))

	LeftMargin := canvas.NewRectangle(color.Transparent)
	LeftMargin.SetMinSize(fyne.NewSize(10, 0))

	return container.NewBorder(Margin, Margin, LeftMargin, nil, sideBox)
}