package modifiedFile

import (
	"fmt"
	gitaddremove "gocmd/testfiles/GitAddRemove"
	gitCurrent "gocmd/testfiles/GitCurrent"
	gitobject "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FileStatus struct {
	Name   string // File Name
	Status string // "MODIFIED" or "ADDED" or "DELETED"
}

func getFileStatuses(repoPath string) ([]FileStatus, *gitpath.GitRepository) {
	var result []FileStatus

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in getFileStatuses:", r)
		}
	}()

	repo, err := gitpath.Repo_find(repoPath, false)
	if err != nil || repo == nil {
		fmt.Println("No repo found at:", repoPath)
		return result, nil
	}

	index, err := gitobject.Index_Read2(*repo)
	if err != nil {
		fmt.Println("Failed to read index:", err)
		return result, repo
	}

	// Build a set of files staged as DELETED
	deletedFromIndex := map[string]bool{}
	headStatus, headErr := gitCurrent.StatusHeadIndex(*repo, *index)
	if headErr == nil {
		for _, f := range headStatus.Deleted {
			normalized := filepath.ToSlash(f)
			deletedFromIndex[normalized] = true
			result = append(result, FileStatus{Name: f, Status: "DELETED"})
		}
	}

	// Check worktree vs index
	status_files, err := gitCurrent.StatusIndexWorktree(*repo, *index)
	if err == nil {
		for _, f := range status_files.Modified {
			result = append(result, FileStatus{Name: f, Status: "MODIFIED"})
		}
		for _, f := range status_files.Untracked {
			normalized := filepath.ToSlash(f)
			// Skip files already shown as DELETED
			if deletedFromIndex[normalized] {
				continue
			}
			result = append(result, FileStatus{Name: f, Status: "ADDED"})
		}
	}

	return result, repo
}

func reloadFiles(repo *gitpath.GitRepository, files *[]FileStatus) {
    index, err := gitobject.Index_Read2(*repo)
    if err != nil || index == nil {
        return
    }

    *files = []FileStatus{}

    headStatus, headErr := gitCurrent.StatusHeadIndex(*repo, *index)
    if headErr == nil {
        for _, f := range headStatus.Deleted {
            *files = append(*files, FileStatus{Name: f, Status: "DELETED"})
        }
    }

    status_files, err := gitCurrent.StatusIndexWorktree(*repo, *index)
    if err == nil {
        for _, f := range status_files.Modified {
            *files = append(*files, FileStatus{Name: f, Status: "MODIFIED"})
        }
        for _, f := range status_files.Untracked {
            *files = append(*files, FileStatus{Name: f, Status: "ADDED"})
        }
    }
}

func FileDirectory(repoPath string, window fyne.Window) fyne.CanvasObject {
	title := canvas.NewText("File Directory", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Manage your file in your directory", color.Gray{Y: 150})
	subtitle.TextSize = 15

	// Load files from git
	files, repo := getFileStatuses(repoPath)

	modifiedBox, updateModifiedList := modifiedListBox(&files, repo, window)
	updateModifiedList()

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	modifiedFileContent := container.NewVBox(heightMargin, title, subtitle, heightMargin, modifiedBox)

	return container.NewBorder(nil, nil, widthMargin, widthMargin, container.NewPadded(modifiedFileContent))
}

func modifiedListBox(files *[]FileStatus, repo *gitpath.GitRepository, window fyne.Window) (fyne.CanvasObject, func()) {
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

	// Checked Files
	checkedFiles := map[string]bool{}

	var updateFunction func()

	addButton := widget.NewButton("Add", func() {
		if repo == nil {
			fmt.Println("No repo available")
			return
		}

		var selectedFiles []string
		for name, checked := range checkedFiles {
			if checked {
				absolutePath := filepath.Join(repo.WorkTree, name)
				selectedFiles = append(selectedFiles, absolutePath)
			}
		}

		if len(selectedFiles) == 0 {
			dialog.ShowInformation("No Files Selected", "Please select at least one file to add.", window)
			return
		}

		_, err := gitaddremove.Add(repo, selectedFiles, gitaddremove.Options{All: false})
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		reloadFiles(repo, files)
		if updateFunction != nil {
			updateFunction() // refresh UI
		}

		dialog.ShowInformation(
        "Success",
        fmt.Sprintf("%d file(s) added to save list!", len(selectedFiles)),
        window,
    )
	})
	addButton.Importance = widget.HighImportance

	ignoreButton := widget.NewButton("Ignore", func() {
		if repo == nil {
			fmt.Println("No repo available")
			return
		}

		var selectedFiles []string
		var cannotIgnore []string

		for name, checked := range checkedFiles {
			if !checked {
				continue
			}
			// find the file status
			for _, file := range *files {
				if file.Name == name {
					if file.Status == "MODIFIED" {
						cannotIgnore = append(cannotIgnore, name) 
					} else {
						selectedFiles = append(selectedFiles, name) 
					}
				}
			}
		}

		gitignorePath := filepath.Join(repo.WorkTree, ".gitignore")
		file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to open .gitignore: %w", err), window)
			return
		}
		defer file.Close()

		for _, name := range selectedFiles {
			fmt.Fprintln(file, name)  // write each file name as a new line
		}

		if len(cannotIgnore) > 0 {
			dialog.ShowInformation("Ignore Failed",
				fmt.Sprintf("Modified file(s) cannot be ignored:\n%s\n\nRemove them from save list first.",
					strings.Join(cannotIgnore, "\n")),
				window)
		}

		if len(selectedFiles) == 0 {
			dialog.ShowInformation("No Files Selected", "Please select at least one file to ignore.", window)
			return
		}

		reloadFiles(repo, files)
		if updateFunction != nil {
			updateFunction()
		}

		dialog.ShowInformation("Ignored", fmt.Sprintf("%d file(s) added as ignored file(s)!", len(selectedFiles)), window)
	})
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
	background.SetMinSize(fyne.NewSize(0, 500))

	content := container.NewBorder(modifiedListHeader, button, nil, nil, scrollableFileList)
	box := container.NewStack(background, container.NewPadded(content))

	update := func() {
		modifiedListTitle.Text = fmt.Sprintf("File List (%d)", len(*files))
		modifiedListTitle.Refresh()

		fileList.Objects = nil
		checkedFiles = map[string]bool{} // reset checked files

		for _, file := range *files {
			checkedFiles[file.Name] = false
			checkbox := widget.NewCheck("", func(checked bool) {
				checkedFiles[file.Name] = checked
			})

			fileName := canvas.NewText(file.Name, color.White)
			fileName.TextSize = 14

			fileScroll := container.NewHScroll(fileName)

			row := container.NewBorder(
				nil,
				nil,
				checkbox,
				nil,
				fileScroll,
			)

			// Different color for ADDED vs MODIFIED
			var statusColor color.Color
			switch file.Status {
			case "ADDED":
				statusColor = color.RGBA{R: 100, G: 200, B: 100, A: 255} // green
			case "MODIFIED":
				statusColor = color.RGBA{R: 255, G: 200, B: 0, A: 255} // yellow
			case "DELETED":
				statusColor = color.RGBA{R: 220, G: 50, B: 50, A: 255} // red
			default:
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

	updateFunction = update

	return box, update
}
