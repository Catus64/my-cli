package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"ezgit/img"
)

func (screen *MyApp) Header() fyne.CanvasObject {
	return screen.BuildHeader(true) //with header line
}

func (screen *MyApp) HeaderNoneLine() fyne.CanvasObject {
	return screen.BuildHeader(false) //with no header line
}

func (screen *MyApp) BuildHeader(line bool) fyne.CanvasObject {
	headerTitle := canvas.NewText("EzGit", color.White)
	headerTitle.TextSize = 24
	headerTitle.TextStyle = fyne.TextStyle{Bold: true}

	headerSub := canvas.NewText("A Single-User Version Control System", color.Gray{Y: 150})
	headerSub.TextSize = 12

	HeaderText := container.NewVBox(headerTitle, headerSub)

	iconBg := canvas.NewRectangle(color.White)
	iconBg.SetMinSize(fyne.NewSize(50, 50))
	iconBg.CornerRadius = 8

	data, _ := img.Assets.ReadFile("icon.png")
    res := fyne.NewStaticResource("icon.png", data)
    iconImg := canvas.NewImageFromResource(res)
    iconImg.FillMode = canvas.ImageFillContain

	icon := container.NewStack(iconBg, container.NewPadded(iconImg))

	header := container.NewHBox(icon, HeaderText)

	headerContent := container.NewPadded(header)

	if line {
		HeaderLine := canvas.NewRectangle(color.Gray{Y: 100})
		HeaderLine.SetMinSize(fyne.NewSize(0, 2))
		return container.NewVBox(headerContent, HeaderLine)
	}
	
	return headerContent
}