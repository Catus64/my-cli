package homepage

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func HomePageContent() fyne.CanvasObject {
	title := canvas.NewText("Home Page", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Manage your save list", color.Gray{Y: 150})
	subtitle.TextSize = 15

	// Save List and Preview Box
	files := []string{}

	saveListBox, updateSaveList := saveListBox(&files)
	previewBox, updatePreview := saveListPreview(&files)

	onFilesChanged := func() {
        updateSaveList()
        updatePreview()
    }
    _ = onFilesChanged

	// Commit Box
	commitBox := commitBox()

	columnGap := canvas.NewRectangle(color.Transparent)
	columnGap.SetMinSize(fyne.NewSize(10, 0))

	rowGap := canvas.NewRectangle(color.Transparent)
	rowGap.SetMinSize(fyne.NewSize(0, 10))

	rightColumn := container.NewBorder(nil, container.NewVBox(rowGap, previewBox), nil, nil, commitBox)

	paddedRight := container.NewBorder(nil, nil, columnGap, nil, rightColumn)

	content := container.NewGridWithColumns(2, saveListBox, paddedRight)

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	homeContent := container.NewVBox(heightMargin, title, subtitle, heightMargin, content)

	return container.NewBorder(nil, nil, widthMargin, widthMargin, container.NewPadded(homeContent))
}

func saveListBox(file *[]string) (fyne.CanvasObject, func()) {
	saveListTitle := canvas.NewText(fmt.Sprintf("Save List (%d)", len(*file)), color.RGBA{R: 208, G: 200, B: 200, A: 255})
	saveListTitle.TextSize = 20
	saveListTitle.TextStyle = fyne.TextStyle{Bold: true}

	titleLine := canvas.NewRectangle(color.RGBA{R: 208, G: 200, B: 200, A: 255})
    titleLine.SetMinSize(fyne.NewSize(0, 1))

	LRMargin := canvas.NewRectangle(color.Transparent)
   	LRMargin.SetMinSize(fyne.NewSize(10, 0))

	TDMargin := canvas.NewRectangle(color.Transparent)
   	TDMargin.SetMinSize(fyne.NewSize(0, 5))

	title := container.NewHBox(LRMargin, saveListTitle)

	saveListHeader := container.NewVBox(TDMargin, title, TDMargin, titleLine)

	fileList := container.NewVBox()

	removeButton := widget.NewButton("Remove", func() {

	})
	removeButton.Importance = widget.DangerImportance
	removeButtonRow := container.NewHBox(layout.NewSpacer(), removeButton, layout.NewSpacer())
	removeBtn := container.NewVBox(removeButtonRow, TDMargin)

	background := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})
	background.StrokeColor = color.RGBA{R: 208, G: 200, B: 200, A: 255}
	background.StrokeWidth = 1
	background.CornerRadius = 8
	background.SetMinSize(fyne.NewSize(0, 350))

	content := container.NewBorder(saveListHeader, removeBtn, nil, nil, fileList)

	box := container.NewStack(background, container.NewPadded(content))

	// Update file
	update := func() {
		saveListTitle.Text = fmt.Sprintf("Save List (%d)", len(*file))
        saveListTitle.Refresh()

		fileList.Objects = nil
		for _, files:= range *file {
			fileName := files
			checkbox := widget.NewCheck(fileName, func(checked bool) {})
			status := canvas.NewText("ADDED", color.Gray{Y: 150})
			status.TextSize = 10
			statusPosition := container.NewBorder(nil, nil, nil, status, checkbox)
			fileList.Add(statusPosition)
		}
		fileList.Refresh()
	}
	return box, update
}

func commitBox() fyne.CanvasObject {
	commitMessageEntry := widget.NewMultiLineEntry()
	commitMessageEntry.Wrapping = fyne.TextWrapWord

	placeholder := canvas.NewText("Summary (required)", color.Gray{Y: 150})
    placeholder.Alignment = fyne.TextAlignCenter
    placeholder.TextSize = 14

    // hide placeholder when user types
    commitMessageEntry.OnChanged = func(text string) {
        if text == "" {
            placeholder.Show()
        } else {
            placeholder.Hide()
        }
        placeholder.Refresh()
    }

	heightMargin := canvas.NewRectangle(color.Transparent)
    heightMargin.SetMinSize(fyne.NewSize(0, 5))

	widthMargin := canvas.NewRectangle(color.Transparent)
    widthMargin.SetMinSize(fyne.NewSize(5, 0))

	placeholderContainer := container.NewBorder(container.NewPadded(placeholder), nil, nil, nil, nil)
	placeholderPosition := container.NewVBox(heightMargin, placeholderContainer)

	commitBackground := canvas.NewRectangle(color.Transparent)
	commitBackground.CornerRadius = 8
	commitBackground.SetMinSize(fyne.NewSize(0, 220))

	commitBox := container.NewStack(commitBackground, container.NewPadded(commitMessageEntry), placeholderPosition)

	saveButton := widget.NewButton("Save", func() {})
	saveButton.Importance = widget.HighImportance
	saveButtonRow := container.NewHBox(layout.NewSpacer(), saveButton, layout.NewSpacer())

	background := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})
	background.StrokeColor = color.RGBA{R: 208, G: 200, B: 200, A: 255}
	background.StrokeWidth = 1
	background.CornerRadius = 8

	saveBtn := container.NewVBox(heightMargin, saveButtonRow, heightMargin)

	innerContent := container.NewBorder(heightMargin, saveBtn, widthMargin, widthMargin, commitBox)

    return container.NewStack(background, container.NewPadded(innerContent))
}

func saveListPreview(file *[]string) (fyne.CanvasObject, func()) {
	previewTitle := canvas.NewText("Save List Preview", color.White)
	previewTitle.TextSize = 16
	previewTitle.TextStyle = fyne.TextStyle{Bold: true}

	previewSubTitle := canvas.NewText(fmt.Sprintf("%d file(s) ready to be saved.", len(*file)), color.Gray{Y: 150})
	previewSubTitle.TextSize = 12

	previewList := container.NewVBox()

	background := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})
	background.StrokeColor = color.RGBA{R: 208, G: 200, B: 200, A: 255}
	background.StrokeWidth = 1
	background.CornerRadius = 8
	background.SetMinSize(fyne.NewSize(0, 120))

	leftMargin := canvas.NewRectangle(color.Transparent)
    leftMargin.SetMinSize(fyne.NewSize(10, 0))

	topMargin := canvas.NewRectangle(color.Transparent)
    topMargin.SetMinSize(fyne.NewSize(0, 5))

	content := container.NewVBox(topMargin, previewTitle, previewSubTitle, previewList)
	fullcontent := container.NewHBox(leftMargin, content)
	box := container.NewStack(background, container.NewPadded(fullcontent))

	// Update file
	update := func ()  {
		previewSubTitle.Text = fmt.Sprintf("Save List (%d)", len(*file))
        previewSubTitle.Refresh()

		previewList.Objects = nil
		for _, files:= range *file {
			bullet := canvas.NewText("• "+files, color.Gray{Y: 200})
			bullet.TextSize = 12
			previewList.Add(bullet)
		}
		previewList.Refresh()
	}

	return box, update
}
