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
		fmt.Println("  ./clipboard-manager        - Start clipboard watcher")
		fmt.Println("  ./clipboard-manager show   - Show GUI history")
		fmt.Println("  ./clipboard-manager list   - Show terminal history")
		fmt.Println("  ./clipboard-manager help   - Show this help")
		return
	}

	fmt.Println("Clipboard Manager started. Press Ctrl+C to stop.")

	// graceful exit (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nSaving clipboard history...")
		saveHistory()
		os.Exit(0)
	}()

	watchClipboard() // start daemon
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
	
	for {
		text, err := clipboard.ReadAll()
		if err != nil {
			fmt.Printf("Clipboard read error: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}
		
		// Clean up the text
		text = strings.TrimSpace(text)
		
		if text != "" && text != last && len(text) > 0 {
			addToHistory(text)
			last = text
			fmt.Printf("Copied: %.50s", text)
			if len(text) > 50 {
				fmt.Print("...")
			}
			fmt.Println()
		}
		time.Sleep(1 * time.Second) // Reduced frequency to be less resource intensive
	}
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
