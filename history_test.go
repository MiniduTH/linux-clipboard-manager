package main

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func TestRemoveHistoryItem(t *testing.T) {
	// Setup: Clear history and add test items
	history = []ClipboardItem{
		{Type: ItemTypeText, Content: "item1", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item2", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item3", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item4", Timestamp: time.Now()},
	}
	
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
			history = []ClipboardItem{
				{Type: ItemTypeText, Content: "item1", Timestamp: time.Now()},
				{Type: ItemTypeText, Content: "item2", Timestamp: time.Now()},
				{Type: ItemTypeText, Content: "item3", Timestamp: time.Now()},
				{Type: ItemTypeText, Content: "item4", Timestamp: time.Now()},
			}
			
			removeHistoryItem(tt.index)
			
			if len(history) != tt.expectedLen {
				t.Errorf("Expected history length %d, got %d", tt.expectedLen, len(history))
			}
			
			for i, expected := range tt.expectedItems {
				if i >= len(history) || history[i].Content != expected {
					t.Errorf("Expected item at index %d to be %s, got %s", i, expected, history[i].Content)
				}
			}
		})
	}
}

func TestRemoveHistoryItemInvalidIndex(t *testing.T) {
	// Setup: Clear history and add test items
	history = []ClipboardItem{
		{Type: ItemTypeText, Content: "item1", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item2", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item3", Timestamp: time.Now()},
	}
	originalLen := len(history)
	
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
			history = []ClipboardItem{
				{Type: ItemTypeText, Content: "item1", Timestamp: time.Now()},
				{Type: ItemTypeText, Content: "item2", Timestamp: time.Now()},
				{Type: ItemTypeText, Content: "item3", Timestamp: time.Now()},
			}
			
			removeHistoryItem(tt.index)
			
			// History should remain unchanged
			if len(history) != originalLen {
				t.Errorf("Expected history length to remain %d, got %d", originalLen, len(history))
			}
		})
	}
}

func TestClearHistory(t *testing.T) {
	// Setup: Add test items
	history = []ClipboardItem{
		{Type: ItemTypeText, Content: "item1", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item2", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item3", Timestamp: time.Now()},
	}
	
	clearHistory()
	
	if len(history) != 0 {
		t.Errorf("Expected history to be empty after clear, got length %d", len(history))
	}
}

func TestClearHistoryWithCallback(t *testing.T) {
	// Setup: Add test items
	history = []ClipboardItem{
		{Type: ItemTypeText, Content: "item1", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item2", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item3", Timestamp: time.Now()},
	}
	
	callbackCalled := false
	callback := func() {
		callbackCalled = true
	}
	
	clearHistory(callback)
	
	if len(history) != 0 {
		t.Errorf("Expected history to be empty after clear, got length %d", len(history))
	}
	
	if !callbackCalled {
		t.Error("Expected callback to be called")
	}
}

func TestClearHistoryWithMultipleCallbacks(t *testing.T) {
	// Setup: Add test items
	history = []ClipboardItem{
		{Type: ItemTypeText, Content: "item1", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item2", Timestamp: time.Now()},
		{Type: ItemTypeText, Content: "item3", Timestamp: time.Now()},
	}
	
	callback1Called := false
	callback2Called := false
	
	callback1 := func() {
		callback1Called = true
	}
	
	callback2 := func() {
		callback2Called = true
	}
	
	clearHistory(callback1, callback2)
	
	if len(history) != 0 {
		t.Errorf("Expected history to be empty after clear, got length %d", len(history))
	}
	
	if !callback1Called {
		t.Error("Expected callback1 to be called")
	}
	
	if !callback2Called {
		t.Error("Expected callback2 to be called")
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Setup: Clear history
	history = []ClipboardItem{}
	
	// Add initial items
	for i := 0; i < 10; i++ {
		history = append(history, ClipboardItem{
			Type:      ItemTypeText,
			Content:   fmt.Sprintf("item%d", i),
			Timestamp: time.Now(),
		})
	}
	
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
	currentLen := getHistoryLength()
	if currentLen < 0 || currentLen > 10 {
		t.Errorf("History length after concurrent deletion is invalid: %d", currentLen)
	}
}

func TestConcurrentClearAndRemove(t *testing.T) {
	// Setup: Add test items
	history = []ClipboardItem{}
	for i := 0; i < 20; i++ {
		history = append(history, ClipboardItem{
			Type:      ItemTypeText,
			Content:   fmt.Sprintf("item%d", i),
			Timestamp: time.Now(),
		})
	}
	
	var wg sync.WaitGroup
	
	// Start concurrent operations
	wg.Add(2)
	
	// Goroutine 1: Clear history
	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond) // Small delay
		clearHistory()
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
	finalLen := getHistoryLength()
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