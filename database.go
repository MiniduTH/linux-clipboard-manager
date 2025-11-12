package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// initDatabase initializes the SQLite database and creates tables if they don't exist
func initDatabase() error {
	dbPath := getDatabasePath()
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %v", err)
	}
	
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	
	// Test connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}
	
	// Create tables
	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}
	
	return nil
}

// getDatabasePath returns the path to the SQLite database file
func getDatabasePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to /tmp if home directory is not available
		return "/tmp/clipboard-history.db"
	}
	
	dir := filepath.Join(home, ".local", "share", "clipboard-manager")
	return filepath.Join(dir, "history.db")
}

// createTables creates the necessary database tables
func createTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS clipboard_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL CHECK(type IN ('text', 'image')),
		content TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		image_format TEXT,
		image_width INTEGER,
		image_height INTEGER,
		image_size INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_timestamp ON clipboard_history(timestamp DESC);
	CREATE INDEX IF NOT EXISTS idx_type ON clipboard_history(type);
	`
	
	_, err := db.Exec(createTableSQL)
	return err
}

// closeDatabase closes the database connection
func closeDatabase() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// saveClipboardItem saves a clipboard item to the database
func saveClipboardItem(item ClipboardItem) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	
	// Check for duplicates (same content and type)
	var count int
	checkSQL := "SELECT COUNT(*) FROM clipboard_history WHERE content = ? AND type = ?"
	err := db.QueryRow(checkSQL, item.Content, string(item.Type)).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check for duplicates: %v", err)
	}
	
	// If duplicate exists, delete it first (we'll add the new one at the end)
	if count > 0 {
		deleteSQL := "DELETE FROM clipboard_history WHERE content = ? AND type = ?"
		_, err = db.Exec(deleteSQL, item.Content, string(item.Type))
		if err != nil {
			return fmt.Errorf("failed to delete duplicate: %v", err)
		}
	}
	
	// Insert new item
	insertSQL := `
	INSERT INTO clipboard_history (type, content, timestamp, image_format, image_width, image_height, image_size)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	var imageFormat sql.NullString
	var imageWidth, imageHeight, imageSize sql.NullInt64
	
	if item.ImageMeta != nil {
		imageFormat = sql.NullString{String: item.ImageMeta.Format, Valid: true}
		imageWidth = sql.NullInt64{Int64: int64(item.ImageMeta.Width), Valid: true}
		imageHeight = sql.NullInt64{Int64: int64(item.ImageMeta.Height), Valid: true}
		imageSize = sql.NullInt64{Int64: int64(item.ImageMeta.Size), Valid: true}
	}
	
	_, err = db.Exec(insertSQL, string(item.Type), item.Content, item.Timestamp,
		imageFormat, imageWidth, imageHeight, imageSize)
	if err != nil {
		return fmt.Errorf("failed to insert clipboard item: %v", err)
	}
	
	// Maintain max history size
	return maintainHistoryLimit()
}

// loadClipboardHistory loads clipboard history from the database
func loadClipboardHistory() ([]ClipboardItem, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	
	query := `
	SELECT type, content, timestamp, image_format, image_width, image_height, image_size
	FROM clipboard_history
	ORDER BY timestamp ASC
	`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query clipboard history: %v", err)
	}
	defer rows.Close()
	
	var items []ClipboardItem
	
	for rows.Next() {
		var item ClipboardItem
		var itemType string
		var imageFormat sql.NullString
		var imageWidth, imageHeight, imageSize sql.NullInt64
		
		err := rows.Scan(&itemType, &item.Content, &item.Timestamp,
			&imageFormat, &imageWidth, &imageHeight, &imageSize)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		
		item.Type = ClipboardItemType(itemType)
		
		// Set image metadata if available
		if imageFormat.Valid {
			item.ImageMeta = &ImageMetadata{
				Format: imageFormat.String,
				Width:  int(imageWidth.Int64),
				Height: int(imageHeight.Int64),
				Size:   int(imageSize.Int64),
			}
		}
		
		items = append(items, item)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}
	
	return items, nil
}

