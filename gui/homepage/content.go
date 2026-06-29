package homepage

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	githashread "gocmd/testfiles/GitHashRead"
	gitlog "gocmd/testfiles/GitLog"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type transparentButtonTheme struct{ fyne.Theme }

func (t *transparentButtonTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
    if n == theme.ColorNameButton {
        return color.Transparent  // no background
    }
    if n == theme.ColorNameForeground {
        return color.White  // force white text
    }
    if n == theme.ColorNameHover {
        return color.RGBA{R: 255, G: 255, B: 255, A: 20}  // subtle hover
    }
    return t.Theme.Color(n, v)
}
func (t *transparentButtonTheme) Font(s fyne.TextStyle) fyne.Resource { return t.Theme.Font(s) }
func (t *transparentButtonTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return t.Theme.Icon(n) }
func (t *transparentButtonTheme) Size(n fyne.ThemeSizeName) float32 { return t.Theme.Size(n) }

func HomePageContent(pathName string, window fyne.Window) fyne.CanvasObject {
	title := canvas.NewText(filepath.Base(pathName), color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("View your file(s) in repository", color.Gray{Y: 150})
	subtitle.TextSize = 15

	// Load repo
	repo, _ := gitpath.Repo_find(pathName, false)

	var rootTreeSHA string
	var commitSHA, commitAuthor, commitDate, commitMessage string

	if repo != nil {
		headSHA, err := gitobj.Ref_Resolve(*repo, "HEAD")
		if err == nil && headSHA != nil {
			commitSHA = *headSHA
			obj, err := githashread.Object_Read(*repo, *headSHA)
			if err == nil {
				obj.Deserialize()
				commit, ok := obj.(*gitobj.GitCommit)
				if ok {
					commit.Deserialize()
					rootTreeSHA = strings.TrimSpace(string(commit.KvlmDict.Dict["tree"]))
					commitDate, commitAuthor = gitlog.Format_Date_Author(string(commit.KvlmDict.Dict["author"]))
					commitMessage = strings.TrimSpace(string(commit.KvlmDict.Dict["data"]))
				}
			}
		}
	}

	// Different Used Margin
	LRMargin := canvas.NewRectangle(color.Transparent)
	LRMargin.SetMinSize(fyne.NewSize(10, 0))

	TDMargin := canvas.NewRectangle(color.Transparent)
	TDMargin.SetMinSize(fyne.NewSize(0, 1))

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	fileListTitle := canvas.NewText("File List (0)", color.RGBA{R: 208, G: 200, B: 200, A: 255})
	fileListTitle.TextSize = 20
	fileListTitle.TextStyle = fyne.TextStyle{Bold: true}

	// File list
	fileList := container.NewVBox()
	scrollableFileList := container.NewScroll(fileList)

	var treeStack []string // navigate to folder content

	backContainer := container.NewVBox() // To be used for back button later

	var showDirectory func(treeSHA string)
	showDirectory = func(treeSHA string) {
		fileList.Objects = nil
		backContainer.Objects = nil

		if repo == nil || treeSHA == "" {
			empty := canvas.NewText("No saved files found.", color.Gray{Y: 150})
			fileList.Add(empty)
			fileList.Refresh()
			return
		}

		treeObj, err := githashread.Object_Read(*repo, treeSHA)
		if err != nil {
			fileList.Refresh()
			return
		}

		treeContent := gitobj.Tree_Parse(treeObj.Deserialize())

		// update count
		fileListTitle.Text = fmt.Sprintf("File List (%d)", len(treeContent))
		fileListTitle.Refresh()

		// Folders only 
		for _, file := range treeContent {
			if !strings.HasPrefix(string(file.Mode), "04") {
				continue // skip if not folder
			}

			iconWidget := widget.NewIcon(theme.FolderIcon())
			typeText := canvas.NewText("FOLDER", color.RGBA{R: 160, G: 160, B: 160, A: 255})
			typeText.TextSize = 14

			button := &widget.Button{
				Text:       string(file.Path),
				Alignment:  widget.ButtonAlignLeading,
				OnTapped: func() {
					treeStack = append(treeStack, treeSHA)
					showDirectory(file.Sha)
				},
			}

			themedButton := container.NewThemeOverride(button, &transparentButtonTheme{theme.DefaultTheme()})
			nameWithMargin := container.NewBorder(nil, nil, LRMargin, nil, themedButton)

			row := container.NewBorder(
				nil, nil,
				container.NewHBox(LRMargin, iconWidget),
				container.NewHBox(typeText, LRMargin),
				nameWithMargin,
			)
			fileList.Add(row)
		}

		// Files only
		for _, file := range treeContent {
			if strings.HasPrefix(string(file.Mode), "04") {
				continue // skip if folder
			}

			iconWidget := widget.NewIcon(theme.FileIcon())
			typeText := canvas.NewText("FILE", color.RGBA{R: 160, G: 160, B: 160, A: 255})
			typeText.TextSize = 14

			fileBtn := &widget.Button{
				Text:      string(file.Path),
				Alignment: widget.ButtonAlignLeading,
				OnTapped: func() {
					// read file content from git object
					obj, err := githashread.Object_Read(*repo, file.Sha)
					if err != nil {
						dialog.ShowError(fmt.Errorf("failed to read file: %w", err), window)
						return
					}

					content := string(obj.Deserialize())

					// show content in a scrollable dialog
					contentLabel := widget.NewLabel(content)
					contentLabel.Wrapping = fyne.TextWrapWord

					scroll := container.NewScroll(contentLabel)
					scroll.SetMinSize(fyne.NewSize(500, 350))

					dialog.ShowCustom(string(file.Path), "Close", scroll, window)
				},
			}
			themedFileBtn := container.NewThemeOverride(fileBtn, &transparentButtonTheme{theme.DefaultTheme()})

			nameWithMargin := container.NewBorder(nil, nil, LRMargin, nil, themedFileBtn)

			row := container.NewBorder(
				nil, nil,
				container.NewHBox(LRMargin, iconWidget),
				container.NewHBox(typeText, LRMargin),
				nameWithMargin,
			)
			fileList.Add(row)
		}

		// Back button shown only when inside a subfolder
		if len(treeStack) > 0 {
			backButton := widget.NewButton("◀  Back", func() {
				if len(treeStack) > 0 {
					prev := treeStack[len(treeStack)-1]
					treeStack = treeStack[:len(treeStack)-1]
					showDirectory(prev)
				}
			})
			backContainer.Add(backButton)
		}

		backContainer.Refresh()
		fileList.Refresh()
		scrollableFileList.Refresh()
	}

	showDirectory(rootTreeSHA)

	// Current Version button
	currentVersionBtn := widget.NewButton("Current Details", func() {
		if commitSHA == "" {
			dialog.ShowInformation("No Details", "No saved files found in this repository.", window)
			return
		}

		gap := canvas.NewRectangle(color.Transparent)
		gap.SetMinSize(fyne.NewSize(0, 15))

		makeLabel := func(text string) *canvas.Text {
			t := canvas.NewText(text, color.Black)
			t.TextSize = 16
			t.TextStyle = fyne.TextStyle{Bold: true}
			return t
		}

		popupContent := container.NewVBox(
			makeLabel("Hash: "+commitSHA),
			gap,
			makeLabel("Author: "+commitAuthor),
			gap,
			makeLabel("Date: "+commitDate),
			gap,
			makeLabel("Saved Message: "+commitMessage),
		)

		dialog.ShowCustom("Current Details", "Close", container.NewPadded(popupContent), window)
	})
	currentVersionBtn.Importance = widget.HighImportance

	headerContent := container.NewBorder(nil, nil,
		container.NewHBox(LRMargin, fileListTitle),
		currentVersionBtn,
		nil,
	)

	titleLine := canvas.NewRectangle(color.RGBA{R: 208, G: 200, B: 200, A: 255})
	titleLine.SetMinSize(fyne.NewSize(0, 1))

	header := container.NewVBox(TDMargin, headerContent, TDMargin, titleLine)

	// Background
	background := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})
	background.StrokeColor = color.RGBA{R: 208, G: 200, B: 200, A: 255}
	background.StrokeWidth = 1
	background.CornerRadius = 8
	background.SetMinSize(fyne.NewSize(0, 500))

	bottomBackButton := container.NewHBox(backContainer) // Back Button
	content := container.NewBorder(header, bottomBackButton, nil, nil, scrollableFileList)
	box := container.NewStack(background, container.NewPadded(content))

	pageContent := container.NewVBox(heightMargin, title, subtitle, heightMargin, box)

	return container.NewBorder(nil, nil, widthMargin, widthMargin, container.NewPadded(pageContent))
}
