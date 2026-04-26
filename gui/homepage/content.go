package homepage

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/layout"
	// "fyne.io/fyne/v2/widget"
)

func HomePageContent() fyne.CanvasObject {
	title := canvas.NewText("Home Page", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("View your file in repository", color.Gray{Y: 150})
	subtitle.TextSize = 15

	header := container.NewVBox(title, subtitle)

	fileListLabel := canvas.NewText("File List", color.RGBA{R: 208, G: 200, B: 200, A: 255})
	fileListLabel.TextSize = 20
	fileListLabel.TextStyle = fyne.TextStyle{Bold: true}

	return header
}
