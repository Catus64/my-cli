package history

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	gitCurrent "gocmd/testfiles/GitCurrent"
	gitcheckout "gocmd/testfiles/GitCheckout"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowSaveHistoryWindow(app fyne.App, clicked Commit, files []string, repo *gitpath.GitRepository) {	
	window := app.NewWindow(fmt.Sprintf("Commit: %s", clicked.SHA))
	window.Resize(fyne.NewSize(900, 600))

	topHash := clicked.ShortSHA

	// Commit Hash
	hashTitle := canvas.NewText(topHash, color.White)
	hashTitle.TextSize = 40
	hashTitle.TextStyle = fyne.TextStyle{Bold: true}

	hashUnderline := canvas.NewRectangle(color.White)
	hashUnderline.SetMinSize(fyne.NewSize(0, 2))

	hashHeader := container.NewVBox(hashTitle, hashUnderline)
	headerArea := container.NewHBox(hashHeader, layout.NewSpacer())

	// get all files for current commit
	allFiles := []string{}
	if repo != nil && clicked.TreeSHA != "" {
		treeMap := map[string]string{}
		gitCurrent.TreeToMap(*repo, clicked.TreeSHA, "", treeMap)
		for path := range treeMap {
			allFiles = append(allFiles, path)
		}
	}

	// File List
	fileListTitle := canvas.NewText(fmt.Sprintf("File List (%d)", len(allFiles)), color.White)
	fileListTitle.TextSize = 16
	fileListTitle.TextStyle = fyne.TextStyle{Bold: true}

	currentChangeButton := widget.NewButton("Current Change", func() {
		if len(files) == 0 {
			dialog.ShowInformation("No Changes", "No changed files in this commit.", window)
			return
		}

		changeList := container.NewVBox()
		for _, file := range files {
			bullet := canvas.NewText("• "+file, color.Black)
			bullet.TextSize = 13
			changeList.Add(bullet)
		}

		scroll := container.NewScroll(changeList)
		scroll.SetMinSize(fyne.NewSize(400, 300))
		dialog.ShowCustom(fmt.Sprintf("Changed Files (%d)", len(files)), "Close", scroll, window)
	})
	currentChangeButton.Importance = widget.HighImportance

	titleLine := canvas.NewRectangle(color.RGBA{R: 208, G: 200, B: 200, A: 255})
	titleLine.SetMinSize(fyne.NewSize(0, 1))

	fileList := container.NewVBox()

	topSpacing := canvas.NewRectangle(color.Transparent)
	topSpacing.SetMinSize(fyne.NewSize(0, 5)) 
	fileList.Add(topSpacing)

	for _, f := range allFiles {
		bullet := canvas.NewText("• "+f, color.White)
		bullet.TextSize = 13

		leftSpacing := canvas.NewRectangle(color.Transparent)
		leftSpacing.SetMinSize(fyne.NewSize(10, 0))

		lrRow := container.NewHBox(leftSpacing, bullet)

		itemRow := container.NewVBox(topSpacing, lrRow)

		fileList.Add(itemRow)
	}

	scrollFileList := container.NewScroll(fileList)
	scrollFileList.SetMinSize(fyne.NewSize(250, 300))

	LRMargin := canvas.NewRectangle(color.Transparent)
	LRMargin.SetMinSize(fyne.NewSize(10, 0))

	TDMargin := canvas.NewRectangle(color.Transparent)
	TDMargin.SetMinSize(fyne.NewSize(0, 5))

	titleRow := container.NewBorder(nil, nil, container.NewHBox(LRMargin, fileListTitle), currentChangeButton, nil)
	header := container.NewVBox(titleRow, titleLine)

	// Load button at bottom
	loadBtn := widget.NewButton("Load", func() {
		if repo == nil {
			dialog.ShowError(fmt.Errorf("no repository found"), window)
			return
		}

		// Require user to enter the folder name to load
		nameEntry := widget.NewEntry()
    	nameEntry.SetText(clicked.ShortSHA)

		form := []*widget.FormItem{
			widget.NewFormItem("Folder Name", nameEntry),
		}

		loadDialog := dialog.NewForm("Load Version", "Load", "Cancel", form, func(confirmed bool) {
			if !confirmed {
				return
			}

			folderName := strings.TrimSpace(nameEntry.Text)
			if folderName == "" {
				dialog.ShowInformation("Error", "Please enter a folder name.", window)
				return
			}

			// check if directory already exists
			outputPath := filepath.Join(repo.WorkTree, "..", folderName)
			_, statErr := os.Stat(outputPath)

			doLoad := func() {
				loadedPath, err := gitcheckout.Load(clicked.SHA, folderName, *repo)
				if err != nil {
					dialog.ShowError(fmt.Errorf("failed to load: %w", err), window)
					return
				}
				dialog.ShowInformation("Loaded Successfully",
					fmt.Sprintf("Files loaded to:\n%s", loadedPath), window)
			}

			if statErr == nil {
				// directory exists, ask to overwrite
				dialog.ShowConfirm("Directory Exists",
					fmt.Sprintf("Folder '%s' already exists. Overwrite?", folderName),
					func(overwrite bool) {
						if !overwrite {
							return
						}
						os.RemoveAll(outputPath)
						doLoad()
					}, window)
			} else {
				doLoad()
			}

		}, window)

		loadDialog.Resize(fyne.NewSize(400, 200))

		loadDialog.Show()
	})
	loadBtn.Importance = widget.HighImportance
	loadBtnRow := container.NewHBox(layout.NewSpacer(), loadBtn, layout.NewSpacer())
	loadBtnContainer := container.NewVBox(TDMargin, loadBtnRow, TDMargin)

	fileBackground := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})
	fileBackground.StrokeColor = color.RGBA{R: 208, G: 200, B: 200, A: 255}
	fileBackground.StrokeWidth = 1
	fileBackground.CornerRadius = 8

	fileInner := container.NewBorder(header, loadBtnContainer, nil, nil, scrollFileList)
	filePanel := container.NewStack(fileBackground, container.NewPadded(fileInner))

	// Commit Details
	detailTitle := canvas.NewText("Save Details", color.White)
	detailTitle.TextSize = 20
	detailTitle.TextStyle = fyne.TextStyle{Bold: true}

	hashLabel := canvas.NewText(fmt.Sprintf("Hash: %s", clicked.SHA), color.RGBA{R: 120, G: 200, B: 255, A: 255})
	hashLabel.TextSize = 14
	hashLabel.TextStyle = fyne.TextStyle{Bold: true}

	parentStr := "none"
	if len(clicked.Parents) > 0 {
		parts := []string{}
		for _, p := range clicked.Parents {
			short := p
			// If the parent hash is longer than 5 characters, cut it down
			if len(p) > 5 {
				short = p[:5]
			}
			parts = append(parts, short)
		}
		// Join the shortened hashes together with a comma
		parentStr = strings.Join(parts, ", ")
	}
	parentsLabel := canvas.NewText(fmt.Sprintf("Parent(s): %s", parentStr), color.RGBA{R: 180, G: 180, B: 180, A: 255})
	parentsLabel.TextSize = 13

	authorLabel := canvas.NewText(fmt.Sprintf("Author: %s", clicked.Author), color.RGBA{R: 180, G: 180, B: 180, A: 255})
	authorLabel.TextSize = 13

	dateLabel := canvas.NewText(fmt.Sprintf("Date: %s", clicked.Date), color.RGBA{R: 180, G: 180, B: 180, A: 255})
	dateLabel.TextSize = 13

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	// Commit message centered
	msgLabel := widget.NewLabel(clicked.Message)
	msgLabel.TextStyle = fyne.TextStyle{Bold: true}
	msgLabel.Alignment = fyne.TextAlignCenter
	msgLabel.Wrapping = fyne.TextWrapWord

	msgContainer := container.NewPadded(msgLabel)

	detailContent := container.NewVBox(
		detailTitle,
		hashLabel,
		parentsLabel,
		authorLabel,
		dateLabel,
		heightMargin,
		msgContainer, // centered message
	)

	divider := canvas.NewRectangle(color.RGBA{R: 60, G: 60, B: 60, A: 255})
	divider.SetMinSize(fyne.NewSize(2, 0))

	detailPanel := container.NewBorder(nil, nil, divider, nil, container.NewPadded(detailContent))

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	panelSpacing := canvas.NewRectangle(color.Transparent)
	panelSpacing.SetMinSize(fyne.NewSize(20, 0)) 
	filePanelWithSpace := container.NewBorder(nil, nil, nil, panelSpacing, filePanel)

	bottomMargin := canvas.NewRectangle(color.Transparent)
	bottomMargin.SetMinSize(fyne.NewSize(0, 30))

	content := container.NewGridWithColumns(2, filePanelWithSpace, detailPanel)

	fullContent := container.NewBorder(
		container.NewPadded(headerArea), 
		bottomMargin,
		widthMargin, 
		widthMargin,
		container.NewPadded(content),
	)

	window.SetContent(fullContent)
	window.Show()
}