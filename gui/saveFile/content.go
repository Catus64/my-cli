package saveFile

import (
	"fmt"
	gitaddremove "gocmd/testfiles/GitAddRemove"
	gitCurrent "gocmd/testfiles/GitCurrent"
	gitobject "gocmd/testfiles/GitObject"
	gitsave "gocmd/testfiles/GitSave"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"image/color"
	"net/mail"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Helper to get repo and index
func getRepoAndIndex(repoPath string) (*gitpath.GitRepository, *gitobject.GitIndex) {
	repo := gitpath.MakeRepo(repoPath, false)
	if repo == nil {
		fmt.Println("No repo found at:", repoPath)
		return nil, nil
	}

	index, err := gitobject.Index_Read2(*repo)
	if err != nil || index == nil {
		fmt.Println("Failed to read index:", err)
		return repo, nil
	}

	return repo, index
}

// Load files from index
func getSaveListFiles(repoPath string) []string {
	var result []string

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in getSaveListFiles:", r)
		}
	}()

	repo, index := getRepoAndIndex(repoPath)
	if repo == nil || index == nil {
		return result
	}

	// Show ALL index entries
	for _, entry := range index.Entries {
		result = append(result, entry.Name)
		println("File: ", entry.Name, " Mode: ", entry.ModePerms, " SHA: ", entry.SHA)
	}

	fmt.Println("Staged files:", len(result))
	return result
}

func getStagedFiles(repoPath string) []string {
	var result []string

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in getStagedFiles:", r)
		}
	}()

	repo, index := getRepoAndIndex(repoPath)
	if repo == nil || index == nil {
		return result
	}

	// Trigger when zero commit in the repo
	headResolved, headErr := gitobject.Ref_Resolve(*repo, "HEAD")
	if headErr != nil || headResolved == nil {
		for _, entry := range index.Entries {
			result = append(result, entry.Name)
		}
		return result
	}

	// Compare index vs HEAD
	head := make(map[string]string)
	treeSHA, treeErr := gitCurrent.Get_Tree_SHA(*repo, "HEAD")
	if treeErr == nil {
		gitCurrent.TreeToMap(*repo, treeSHA, "", head)
	}

	for _, entry := range index.Entries {
		normalizedName := filepath.ToSlash(entry.Name)
		if headSHA, exists := head[normalizedName]; exists {
			if headSHA != entry.SHA {
				result = append(result, entry.Name)
			}
		} else {
			result = append(result, entry.Name)
		}
	}

	fmt.Println("Files ready to commit:", len(result))
	return result
}

func saveListBox(file *[]string, repo *gitpath.GitRepository, window fyne.Window, onFileChange func()) (fyne.CanvasObject, func()) {
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

	// Add scroll
	scrollableFileList := container.NewScroll(fileList)
	scrollableFileList.Direction = container.ScrollBoth

	// Track checked files
	checkedFiles := map[string]bool{}

	var updateFunction func()

	removeButton := widget.NewButton("Remove", func() {
		var selectedFiles []string
		for name, checked := range checkedFiles {
			if checked {
				absolutePath := filepath.Join(repo.WorkTree, name)
				selectedFiles = append(selectedFiles, absolutePath)
			}
		}

		if len(selectedFiles) == 0 {
			dialog.ShowInformation("No Files Selected", "Please select at least one file to remove.", window)
			return
		}

		_, err := gitaddremove.Remove(repo, selectedFiles, gitaddremove.RemoveOptions{
			Delete:          false,
			SkipMissingFile: false,
		})

		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowInformation(
			"Success",
			fmt.Sprintf("%d file(s) removed from save list!", len(selectedFiles)),
			window,
		)

		// Reload files after removing
		newIndex, err := gitobject.Index_Read2(*repo)
		if err == nil && newIndex != nil {
			*file = []string{}
			for _, entry := range newIndex.Entries {
				*file = append(*file, entry.Name)
			}
		}

		if updateFunction != nil {
			updateFunction()
		}

		if onFileChange != nil {
			onFileChange()
		}
	})
	removeButton.Importance = widget.DangerImportance

	removeButtonRow := container.NewHBox(layout.NewSpacer(), removeButton, layout.NewSpacer())
	removeBtn := container.NewVBox(removeButtonRow, TDMargin)

	background := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})
	background.StrokeColor = color.RGBA{R: 208, G: 200, B: 200, A: 255}
	background.StrokeWidth = 1
	background.CornerRadius = 8
	background.SetMinSize(fyne.NewSize(0, 425))

	content := container.NewBorder(saveListHeader, removeBtn, nil, nil, scrollableFileList)

	box := container.NewStack(background, container.NewPadded(content))

	// Update file
	update := func() {
		saveListTitle.Text = fmt.Sprintf("Save List (%d)", len(*file))
		saveListTitle.Refresh()

		fileList.Objects = nil
		checkedFiles = map[string]bool{}

		for _, files := range *file {
			checkedFiles[files] = false
			checkbox := widget.NewCheck("", func(checked bool) {
				checkedFiles[files] = checked
			})

			fileName := canvas.NewText(files, color.White)
			fileName.TextSize = 14

			fileScroll := container.NewHScroll(fileName)

			row := container.NewBorder(nil, nil, checkbox, LRMargin, fileScroll)
			fileList.Add(row)
		}
		fileList.Refresh()
		scrollableFileList.Refresh()
	}

	updateFunction = update

	return box, update
}

