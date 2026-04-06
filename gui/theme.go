package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyDarkTheme struct{}

func (MyDarkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		return color.Black // force black background
	}
	if name == theme.ColorNameDisabled { //create and open button after select dashboard option
		return color.Gray{Y: 150} 
	}
	if name == theme.ColorNameHover {
		return color.RGBA{R: 211, G: 211, B: 211, A: 40} //light gray hover
	}
	return theme.DefaultTheme().Color(name, variant)
}
func (MyDarkTheme) Font(style fyne.TextStyle) fyne.Resource { return theme.DefaultTheme().Font(style) }
func (MyDarkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (MyDarkTheme) Size(name fyne.ThemeSizeName) float32 { return theme.DefaultTheme().Size(name) }