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
	repoBullet2 := bullet("Current Details - When user click this button, user can view the detail of the current version in the repository.")
	repoBullet3 := bullet("User can click the file in the file list box to view the content of the file in the text area.")

	section1 := container.NewVBox(section1Title, repoBullet1, repoBullet2, repoBullet3)

	// 2. File to be Saved Page
	section2Title := sectionTitle("2. File to be Saved Page")
	fileToBeSavedBullet1 := bullet("The file(s) showed in the Files to be Saved box are the file(s) that user want to save to repository.")
	fileToBeSavedBullet2 := bullet("The file(s) showed in the \"File to be Save\" box are the file(s) that ready to be save to repository.")
	fileToBeSavedBullet3 := bullet("Remove Button - When user click this button, the file in File(s) to be Saved box will be removed to file list in File Directory Page.")
	fileToBeSavedBullet4 := bullet("Save Button - When user click this button after fill in messages, the file in File to be Save box will save with the message.")

	section2 := container.NewVBox(section2Title, fileToBeSavedBullet1, fileToBeSavedBullet2, fileToBeSavedBullet3, fileToBeSavedBullet4)

	// 3. Working Directory Page
	section3Title := sectionTitle("3. Working Directory Page")
	workingDirectoryBullet1 := bullet("The file(s) showed in the file list box are the file(s) that added or modified by user but haven't save to repository.")
	workingDirectoryBullet2 := bullet("Add Button - When user click this button, the file(s) in file list will add to File(s) to be Saved box in File to be Saved Page.")

	section3 := container.NewVBox(section3Title, workingDirectoryBullet1, workingDirectoryBullet2)

	// 4. Ignored File Page
	section4Title := sectionTitle("4. Ignored File Page")
	ignoredFileBullet1 := bullet("The file(s) showed in the ignored file list box are the file(s) that user want to ignore.")
	ignoredFileBullet2 := bullet("View File Type Button - When user click this button, user can view the file type that will be ignored in the ignored file list box and edit it.")

	section4 := container.NewVBox(section4Title, ignoredFileBullet1, ignoredFileBullet2)

	// 5. Branch Management Page
	section5Title := sectionTitle("5. Branch Management Page")
	branchBullet1 := bullet("The file(s) showed in the table are the Save States of the repository.")
	branchBullet2 := bullet("Add Button - When user click this button, user are able to add the branch into the repository.")
	branchBullet3 := bullet("Switch Button - Before user click this button, user need to select a branch in the table to switch.")
	branchBullet4 := bullet("Merge Button - Before user click this button, user need to select a branch in the table to merge.")

	section5 := container.NewVBox(section5Title, branchBullet1, branchBullet2, branchBullet3, branchBullet4)

	// 6. Version History Page
	section6Title := sectionTitle("6. Version History Page")
	versionHistoryBullet1 := bullet("The page will display history and past changes by user to enable them to browse and compare changes across files.")
	versionHistoryBullet2 := bullet("User can search the history bubbles in the search bar with history bubbles's hash key to find the specific history bubble.")
	versionHistoryBullet3 := bullet("User can click the history bubble to view the detail.")

	green  := color.RGBA{R: 100, G: 220, B: 100, A: 255}
	orange := color.RGBA{R: 255, G: 160, B: 60,  A: 255}
	white  := color.White

	versionHistoryBullet4 := coloredBullet([]struct{ text string; color color.Color }{
		{" a) ", white}, {"Merge→xxxxx", green}, {" = Merge branch to main", white},
	})
	versionHistoryBullet5 := coloredNoBullet([]struct{ text string; color color.Color }{
		{"b) ", white}, {"Merge→xxxxx", orange}, {" = Merge branch to branch", white},
	})

	section6 := container.NewVBox(section6Title, versionHistoryBullet1, versionHistoryBullet2, versionHistoryBullet3, versionHistoryBullet4, versionHistoryBullet5)

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