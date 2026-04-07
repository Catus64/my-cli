package createRepo

import (
	"gocmd/gui"
	"gocmd/gui/homepage"
	"image/color"
	"strings"

	"fyne.io/fyne/v2" // Core Fyne types (like window sizes)
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog" // Pop-up windows/file pickers
	"fyne.io/fyne/v2/widget" // UI components (buttons, text boxes)
)

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
        window.Hide()
        homepage.Show(gui, path)
    })

    createButton.Disable()

    entry.OnChanged = func(text string) {
        cleanText := strings.TrimSpace(text)
        
        if cleanText == "" {
            createButton.Disable()
        }else {
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