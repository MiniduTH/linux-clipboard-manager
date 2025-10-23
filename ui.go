package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

func showPopup() error {
	historyMu.RLock()
	historyLen := len(history)
	historyCopy := make([]string, historyLen)
	copy(historyCopy, history)
	historyMu.RUnlock()
	
	if historyLen == 0 {
		fmt.Println("No clipboard history yet.")
		return nil
	}

	// Try to create the app with error handling
	a := app.New()
	if a == nil {
		return fmt.Errorf("failed to create GUI application")
	}
	
	a.SetIcon(nil) // Avoid icon loading issues
	w := a.NewWindow("Clipboard History")
	w.Resize(fyne.NewSize(600, 500))
	w.CenterOnScreen()

	// Create a scrollable list
	list := widget.NewList(
		func() int { return historyLen },
		func() fyne.CanvasObject { 
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextWrapWord
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			item := historyCopy[historyLen-1-i]
			// Clean up the text for display
			item = strings.ReplaceAll(item, "\n", " ")
			item = strings.ReplaceAll(item, "\t", " ")
			
			// Truncate long text for display
			if len(item) > 120 {
				item = item[:120] + "..."
			}
			
			label := o.(*widget.Label)
			label.SetText(fmt.Sprintf("%d: %s", historyLen-i, item))
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		item := historyCopy[historyLen-1-id]
		if err := clipboard.WriteAll(item); err != nil {
			fmt.Printf("Error writing to clipboard: %v\n", err)
			return
		}
		fmt.Printf("Restored to clipboard: %.50s", item)
		if len(item) > 50 {
			fmt.Print("...")
		}
		fmt.Println()
		w.Close()
	}

	// Create buttons
	clearBtn := widget.NewButton("Clear History", func() {
		clearHistory()
		w.Close()
	})
	
	closeBtn := widget.NewButton("Close", func() {
		w.Close()
	})

	buttonContainer := container.NewHBox(clearBtn, closeBtn)

	content := container.NewVBox(
		widget.NewLabel("Clipboard History (click to restore):"),
		widget.NewSeparator(),
		list,
		widget.NewSeparator(),
		buttonContainer,
	)

	w.SetContent(content)
	
	// Handle window close
	w.SetCloseIntercept(func() {
		w.Close()
	})

	w.ShowAndRun()
	return nil
}
