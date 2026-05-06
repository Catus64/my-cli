package openRepo

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2" // Core Fyne types (like window sizes)
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog" // Pop-up windows/file pickers
	"fyne.io/fyne/v2/widget" // UI components (buttons, text boxes)

	"gocmd/gui"
	"gocmd/gui/help"
	"gocmd/gui/history"
	"gocmd/gui/homepage"
	"gocmd/gui/ignoredFile"
	"gocmd/gui/modifiedFile"
	"gocmd/gui/saveFile"
	gitpath "gocmd/testfiles/Gitrepostruct"
)

func OpenRepo(path string) error{
    _, err := gitpath.Repo_find(path, false)
	if err != nil {
		return fmt.Errorf("Repository doesn't exist!")
	}
    return nil
}

func Show(gui *gui.MyApp) {
	myWindow := gui.App.NewWindow(gui.Window.Title())

    selectedPath := "" 
    entry := path(&selectedPath)
    browseButton := browse(myWindow, entry, &selectedPath)
    openButton := openButton(entry, myWindow, gui)
    quitButton := quitButton(myWindow, gui)
    content := OpenRepoContent(entry, browseButton, openButton, quitButton)

    fullcontent := container.NewBorder(gui.Header(), nil, gui.SideRect(), nil, content)

	myWindow.SetContent(fullcontent)
	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.Show()
}

func path(selectedPath *string) *widget.Entry {
    entry := widget.NewEntry()
    entry.SetText(*selectedPath) 
    entry.Disable()
    return entry
}

func browse(window fyne.Window, entry *widget.Entry, selectedPath *string) *fyne.Container {
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

func openButton(entry *widget.Entry, window fyne.Window, gui *gui.MyApp) *widget.Button {
    openButton := widget.NewButton("Open", func() {
        path := strings.TrimSpace(entry.Text)

        err := OpenRepo(path)
        if err != nil{
            dialog.ShowError(err, window)
            entry.SetText("")

            return
        }

        // Window after Open
        mainWindow := gui.App.NewWindow(gui.Window.Title())
        mainWindow.Resize(fyne.NewSize(1000, 600))
        
        // define all navigation function
        var showHome, showSave, showFileDirectory, showIgnoredFile, showHistory, showHelp func()

		showHome = func() {
			homepage.Show(gui, path, mainWindow, showHome, showSave, showFileDirectory, showIgnoredFile, showHistory, showHelp)
		}
		showSave = func() {
			saveFile.Show(gui, path, mainWindow, showHome, showSave, showFileDirectory, showIgnoredFile, showHistory, showHelp)
		}
		showFileDirectory = func() {
			modifiedFile.Show(gui, path, mainWindow, showHome, showSave, showFileDirectory, showIgnoredFile, showHistory, showHelp)
		}
		showIgnoredFile = func() {
			ignoredFile.Show(gui, path, mainWindow, showHome, showSave, showFileDirectory, showIgnoredFile, showHistory, showHelp)
		}
		showHistory = func() {
			history.Show(gui, path, mainWindow, showHome, showSave, showFileDirectory, showIgnoredFile, showHistory, showHelp)
		}
		showHelp = func() {
			help.Show(gui, path, mainWindow, showHome, showSave, showFileDirectory, showIgnoredFile, showHistory, showHelp)
		}

        window.Hide()
        mainWindow.Show()
        showHome() // show homepage after create
    })
    openButton.Importance = widget.HighImportance // blue

    openButton.Disable()

    entry.OnChanged = func(text string) {
        cleanText := strings.TrimSpace(text)
        
        if cleanText == "" {
            openButton.Disable()
        }else {
            openButton.Enable()
        }
    }

    return openButton
}

func quitButton(window fyne.Window, gui *gui.MyApp) *widget.Button {
    quitButton := widget.NewButton("Quit", func() {
        window.Close()
        gui.Window.Show()
    })
    quitButton.Importance = widget.DangerImportance // red
    return quitButton
}

func OpenRepoContent(entry *widget.Entry, browseButton *fyne.Container, openButton *widget.Button, quitButton *widget.Button) fyne.CanvasObject {
    title := canvas.NewText("Open Repository", color.White)
    title.TextSize = 56
    title.TextStyle = fyne.TextStyle{Bold: true}

    subtitle := canvas.NewText("Select your file", color.Gray{Y: 150})
    subtitle.TextSize = 16

    Label := canvas.NewText("Directory:", color.White)
    Label.TextSize = 28
    Label.TextStyle = fyne.TextStyle{Bold: true}

    bottomButton := container.NewBorder(nil, nil, openButton, quitButton, nil)

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