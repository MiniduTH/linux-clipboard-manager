package main

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func TestRemoveHistoryItem(t *testing.T) {
	// Setup test database
	setupTestDB(t)
	defer teardownTestDB(t)
	
	// Setup: Clear history and add test items
	testItems := []ClipboardItem{
		{Type: ItemTypeText, Content: "item1", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item2", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item3", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item4", Timestamp: time.Now()},
	}
	addTestItems(t, testItems)
	
	tests := []struct {
		name          string
		index         int
		expectedLen   int
		expectedItems []string
	}{
		{
			name:          "Remove first item",
			index:         0,
			expectedLen:   3,
			expectedItems: []string{"item2", "item3", "item4"},
		},
		{
			name:          "Remove middle item",
			index:         1,
			expectedLen:   3,
			expectedItems: []string{"item1", "item3", "item4"},
		},
		{
			name:          "Remove last item",
			index:         3,
			expectedLen:   3,
			expectedItems: []string{"item1", "item2", "item3"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset history for each test
			clearTestHistory(t)
			addTestItems(t, testItems)
			
			removeHistoryItem(tt.index)
			
			actualLen := getTestHistoryLength()
			if actualLen != tt.expectedLen {
				t.Errorf("Expected history length %d, got %d", tt.expectedLen, actualLen)
			}
			
			for i, expected := range tt.expectedItems {
				if i >= actualLen {
					t.Errorf("Expected item at index %d, but history is too short", i)
					continue
				}
				item := getTestHistoryItem(i)
				if item.Content != expected {
					t.Errorf("Expected item at index %d to be %s, got %s", i, expected, item.Content)
				}
			}
		})
	}
}

func TestRemoveHistoryItemInvalidIndex(t *testing.T) {
	// Setup test database
	setupTestDB(t)
	defer teardownTestDB(t)
	
	// Setup: Clear history and add test items
	testItems := []ClipboardItem{
		{Type: ItemTypeText, Content: "item1", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item2", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item3", Timestamp: time.Now()},
	}
	addTestItems(t, testItems)
	originalLen := getTestHistoryLength()
	
	tests := []struct {
		name  string
		index int
	}{
		{
			name:  "Negative index",
			index: -1,
		},
		{
			name:  "Index too large",
			index: 5,
		},
		{
			name:  "Index equal to length",
			index: 3,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset history for each test
			clearTestHistory(t)
			addTestItems(t, testItems)
			
			removeHistoryItem(tt.index)
			
			// History should remain unchanged
			actualLen := getTestHistoryLength()
			if actualLen != originalLen {
				t.Errorf("Expected history length to remain %d, got %d", originalLen, actualLen)
			}
		})
	}
}

func TestClearHistory(t *testing.T) {
	// Setup test database
	setupTestDB(t)
	defer teardownTestDB(t)
	
	// Setup: Add test items
	testItems := []ClipboardItem{
		{Type: ItemTypeText, Content: "item1", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item2", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item3", Timestamp: time.Now()},
	}
	addTestItems(t, testItems)
	
	err := clearHistory()
	if err != nil {
		t.Fatalf("clearHistory() failed: %v", err)
	}
	
	actualLen := getTestHistoryLength()
	if actualLen != 0 {
		t.Errorf("Expected history to be empty after clear, got length %d", actualLen)
	}
}

func TestClearHistoryWithCallback(t *testing.T) {
	// Setup test database
	setupTestDB(t)
	defer teardownTestDB(t)
	
	// Setup: Add test items
	testItems := []ClipboardItem{
		{Type: ItemTypeText, Content: "item1", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item2", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item3", Timestamp: time.Now()},
	}
	addTestItems(t, testItems)
	
	callbackCalled := false
	callback := func() {
		callbackCalled = true
	}
	
	err := clearHistory(callback)
	if err != nil {
		t.Fatalf("clearHistory() failed: %v", err)
	}
	
	actualLen := getTestHistoryLength()
	if actualLen != 0 {
		t.Errorf("Expected history to be empty after clear, got length %d", actualLen)
	}
	
	if !callbackCalled {
		t.Error("Expected callback to be called")
	}
}

func TestClearHistoryWithMultipleCallbacks(t *testing.T) {
	// Setup test database
	setupTestDB(t)
	defer teardownTestDB(t)
	
	// Setup: Add test items
	testItems := []ClipboardItem{
		{Type: ItemTypeText, Content: "item1", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item2", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item3", Timestamp: time.Now()},
	}
	addTestItems(t, testItems)
	
	callback1Called := false
	callback2Called := false
	
	callback1 := func() {
		callback1Called = true
	}
	
	callback2 := func() {
		callback2Called = true
	}
	
	err := clearHistory(callback1, callback2)
	if err != nil {
		t.Fatalf("clearHistory() failed: %v", err)
	}
	
	actualLen := getTestHistoryLength()
	if actualLen != 0 {
		t.Errorf("Expected history to be empty after clear, got length %d", actualLen)
	}
	
	if !callback1Called {
		t.Error("Expected callback1 to be called")
	}
	
	if !callback2Called {
		t.Error("Expected callback2 to be called")
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Setup test database
	setupTestDB(t)
	defer teardownTestDB(t)
	
	// Setup: Clear history and add initial items
	var testItems []ClipboardItem
	for i := 0; i < 10; i++ {
		testItems = append(testItems, ClipboardItem{
			Type:      ItemTypeText,
			Content:   fmt.Sprintf("item%d", i),
			Timestamp: time.Now(),
		})
	}
	addTestItems(t, testItems)
	
	var wg sync.WaitGroup
	numGoroutines := 5
	
	// Test concurrent deletion
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			// Try to remove items, some may fail due to invalid indices
			removeHistoryItem(index)
		}(i)
	}
	
	wg.Wait()
	
	// History should still be valid (no panics or corruption)
	currentLen := getTestHistoryLength()
	if currentLen < 0 || currentLen > 10 {
		t.Errorf("History length after concurrent deletion is invalid: %d", currentLen)
	}
}

func TestConcurrentClearAndRemove(t *testing.T) {
	// Setup test database
	setupTestDB(t)
	defer teardownTestDB(t)
	
	// Setup: Add test items
	var testItems []ClipboardItem
	for i := 0; i < 20; i++ {
		testItems = append(testItems, ClipboardItem{
			Type:      ItemTypeText,
			Content:   fmt.Sprintf("item%d", i),
			Timestamp: time.Now(),
		})
	}
	addTestItems(t, testItems)
	
	var wg sync.WaitGroup
	
	// Start concurrent operations
	wg.Add(2)
	
	// Goroutine 1: Clear history
	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond) // Small delay
		if err := clearHistory(); err != nil {
			t.Errorf("clearHistory() failed: %v", err)
		}
	}()
	
	// Goroutine 2: Try to remove items
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			removeHistoryItem(i)
			time.Sleep(5 * time.Millisecond)
		}
	}()
	
	wg.Wait()
	
	// Final state should be consistent (empty due to clear)
	finalLen := getTestHistoryLength()
	if finalLen != 0 {
		t.Errorf("Expected final history length to be 0, got %d", finalLen)
	}
}

// Setup and teardown for tests
func TestMain(m *testing.M) {
	// Save original history file path to restore later
	originalHistory := make([]ClipboardItem, len(history))
	copy(originalHistory, history)
	
	// Run tests
	code := m.Run()
	
	// Restore original history
	history = originalHistory
	
	os.Exit(code)
}