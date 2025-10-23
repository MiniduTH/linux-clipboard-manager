package main

import (
	"encoding/base64"
	"encoding/json"
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
	
	// Remove duplicates from history
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Content == text && history[i].Type == ItemTypeText {
			history = append(history[:i], history[i+1:]...)
			break
		}
	}
	
	// Add to end
	history = append(history, newItem)
	
	// Maintain max size
	if len(history) > maxHistory {
		history = history[len(history)-maxHistory:]
	}
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
	
	// Remove duplicates from history
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Content == base64Data && history[i].Type == ItemTypeImage {
			history = append(history[:i], history[i+1:]...)
			break
		}
	}
	
	// Add to end
	history = append(history, newItem)
	
	// Maintain max size
	if len(history) > maxHistory {
		history = history[len(history)-maxHistory:]
	}
}

// Remove specific item from history by index
func removeHistoryItem(index int) {
	historyMu.Lock()
	
	if index < 0 || index >= len(history) {
		historyMu.Unlock()
		fmt.Printf("Invalid index %d for history removal\n", index)
		return
	}
	
	// Remove item at index
	history = append(history[:index], history[index+1:]...)
	historyMu.Unlock()
	
	// Save updated history after releasing the lock to avoid deadlock
	saveHistory()
	fmt.Printf("Removed history item at index %d\n", index)
}

// Clear all history with optional UI callback
func clearHistory(onComplete ...func()) {
	historyMu.Lock()
	history = []ClipboardItem{}
	historyMu.Unlock()
	
	// Save history after releasing the lock to avoid deadlock
	saveHistory()
	fmt.Println("Clipboard history cleared.")
	
	// Call UI callback if provided
	for _, callback := range onComplete {
		if callback != nil {
			callback()
		}
	}
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
	historyMu.RLock()
	historyCopy := make([]ClipboardItem, len(history))
	copy(historyCopy, history)
	historyMu.RUnlock()
	
	fpath := getHistoryFile()
	f, err := os.Create(fpath)
	if err != nil {
		fmt.Printf("Save error: %v\n", err)
		return
	}
	defer f.Close()
	
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(historyCopy); err != nil {
		fmt.Printf("Encode error: %v\n", err)
	}
}

func loadHistory() {
	fpath := getHistoryFile()
	f, err := os.Open(fpath)
	if err != nil {
		// File doesn't exist yet, that's okay
		return
	}
	defer f.Close()
	
	historyMu.Lock()
	defer historyMu.Unlock()
	
	// Try to load new format first
	var loadedHistory []ClipboardItem
	if err := json.NewDecoder(f).Decode(&loadedHistory); err != nil {
		// If new format fails, try legacy format ([]string)
		f.Seek(0, 0) // Reset file position
		var legacyHistory []string
		if err := json.NewDecoder(f).Decode(&legacyHistory); err != nil {
			// Reset to empty on error, don't print error for new installations
			history = []ClipboardItem{}
			return
		}
		
		// Convert legacy format to new format
		history = []ClipboardItem{}
		for _, item := range legacyHistory {
			if strings.TrimSpace(item) != "" && utf8.ValidString(item) {
				history = append(history, ClipboardItem{
					Type:      ItemTypeText,
					Content:   item,
					Timestamp: time.Now(),
				})
			}
		}
		return
	}
	
	// Validate loaded history (new format)
	history = []ClipboardItem{}
	for _, item := range loadedHistory {
		// Validate text items
		if item.Type == ItemTypeText {
			if strings.TrimSpace(item.Content) != "" && utf8.ValidString(item.Content) {
				history = append(history, item)
			}
		} else if item.Type == ItemTypeImage {
			// Validate image items - check if base64 content is valid
			if _, err := base64.StdEncoding.DecodeString(item.Content); err == nil {
				history = append(history, item)
			}
		}
	}
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
