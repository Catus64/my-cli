package modifiedFile

import (
	"fmt"
	"image/color"
	gitCurrent "gocmd/testfiles/GitCurrent"
	gitobject "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FileStatus struct {
	Name   string
	Status string // "MODIFIED" or "ADDED"
}

func getFileStatuses(repoPath string) (result []FileStatus) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in getFileStatuses:", r)
		}
	}()

	repo, err := gitpath.Repo_find(repoPath, false)
	if err != nil || repo == nil {
		fmt.Println("No repo found at:", repoPath)
		return
	}

	index, err := gitobject.Index_Read2(*repo)
	if err != nil {
		fmt.Println("Failed to read index:", err)
		return
	}

	// --- MODIFIED: only run if HEAD exists ---
	modified, untracked, err := gitCurrent.StatusIndexWorktree(*repo, *index)
	if err == nil {
		for _, f := range modified {
			result = append(result, FileStatus{Name: f, Status: "MODIFIED"})
		}
		for _, f := range untracked {
			result = append(result, FileStatus{Name: f, Status: "ADDED"})
		}
	}

	return
}

func FolderDirectory(repoPath string) fyne.CanvasObject {
	title := canvas.NewText("File Directory", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Manage your file in your directory", color.Gray{Y: 150})
	subtitle.TextSize = 15

	// Load files from git
	files := getFileStatuses(repoPath)

	modifiedBox, updateModifiedList := modifiedListBox(&files)
	updateModifiedList()

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	modifiedFileContent := container.NewVBox(heightMargin, title, subtitle, heightMargin, modifiedBox)

	return container.NewBorder(nil, nil, widthMargin, widthMargin, container.NewPadded(modifiedFileContent))
}

func modifiedListBox(files *[]FileStatus) (fyne.CanvasObject, func()) {
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

	scrollableFileList := container.NewVScroll(fileList)
	scrollableFileList.SetMinSize(fyne.NewSize(0, 300))

	addButton := widget.NewButton("Add", func() {})
	addButton.Importance = widget.HighImportance

	ignoreButton := widget.NewButton("Ignore", func() {})
	ignoreButton.Importance = widget.DangerImportance

	buttonWidth := canvas.NewRectangle(color.Transparent)
	buttonWidth.SetMinSize(fyne.NewSize(100, 0))

	buttonRow := container.NewHBox(
		layout.NewSpacer(),
		container.NewStack(buttonWidth, addButton),
		layout.NewSpacer(),
		container.NewStack(buttonWidth, ignoreButton),
		layout.NewSpacer(),
	)
	button := container.NewVBox(buttonRow, TDMargin)

	background := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})
	background.StrokeColor = color.RGBA{R: 208, G: 200, B: 200, A: 255}
	background.StrokeWidth = 1
	background.CornerRadius = 8
	background.SetMinSize(fyne.NewSize(0, 420))

	content := container.NewBorder(modifiedListHeader, button, nil, nil, scrollableFileList)
	box := container.NewStack(background, container.NewPadded(content))

	update := func() {
		modifiedListTitle.Text = fmt.Sprintf("File List (%d)", len(*files))
		modifiedListTitle.Refresh()

		fileList.Objects = nil
		for _, file := range *files {
			checkbox := widget.NewCheck("", func(checked bool) {})

			fileName := canvas.NewText(file.Name, color.White)
			fileName.TextSize = 14

			row := container.NewHBox(
				checkbox,
				fileName,
			)

			// Different color for ADDED vs MODIFIED
			var statusColor color.Color
			if file.Status == "ADDED" {
				statusColor = color.RGBA{R: 100, G: 200, B: 100, A: 255} // green
			} else {
				statusColor = color.Gray{Y: 150} // grey
			}

			status := canvas.NewText(file.Status, statusColor)
			status.TextSize = 10
			statusWithSpace := container.NewHBox(status, LRMargin)
			statusPosition := container.NewBorder(nil, nil, nil, statusWithSpace, row)
			fileList.Add(statusPosition)
		}
		fileList.Refresh()
		scrollableFileList.Refresh()
	}

	return box, update
}