// updateClipboardItem updates the content of an existing clipboard item
func updateClipboardItem(oldContent string, newContent string, itemType ClipboardItemType) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	
	// Update the item's content and timestamp
	updateSQL := "UPDATE clipboard_history SET content = ?, timestamp = ? WHERE content = ? AND type = ?"
	result, err := db.Exec(updateSQL, newContent, time.Now(), oldContent, string(itemType))
	if err != nil {
		return fmt.Errorf("failed to update clipboard item: %v", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("no item found to update")
	}
	
	return nil
}

// deleteClipboardItem deletes a clipboard item by its content and type
func deleteClipboardItem(content string, itemType ClipboardItemType) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	
	deleteSQL := "DELETE FROM clipboard_history WHERE content = ? AND type = ?"
	result, err := db.Exec(deleteSQL, content, string(itemType))
	if err != nil {
		return fmt.Errorf("failed to delete clipboard item: %v", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("no item found to delete")
	}
	
	return nil
}

// clearClipboardHistory deletes all clipboard history
func clearClipboardHistory() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	
	_, err := db.Exec("DELETE FROM clipboard_history")
	if err != nil {
		return fmt.Errorf("failed to clear clipboard history: %v", err)
	}
	
	return nil
}

// getClipboardHistoryCount returns the number of items in clipboard history
func getClipboardHistoryCount() (int, error) {
	if db == nil {
		return 0, fmt.Errorf("database not initialized")
	}
	
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM clipboard_history").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get history count: %v", err)
	}
	
	return count, nil
}

// maintainHistoryLimit ensures the history doesn't exceed maxHistory items
func maintainHistoryLimit() error {
	count, err := getClipboardHistoryCount()
	if err != nil {
		return err
	}
	
	if count > maxHistory {
		// Delete oldest items
		deleteSQL := `
		DELETE FROM clipboard_history 
		WHERE id IN (
			SELECT id FROM clipboard_history 
			ORDER BY timestamp ASC 
			LIMIT ?
		)
		`
		_, err = db.Exec(deleteSQL, count-maxHistory)
		if err != nil {
			return fmt.Errorf("failed to maintain history limit: %v", err)
		}
	}
	
	return nil
}

// migrateFromJSON migrates existing JSON history to SQLite database
func migrateFromJSON() error {
	// Check if JSON file exists
	jsonPath := getHistoryFile()
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		// No JSON file to migrate
		return nil
	}
	
	fmt.Println("Migrating clipboard history from JSON to SQLite...")
	
	// Load existing JSON history using the old method
	f, err := os.Open(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to open JSON history file: %v", err)
	}
	defer f.Close()
	
	// Try to load new format first
	var jsonHistory []ClipboardItem
	if err := json.NewDecoder(f).Decode(&jsonHistory); err != nil {
		// If new format fails, try legacy format ([]string)
		f.Seek(0, 0) // Reset file position
		var legacyHistory []string
		if err := json.NewDecoder(f).Decode(&legacyHistory); err != nil {
			return fmt.Errorf("failed to decode JSON history: %v", err)
		}
		
		// Convert legacy format to new format
		jsonHistory = []ClipboardItem{}
		for _, item := range legacyHistory {
			if strings.TrimSpace(item) != "" {
				jsonHistory = append(jsonHistory, ClipboardItem{
					Type:      ItemTypeText,
					Content:   item,
					Timestamp: time.Now(),
				})
			}
		}
	}
	
	// Save each item to SQLite
	for _, item := range jsonHistory {
		if err := saveClipboardItem(item); err != nil {
			fmt.Printf("Warning: failed to migrate item: %v\n", err)
		}
	}
	
	// Backup the JSON file
	backupPath := jsonPath + ".backup"
	if err := os.Rename(jsonPath, backupPath); err != nil {
		fmt.Printf("Warning: failed to backup JSON file: %v\n", err)
	} else {
		fmt.Printf("JSON history backed up to: %s\n", backupPath)
	}
	
	fmt.Printf("Successfully migrated %d items to SQLite database\n", len(jsonHistory))
	return nil
}