package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
)

func main() {
	// Check if we're in a proper environment
	if !checkEnvironment() {
		log.Fatal("This application requires a graphical environment (X11 or Wayland) and clipboard utilities")
	}

	loadHistory() // load previous data

	if len(os.Args) > 1 && os.Args[1] == "show" {
		// popup mode
		if err := showPopup(); err != nil {
			fmt.Printf("Error showing popup: %v\n", err)
			// Fallback to terminal mode
			showTerminalHistory()
		}
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "list" {
		showTerminalHistory()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "help" {
		fmt.Println("Clipboard Manager for Linux")
		fmt.Println("Usage:")
		fmt.Println("  ./clipboard-manager        - Start with system integration")
		fmt.Println("  ./clipboard-manager show   - Show GUI history")
		fmt.Println("  ./clipboard-manager list   - Show terminal history")
		fmt.Println("  ./clipboard-manager tray   - Start with system tray")
		fmt.Println("  ./clipboard-manager daemon - Start in background (no GUI)")
		fmt.Println("  ./clipboard-manager help   - Show this help")
		fmt.Println()
		fmt.Println("System Integration:")
		fmt.Println("  The app will try to set up Super+Z hotkey automatically")
		fmt.Println("  Use 'tray' mode for system tray integration")
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "tray" {
		fmt.Println("Starting Clipboard Manager with system tray...")
		runWithSystemTray()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "daemon" {
		fmt.Println("Clipboard Manager started in daemon mode (no hotkeys).")
		fmt.Println("Press Ctrl+C to stop.")
		
		// graceful exit (Ctrl+C)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			fmt.Println("\nSaving clipboard history...")
			saveHistory()
			os.Exit(0)
		}()

		watchClipboard() // start daemon without hotkeys
		return
	}

	fmt.Println("Clipboard Manager started with system integration.")

	// Start clipboard monitoring in background
	go watchClipboard()

	// graceful exit (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nSaving clipboard history...")
		saveHistory()
		os.Exit(0)
	}()

	// Setup system hotkeys and keep running
	setupLinuxHotkeys()
}

// Check if we have the necessary environment and tools
func checkEnvironment() bool {
	// Check for display server
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		return false
	}

	// Check for clipboard utilities
	tools := []string{"xclip", "xsel", "wl-clipboard"}
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err == nil {
			return true
		}
	}
	
	// If no external tools, try the Go clipboard library
	_, err := clipboard.ReadAll()
	return err == nil
}

// watches clipboard continuously
func watchClipboard() {
	last := ""
	saveTimer := time.NewTicker(10 * time.Second) // save every 10 seconds
	defer saveTimer.Stop()
	
	go func() {
		for range saveTimer.C {
			saveHistory()
		}
	}()
	
	errorCount := 0
	maxErrors := 5
	
	for {
		text, err := clipboard.ReadAll()
		if err != nil {
			errorCount++
			if errorCount <= maxErrors {
				fmt.Printf("Clipboard read error (%d/%d): %v\n", errorCount, maxErrors, err)
			}
			if errorCount >= maxErrors {
				fmt.Println("Too many clipboard errors, reducing check frequency...")
				time.Sleep(5 * time.Second)
			} else {
				time.Sleep(2 * time.Second)
			}
			continue
		}
		
		// Reset error count on successful read
		errorCount = 0
		
		// Clean up the text
		text = strings.TrimSpace(text)
		
		// Skip empty or very short text
		if len(text) < 2 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		
		// Skip if same as last
		if text == last {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		
		// Skip if it looks like system clipboard noise
		if isSystemNoise(text) {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		
		addToHistory(text)
		last = text
		
		// Only show notification for meaningful content
		if len(text) > 5 {
			displayText := text
			if len(displayText) > 60 {
				displayText = displayText[:60] + "..."
			}
			// Replace newlines for cleaner output
			displayText = strings.ReplaceAll(displayText, "\n", " ")
			fmt.Printf("📋 Copied: %s\n", displayText)
		}
		
		time.Sleep(500 * time.Millisecond)
	}
}

// Check if text is likely system noise
func isSystemNoise(text string) bool {
	// Skip very short text
	if len(text) < 3 {
		return true
	}
	
	// Skip common system clipboard noise patterns
	noisePatterns := []string{
		"signal\"",
		"syscall\"",
		"time\"",
		"import",
		"package",
		"func",
	}
	
	for _, pattern := range noisePatterns {
		if strings.Contains(text, pattern) && len(text) < 50 {
			return true
		}
	}
	
	return false
}

// Terminal-based history viewer as fallback
func showTerminalHistory() {
	historyMu.RLock()
	defer historyMu.RUnlock()
	
	if len(history) == 0 {
		fmt.Println("No clipboard history yet.")
		return
	}

	fmt.Println("\nClipboard History (newest first):")
	fmt.Println(strings.Repeat("-", 50))
	
	for i := len(history) - 1; i >= 0; i-- {
		item := history[i]
		if len(item) > 80 {
			item = item[:80] + "..."
		}
		// Replace newlines with spaces for better terminal display
		item = strings.ReplaceAll(item, "\n", " ")
		fmt.Printf("%2d: %s\n", len(history)-i, item)
	}
	
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Total items: %d\n", len(history))
}
