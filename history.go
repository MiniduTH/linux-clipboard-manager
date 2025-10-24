package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// ClipboardItemType represents the type of clipboard content
type ClipboardItemType string

const (
	ItemTypeText  ClipboardItemType = "text"
	ItemTypeImage ClipboardItemType = "image"
)

// ClipboardItem represents a single clipboard entry that can be text or image
type ClipboardItem struct {
	Type      ClipboardItemType `json:"type"`
	Content   string            `json:"content"`   // Text content or base64 encoded image
	Timestamp time.Time         `json:"timestamp"`
	ImageMeta *ImageMetadata    `json:"image_meta,omitempty"` // Metadata for images
}

// ImageMetadata contains metadata about image clipboard items
type ImageMetadata struct {
	Format string `json:"format"` // "png", "jpeg", etc.
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int    `json:"size"` // Size in bytes of original image
}

var (
	history   []ClipboardItem
	historyMu sync.RWMutex
)

const maxHistory = 50

// addToHistory adds a new text item to clipboard history
func addToHistory(text string) {
	historyMu.Lock()
	defer historyMu.Unlock()
	
	// Clean and validate the text
	text = strings.TrimSpace(text)
	if text == "" || !utf8.ValidString(text) {
		return
	}
	
	// Create new text item
	newItem := ClipboardItem{
		Type:      ItemTypeText,
		Content:   text,
		Timestamp: time.Now(),
	}
	
	// Skip if it's the same as the last item
	if len(history) > 0 && history[len(history)-1].Content == text {
		return
	}
	
	// Save to database
	if err := saveClipboardItem(newItem); err != nil {
		fmt.Printf("Error saving text to database: %v\n", err)
		return
	}
	
	// Update in-memory history
	refreshHistoryFromDB()
}

// addImageToHistory adds a new image item to clipboard history
func addImageToHistory(imageData []byte, format string) {
	historyMu.Lock()
	defer historyMu.Unlock()
	
	if len(imageData) == 0 {
		return
	}
	
	// Decode image to get metadata
	var img image.Image
	var err error
	
	switch format {
	case "png":
		img, err = png.Decode(strings.NewReader(string(imageData)))
	case "jpeg", "jpg":
		img, err = jpeg.Decode(strings.NewReader(string(imageData)))
	default:
		// Try to decode as PNG first, then JPEG
		if img, err = png.Decode(strings.NewReader(string(imageData))); err != nil {
			img, err = jpeg.Decode(strings.NewReader(string(imageData)))
			if err == nil {
				format = "jpeg"
			}
		} else {
			format = "png"
		}
	}
	
	if err != nil {
		fmt.Printf("Error decoding image: %v\n", err)
		return
	}
	
	// Create image metadata
	bounds := img.Bounds()
	imageMeta := &ImageMetadata{
		Format: format,
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
		Size:   len(imageData),
	}
	
	// Encode image data as base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)
	
	// Create new image item
	newItem := ClipboardItem{
		Type:      ItemTypeImage,
		Content:   base64Data,
		Timestamp: time.Now(),
		ImageMeta: imageMeta,
	}
	
	// Skip if it's the same as the last item (compare base64 content)
	if len(history) > 0 && history[len(history)-1].Content == base64Data {
		return
	}
	
	// Save to database
	if err := saveClipboardItem(newItem); err != nil {
		fmt.Printf("Error saving image to database: %v\n", err)
		return
	}
	
	// Update in-memory history
	refreshHistoryFromDB()
}

// Remove specific item from history by index
func removeHistoryItem(index int) {
	historyMu.Lock()
	defer historyMu.Unlock()
	
	if index < 0 || index >= len(history) {
		fmt.Printf("Invalid index %d for history removal\n", index)
		return
	}
	
	// Get the item to remove
	item := history[index]
	
	// Remove from database
	if err := deleteClipboardItem(item.Content, item.Type); err != nil {
		fmt.Printf("Error removing item from database: %v\n", err)
		return
	}
	
	// Update in-memory history
	refreshHistoryFromDB()
	fmt.Printf("Removed history item at index %d\n", index)
}

// Clear all history with optional UI callback
func clearHistory(onComplete ...func()) error {
	historyMu.Lock()
	defer historyMu.Unlock()
	
	// Clear from database
	if err := clearClipboardHistory(); err != nil {
		fmt.Printf("Error clearing history from database: %v\n", err)
		return err
	}
	
	// Clear in-memory history
	history = []ClipboardItem{}
	fmt.Println("Clipboard history cleared.")
	
	// Call UI callback if provided
	for _, callback := range onComplete {
		if callback != nil {
			callback()
		}
	}
	
	return nil
}

func getHistoryFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to /tmp if home directory is not available
		return "/tmp/clipboard-history.json"
	}
	
	dir := filepath.Join(home, ".local", "share", "clipboard-manager")
	if err := os.MkdirAll(dir, 0755); err != nil {
		// Fallback to /tmp if we can't create the directory
		return "/tmp/clipboard-history.json"
	}
	
	return filepath.Join(dir, "history.json")
}

func saveHistory() {
	// This function is now deprecated as we use SQLite database
	// Keeping it for backward compatibility but it does nothing
	// The database is automatically saved when items are added/removed
}

func loadHistory() {
	// Initialize database
	if err := initDatabase(); err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		// Fallback to empty history
		historyMu.Lock()
		history = []ClipboardItem{}
		historyMu.Unlock()
		return
	}
	
	// Try to migrate from JSON if it exists
	if err := migrateFromJSON(); err != nil {
		fmt.Printf("Warning: failed to migrate from JSON: %v\n", err)
	}
	
	// Load history from database
	refreshHistoryFromDB()
}

// Get history length safely
func getHistoryLength() int {
	historyMu.RLock()
	defer historyMu.RUnlock()
	return len(history)
}

// Get history copy safely for UI operations
func getHistoryCopy() []ClipboardItem {
	historyMu.RLock()
	defer historyMu.RUnlock()
	
	historyCopy := make([]ClipboardItem, len(history))
	copy(historyCopy, history)
	return historyCopy
}

// restoreImageToClipboard restores an image from base64 storage to clipboard
// This function is kept for backward compatibility but the actual restoration
// is now handled by restoreImageToSystemClipboard in image_clipboard.go
func restoreImageToClipboard(base64Data string, format string) error {
	// Decode base64 data
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("failed to decode base64 image: %v", err)
	}
	
	// Use the system clipboard restoration function
	return restoreImageToSystemClipboard(imageData, format)
}
// refreshHistoryFromDB loads the current history from database into memory
func refreshHistoryFromDB() {
	loadedHistory, err := loadClipboardHistory()
	if err != nil {
		fmt.Printf("Error loading history from database: %v\n", err)
		return
	}
	
	// Update in-memory history (lock should already be held by caller)
	history = loadedHistory
}