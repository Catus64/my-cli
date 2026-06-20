package branchPage

import (
	"fmt"
	"strings"
	gitCurrent "gocmd/testfiles/GitCurrent"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	alternateversions "gocmd/testfiles/alternateVersions"
	githashread "gocmd/testfiles/GitHashRead"
	gitlog "gocmd/testfiles/GitLog"
	"image/color"
	"io/fs"
	"path/filepath"
	"gocmd/gui/history"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type branchRow struct {
	name      string
	latestSHA string
	isCurrent bool
}

func loadBranches(repoPath string) ([]branchRow, *gitpath.GitRepository) {
	repo, err := gitpath.Repo_find(repoPath, false)
	if err != nil || repo == nil {
		return nil, nil
	}

	activeBranch, _ := gitCurrent.Get_Active_Branch(*repo)
	headsDirectory := gitpath.Repo_Path(*repo, "refs", "heads")

	var rows []branchRow
	filepath.WalkDir(headsDirectory, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return nil
		}
		name, err := filepath.Rel(headsDirectory, path)
		if err != nil {
			return nil
		}
		name = filepath.ToSlash(name)

		sha, err := gitobj.Ref_Resolve(*repo, "refs/heads/"+name)
		shortSHA := ""
		if err == nil && sha != nil && len(*sha) >= 5 {
			shortSHA = (*sha)[:5]
		}

		rows = append(rows, branchRow{
			name:      name,
			latestSHA: shortSHA,
			isCurrent: name == activeBranch,
		})
		return nil
	})
	
	return rows, repo
}	

