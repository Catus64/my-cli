package sidebar

import (
	"gocmd/gui"
	"image/color"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func NavButton(label string, isActive bool, action func()) fyne.CanvasObject {
	button := widget.NewButton("", action)
	button.Importance = widget.LowImportance

	var backgroundColor color.Color
	if isActive {
		backgroundColor = color.RGBA{R: 88, G: 84, B: 84, A: 255}
	} else {
		backgroundColor = color.Black
	}
	
	buttonBg := canvas.NewRectangle(backgroundColor)
	buttonBg.StrokeColor = color.RGBA{R: 204, G: 200, B: 200, A: 255}
	buttonBg.StrokeWidth = 1
	buttonBg.CornerRadius = 8

	text := canvas.NewText(label, color.White)
	text.Alignment = fyne.TextAlignCenter

	buttonHeightSize := canvas.NewRectangle(color.Transparent)
	buttonHeightSize.SetMinSize(fyne.NewSize(0, 40))

	buttonStack := container.NewStack(buttonBg, buttonHeightSize, button, text)

	LRmargin := canvas.NewRectangle(color.Transparent)
	LRmargin.SetMinSize(fyne.NewSize(15, 0))

	return container.NewBorder(nil, nil, LRmargin, LRmargin, buttonStack)
}

func SideBar(gui *gui.MyApp, window fyne.Window, pathName string, activePage string) fyne.CanvasObject {
	homeButton := NavButton("Home Page", activePage == "home", func()  {
		
	})
	modifiedFileButton := NavButton("Modified File", activePage == "modified", func()  {
		
	})
	historyButton := NavButton("History", activePage == "history", func()  {
		
	})
	helpButton := NavButton("Help", activePage == "help", func()  {
		
	})

	// Current repository box
	repoLabel := canvas.NewText("Current Repository", color.Gray{Y: 150})
	repoLabel.TextSize = 16
	repoLabel.Alignment = fyne.TextAlignCenter

	repoName := canvas.NewText(filepath.Base(pathName), color.White)
	repoName.TextSize = 20
	repoName.TextStyle = fyne.TextStyle{Bold: true}
	repoName.Alignment = fyne.TextAlignCenter

	repoBackground := canvas.NewRectangle(color.Black)
	repoBackground.StrokeColor = color.RGBA{R: 124, G: 115, B: 115, A: 255}
	repoBackground.StrokeWidth = 1
	repoBackground.CornerRadius = 2

	repoHeightMargin := canvas.NewRectangle(color.Transparent)
	repoHeightMargin.SetMinSize(fyne.NewSize(0, 80)) 

	textHeightMargin := canvas.NewRectangle(color.Transparent)
	textHeightMargin.SetMinSize(fyne.NewSize(0, 10)) 

	textContent := container.NewCenter(container.NewVBox(repoLabel, textHeightMargin, repoName))

	repoBox := container.NewStack(repoBackground, repoHeightMargin, textContent)

	// Quit Button
	quitButton := widget.NewButton("Quit", func ()  {
		window.Close()
        gui.Window.Show()
	})
	quitButton.Importance = widget.DangerImportance

	BoxLRmargin := canvas.NewRectangle(color.Transparent)
	BoxLRmargin.SetMinSize(fyne.NewSize(40, 0))

	QuitLRmargin := canvas.NewRectangle(color.Transparent)
	QuitLRmargin.SetMinSize(fyne.NewSize(80, 0))

	RepoBoxPosition := container.NewBorder(nil, nil, BoxLRmargin, BoxLRmargin, repoBox)
	QuitButtonPosition := container.NewBorder(nil, nil, QuitLRmargin, QuitLRmargin, quitButton)

	// Side Bar Background
	sideBackground := canvas.NewRectangle(color.RGBA{R: 3, G: 36, B: 63, A: 255})
	sideBackground.SetMinSize(fyne.NewSize(200, 0))
	sideBackground.StrokeColor = color.RGBA{R: 124, G: 115, B: 115, A: 255}
	sideBackground.StrokeWidth = 1

	heightSize := canvas.NewRectangle(color.Transparent)
	heightSize.SetMinSize(fyne.NewSize(0, 20))

	sideBarContent := container.NewVBox(
		gui.HeaderNoneLine(),
		heightSize,
		homeButton,
		heightSize,
		modifiedFileButton,
		heightSize,
		historyButton,
		heightSize,
		helpButton,
		layout.NewSpacer(),
		RepoBoxPosition,
		heightSize,
		QuitButtonPosition,
		heightSize,
	)

	return container.NewStack(sideBackground, container.NewPadded(sideBarContent))
}