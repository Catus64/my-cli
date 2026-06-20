package branchPage

import (
	gitpath "gocmd/testfiles/Gitrepostruct"
	alternateversions "gocmd/testfiles/alternateVersions"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func showMergeMessageDialog(repo *gitpath.GitRepository, targetBranch string, resolutions map[string]string, window fyne.Window, onResolved func()) {
	messageEntry := widget.NewMultiLineEntry()
	messageEntry.SetPlaceHolder("Merge save message")
	messageEntry.SetText("Merge branch '" + targetBranch + "'")
	messageEntry.Wrapping = fyne.TextWrapWord

	sizedMsgEntry := container.NewGridWrap(fyne.NewSize(400, 100), messageEntry)

	formItems := []*widget.FormItem{
		widget.NewFormItem("Message", sizedMsgEntry),
	}
	dialog.ShowForm("Merge Save Message", "Save", "Cancel", formItems, func(submitted bool) {
		if !submitted || strings.TrimSpace(messageEntry.Text) == "" {
			// user cancelled — nothing was ever written to disk, so nothing to undo
			return
		}
		if repo != nil {
			for filePath, content := range resolutions {
				if err := alternateversions.ApplyConflictResolution(*repo, filePath, content); err != nil {
					dialog.ShowError(err, window)
					return
				}
			}
			err := alternateversions.CompleteMerge(*repo, targetBranch, strings.TrimSpace(messageEntry.Text), window)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
		}
		info := dialog.NewInformation("Merge Complete", "Merged successfully. All conflicts resolved.", window)
		info.SetOnClosed(func() {
			window.Close()
			if onResolved != nil {
				onResolved()
			}
		})
		info.Show()
	}, window)
}

func showMergeConflictWindow(app fyne.App, conflicts []alternateversions.MergeConflict, repo *gitpath.GitRepository, currentBranch string, targetBranch string, onResolved func()) {
	if len(conflicts) == 0 {
		return
	}

	conflictIndex := 0
	resolutions := make(map[string]string)
	var window fyne.Window
	window = app.NewWindow("Merge Conflict")
	window.Resize(fyne.NewSize(1000, 600))

	var renderConflict func(index int) fyne.CanvasObject
	renderConflict = func(index int) fyne.CanvasObject {
		conflict := conflicts[index]

		titleText := canvas.NewText("Merge Conflict", color.White)
		titleText.TextSize = 40
		titleText.TextStyle = fyne.TextStyle{Bold: true}

		subtitleText := canvas.NewText("Select the version you want", color.Gray{Y: 150})
		subtitleText.TextSize = 15

		fileText := canvas.NewText("Conflict File: "+conflict.FilePath, color.RGBA{R: 120, G: 200, B: 255, A: 255})
		fileText.TextSize = 13

		widthMargin := canvas.NewRectangle(color.Transparent)
		widthMargin.SetMinSize(fyne.NewSize(30, 0))

		headerBlock := container.NewBorder(nil, nil, widthMargin, nil,
			container.NewVBox(
				container.NewPadded(titleText),
				container.NewPadded(subtitleText),
				container.NewPadded(fileText),
			),
		)

		// Left panel — current branch
		currentLabel := canvas.NewText(currentBranch, color.Black)
		currentLabel.TextSize = 22
		currentLabel.TextStyle = fyne.TextStyle{Bold: true}

		currentUnderline := canvas.NewRectangle(color.Black)
		currentUnderline.SetMinSize(fyne.NewSize(200, 2))

		currentContent := widget.NewMultiLineEntry()
		currentContent.SetText(conflict.CurrentContent)
		currentContent.Wrapping = fyne.TextWrapWord
		currentContent.Disable()

		leftMargin := canvas.NewRectangle(color.Transparent)
		leftMargin.SetMinSize(fyne.NewSize(10, 0))

		currentHeader := container.NewBorder(nil, currentUnderline, leftMargin, nil, currentLabel)

		chooseCurrentBtn := widget.NewButton("Use "+currentBranch, func() {
			resolutions[conflict.FilePath] = conflict.CurrentContent
			if index+1 < len(conflicts) {
				conflictIndex++
				window.SetContent(renderConflict(conflictIndex))
			} else {
				showMergeMessageDialog(repo, targetBranch, resolutions, window, onResolved)
			}
		})
		chooseCurrentBtn.Importance = widget.HighImportance
		currentBtnRow := container.NewHBox(layout.NewSpacer(), chooseCurrentBtn, layout.NewSpacer())

		btnMargin := canvas.NewRectangle(color.Transparent)
		btnMargin.SetMinSize(fyne.NewSize(0, 10))

		currentBg := canvas.NewRectangle(color.White)
		currentBg.CornerRadius = 12
		currentInner := container.NewBorder(
			currentHeader,
			container.NewVBox(btnMargin, currentBtnRow, btnMargin),
			nil, nil,
			container.NewPadded(currentContent),
		)
		currentBox := container.NewStack(currentBg, container.NewPadded(currentInner))

		// Right panel — target branch
		targetLabel := canvas.NewText(targetBranch, color.Black)
		targetLabel.TextSize = 22
		targetLabel.TextStyle = fyne.TextStyle{Bold: true}

		targetUnderline := canvas.NewRectangle(color.Black)
		targetUnderline.SetMinSize(fyne.NewSize(200, 2))

		targetContent := widget.NewMultiLineEntry()
		targetContent.SetText(conflict.IncomingContent)
		targetContent.Wrapping = fyne.TextWrapWord
		targetContent.Disable()

		targetHeader := container.NewBorder(nil, targetUnderline, leftMargin, nil, targetLabel)

		chooseTargetBtn := widget.NewButton("Use "+targetBranch, func() {
			resolutions[conflict.FilePath] = conflict.IncomingContent
			if index+1 < len(conflicts) {
				conflictIndex++
				window.SetContent(renderConflict(conflictIndex))
			} else {
				showMergeMessageDialog(repo, targetBranch, resolutions, window, onResolved)
			}
		})
		chooseTargetBtn.Importance = widget.HighImportance
		targetBtnRow := container.NewHBox(layout.NewSpacer(), chooseTargetBtn, layout.NewSpacer())

		targetBg := canvas.NewRectangle(color.White)
		targetBg.CornerRadius = 12
		targetInner := container.NewBorder(
			targetHeader,
			container.NewVBox(btnMargin, targetBtnRow, btnMargin),
			nil, nil,
			container.NewPadded(targetContent),
		)
		targetBox := container.NewStack(targetBg, container.NewPadded(targetInner))

		panels := container.NewGridWithColumns(2, currentBox, targetBox)

		heightMargin := canvas.NewRectangle(color.Transparent)
		heightMargin.SetMinSize(fyne.NewSize(0, 20))

		paddedPanels := container.NewBorder(
			nil,
			heightMargin,
			widthMargin,
			widthMargin,
			panels,
		)

		return container.NewBorder(
			container.NewVBox(
				heightMargin,
				headerBlock,
				heightMargin,
			),
			nil, nil, nil,
			paddedPanels,
		)
	}

	window.SetContent(renderConflict(conflictIndex))
	window.Show()
}