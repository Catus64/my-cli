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

func coloredBullet(parts []struct{ text string; color color.Color }) fyne.CanvasObject {
    bulletDot := canvas.NewText("•", color.White)  
    bulletDot.TextSize = 20

    hbox := container.NewHBox()
    for _, p := range parts {
        t := canvas.NewText(p.text, p.color)
        t.TextSize = 14
        hbox.Add(t)
    }

    leftMargin := canvas.NewRectangle(color.Transparent)
    leftMargin.SetMinSize(fyne.NewSize(20, 0))

    return container.NewBorder(nil, nil, container.NewHBox(leftMargin, bulletDot), nil, container.NewPadded(hbox))
}

// without bullet dot
func coloredNoBullet(parts []struct{ text string; color color.Color }) fyne.CanvasObject {
    hbox := container.NewHBox()
    for _, p := range parts {
        t := canvas.NewText(p.text, p.color)
        t.TextSize = 14
        hbox.Add(t)
    }

    leftMargin := canvas.NewRectangle(color.Transparent)
    leftMargin.SetMinSize(fyne.NewSize(35, 0)) 

    return container.NewBorder(nil, nil, leftMargin, nil, container.NewPadded(hbox))
}

func HelpContent() fyne.CanvasObject {
	title := canvas.NewText("Help", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Description of GUI action (Scrollable)", color.Gray{Y: 150})
	subtitle.TextSize = 15

	// 1. Repository Page
	section1Title := sectionTitle("1. Repository Page")
	repoBullet1 := bullet("The file(s) showed in the file list box are the file(s) that saved by user for current version in the repository.")
	repoBullet2 := bullet("Current Version - When user click this button, user can view the detail of the current version in the repository.")
	repoBullet3 := bullet("User can click the file in the file list box to view the content of the file in the text area.")

	section1 := container.NewVBox(section1Title, repoBullet1, repoBullet2, repoBullet3)

	// 2. Save List Page
	section2Title := sectionTitle("2. Save List Page")
	saveListBullet1 := bullet("The file(s) showed in the save list box are the file(s) that user want to save to repository.")
	saveListBullet2 := bullet("The file(s) showed in the \"File to Save\" box are the file(s) that ready to be saved to repository.")
	saveListBullet3 := bullet("Remove Button - When user click this button, the file in save list will be removed to file list in File Directory Page.")
	saveListBullet4 := bullet("Save Button - When user click this button after fill in messages, the file in save list will save with the message.")

	section2 := container.NewVBox(section2Title, saveListBullet1, saveListBullet2, saveListBullet3, saveListBullet4)

	// 3. File Directory Page
	section3Title := sectionTitle("3. File Directory Page")
	fileDirectoryBullet1 := bullet("The file(s) showed in the file list box are the file(s) that added or modified by user but haven't save to repository.")
	fileDirectoryBullet2 := bullet("Add Button - When user click this button, the file(s) in file list will add to save list in Save List Page.")

	section3 := container.NewVBox(section3Title, fileDirectoryBullet1, fileDirectoryBullet2)

	// 4. Ignored File Page
	section4Title := sectionTitle("4. Ignored File Page")
	ignoredFileBullet1 := bullet("The file(s) showed in the ignored file list box are the file(s) that user want to ignore.")
	ignoredFileBullet2 := bullet("View File Type Button - When user click this button, user can view the file type that will be ignored in the ignored file list box and edit it.")

	section4 := container.NewVBox(section4Title, ignoredFileBullet1, ignoredFileBullet2)

	// 5. Save File Page
	section5Title := sectionTitle("5. Save File Page")
	saveFileBullet1 := bullet("The file(s) showed in the save file list box are the Save States of the project.")
	saveFileBullet2 := bullet("Add Button - When user click this button, user are able to add the save state in the project.")
	saveFileBullet3 := bullet("Switch Button - Before user click this button, user need to select the save state in the row box to switch.")
	saveFileBullet4 := bullet("Merge Button - Before user click this button, user need to select the save state in the row box to merge.")

	section5 := container.NewVBox(section5Title, saveFileBullet1, saveFileBullet2, saveFileBullet3, saveFileBullet4)

	// 6. Save History Page
	section6Title := sectionTitle("6. Save History Page")
	saveHistoryBullet1 := bullet("The page will display history and past changes by user to enable them to browse and compare changes across files.")
	saveHistoryBullet2 := bullet("User can search the history bubbles in the search bar with history bubbles's hash key to find the specific history bubble.")
	saveHistoryBullet3 := bullet("User can click the history bubble to view the detail.")

	green  := color.RGBA{R: 100, G: 220, B: 100, A: 255}
	orange := color.RGBA{R: 255, G: 160, B: 60,  A: 255}
	white  := color.White

	saveHistoryBullet4 := coloredBullet([]struct{ text string; color color.Color }{
		{" a) ", white}, {"Merge→xxxxx", green}, {" = Merge branch to main", white},
	})
	saveHistoryBullet5 := coloredNoBullet([]struct{ text string; color color.Color }{
		{"b) ", white}, {"Merge→xxxxx", orange}, {" = Merge branch to branch", white},
	})

	section6 := container.NewVBox(section6Title, saveHistoryBullet1, saveHistoryBullet2, saveHistoryBullet3, saveHistoryBullet4, saveHistoryBullet5)

	// 7. Quit Button
	section7 := sectionTitle("7. Quit button allows user to quit to dashboard page.")

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
		section4,
		heightMargin, 
		section5,
		heightMargin,
		section6,
		heightMargin,
		section7,
		heightMargin,
	)

	scrollContent := container.NewScroll(fullcontent)

	return container.NewBorder(nil, nil, widthMargin, widthMargin, container.NewPadded(scrollContent))
}