func branchContent(repoPath string, window fyne.Window, app fyne.App) fyne.CanvasObject {
	title := canvas.NewText("Save File", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Manage your save file", color.Gray{Y: 150})
	subtitle.TextSize = 15

	// Table
	headerbackground := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})

	makeHeaderCell := func(text string, align fyne.TextAlign) fyne.CanvasObject {
		t := canvas.NewText(text, color.White)
		t.TextSize = 15
		t.TextStyle = fyne.TextStyle{Bold: true}
		t.Alignment = align
		return container.NewPadded(t)
	}

	headerRow := container.NewGridWithColumns(4,
		makeHeaderCell("No", fyne.TextAlignCenter),
		makeHeaderCell("Name", fyne.TextAlignLeading),
		makeHeaderCell("Latest Save", fyne.TextAlignCenter),
		makeHeaderCell("Status", fyne.TextAlignCenter),
	)
	header := container.NewStack(headerbackground, headerRow)

	selectedBranch := ""

	rows, repo := loadBranches(repoPath)

	showBranchHistory := func(branchName string, branchRepo *gitpath.GitRepository) {
		if branchRepo == nil {
			return
		}

		branchSHA, err := gitobj.Ref_Resolve(*branchRepo, "refs/heads/"+branchName)
		if err != nil || branchSHA == nil {
			return
		}

		var commits []history.Commit
		seen := map[string]bool{}
		queue := []string{*branchSHA}

		for len(queue) > 0 {
			sha := queue[0]
			queue = queue[1:]
			if seen[sha] || sha == "" {
				continue
			}
			seen[sha] = true

			obj, err := githashread.Object_Read(*branchRepo, sha)
			if err != nil {
				continue
			}
			gitCommit, ok := obj.(*gitobj.GitCommit)
			if !ok {
				continue
			}
			gitCommit.Deserialize()

			date, author := gitlog.Format_Date_Author(string(gitCommit.KvlmDict.Dict["author"]))
			treeSHA := strings.TrimSpace(string(gitCommit.KvlmDict.Dict["tree"]))
			parents := gitobj.GetKvlmValues(gitCommit.KvlmDict.Dict, "parent")
			message := strings.TrimSpace(string(gitCommit.KvlmDict.Dict["data"]))

			commits = append(commits, history.Commit{
				SHA:      sha,
				ShortSHA: sha[:5],
				Author:   author,
				Date:     date,
				Message:  message,
				TreeSHA:  treeSHA,
				Parents:  parents,
			})

			for _, p := range parents {
				if !seen[p] {
					queue = append(queue, p)
				}
			}
		}

		if len(commits) == 0 {
			dialog.ShowInformation("No History", "No commits found for: "+branchName, window)
			return
		}

		historyWindow := app.NewWindow("History")
		historyWindow.Resize(fyne.NewSize(900, 300))

		hCanvas := history.NewMainChainCanvas(commits, func(c history.Commit) {
			history.ShowCommitDetailWindow(app, c)
		})

		titleText := canvas.NewText("History: "+branchName, color.White)
		titleText.TextSize = 22
		titleText.TextStyle = fyne.TextStyle{Bold: true}

		graphScroll := container.NewHScroll(hCanvas)

		margin := canvas.NewRectangle(color.Transparent)
		margin.SetMinSize(fyne.NewSize(0, 20))

		content := container.NewBorder(
			container.NewVBox(margin, container.NewPadded(titleText), margin),
			nil, nil, nil,
			container.NewPadded(graphScroll),
		)

		historyWindow.SetContent(content)
		historyWindow.Show()
	}

	tableRows := container.NewVBox()
	var rowObjects []fyne.CanvasObject
	selectedBg := make([]*canvas.Rectangle, len(rows))

	rebuildTable := func(newRows []branchRow, _ *gitpath.GitRepository) {
		tableRows.Objects = nil
		rowObjects = nil
		selectedBg = make([]*canvas.Rectangle, len(newRows))
		
		for i, row := range newRows {
			bg := canvas.NewRectangle(color.Transparent)
			selectedBg[i] = bg

			divider := canvas.NewRectangle(color.RGBA{R: 60, G: 70, B: 90, A: 255})
			divider.SetMinSize(fyne.NewSize(0, 1))

			noText := canvas.NewText(fmt.Sprintf("%d", i+1), color.RGBA{R: 180, G: 180, B: 180, A: 255})
			noText.TextSize = 14
			noText.Alignment = fyne.TextAlignCenter

			nameText := canvas.NewText(row.name, color.White)
			nameText.TextSize = 14
			nameText.Alignment = fyne.TextAlignLeading

			shaText := canvas.NewText(row.latestSHA, color.White)
			shaText.TextSize = 14
			shaText.Alignment = fyne.TextAlignCenter

			var statusBranch fyne.CanvasObject
			if row.isCurrent {
				statusText := canvas.NewText("CURRENT", color.RGBA{R: 100, G: 220, B: 100, A: 255})
				statusText.TextSize = 14
				statusText.TextStyle = fyne.TextStyle{Bold: true}
				statusText.Alignment = fyne.TextAlignCenter
				statusBranch = container.NewPadded(statusText)
			} else {
				statusBranch = canvas.NewText("", color.White)
				statusBranch.(*canvas.Text).Alignment = fyne.TextAlignCenter
			}

			tableContent := container.NewGridWithColumns(4,
				container.NewPadded(noText),
				container.NewPadded(nameText),
				container.NewPadded(shaText),
				statusBranch,
			)

			selectedTarget := widget.NewButton("", func() {
				if row.isCurrent {
					return
				}
				selectedBranch = row.name
				for j, bg := range selectedBg {
					if j == i {
						bg.FillColor = color.RGBA{R: 50, G: 100, B: 150, A: 255}
					} else {
						bg.FillColor = color.Transparent
					}
					bg.Refresh()
				}
				showBranchHistory(row.name, repo)
			})
			selectedTarget.Importance = widget.LowImportance

			rowObject := container.NewStack(bg, selectedTarget, tableContent)
			tableRows.Add(divider)
			tableRows.Add(rowObject)
			rowObjects = append(rowObjects, rowObject)
		}
		tableRows.Refresh()
	}
	rebuildTable(rows, repo)

	tableBackground := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})
	tableBackground.StrokeColor = color.RGBA{R: 208, G: 200, B: 200, A: 255}
	tableBackground.StrokeWidth = 1
	tableBackground.CornerRadius = 8

	tableScroll := container.NewScroll(tableRows)
	tableScroll.Direction = container.ScrollBoth

	headerDivider := canvas.NewRectangle(color.RGBA{R: 208, G: 200, B: 200, A: 255})
	headerDivider.SetMinSize(fyne.NewSize(0, 1))

	tableContent := container.NewBorder(container.NewVBox(header, headerDivider), nil, nil, nil, tableScroll)
	tableBox := container.NewStack(tableBackground, container.NewPadded(tableContent))

	addButton := widget.NewButton("Add", func() {
		if repo == nil {
			dialog.ShowInformation("Error", "No repository found.", window)
			return
		}
		nameEntry := widget.NewEntry()
		nameEntry.SetPlaceHolder("New save file name")
		sizedEntry := container.NewGridWrap(fyne.NewSize(300, nameEntry.MinSize().Height), nameEntry)
		formItems := []*widget.FormItem{
			widget.NewFormItem("Name", sizedEntry),
		}
		dialog.ShowForm("Add Save File", "Add", "Cancel", formItems, func(submitted bool) {
			if !submitted || nameEntry.Text == "" {
				return
			}
			err := alternateversions.CreateAltVer(*repo, nameEntry.Text)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			dialog.ShowInformation("Success", fmt.Sprintf("Save file '%s' created.", nameEntry.Text), window)
			newRows, newRepo := loadBranches(repoPath)
			rows = newRows
			repo = newRepo
			selectedBranch = ""
			rebuildTable(rows, repo)
		}, window)
	})
	addButton.Importance = widget.HighImportance

	doSwitch := func(branchName string) {
		err := alternateversions.SwitchAltVer(*repo, branchName)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		dialog.ShowInformation("Switched", fmt.Sprintf("Switched to '%s'.", branchName), window)
		selectedBranch = ""
		newRows, newRepo := loadBranches(repoPath)
		rows = newRows
		repo = newRepo
		rebuildTable(rows, repo)
	}

	switchButton := widget.NewButton("Switch", func() {
		if selectedBranch == "" {
			dialog.ShowInformation("No Selection", "Please select a save file to switch.", window)
			return
		}
		if repo == nil {
			dialog.ShowInformation("Error", "No repository found.", window)
			return
		}

		dirty, _ := alternateversions.IsDirty(*repo)
		if dirty {
			dialog.ShowConfirm("Unsaved Changes",
				"You have unsaved changes that will be lost if you switch. Continue anyway?",
				func(confirmed bool) {
					if !confirmed {
						return
					}
					doSwitch(selectedBranch)
				}, window)
			return
		}

		doSwitch(selectedBranch)
	})
	switchButton.Importance = widget.HighImportance

	mergeButton := widget.NewButton("Merge", func() {

	})
	mergeButton.Importance = widget.HighImportance

	buttonRow := container.NewHBox(
		layout.NewSpacer(),
		addButton,
		layout.NewSpacer(),
		switchButton,
		layout.NewSpacer(),
		mergeButton,
		layout.NewSpacer(),
	)

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	branchContent := container.NewBorder(
		container.NewVBox(heightMargin, title, subtitle, heightMargin),
		container.NewVBox(heightMargin, buttonRow, heightMargin),
		nil, nil,
		tableBox,
	)

	return container.NewBorder(nil, nil, widthMargin, widthMargin, container.NewPadded(branchContent))
}