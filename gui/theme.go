package gui

import (
	"image/color"
	"os"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var persianFont fyne.Resource
var persianFontBold fyne.Resource

// loadPersianFont tries to load a Persian-supporting font from system
func loadPersianFont() {
	// Try to find and load Persian fonts from common locations
	fontPaths := []string{
		"/usr/share/fonts/truetype/noto/NotoSansArabic-Regular.ttf",
		"/usr/share/fonts/truetype/noto/NotoSansArabicUI-Regular.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/TTF/NotoSansArabic-Regular.ttf",
		"/usr/share/fonts/TTF/DejaVuSans.ttf",
	}
	
	for _, path := range fontPaths {
		if _, err := os.Stat(path); err == nil {
			resource, err := fyne.LoadResourceFromPath(path)
			if err == nil {
				persianFont = resource
				break
			}
		}
	}
	
	// Try bold version
	boldPaths := []string{
		"/usr/share/fonts/truetype/noto/NotoSansArabic-Bold.ttf",
		"/usr/share/fonts/truetype/noto/NotoSansArabicUI-Bold.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
	}
	
	for _, path := range boldPaths {
		if _, err := os.Stat(path); err == nil {
			resource, err := fyne.LoadResourceFromPath(path)
			if err == nil {
				persianFontBold = resource
				break
			}
		}
	}
}

// PersianTheme provides Persian font support
type PersianTheme struct {
	fyne.Theme
}

// NewPersianTheme creates a theme with Persian font support
func NewPersianTheme() fyne.Theme {
	// Load Persian font on first call
	if persianFont == nil {
		loadPersianFont()
	}
	return &PersianTheme{Theme: theme.DefaultTheme()}
}

// Font returns a font that supports Persian characters
func (t *PersianTheme) Font(style fyne.TextStyle) fyne.Resource {
	// If we have a loaded Persian font, use it
	if style.Bold && persianFontBold != nil {
		return persianFontBold
	}
	if persianFont != nil {
		return persianFont
	}
	// Fallback to default theme font
	return theme.DefaultTheme().Font(style)
}

// Size returns the size for a theme size name
func (t *PersianTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// Color returns the color for a theme color name
func (t *PersianTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

// Icon returns the icon for a theme icon name
func (t *PersianTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

