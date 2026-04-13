package help

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type whiteTextTheme struct{}

func (whiteTextTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground {
		return color.White // Force only the text to be pure white
	}
	return theme.DefaultTheme().Color(name, variant)
}
func (whiteTextTheme) Font(style fyne.TextStyle) fyne.Resource { return theme.DefaultTheme().Font(style) }
func (whiteTextTheme) Icon(name fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(name) }
func (whiteTextTheme) Size(name fyne.ThemeSizeName) float32 { return theme.DefaultTheme().Size(name) }


func HelpContent() fyne.CanvasObject {
	title := canvas.NewText("Help", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Description of GUI action (Scrollable)", color.Gray{Y: 150})
	subtitle.TextSize = 15

	// 1. Home Page
	section1Title := sectionTitle("1. Home Page")
	homeBullet1 := bullet("Remove Button - When user click this button, the file in save list will be removed to file list in Modified File page.")
	homeBullet2 := bullet("Save Button - When user click this button after fill in messages, the file in save list will save with the message.")

	section1 := container.NewVBox(section1Title, homeBullet1, homeBullet2)

	// 2. Modified File
	section2Title := sectionTitle("2. Modified File Page")
	modifiedBullet1 := bullet("The file(s) showed in the file list box are the file(s) that modified by user but haven't save to repository.")
	modifiedBullet2 := bullet("Add Button - When user click this button, the file in file list will add to save list in Home Page.")

	section2 := container.NewVBox(section2Title, modifiedBullet1, modifiedBullet2)

	// 3. History
	section3Title := sectionTitle("3. History")
	historyBullet1 := bullet("The page will display history and past changes by user to enable them to browse and compare changes across files.")
	historyBullet2 := bullet("If user click any version button in the page, the interface will direct user to a new tab with column view format to search for history.")

	section3 := container.NewVBox(section3Title, historyBullet1, historyBullet2)

	// 4. Quit Button
	section4Title := sectionTitle("4. Quit button allows user to quit to dashboard page.")

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	fullcontent := container.NewVBox(
		heightMargin, 
		title,
		subtitle, 
		heightMargin, 
		section1,
		heightMargin, 
		section2,
		heightMargin,
		section3,
		heightMargin,
		section4Title,
		heightMargin,
	)

	scrollContent := container.NewScroll(fullcontent)

	return container.NewBorder(nil, nil, widthMargin, widthMargin, container.NewPadded(scrollContent))
}

func sectionTitle(text string) fyne.CanvasObject {
    t := canvas.NewText(text, color.White)
    t.TextSize = 22
    t.TextStyle = fyne.TextStyle{Bold: true}
    return t
}

func bullet(text string) fyne.CanvasObject {
	bullet := canvas.NewText("•", color.White)
    bullet.TextSize = 20

    content := widget.NewLabel(text)
    content.Wrapping = fyne.TextWrapWord

	// Change the default color of label text to white color
	fullContent := container.NewThemeOverride(content, &whiteTextTheme{})

	leftMargin := canvas.NewRectangle(color.Transparent)
	leftMargin.SetMinSize(fyne.NewSize(20, 0))

	return container.NewBorder(nil, nil, container.NewHBox(leftMargin, bullet), nil, container.NewPadded(fullContent))
}