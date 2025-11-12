package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

// refreshUI updates the window content with current history state
func refreshUI(w fyne.Window) {
	historyLen := getHistoryLength()
	historyCopy := getHistoryCopy()
	
	if historyLen == 0 {
		// Show empty state
		emptyLabel := widget.NewLabel("No clipboard history yet.\nStart copying text to see it here!")
		emptyLabel.Alignment = fyne.TextAlignCenter
		
		closeBtn := widget.NewButton("Close", func() {
			w.Close()
		})
		
		emptyContent := container.NewVBox(
			widget.NewLabel("Clipboard History"),
			widget.NewSeparator(),
			container.NewCenter(emptyLabel),
			widget.NewSeparator(),
			container.NewCenter(closeBtn),
		)
		
		w.SetContent(emptyContent)
		w.Canvas().Refresh(emptyContent)
		return
	}
	
	// Create custom history list items
	var historyItems []fyne.CanvasObject
	
	// Define handlers for item selection, deletion, and editing
	onSelect := func(index int) {
		// Convert UI index to history array index (newest first display)
		historyIndex := historyLen - 1 - index
		originalItem := historyCopy[historyIndex]
		
		if originalItem.Type == ItemTypeText {
			// Restore text to clipboard
			if err := clipboard.WriteAll(originalItem.Content); err != nil {
				fmt.Printf("Error writing text to clipboard: %v\n", err)
				return
			}
			fmt.Printf("Restored text to clipboard: %.50s", originalItem.Content)
			if len(originalItem.Content) > 50 {
				fmt.Print("...")
			}
			fmt.Println()
		} else if originalItem.Type == ItemTypeImage {
			// Restore image to clipboard
			imageData, err := base64.StdEncoding.DecodeString(originalItem.Content)
			if err != nil {
				fmt.Printf("Error decoding image: %v\n", err)
				return
			}
			
			format := "png"
			if originalItem.ImageMeta != nil {
				format = originalItem.ImageMeta.Format
			}
			
			if err := restoreImageToSystemClipboard(imageData, format); err != nil {
				fmt.Printf("Error restoring image to clipboard: %v\n", err)
				// Fallback: show message about image restoration limitation
				fmt.Println("Note: Image clipboard restoration may have limited support on this system.")
			} else {
				fmt.Printf("Restored %s image to clipboard", strings.ToUpper(format))
				if originalItem.ImageMeta != nil {
					fmt.Printf(" (%dx%d)", originalItem.ImageMeta.Width, originalItem.ImageMeta.Height)
				}
				fmt.Println()
			}
		}
		
		w.Close()
	}
	
	onDelete := func(index int) {
		// Convert UI index to history array index (newest first display)
		historyIndex := historyLen - 1 - index
		removeHistoryItem(historyIndex)
		
		// Refresh the UI in place instead of closing and reopening
		refreshUI(w)
	}
	
	onEdit := func(index int) {
		// Convert UI index to history array index (newest first display)
		historyIndex := historyLen - 1 - index
		originalItem := historyCopy[historyIndex]
		
		// Only allow editing text items
		if originalItem.Type != ItemTypeText {
			return
		}
		
		// Show edit dialog
		showEditDialog(w, historyIndex, originalItem.Content)
	}
	
	// Create custom list items (newest first)
	for i := 0; i < historyLen; i++ {
		item := historyCopy[historyLen-1-i]
		historyItem := NewHistoryListItem(item, i, onDelete, onSelect, onEdit)
		historyItems = append(historyItems, historyItem)
	}
	
	// Create scrollable container for the history items using VBox layout
	listContainer := container.NewVBox(historyItems...)
	paddedContainer := container.NewPadded(listContainer)
	scrollContainer := container.NewScroll(paddedContainer)
	scrollContainer.SetMinSize(fyne.NewSize(680, 400))
	scrollContainer.Direction = container.ScrollVerticalOnly

	// Create buttons with improved styling
	closeBtn := widget.NewButton("Close", func() {
		w.Close()
	})
	closeBtn.Importance = widget.HighImportance

	buttonContainer := container.NewHBox(closeBtn)
	headerLabel := widget.NewLabel(fmt.Sprintf("Clipboard History (%d items) - Click to copy, Edit/Delete with buttons", historyLen))
	headerLabel.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		headerLabel,
		widget.NewSeparator(),
		scrollContainer,
		widget.NewSeparator(),
		buttonContainer,
	)

	w.SetContent(content)
	w.Canvas().Refresh(content)
}

// showEditDialog displays a dialog for editing text content
func showEditDialog(parent fyne.Window, historyIndex int, currentContent string) {
	// Create a multi-line entry for editing
	entry := widget.NewMultiLineEntry()
	entry.SetText(currentContent)
	entry.Wrapping = fyne.TextWrapWord
	entry.SetMinRowsVisible(10)
	
	// Create dialog variable that we'll populate
	var dialog *widget.PopUp
	
	// Create save button handler
	saveHandler := func() {
		newContent := strings.TrimSpace(entry.Text)
		if newContent == "" {
			// Show error if content is empty
			errDialog := widget.NewModalPopUp(
				container.NewVBox(
					widget.NewLabel("Error: Content cannot be empty"),
					widget.NewButton("OK", func() {
						// Just close the error dialog, keep edit dialog open
					}),
				),
				parent.Canvas(),
			)
			errDialog.Show()
			return
		}
		
		// Update the item
		if err := editHistoryItem(historyIndex, newContent); err != nil {
			fmt.Printf("Error editing item: %v\n", err)
			errDialog := widget.NewModalPopUp(
				container.NewVBox(
					widget.NewLabel(fmt.Sprintf("Error: %v", err)),
					widget.NewButton("OK", func() {
						// Just close the error dialog, keep edit dialog open
					}),
				),
				parent.Canvas(),
			)
			errDialog.Show()
			return
		}
		
		// Close dialog and refresh UI
		dialog.Hide()
		refreshUI(parent)
	}
	
	// Create cancel button handler
	cancelHandler := func() {
		dialog.Hide()
	}
	
	// Create dialog with save and cancel buttons
	dialog = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Edit Clipboard Item"),
			widget.NewSeparator(),
			container.NewScroll(entry),
			widget.NewSeparator(),
			container.NewHBox(
				widget.NewButton("Save", saveHandler),
				widget.NewButton("Cancel", cancelHandler),
			),
		),
		parent.Canvas(),
	)
	
	dialog.Resize(fyne.NewSize(600, 400))
	dialog.Show()
}

func showPopup() error {
	// Try to create the app with error handling with proper ID
	a := app.NewWithID("com.clipboard.manager")
	if a == nil {
		return fmt.Errorf("failed to create GUI application")
	}
	
	// Apply custom theme for better hover states and contrast
	a.Settings().SetTheme(NewCustomTheme())
	a.SetIcon(nil) // Avoid icon loading issues
	w := a.NewWindow("Clipboard History")
	w.Resize(fyne.NewSize(700, 600))
	w.CenterOnScreen()

	// Use the refreshUI function to set initial content
	refreshUI(w)
	
	w.ShowAndRun()
	return nil
}