func commitBox(repoPath string, window fyne.Window, onFileChange func()) fyne.CanvasObject {
	commitMessageEntry := widget.NewMultiLineEntry()
	commitMessageEntry.Wrapping = fyne.TextWrapWord

	placeholder := canvas.NewText("Version Message (required)", color.Gray{Y: 150})
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
	commitBackground.SetMinSize(fyne.NewSize(0, 270))

	commitBox := container.NewStack(commitBackground, container.NewPadded(commitMessageEntry), placeholderPosition)

	doSave := func(repo *gitpath.GitRepository, index *gitobject.GitIndex, treeSHA string, parents []string, author string, message string) {
		// Create Commit
		commitSHA, err := gitsave.Version_Create(*repo, treeSHA, parents, author, time.Now(), message)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to create version: %w", err), window)
			return
		}

		// Update branch ref
		branchName, err := gitsave.Update_Branch_Ref(*repo, commitSHA)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to update branch: %w", err), window)
			return
		}

		// Refresh Index
		err = gitsave.RefreshIndex(*repo, index)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to refresh index: %w", err), window)
			return
		}

		// Clear message box after successful save
		commitMessageEntry.SetText("")
		placeholder.Show()
		placeholder.Refresh()

		dialog.ShowInformation(
			"Saved Successfully",
			fmt.Sprintf("Version saved on '%s'\n%s", branchName, commitSHA[:7]),
			window,
		)

		if onFileChange != nil {
			onFileChange()
		}
	}

	setupConfigForm := func(repo *gitpath.GitRepository, index *gitobject.GitIndex, treeSHA string, parents []string, message string) {
		nameEntry := widget.NewEntry()
		nameEntry.SetPlaceHolder("Your Name")

		emailEntry := widget.NewEntry()
		emailEntry.SetPlaceHolder("Your Email")

		formItems := []*widget.FormItem{
			widget.NewFormItem("Name", nameEntry),
			widget.NewFormItem("Email", emailEntry),
		}

		dialog.ShowForm(
			"Setup Required",
			"Save", "Cancel",
			formItems,
			func(submitted bool) {
				if !submitted {
					return
				}

				name := strings.TrimSpace(nameEntry.Text)
				email := strings.ToLower(strings.TrimSpace(emailEntry.Text))

				if name == "" || email == "" {
					dialog.ShowInformation("Required", "Please enter both name and email.", window)
					return
				}

				// Validate email format
				_, err := mail.ParseAddress(email)
				if err != nil {
					dialog.ShowInformation("Invalid Email", "Please enter a valid email address.", window)
					return
				}

				// Save to disk so next time it loads automatically
				newConfig := &gitpath.EzGitConfig{Name: name, Email: email}
				if err := gitpath.Save(newConfig); err != nil {
					dialog.ShowError(fmt.Errorf("failed to save config: %w", err), window)
					return
				}

				doSave(repo, index, treeSHA, parents, newConfig.Format(), message)
			},
			window,
		)
	}

	saveButton := widget.NewButton("Save", func() {
		message := strings.TrimSpace(commitMessageEntry.Text)
		if message == "" {
			dialog.ShowInformation("Message Required", "Please enter a version message before saving.", window)
			return
		}

		// Find repo
		repo, err := gitpath.Repo_find(repoPath, true)
		if err != nil || repo == nil {
			dialog.ShowError(fmt.Errorf("could not find repository: %w", err), window)
			return
		}

		// Read index
		index, err := gitobject.Index_Read2(*repo)
		if err != nil || index == nil {
			dialog.ShowError(fmt.Errorf("could not read index: %w", err), window)
			return
		}

		// Check if there is anything to commit
		headResult, err := gitCurrent.StatusHeadIndex(*repo, *index)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if !headResult.HasChanges() {
			dialog.ShowInformation("Nothing to Save", "Your save list is already up to date.", window)
			return
		}

		// Build tree from index
		treeSHA, err := gitobject.TreeFromIndex(*repo, *index)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to build tree: %w", err), window)
			return
		}

		// Get parent commit SHA
		var parents []string
		parentSHA, err := gitobject.Ref_Resolve(*repo, "HEAD")
		if err == nil && parentSHA != nil {
			parents = []string{*parentSHA}
		}

		// Get author from config
		userConfig, err := gitpath.Load()
		if err != nil {
			// Config not found — show GUI form to collect name and email
			setupConfigForm(repo, index, treeSHA, parents, message)
			return
		}
		// Config exists, commit directly
		doSave(repo, index, treeSHA, parents, userConfig.Format(), message)
	})
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

