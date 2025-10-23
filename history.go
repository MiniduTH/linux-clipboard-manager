package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"
)

var (
	history   []string
	historyMu sync.RWMutex
)

const maxHistory = 50

// add new item
func addToHistory(text string) {
	historyMu.Lock()
	defer historyMu.Unlock()
	
	// Clean and validate the text
	text = strings.TrimSpace(text)
	if text == "" || !utf8.ValidString(text) {
		return
	}
	
	// Skip if it's the same as the last item
	if len(history) > 0 && history[len(history)-1] == text {
		return
	}
	
	// Remove duplicates from history
	for i := len(history) - 1; i >= 0; i-- {
		if history[i] == text {
			history = append(history[:i], history[i+1:]...)
			break
		}
	}
	
	// Add to end
	history = append(history, text)
	
	// Maintain max size
	if len(history) > maxHistory {
		history = history[len(history)-maxHistory:]
	}
}

// Clear all history
func clearHistory() {
	historyMu.Lock()
	defer historyMu.Unlock()
	history = []string{}
	saveHistory()
	fmt.Println("Clipboard history cleared.")
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
	historyCopy := make([]string, len(history))
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
	
	var loadedHistory []string
	if err := json.NewDecoder(f).Decode(&loadedHistory); err != nil {
		// Reset to empty on error, don't print error for new installations
		history = []string{}
		return
	}
	
	// Validate loaded history
	history = []string{}
	for _, item := range loadedHistory {
		if strings.TrimSpace(item) != "" && utf8.ValidString(item) {
			history = append(history, item)
		}
	}
}
