package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

// System tray integration
func setupSystemTray() {
	a := app.NewWithID("com.clipboard.manager")
	a.SetIcon(nil)
	
	if desk, ok := a.(desktop.App); ok {
		menu := fyne.NewMenu("Clipboard Manager",
			fyne.NewMenuItem("Show History", func() {
				go func() {
					if err := showPopup(); err != nil {
						fmt.Printf("GUI failed: %v\n", err)
						showTerminalHistory()
					}
				}()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Clear History", func() {
				if err := clearHistory(); err != nil {
					fmt.Printf("Error clearing history: %v\n", err)
				}
			}),
			fyne.NewMenuItem("Quit", func() {
				saveHistory()
				a.Quit()
			}),
		)
		desk.SetSystemTrayMenu(menu)
		
		fmt.Println("✓ System tray icon created")
		fmt.Println("Right-click the tray icon to access clipboard history")
	}
	
	// Keep the app running
	a.Run()
}

// Run with system tray (non-blocking)
func runWithSystemTray() {
	go func() {
		// Start clipboard monitoring
		watchClipboard()
	}()
	
	// Setup system tray (this blocks)
	setupSystemTray()
}