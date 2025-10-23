package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// HistoryListItem is a custom widget for displaying clipboard history items
// with delete functionality and proper text wrapping
type HistoryListItem struct {
	widget.BaseWidget
	item     ClipboardItem
	index    int
	onDelete func(int)
	onSelect func(int)
	
	// Internal widgets
	textWidget     *widget.RichText
	imageWidget    *canvas.Image
	deleteButton   *widget.Button
	container      *fyne.Container
	background     *canvas.Rectangle
	
	// Hover state management
	isHovered       bool
	deleteHovered   bool
}

// NewHistoryListItem creates a new history list item widget
func NewHistoryListItem(clipboardItem ClipboardItem, index int, onDelete func(int), onSelect func(int)) *HistoryListItem {
	item := &HistoryListItem{
		item:     clipboardItem,
		index:    index,
		onDelete: onDelete,
		onSelect: onSelect,
	}
	
	item.ExtendBaseWidget(item)
	item.createContent()
	return item
}

// createContent builds the internal widget structure
func (h *HistoryListItem) createContent() {
	// Create background rectangle for hover effects
	h.background = canvas.NewRectangle(color.Transparent)
	h.background.CornerRadius = 4 // Rounded corners for modern look
	
	// Create delete button with X icon and custom hover handling
	h.deleteButton = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		if h.onDelete != nil {
			h.onDelete(h.index)
		}
	})
	
	// Style the delete button with consistent sizing and improved accessibility
	h.deleteButton.Resize(fyne.NewSize(28, 28)) // Slightly larger for better touch targets
	h.deleteButton.Importance = widget.LowImportance
	
	var contentWidget fyne.CanvasObject
	
	if h.item.Type == ItemTypeText {
		// Create the text widget with proper wrapping
		h.textWidget = widget.NewRichText()
		h.textWidget.Wrapping = fyne.TextWrapWord
		
		// Set maximum height to prevent excessive vertical space usage
		// This ensures consistent item heights in the list
		maxHeight := float32(80) // Approximately 4 lines of text
		h.textWidget.Resize(fyne.NewSize(0, maxHeight))
		
		// Set the text content with proper formatting
		displayText := h.prepareDisplayText()
		h.textWidget.ParseMarkdown(displayText)
		
		contentWidget = h.textWidget
	} else if h.item.Type == ItemTypeImage {
		// Create image preview widget
		if img, err := h.createImageFromBase64(); err == nil {
			h.imageWidget = canvas.NewImageFromImage(img)
			h.imageWidget.FillMode = canvas.ImageFillContain
			h.imageWidget.SetMinSize(fyne.NewSize(100, 60))
			
			// Create a container with image and metadata
			metaText := h.createImageMetaText()
			metaLabel := widget.NewLabel(metaText)
			metaLabel.TextStyle = fyne.TextStyle{Italic: true}
			
			imageContainer := container.NewHBox(
				container.NewBorder(nil, nil, nil, nil, h.imageWidget),
				widget.NewSeparator(),
				metaLabel,
			)
			
			contentWidget = imageContainer
		} else {
			// Fallback to text display if image can't be loaded
			h.textWidget = widget.NewRichText()
			h.textWidget.ParseMarkdown(fmt.Sprintf("**[IMAGE ERROR]** %v", err))
			contentWidget = h.textWidget
		}
	}
	
	// Create content container with content and delete button
	contentContainer := container.NewBorder(
		nil, nil, nil, h.deleteButton,
		contentWidget,
	)
	
	// Create main container with background and content layered
	h.container = container.NewStack(
		h.background,
		container.NewPadded(contentContainer),
	)
	
	// Initialize hover states
	h.isHovered = false
	h.deleteHovered = false
	h.updateHoverState()
}