func readySaveList(file *[]string) (fyne.CanvasObject, func()) {
	previewTitle := canvas.NewText("File to Save", color.White)
	previewTitle.TextSize = 16
	previewTitle.TextStyle = fyne.TextStyle{Bold: true}

	previewSubTitle := canvas.NewText(fmt.Sprintf("%d file(s) ready to be saved.", len(*file)), color.Gray{Y: 150})
	previewSubTitle.TextSize = 13

	previewList := container.NewVBox()

	// Add scroll
	scrollablePreviewList := container.NewScroll(previewList)
	scrollablePreviewList.Direction = container.ScrollBoth

	background := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})
	background.StrokeColor = color.RGBA{R: 208, G: 200, B: 200, A: 255}
	background.StrokeWidth = 1
	background.CornerRadius = 8
	background.SetMinSize(fyne.NewSize(0, 160))

	leftMargin := canvas.NewRectangle(color.Transparent)
	leftMargin.SetMinSize(fyne.NewSize(10, 0))

	topMargin := canvas.NewRectangle(color.Transparent)
	topMargin.SetMinSize(fyne.NewSize(0, 5))

	downMargin := canvas.NewRectangle(color.Transparent)
	downMargin.SetMinSize(fyne.NewSize(0, 3))

	content := container.NewVBox(topMargin, previewTitle, previewSubTitle, downMargin)
	fullcontent := container.NewHBox(leftMargin, content)
	box := container.NewStack(background, container.NewPadded(container.NewBorder(fullcontent, nil, nil, nil, scrollablePreviewList)))

	// Update file
	update := func() {
		previewSubTitle.Text = fmt.Sprintf("%d file(s) ready to be saved.", len(*file))
		previewSubTitle.Refresh()

		previewList.Objects = nil
		for _, files := range *file {
			bullet := canvas.NewText("    •  "+files, color.Gray{Y: 200})
			bullet.TextSize = 12
			previewList.Add(bullet)
		}
		previewList.Refresh()
		scrollablePreviewList.Refresh()
	}

	return box, update
}

func SaveFileContent(repoPath string, window fyne.Window) fyne.CanvasObject {
	title := canvas.NewText("Save File", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Manage your save list", color.Gray{Y: 150})
	subtitle.TextSize = 15

	repo := gitpath.MakeRepo(repoPath, false)

	// All index files for save list box
	allFiles := getSaveListFiles(repoPath)

	// Only staged files for staged box
	stagedFiles := getStagedFiles(repoPath)

	var updateStagedFiles func()

	onFileChange := func() {
		// Reload staged files from index and refresh the readySaveBox
		stagedFiles = getStagedFiles(repoPath)
		if updateStagedFiles != nil {
			updateStagedFiles()
		}
	}

	saveListBox, updateSaveList := saveListBox(&allFiles, repo, window, onFileChange)
	readySaveBox, updateStagedFile := readySaveList(&stagedFiles)
	updateStagedFiles = updateStagedFile

	updateSaveList()
	updateStagedFiles()

	// Commit Box
	commitBox := commitBox(repoPath, window, onFileChange)

	columnGap := canvas.NewRectangle(color.Transparent)
	columnGap.SetMinSize(fyne.NewSize(10, 0))

	rowGap := canvas.NewRectangle(color.Transparent)
	rowGap.SetMinSize(fyne.NewSize(0, 10))

	rightColumn := container.NewBorder(nil, container.NewVBox(rowGap, readySaveBox), nil, nil, commitBox)

	paddedRight := container.NewBorder(nil, nil, columnGap, nil, rightColumn)

	content := container.NewGridWithColumns(2, saveListBox, paddedRight)

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	homeContent := container.NewVBox(heightMargin, title, subtitle, heightMargin, content)

	return container.NewBorder(nil, nil, widthMargin, widthMargin, container.NewPadded(homeContent))
}