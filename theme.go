package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CustomTheme extends the default theme with better hover colors and contrast
type CustomTheme struct {
	fyne.Theme
}

// NewCustomTheme creates a new custom theme based on the default theme
func NewCustomTheme() fyne.Theme {
	return &CustomTheme{
		Theme: theme.DefaultTheme(),
	}
}

// Color returns theme colors with enhanced hover support
func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Use default theme colors as base
	baseColor := t.Theme.Color(name, variant)
	
	// Enhance specific colors for better contrast and visibility
	switch name {
	case theme.ColorNameButton:
		// Improve button colors for better hover visibility
		if variant == theme.VariantDark {
			return &color.RGBA{R: 70, G: 70, B: 70, A: 255}
		}
		return &color.RGBA{R: 250, G: 250, B: 250, A: 255}
		
	case theme.ColorNameHover:
		// Enhanced hover colors for better visibility - use subtle colors
		if variant == theme.VariantDark {
			return &color.RGBA{R: 60, G: 60, B: 60, A: 255}
		}
		return &color.RGBA{R: 240, G: 245, B: 250, A: 255}
		
	case theme.ColorNamePressed:
		// Enhanced pressed state colors
		if variant == theme.VariantDark {
			return &color.RGBA{R: 90, G: 90, B: 90, A: 255}
		}
		return &color.RGBA{R: 200, G: 200, B: 200, A: 255}
	}
	
	return baseColor
}

// GetHoverColor returns appropriate hover color for custom widgets
func GetHoverColor(variant fyne.ThemeVariant) color.Color {
	if variant == theme.VariantDark {
		return &color.RGBA{R: 60, G: 60, B: 60, A: 255}
	}
	// Nice subtle blue-gray hover color - gentle and pleasant
	return &color.RGBA{R: 240, G: 245, B: 250, A: 255}
}

// GetDeleteHoverColor returns appropriate hover color for delete buttons
func GetDeleteHoverColor(variant fyne.ThemeVariant) color.Color {
	if variant == theme.VariantDark {
		return &color.RGBA{R: 80, G: 40, B: 40, A: 255}
	}
	return &color.RGBA{R: 255, G: 200, B: 200, A: 255}
}

// GetHighContrastColor returns a high contrast color for better visibility
func GetHighContrastColor(variant fyne.ThemeVariant) color.Color {
	if variant == theme.VariantDark {
		return &color.RGBA{R: 100, G: 100, B: 100, A: 255}
	}
	return &color.RGBA{R: 180, G: 180, B: 180, A: 255}
}

// DetectThemeVariant automatically detects if the current theme is light or dark
func DetectThemeVariant() fyne.ThemeVariant {
	bg := theme.BackgroundColor()
	if rgba, ok := bg.(*color.RGBA); ok {
		// Calculate luminance to determine if theme is light or dark
		luminance := (0.299*float64(rgba.R) + 0.587*float64(rgba.G) + 0.114*float64(rgba.B)) / 255.0
		if luminance > 0.5 {
			return theme.VariantLight
		}
		return theme.VariantDark
	}
	// Default to light theme if we can't determine
	return theme.VariantLight
}

/*
Alternative hover colors for different preferences:
- Very subtle: RGB(248, 250, 252) - barely noticeable blue-gray tint (current)
- Light gray: RGB(245, 245, 245) - neutral light gray
- Warm tint: RGB(252, 250, 248) - very subtle warm tint
- Cool tint: RGB(248, 250, 252) - very subtle cool tint
*/