// prepareDisplayText formats the text for display with proper wrapping and truncation
func (h *HistoryListItem) prepareDisplayText() string {
	text := strings.TrimSpace(h.item.Content)
	
	// Replace tabs with spaces for better display
	text = strings.ReplaceAll(text, "\t", "    ")
	
	// Limit the maximum length to prevent excessive vertical space usage
	// This ensures items don't take up too much vertical space in the list
	const maxDisplayLength = 250
	if len(text) > maxDisplayLength {
		// Find a good break point near the limit to avoid cutting words
		breakPoint := maxDisplayLength
		for i := maxDisplayLength - 30; i < len(text) && i < maxDisplayLength+10; i++ {
			if text[i] == ' ' || text[i] == '\n' || text[i] == '\t' {
				breakPoint = i
				break
			}
		}
		text = text[:breakPoint] + "..."
	}
	
	// Limit the number of lines to prevent excessive vertical space
	lines := strings.Split(text, "\n")
	const maxLines = 4
	if len(lines) > maxLines {
		text = strings.Join(lines[:maxLines], "\n") + "\n..."
	}
	
	// Escape markdown special characters to prevent formatting issues
	text = strings.ReplaceAll(text, "*", "\\*")
	text = strings.ReplaceAll(text, "_", "\\_")
	text = strings.ReplaceAll(text, "#", "\\#")
	text = strings.ReplaceAll(text, "[", "\\[")
	text = strings.ReplaceAll(text, "]", "\\]")
	
	return text
}

