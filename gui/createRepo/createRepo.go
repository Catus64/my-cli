package createRepo

import (
	"fmt"
	"gocmd/gui"
	"gocmd/gui/help"
	"gocmd/gui/history"
	"gocmd/gui/homepage"
	"gocmd/gui/ignoredFile"
	"gocmd/gui/modifiedFile"
	"gocmd/gui/saveFile"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"image/color"
	"strings"

	"fyne.io/fyne/v2" // Core Fyne types (like window sizes)
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog" // Pop-up windows/file pickers
	"fyne.io/fyne/v2/widget" // UI components (buttons, text boxes)
)

func RepoInit(path string) error {
	_, err := gitpath.Repo_create(path)
	if err != nil {
		return fmt.Errorf("Path does not exist or already exist!")
	}
    return nil
}

func Show(gui *gui.MyApp) {
	myWindow := gui.App.NewWindow(gui.Window.Title())

	selectedPath := ""
	entry := pathEntry(&selectedPath)
	browseButton := browseButton(myWindow, entry, &selectedPath)
	createButton := createButton(entry, myWindow, gui)
	quitButton := QuitButton(myWindow, gui)
	content := Content(entry, browseButton, createButton, quitButton)

	fullcontent := container.NewBorder(gui.Header(), nil, gui.SideRect(), nil, content)

	myWindow.SetContent(fullcontent)
	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.Show()
}

func pathEntry(selectedPath *string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText(*selectedPath)

	return entry
}

func browseButton(window fyne.Window, entry *widget.Entry, selectedPath *string) *fyne.Container {
	browse := widget.NewButton("Browse", func() {
		fileDialog := dialog.NewFolderOpen(
			func(reader fyne.ListableURI, err error) {
				if reader != nil {
					// Update input box
					*selectedPath = reader.Path()
					entry.SetText(*selectedPath)
				}
			},
			window,
		)
		fileDialog.Show()
	})
	return container.NewHBox(browse)
}

func createButton(entry *widget.Entry, window fyne.Window, gui *gui.MyApp) *widget.Button {
	createButton := widget.NewButton("Create", func() {
		path := strings.TrimSpace(entry.Text)

		err := RepoInit(path)
        if err != nil{
            dialog.ShowError(err, window)
            entry.SetText("")

            return
        }

		// Window after creating
		mainWindow := gui.App.NewWindow(gui.Window.Title())
		mainWindow.Resize(fyne.NewSize(1000, 600))

		// define all navigation function
		var showHome, showSave, showFileDirectory, showIgnoredFile, showHistory, showHelp func()

		showHome = func() {
			homepage.Show(gui, path, mainWindow, showSave, showFileDirectory, showIgnoredFile, showHistory, showHelp)
		}
		showSave = func() {
			saveFile.Show(gui, path, mainWindow, showHome, showFileDirectory, showIgnoredFile, showHistory, showHelp)
		}
		showFileDirectory = func() {
			modifiedFile.Show(gui, path, mainWindow, showHome, showSave, showIgnoredFile, showHistory, showHelp)
		}
		showIgnoredFile = func() {
			ignoredFile.Show(gui, path, mainWindow, showHome, showSave, showFileDirectory, showHistory, showHelp)
		}
		showHistory = func() {
			history.Show(gui, path, mainWindow, showHome, showSave, showFileDirectory, showIgnoredFile, showHelp)
		}
		showHelp = func() {
			help.Show(gui, path, mainWindow, showHome, showSave, showFileDirectory, showIgnoredFile, showHistory)
		}

		window.Hide()
		mainWindow.Show()
		showHome() // show homepage after create
	})

	createButton.Disable()

	entry.OnChanged = func(text string) {
		cleanText := strings.TrimSpace(text)

		if cleanText == "" {
			createButton.Disable()
		} else {
			createButton.Enable()
		}
	}
	createButton.Importance = widget.HighImportance // blue
	return createButton
}

func QuitButton(window fyne.Window, gui *gui.MyApp) *widget.Button {
	quitButton := widget.NewButton("Quit", func() {
		window.Close()
		gui.Window.Show()
	})
	quitButton.Importance = widget.DangerImportance // red
	return quitButton
}

func Content(entry *widget.Entry, browseButton *fyne.Container, createButton *widget.Button, quitButton *widget.Button) fyne.CanvasObject {
	title := canvas.NewText("Create Repository", color.White)
	title.TextSize = 56
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Select your file", color.Gray{Y: 150})
	subtitle.TextSize = 16

	Label := canvas.NewText("Directory:", color.White)
	Label.TextSize = 28
	Label.TextStyle = fyne.TextStyle{Bold: true}

	bottomButton := container.NewBorder(nil, nil, createButton, quitButton, nil)

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(50, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 40))

	content := container.NewVBox(
		heightMargin,
		title,
		subtitle,
		heightMargin,
		heightMargin,
		Label,
		entry,
		browseButton,
		heightMargin,
		bottomButton,
	)

	return container.NewBorder(nil, nil, widthMargin, widthMargin, container.NewPadded(content))
}
