package main

import (
	"database/sql"
	"path/filepath"
	"testing"
)

// setupTestDB initializes a temporary SQLite database for testing
func setupTestDB(t *testing.T) {
	// Create a temporary database file
	tempDir := t.TempDir()
	testDBPath := filepath.Join(tempDir, "test_history.db")
	
	var err error
	db, err = sql.Open("sqlite3", testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	
	// Test connection
	if err = db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}
	
	// Create tables
	if err = createTables(); err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}
}

// teardownTestDB closes the test database
func teardownTestDB(t *testing.T) {
	if db != nil {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close test database: %v", err)
		}
		db = nil
	}
}

// clearTestHistory clears all test history data
func clearTestHistory(t *testing.T) {
	if db != nil {
		_, err := db.Exec("DELETE FROM clipboard_history")
		if err != nil {
			t.Fatalf("Failed to clear test history: %v", err)
		}
	}
	
	// Also clear in-memory history
	historyMu.Lock()
	history = []ClipboardItem{}
	historyMu.Unlock()
}

// addTestItems adds test items directly to the database and memory
func addTestItems(t *testing.T, items []ClipboardItem) {
	for _, item := range items {
		if err := saveClipboardItem(item); err != nil {
			t.Fatalf("Failed to add test item: %v", err)
		}
	}
	
	// Refresh in-memory history
	historyMu.Lock()
	refreshHistoryFromDB()
	historyMu.Unlock()
}

// getTestHistoryLength returns the current test history length
func getTestHistoryLength() int {
	historyMu.RLock()
	defer historyMu.RUnlock()
	return len(history)
}

// getTestHistoryItem returns a test history item by index
func getTestHistoryItem(index int) ClipboardItem {
	historyMu.RLock()
	defer historyMu.RUnlock()
	if index < 0 || index >= len(history) {
		return ClipboardItem{}
	}
	return history[index]
}