// createImageFromBase64 creates an image from base64 data
func (h *HistoryListItem) createImageFromBase64() (image.Image, error) {
	// Decode base64 data
	imageData, err := base64.StdEncoding.DecodeString(h.item.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %v", err)
	}
	
	// Create a reader from the image data
	reader := strings.NewReader(string(imageData))
	
	// Try to decode based on the format
	var img image.Image
	if h.item.ImageMeta != nil {
		switch h.item.ImageMeta.Format {
		case "png":
			img, err = png.Decode(reader)
		case "jpeg", "jpg":
			img, err = jpeg.Decode(reader)
		default:
			// Try PNG first, then JPEG
			if img, err = png.Decode(reader); err != nil {
				reader = strings.NewReader(string(imageData)) // Reset reader
				img, err = jpeg.Decode(reader)
			}
		}
	} else {
		// Try PNG first, then JPEG
		if img, err = png.Decode(reader); err != nil {
			reader = strings.NewReader(string(imageData)) // Reset reader
			img, err = jpeg.Decode(reader)
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}
	
	return img, nil
}

// createImageMetaText creates metadata text for image items
func (h *HistoryListItem) createImageMetaText() string {
	if h.item.ImageMeta == nil {
		return "Image (unknown format)"
	}
	
	sizeKB := h.item.ImageMeta.Size / 1024
	if sizeKB == 0 {
		sizeKB = 1 // Show at least 1KB
	}
	
	return fmt.Sprintf("%s Image\n%dx%d\n%d KB", 
		strings.ToUpper(h.item.ImageMeta.Format),
		h.item.ImageMeta.Width,
		h.item.ImageMeta.Height,
		sizeKB)
}

// getHoverColors returns appropriate hover colors based on the current theme
func (h *HistoryListItem) getHoverColors() (itemHover, deleteHover color.Color) {
	// Use improved theme detection
	variant := DetectThemeVariant()
	
	// Use custom theme functions for consistent hover colors
	itemHover = GetHoverColor(variant)
	deleteHover = GetDeleteHoverColor(variant)
	
	return itemHover, deleteHover
}

// isLightTheme determines if the current theme is light based on background color
func (h *HistoryListItem) isLightTheme(bg color.Color) bool {
	if rgba, ok := bg.(*color.RGBA); ok {
		// Calculate luminance to determine if theme is light or dark
		luminance := (0.299*float64(rgba.R) + 0.587*float64(rgba.G) + 0.114*float64(rgba.B)) / 255.0
		return luminance > 0.5
	}
	// Default to light theme if we can't determine
	return true
}

// updateHoverState updates the visual appearance based on current hover state
func (h *HistoryListItem) updateHoverState() {
	itemHover, deleteHover := h.getHoverColors()
	
	if h.deleteHovered {
		// Delete button is hovered - use delete hover color with higher opacity
		h.background.FillColor = deleteHover
		h.background.StrokeColor = &color.RGBA{R: 200, G: 100, B: 100, A: 100}
		h.background.StrokeWidth = 1
	} else if h.isHovered {
		// Item is hovered - use item hover color
		h.background.FillColor = itemHover
		h.background.StrokeColor = theme.PrimaryColor()
		h.background.StrokeWidth = 1
	} else {
		// No hover - transparent background, no stroke
		h.background.FillColor = color.Transparent
		h.background.StrokeColor = color.Transparent
		h.background.StrokeWidth = 0
	}
	
	// Apply additional styling to delete button when hovered
	if h.deleteHovered {
		// Make delete button more prominent when hovered
		h.deleteButton.Importance = widget.HighImportance
	} else {
		// Reset to normal importance
		h.deleteButton.Importance = widget.LowImportance
	}
	
	h.background.Refresh()
	h.deleteButton.Refresh()
}

// CreateRenderer creates the renderer for this widget
func (h *HistoryListItem) CreateRenderer() fyne.WidgetRenderer {
	return &historyListItemRenderer{
		item:      h,
		container: h.container,
	}
}

// Tapped handles tap events for item selection
func (h *HistoryListItem) Tapped(*fyne.PointEvent) {
	if h.onSelect != nil {
		h.onSelect(h.index)
	}
}

// TappedSecondary handles secondary tap events (right-click)
func (h *HistoryListItem) TappedSecondary(*fyne.PointEvent) {
	// Could be used for context menu in the future
}

// MouseIn handles mouse enter events for hover effects
func (h *HistoryListItem) MouseIn(event *desktop.MouseEvent) {
	h.isHovered = true
	h.checkDeleteButtonHover(event.Position)
	h.updateHoverState()
}

// MouseOut handles mouse leave events
func (h *HistoryListItem) MouseOut() {
	h.isHovered = false
	h.deleteHovered = false
	h.updateHoverState()
}

// MouseMoved handles mouse movement for hover tracking
func (h *HistoryListItem) MouseMoved(event *desktop.MouseEvent) {
	// Check if mouse is over delete button area
	h.checkDeleteButtonHover(event.Position)
	h.updateHoverState()
}

// checkDeleteButtonHover determines if the mouse is over the delete button
func (h *HistoryListItem) checkDeleteButtonHover(pos fyne.Position) {
	if h.deleteButton == nil {
		h.deleteHovered = false
		return
	}
	
	// Get delete button position and size
	buttonPos := h.deleteButton.Position()
	buttonSize := h.deleteButton.Size()
	
	// Check if mouse position is within delete button bounds
	// Add generous padding for easier interaction and better accessibility
	padding := float32(8)
	h.deleteHovered = pos.X >= buttonPos.X-padding &&
		pos.X <= buttonPos.X+buttonSize.Width+padding &&
		pos.Y >= buttonPos.Y-padding &&
		pos.Y <= buttonPos.Y+buttonSize.Height+padding
}

// UpdateItem updates the clipboard item content
func (h *HistoryListItem) UpdateItem(newItem ClipboardItem) {
	h.item = newItem
	h.createContent() // Recreate content with new item
	h.Refresh()
}

// UpdateIndex updates the index of the item
func (h *HistoryListItem) UpdateIndex(newIndex int) {
	h.index = newIndex
}

// historyListItemRenderer renders the HistoryListItem widget
type historyListItemRenderer struct {
	item      *HistoryListItem
	container *fyne.Container
}

// Layout arranges the child widgets
func (r *historyListItemRenderer) Layout(size fyne.Size) {
	r.container.Resize(size)
	
	// Ensure background fills the entire widget area
	if r.item.background != nil {
		r.item.background.Resize(size)
	}
}

// MinSize returns the minimum size required with height constraints
func (r *historyListItemRenderer) MinSize() fyne.Size {
	minSize := r.container.MinSize()
	
	// Ensure minimum height for readability but cap maximum height
	const minHeight = 40
	const maxHeight = 80
	
	if minSize.Height < minHeight {
		minSize.Height = minHeight
	} else if minSize.Height > maxHeight {
		minSize.Height = maxHeight
	}
	
	// Add some padding for better visual appearance
	minSize.Width += 8  // Extra width for padding
	minSize.Height += 4 // Extra height for padding
	
	return minSize
}

// Refresh updates the visual appearance
func (r *historyListItemRenderer) Refresh() {
	// Refresh all components
	if r.item.background != nil {
		r.item.background.Refresh()
	}
	if r.item.textWidget != nil {
		r.item.textWidget.Refresh()
	}
	if r.item.deleteButton != nil {
		r.item.deleteButton.Refresh()
	}
	r.container.Refresh()
}

// Objects returns the child objects
func (r *historyListItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}

// Destroy cleans up resources
func (r *historyListItemRenderer) Destroy() {
	// Clean up any resources if needed
	r.item.background = nil
	r.item.textWidget = nil
	r.item.deleteButton = nil
}