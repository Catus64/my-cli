package ignoredFile

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	gitcheckignore "gocmd/testfiles/GitCheckIgnore"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func updateIgnoredList(repoPath string, fileList *fyne.Container, fileListTitle *canvas.Text, LRMargin *canvas.Rectangle) {
    ignoredFiles := []string{}

    repo, err := gitpath.Repo_find(repoPath, false)
    if err == nil && repo != nil {
        rules, err := gitcheckignore.ReadGitIgnore(*repo)
        if err == nil {
            filepath.WalkDir(repoPath, func(path string, directory os.DirEntry, err error) error {
                if err != nil || directory.IsDir() {
                    return nil
                }
                if strings.Contains(filepath.ToSlash(path), "/.git/") {
                    return nil
                }
                rel, err := filepath.Rel(repo.WorkTree, path)
                if err != nil {
                    return nil
                }
                relpath := filepath.ToSlash(rel)
                ignored, _ := gitcheckignore.CheckIgnore(rules, relpath)
                if ignored {
                    ignoredFiles = append(ignoredFiles, relpath)
                }
                return nil
            })
        }
    }

    fileList.Objects = nil
    for _, file := range ignoredFiles {
        name := canvas.NewText(file, color.RGBA{R: 150, G: 180, B: 220, A: 255})
        name.TextSize = 14

        status := canvas.NewText("IGNORED", color.RGBA{R: 160, G: 160, B: 160, A: 255})
        status.TextSize = 12
        status.TextStyle = fyne.TextStyle{Bold: true}

        divider := canvas.NewRectangle(color.RGBA{R: 80, G: 80, B: 100, A: 255})
        divider.SetMinSize(fyne.NewSize(0, 1))

        row := container.NewBorder(nil, nil, container.NewHBox(LRMargin, name), container.NewHBox(status, LRMargin), nil)
        fileList.Add(container.NewPadded(row))
        fileList.Add(divider)
    }

    fileListTitle.Text = fmt.Sprintf("File List (%d)", len(ignoredFiles))
    fileListTitle.Refresh()
    fileList.Refresh()
}

func ignoredFileContent(repoPath string, window fyne.Window) fyne.CanvasObject {
	title := canvas.NewText("Ignored File", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("File(s) that will not be included in repository", color.Gray{Y: 150})
	subtitle.TextSize = 15

	LRMargin := canvas.NewRectangle(color.Transparent)
	LRMargin.SetMinSize(fyne.NewSize(10, 0))

	TDMargin := canvas.NewRectangle(color.Transparent)
	TDMargin.SetMinSize(fyne.NewSize(0, 1))

	// File list
	fileList := container.NewVBox()

	scrollableFileList := container.NewScroll(fileList)

	// Header
	fileListTitle := canvas.NewText("File List (0)", color.RGBA{R: 208, G: 200, B: 200, A: 255})
	fileListTitle.TextSize = 20
	fileListTitle.TextStyle = fyne.TextStyle{Bold: true}

	updateIgnoredList(repoPath, fileList, fileListTitle, LRMargin) 

	viewButton := widget.NewButton("View File Type", func() {
		// Read .gitignore content
		gitignorePath := filepath.Join(repoPath, ".gitignore")
		data, err := os.ReadFile(gitignorePath)
		content := ""

		if err != nil {
			dialog.ShowInformation("No Ignored File", "No ignored file found in this repository.", window)
			return
		} else {
			content = string(data)
		}

		editEntry := widget.NewMultiLineEntry()
		editEntry.SetText(content)
		editEntry.Wrapping = fyne.TextWrapWord

		scroll := container.NewScroll(editEntry)
		scroll.SetMinSize(fyne.NewSize(400, 100))

		dialog.ShowCustomConfirm("Ignored File Types", "Update", "Close", scroll, func(confirmed bool) {
			if !confirmed {
				return
			}

			err := os.WriteFile(gitignorePath, []byte(editEntry.Text), 0644)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Failed to save ignored files: %w", err), window)
				return
			}

			updateIgnoredList(repoPath, fileList, fileListTitle, LRMargin)
			dialog.ShowInformation("Updated", "Ignored files updated successfully.", window)
		}, window)
	})
	viewButton.Importance = widget.HighImportance

	headerContent := container.NewBorder(nil, nil, container.NewHBox(LRMargin, fileListTitle), viewButton, nil)

	titleLine := canvas.NewRectangle(color.RGBA{R: 208, G: 200, B: 200, A: 255})
	titleLine.SetMinSize(fyne.NewSize(0, 1))

	header := container.NewVBox(TDMargin, headerContent, TDMargin, titleLine)

	// File list box background
	background := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})
	background.StrokeColor = color.RGBA{R: 208, G: 200, B: 200, A: 255}
	background.StrokeWidth = 1
	background.CornerRadius = 8
	background.SetMinSize(fyne.NewSize(0, 500))

	content := container.NewBorder(header, nil, nil, nil, scrollableFileList)
	box := container.NewStack(background, container.NewPadded(content))

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	ignoredFileContent := container.NewVBox(heightMargin, title, subtitle, heightMargin, box)

	return container.NewBorder(nil, nil, widthMargin, widthMargin, container.NewPadded(ignoredFileContent))
}