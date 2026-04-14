package modifiedFile

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ModifiedFileContent() fyne.CanvasObject {
	title := canvas.NewText("Modified File", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Manage your modified file", color.Gray{Y: 150})
	subtitle.TextSize = 15

	// Modified file list
	files := []string{}

	modifiedListBox, updateModifiedList := modifiedListBox(&files)

	onFilesChanged := func() { updateModifiedList() }
    _ = onFilesChanged

	// Space
	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	modifiedFileContent := container.NewVBox(heightMargin, title, subtitle, heightMargin, modifiedListBox)

	return container.NewBorder(nil, nil, widthMargin, widthMargin, container.NewPadded(modifiedFileContent))
}

func modifiedListBox(files *[]string) (fyne.CanvasObject, func()) {
	modifiedListTitle := canvas.NewText(fmt.Sprintf("File List (%d)", len(*files)), color.RGBA{R: 208, G: 200, B: 200, A: 255})
	modifiedListTitle.TextSize = 20
	modifiedListTitle.TextStyle = fyne.TextStyle{Bold: true}

	titleLine := canvas.NewRectangle(color.RGBA{R: 208, G: 200, B: 200, A: 255})
    titleLine.SetMinSize(fyne.NewSize(0, 1))

	LRMargin := canvas.NewRectangle(color.Transparent)
   	LRMargin.SetMinSize(fyne.NewSize(10, 0))

	TDMargin := canvas.NewRectangle(color.Transparent)
   	TDMargin.SetMinSize(fyne.NewSize(0, 5))

	title := container.NewHBox(LRMargin, modifiedListTitle)

	modifiedListHeader := container.NewVBox(TDMargin, title, TDMargin, titleLine)

	fileList := container.NewVBox()

	addButton := widget.NewButton("Add", func() {

	})
	addButton.Importance = widget.HighImportance

	buttonWidth := canvas.NewRectangle(color.Transparent)
	buttonWidth.SetMinSize(fyne.NewSize(100, 0))

	addButtonRow := container.NewHBox(layout.NewSpacer(), container.NewStack(buttonWidth, addButton), layout.NewSpacer())
	addBtn := container.NewVBox(addButtonRow, TDMargin)

	background := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})
	background.StrokeColor = color.RGBA{R: 208, G: 200, B: 200, A: 255}
	background.StrokeWidth = 1
	background.CornerRadius = 8
	background.SetMinSize(fyne.NewSize(0, 420))

	content := container.NewBorder(modifiedListHeader, addBtn, nil, nil, fileList)

	box := container.NewStack(background, container.NewPadded(content))

	// Update file
	update := func() {
		modifiedListTitle.Text = fmt.Sprintf("File List (%d)", len(*files))
        modifiedListTitle.Refresh()

		fileList.Objects = nil
		for _, files:= range *files {
			fileName := files
			checkbox := widget.NewCheck(fileName, func(checked bool) {})
			status := canvas.NewText("MODIFIED", color.Gray{Y: 150})
			status.TextSize = 10
			statusPosition := container.NewBorder(nil, nil, nil, status, checkbox)
			fileList.Add(statusPosition)
		}
		fileList.Refresh()
	}
	return box, update
}