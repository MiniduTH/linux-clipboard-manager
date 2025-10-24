package main

import (
	"bytes"
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
		// Ensure daemon is running before showing GUI
		ensureDaemonRunning()
		
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
		fmt.Println("  ./clipboard-manager              - Start with system integration")
		fmt.Println("  ./clipboard-manager show         - Show GUI history (auto-starts daemon)")
		fmt.Println("  ./clipboard-manager list         - Show terminal history")
		fmt.Println("  ./clipboard-manager tray         - Start with system tray")
		fmt.Println("  ./clipboard-manager daemon       - Start in background (no GUI)")
		fmt.Println("  ./clipboard-manager daemon-text-only - Start daemon (text only, no image monitoring)")
		fmt.Println("  ./clipboard-manager daemon-minimal   - Start daemon (ultra-minimal polling)")
		fmt.Println("  ./clipboard-manager daemon-passive   - Start daemon (no auto-monitoring)")
		fmt.Println("  ./clipboard-manager daemon-only      - Start daemon only (no hotkeys)")
		fmt.Println("  ./clipboard-manager capture          - Manually capture current clipboard")
		fmt.Println("  ./clipboard-manager status       - Show daemon status")
		fmt.Println("  ./clipboard-manager stop         - Stop daemon")
		fmt.Println("  ./clipboard-manager help         - Show this help")
		fmt.Println()
		fmt.Println("System Integration:")
		fmt.Println("  The app will try to set up Ctrl+Shift+V hotkey automatically")
		fmt.Println("  Use 'tray' mode for system tray integration")
		fmt.Println("  Use 'daemon-only' to run without hotkeys or GUI")
		fmt.Println("  The daemon runs in background to monitor clipboard")
		fmt.Println("  Use 'daemon-minimal' or 'daemon-passive' if you experience clicking/menu issues")
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "status" {
		showDaemonStatus()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "stop" {
		stopDaemon()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "tray" {
		fmt.Println("Starting Clipboard Manager with system tray...")
		runWithSystemTray()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "daemon" {
		// Check if daemon is already running
		if isDaemonRunning() {
			fmt.Println("Clipboard daemon is already running. Use 'clipboard-manager stop' to stop it first.")
			os.Exit(1)
		}
		
		fmt.Println("Clipboard Manager started in daemon mode (no hotkeys).")
		fmt.Println("Press Ctrl+C to stop.")
		
		// graceful exit (Ctrl+C)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			fmt.Println("\nClosing database...")
			closeDatabase()
			os.Exit(0)
		}()

		watchClipboard() // start daemon without hotkeys
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "daemon-text-only" {
		fmt.Println("Clipboard Manager started in text-only daemon mode (no hotkeys, no image monitoring).")
		fmt.Println("Press Ctrl+C to stop.")
		
		// graceful exit (Ctrl+C)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			fmt.Println("\nClosing database...")
			closeDatabase()
			os.Exit(0)
		}()

		watchClipboardTextOnly() // start daemon without hotkeys and image monitoring
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "daemon-minimal" {
		fmt.Println("Clipboard Manager started in minimal daemon mode (ultra-conservative polling).")
		fmt.Println("Press Ctrl+C to stop.")
		
		// graceful exit (Ctrl+C)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			fmt.Println("\nClosing database...")
			closeDatabase()
			os.Exit(0)
		}()

		watchClipboardMinimal() // start daemon with ultra-minimal monitoring
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "daemon-passive" {
		fmt.Println("Clipboard Manager started in passive mode (no automatic monitoring).")
		fmt.Println("Use './clipboard-manager capture' to manually capture current clipboard.")
		fmt.Println("Press Ctrl+C to stop.")
		
		// graceful exit (Ctrl+C)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			fmt.Println("\nClosing database...")
			closeDatabase()
			os.Exit(0)
		}()

		// Just keep the process alive without monitoring
		select {}
	}

	if len(os.Args) > 1 && os.Args[1] == "daemon-only" {
		fmt.Println("Clipboard Manager started in daemon-only mode (no hotkeys, no GUI).")
		fmt.Println("Text-only monitoring to avoid any window creation.")
		fmt.Println("Use './clipboard-manager show' to open GUI manually.")
		fmt.Println("Press Ctrl+C to stop.")
		
		// graceful exit (Ctrl+C)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			fmt.Println("\nClosing database...")
			closeDatabase()
			os.Exit(0)
		}()

		watchClipboardTextOnly() // start daemon without hotkeys, GUI, or image monitoring
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "capture" {
		// Manual clipboard capture
		text, err := clipboard.ReadAll()
		if err != nil {
			fmt.Printf("Error reading clipboard: %v\n", err)
			os.Exit(1)
		}
		
		text = strings.TrimSpace(text)
		if len(text) >= 2 && !isSystemNoise(text) {
			addToHistory(text)
			fmt.Printf("📋 Captured: %s\n", text)
		} else {
			fmt.Println("No meaningful text found in clipboard")
		}
		return
	}

	// Default mode: Start with system integration and hotkeys
	// Check if daemon is already running
	if isDaemonRunning() {
		fmt.Println("Clipboard daemon is already running. Use 'clipboard-manager stop' to stop it first.")
		os.Exit(1)
	}
	
	fmt.Println("Clipboard Manager started with system integration.")
	fmt.Println("Note: This will set up hotkeys but won't show GUI automatically.")
	fmt.Println("Use 'clipboard-manager show' or press Ctrl+Shift+V to open the GUI.")

	// Start clipboard monitoring in background
	go watchClipboard()

	// graceful exit (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nClosing database...")
		closeDatabase()
		os.Exit(0)
	}()

	// Setup system hotkeys and keep running (but don't auto-show GUI)
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
	lastText := ""
	var lastImageData []byte
	
	errorCount := 0
	maxErrors := 5
	imageCheckCounter := 0
	
	// Disable image monitoring if it causes issues (can be made configurable)
	enableImageMonitoring := true
	
	for {
		// Check for text content first
		text, textErr := clipboard.ReadAll()
		
		if textErr == nil {
			// Reset error count on successful read
			errorCount = 0
			
			// Clean up the text
			text = strings.TrimSpace(text)
			
			// Check if we have new meaningful text content
			if len(text) >= 2 && text != lastText && !isSystemNoise(text) {
				addToHistory(text)
				lastText = text
				
				// Show notification for meaningful content
				if len(text) > 5 {
					displayText := text
					if len(displayText) > 60 {
						displayText = displayText[:60] + "..."
					}
					// Replace newlines for cleaner output
					displayText = strings.ReplaceAll(displayText, "\n", " ")
					fmt.Printf("📋 Text copied: %s\n", displayText)
				}
			}
		} else {
			errorCount++
			if errorCount <= maxErrors {
				fmt.Printf("Clipboard text read error (%d/%d): %v\n", errorCount, maxErrors, textErr)
			}
		}
		
		// Check for image content less frequently (every 3rd iteration) and only if enabled
		if enableImageMonitoring {
			imageCheckCounter++
			if imageCheckCounter%3 == 0 {
				if imageData, format, imageErr := detectImageInClipboard(); imageErr == nil {
					// Reset error count on successful image read
					errorCount = 0
					
					// Check if we have new image content (compare first 1KB for efficiency)
					compareSize := 1024
					if len(imageData) < compareSize {
						compareSize = len(imageData)
					}
					
					isNewImage := len(lastImageData) == 0 || 
						len(lastImageData) != len(imageData) ||
						!bytes.Equal(imageData[:compareSize], lastImageData[:compareSize])
					
					if isNewImage {
						addImageToHistory(imageData, format)
						lastImageData = make([]byte, len(imageData))
						copy(lastImageData, imageData)
						
						// Show notification for image
						sizeKB := len(imageData) / 1024
						fmt.Printf("📋 Image copied: %s (%d KB)\n", format, sizeKB)
					}
				}
			}
		}
		
		// Handle errors and sleep timing
		if textErr != nil && errorCount >= maxErrors {
			fmt.Println("Too many clipboard errors, reducing check frequency...")
			time.Sleep(10 * time.Second)
		} else if textErr != nil {
			time.Sleep(3 * time.Second)
		} else {
			// Normal polling - much less aggressive
			time.Sleep(2 * time.Second)
		}
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
		itemNum := len(history) - i
		
		if item.Type == ItemTypeText {
			content := item.Content
			if len(content) > 80 {
				content = content[:80] + "..."
			}
			// Replace newlines with spaces for better terminal display
			content = strings.ReplaceAll(content, "\n", " ")
			fmt.Printf("%2d: [TEXT] %s\n", itemNum, content)
		} else if item.Type == ItemTypeImage {
			if item.ImageMeta != nil {
				fmt.Printf("%2d: [IMAGE] %s %dx%d (%d KB)\n", 
					itemNum, 
					strings.ToUpper(item.ImageMeta.Format),
					item.ImageMeta.Width, 
					item.ImageMeta.Height,
					item.ImageMeta.Size/1024)
			} else {
				fmt.Printf("%2d: [IMAGE] Unknown format\n", itemNum)
			}
		}
	}
	
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Total items: %d\n", len(history))
}
// watchClipboardTextOnly watches clipboard for text content only (no image monitoring)
func watchClipboardTextOnly() {
	lastText := ""
	
	errorCount := 0
	maxErrors := 5
	
	for {
		// Check for text content only
		text, textErr := clipboard.ReadAll()
		
		if textErr == nil {
			// Reset error count on successful read
			errorCount = 0
			
			// Clean up the text
			text = strings.TrimSpace(text)
			
			// Check if we have new meaningful text content
			if len(text) >= 2 && text != lastText && !isSystemNoise(text) {
				addToHistory(text)
				lastText = text
				
				// Show notification for meaningful content
				if len(text) > 5 {
					displayText := text
					if len(displayText) > 60 {
						displayText = displayText[:60] + "..."
					}
					// Replace newlines for cleaner output
					displayText = strings.ReplaceAll(displayText, "\n", " ")
					fmt.Printf("📋 Text copied: %s\n", displayText)
				}
			}
		} else {
			errorCount++
			if errorCount <= maxErrors {
				fmt.Printf("Clipboard text read error (%d/%d): %v\n", errorCount, maxErrors, textErr)
			}
		}
		
		// Handle errors and sleep timing - more conservative for text-only mode
		if textErr != nil && errorCount >= maxErrors {
			fmt.Println("Too many clipboard errors, reducing check frequency...")
			time.Sleep(15 * time.Second)
		} else if textErr != nil {
			time.Sleep(5 * time.Second)
		} else {
			// Even more conservative polling for text-only mode
			time.Sleep(3 * time.Second)
		}
	}
}

// watchClipboardMinimal watches clipboard with ultra-conservative polling to avoid system interference
func watchClipboardMinimal() {
	lastText := ""
	
	errorCount := 0
	maxErrors := 3
	
	fmt.Println("Using ultra-minimal clipboard monitoring (10-second intervals)")
	fmt.Println("This mode minimizes system interference but may miss rapid clipboard changes")
	
	for {
		// Only check clipboard every 10 seconds to minimize system interference
		time.Sleep(10 * time.Second)
		
		// Check for text content only
		text, textErr := clipboard.ReadAll()
		
		if textErr == nil {
			// Reset error count on successful read
			errorCount = 0
			
			// Clean up the text
			text = strings.TrimSpace(text)
			
			// Check if we have new meaningful text content
			if len(text) >= 2 && text != lastText && !isSystemNoise(text) {
				addToHistory(text)
				lastText = text
				
				// Show notification for meaningful content
				if len(text) > 5 {
					displayText := text
					if len(displayText) > 60 {
						displayText = displayText[:60] + "..."
					}
					// Replace newlines for cleaner output
					displayText = strings.ReplaceAll(displayText, "\n", " ")
					fmt.Printf("📋 Text copied: %s\n", displayText)
				}
			}
		} else {
			errorCount++
			if errorCount <= maxErrors {
				fmt.Printf("Clipboard text read error (%d/%d): %v\n", errorCount, maxErrors, textErr)
			}
			
			// If too many errors, wait even longer
			if errorCount >= maxErrors {
				fmt.Println("Too many clipboard errors, waiting 30 seconds...")
				time.Sleep(30 * time.Second)
			}
		}
	}
}