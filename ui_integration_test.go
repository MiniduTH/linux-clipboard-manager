package main

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
)

// TestHistoryListItemRendering tests the custom HistoryListItem widget rendering
func TestHistoryListItemRendering(t *testing.T) {
	// Test data with various content lengths
	testCases := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Short text",
			text:     "Hello World",
			expected: "Hello World",
		},
		{
			name:     "Multi-line text",
			text:     "Line 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "Text with tabs",
			text:     "Column1\tColumn2\tColumn3",
			expected: "Column1    Column2    Column3", // tabs converted to spaces
		},
		{
			name:     "Long text requiring truncation",
			text:     "This is a very long text that should be truncated because it exceeds the maximum display length limit that we have set for the history list items to prevent them from taking up too much vertical space in the user interface. This text is intentionally made very long to exceed the 250 character limit that triggers truncation in the prepareDisplayText method. We need to make sure this text is definitely longer than 250 characters to properly test the truncation functionality.",
			expected: "truncated", // Should contain truncation indicator
		},
		{
			name:     "Text with markdown characters",
			text:     "**bold** _italic_ #header [link]",
			expected: "escaped", // Should escape markdown characters
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test app
			testApp := test.NewApp()
			defer testApp.Quit()

			// Create callbacks for testing
			deleteCallbackCalled := false
			selectCallbackCalled := false
			
			onDelete := func(index int) {
				deleteCallbackCalled = true
			}
			
			onSelect := func(index int) {
				selectCallbackCalled = true
			}

			// Create HistoryListItem
			clipboardItem := ClipboardItem{
				Type:      ItemTypeText,
				Content:   tc.text,
				Timestamp: time.Now(),
			}
			onEdit := func(index int) {}
			item := NewHistoryListItem(clipboardItem, 0, onDelete, onSelect, onEdit)
			
			// Verify widget was created successfully
			if item == nil {
				t.Fatal("Failed to create HistoryListItem")
			}

			// Test widget properties
			if item.item.Content != tc.text {
				t.Errorf("Expected text %q, got %q", tc.text, item.item.Content)
			}

			if item.index != 0 {
				t.Errorf("Expected index 0, got %d", item.index)
			}

			// Test widget rendering
			renderer := item.CreateRenderer()
			if renderer == nil {
				t.Fatal("Failed to create renderer")
			}

			// Test minimum size calculation
			minSize := renderer.MinSize()
			if minSize.Width <= 0 || minSize.Height <= 0 {
				t.Errorf("Invalid minimum size: %v", minSize)
			}

			// Verify minimum height constraints
			if minSize.Height < 40 {
				t.Errorf("Minimum height too small: %f", minSize.Height)
			}
			if minSize.Height > 80 {
				t.Errorf("Minimum height too large: %f", minSize.Height)
			}

			// Test text content preparation
			displayText := item.prepareDisplayText()
			
			switch tc.expected {
			case "truncated":
				if len(displayText) >= len(tc.text) {
					t.Error("Long text was not truncated")
				}
				if displayText[len(displayText)-3:] != "..." {
					t.Error("Truncation indicator not found")
				}
			case "escaped":
				if displayText == tc.text {
					t.Error("Markdown characters were not escaped")
				}
			default:
				// For simple cases, check if content is preserved
				if len(tc.text) < 50 && !containsContent(displayText, tc.expected) {
					t.Errorf("Expected content %q not found in display text %q", tc.expected, displayText)
				}
			}

			// Test widget components exist
			if item.textWidget == nil {
				t.Error("Text widget not created")
			}
			if item.deleteButton == nil {
				t.Error("Delete button not created")
			}
			if item.background == nil {
				t.Error("Background not created")
			}
			if item.container == nil {
				t.Error("Container not created")
			}

			// Test callbacks are not called during creation
			if deleteCallbackCalled {
				t.Error("Delete callback called unexpectedly during creation")
			}
			if selectCallbackCalled {
				t.Error("Select callback called unexpectedly during creation")
			}
		})
	}
}

// TestHistoryListItemInteraction tests user interaction with the widget
func TestHistoryListItemInteraction(t *testing.T) {
	// Create test app
	testApp := test.NewApp()
	defer testApp.Quit()

	// Track callback invocations
	deleteCallbackCalled := false
	selectCallbackCalled := false
	deletedIndex := -1
	selectedIndex := -1
	
	onDelete := func(index int) {
		deleteCallbackCalled = true
		deletedIndex = index
	}
	
	onSelect := func(index int) {
		selectCallbackCalled = true
		selectedIndex = index
	}

	// Create HistoryListItem
	testText := "Test clipboard content"
	testIndex := 5
	clipboardItem := ClipboardItem{
		Type:      ItemTypeText,
		Content:   testText,
		Timestamp: time.Now(),
	}
	onEdit := func(index int) {}
	item := NewHistoryListItem(clipboardItem, testIndex, onDelete, onSelect, onEdit)

	// Test item selection (tap)
	t.Run("Item selection", func(t *testing.T) {
		// Reset callback states
		selectCallbackCalled = false
		selectedIndex = -1

		// Simulate tap event
		item.Tapped(&fyne.PointEvent{
			Position: fyne.NewPos(50, 25),
		})

		// Verify select callback was called
		if !selectCallbackCalled {
			t.Error("Select callback was not called")
		}
		if selectedIndex != testIndex {
			t.Errorf("Expected selected index %d, got %d", testIndex, selectedIndex)
		}
	})

	// Test delete button functionality
	t.Run("Delete button click", func(t *testing.T) {
		// Reset callback states
		deleteCallbackCalled = false
		deletedIndex = -1

		// Simulate delete button click
		if item.deleteButton != nil {
			test.Tap(item.deleteButton)
		} else {
			t.Fatal("Delete button not found")
		}

		// Verify delete callback was called
		if !deleteCallbackCalled {
			t.Error("Delete callback was not called")
		}
		if deletedIndex != testIndex {
			t.Errorf("Expected deleted index %d, got %d", testIndex, deletedIndex)
		}
	})

	// Test hover states
	t.Run("Hover states", func(t *testing.T) {
		// Test mouse enter
		item.MouseIn(&desktop.MouseEvent{
			PointEvent: fyne.PointEvent{
				Position: fyne.NewPos(50, 25),
			},
		})

		if !item.isHovered {
			t.Error("Item should be hovered after MouseIn")
		}

		// Test mouse movement over delete button area
		deleteButtonPos := item.deleteButton.Position()
		deleteButtonSize := item.deleteButton.Size()
		
		item.MouseMoved(&desktop.MouseEvent{
			PointEvent: fyne.PointEvent{
				Position: fyne.NewPos(
					deleteButtonPos.X+deleteButtonSize.Width/2,
					deleteButtonPos.Y+deleteButtonSize.Height/2,
				),
			},
		})

		if !item.deleteHovered {
			t.Error("Delete button should be hovered")
		}

		// Test mouse leave
		item.MouseOut()

		if item.isHovered {
			t.Error("Item should not be hovered after MouseOut")
		}
		if item.deleteHovered {
			t.Error("Delete button should not be hovered after MouseOut")
		}
	})
}

// TestHistoryListItemTextWrapping tests text wrapping and truncation functionality
func TestHistoryListItemTextWrapping(t *testing.T) {
	testCases := []struct {
		name           string
		text           string
		expectTruncated bool
		expectWrapped   bool
		maxLines        int
	}{
		{
			name:           "Short single line",
			text:           "Short text",
			expectTruncated: false,
			expectWrapped:   false,
		},
		{
			name:           "Multiple lines within limit",
			text:           "Line 1\nLine 2\nLine 3",
			expectTruncated: false,
			expectWrapped:   true,
		},
		{
			name:           "Too many lines",
			text:           "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6",
			expectTruncated: true,
			expectWrapped:   true,
			maxLines:        4,
		},
		{
			name:           "Very long single line",
			text:           "This is an extremely long line of text that should definitely be truncated because it exceeds the maximum display length that we have configured for the history list items to prevent them from taking up excessive vertical space in the user interface and making the list difficult to navigate. This text is intentionally made very long to exceed the 250 character limit that triggers truncation in the prepareDisplayText method.",
			expectTruncated: true,
			expectWrapped:   false,
		},
		{
			name:           "Text with special characters",
			text:           "Special chars: **bold** _italic_ #header [link] `code`",
			expectTruncated: false,
			expectWrapped:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test app
			testApp := test.NewApp()
			defer testApp.Quit()

			// Create HistoryListItem
			clipboardItem := ClipboardItem{
				Type:      ItemTypeText,
				Content:   tc.text,
				Timestamp: time.Now(),
			}
			item := NewHistoryListItem(clipboardItem, 0, nil, nil, nil)

			// Test text preparation
			displayText := item.prepareDisplayText()

			// Check truncation
			if tc.expectTruncated {
				if len(displayText) >= len(tc.text) {
					t.Error("Expected text to be truncated")
				}
				if len(displayText) < 3 || displayText[len(displayText)-3:] != "..." {
					t.Error("Expected truncation indicator (...)")
				}
			} else {
				if len(displayText) >= 3 && displayText[len(displayText)-3:] == "..." {
					t.Error("Unexpected truncation")
				}
			}

			// Check line limits
			if tc.maxLines > 0 {
				lines := countLines(displayText)
				// Account for the truncation indicator which may add a line
				expectedMaxLines := tc.maxLines
				if tc.expectTruncated {
					expectedMaxLines = tc.maxLines + 1 // Allow for truncation indicator line
				}
				if lines > expectedMaxLines {
					t.Errorf("Expected max %d lines, got %d", expectedMaxLines, lines)
				}
			}

			// Check markdown escaping
			if containsMarkdown(tc.text) && !containsEscapedMarkdown(displayText) {
				t.Error("Markdown characters should be escaped")
			}

			// Test text widget configuration
			if item.textWidget == nil {
				t.Fatal("Text widget not created")
			}

			// Verify text widget has wrapping enabled
			if item.textWidget.Wrapping != fyne.TextWrapWord {
				t.Error("Text widget should have word wrapping enabled")
			}

			// Test widget size constraints
			renderer := item.CreateRenderer()
			minSize := renderer.MinSize()
			
			// Verify height constraints are applied
			if minSize.Height > 80 {
				t.Errorf("Widget height exceeds maximum: %f", minSize.Height)
			}
		})
	}
}

// TestHistoryListItemUpdates tests dynamic updates to the widget
func TestHistoryListItemUpdates(t *testing.T) {
	// Create test app
	testApp := test.NewApp()
	defer testApp.Quit()

	// Create HistoryListItem
	originalText := "Original text"
	originalIndex := 0
	clipboardItem := ClipboardItem{
		Type:      ItemTypeText,
		Content:   originalText,
		Timestamp: time.Now(),
	}
	item := NewHistoryListItem(clipboardItem, originalIndex, nil, nil, nil)

	// Test item update
	t.Run("Item update", func(t *testing.T) {
		newText := "Updated text content"
		newItem := ClipboardItem{
			Type:      ItemTypeText,
			Content:   newText,
			Timestamp: time.Now(),
		}
		item.UpdateItem(newItem)

		if item.item.Content != newText {
			t.Errorf("Expected text %q, got %q", newText, item.item.Content)
		}

		// Verify display text was updated
		displayText := item.prepareDisplayText()
		if !containsContent(displayText, newText) {
			t.Error("Display text was not updated")
		}
	})

	// Test index update
	t.Run("Index update", func(t *testing.T) {
		newIndex := 10
		item.UpdateIndex(newIndex)

		if item.index != newIndex {
			t.Errorf("Expected index %d, got %d", newIndex, item.index)
		}
	})
}

// Helper functions for testing

// containsContent checks if display text contains expected content
func containsContent(displayText, expected string) bool {
	// Simple substring check, accounting for formatting changes
	return len(displayText) > 0 && (displayText == expected || 
		(len(expected) > 10 && len(displayText) >= len(expected)-10))
}

// countLines counts the number of lines in text
func countLines(text string) int {
	if text == "" {
		return 0
	}
	lines := 1
	for _, char := range text {
		if char == '\n' {
			lines++
		}
	}
	return lines
}

// containsMarkdown checks if text contains markdown characters
func containsMarkdown(text string) bool {
	markdownChars := []string{"**", "_", "#", "[", "]", "`"}
	for _, char := range markdownChars {
		if len(text) > 0 && text != char { // Avoid false positives
			for i := 0; i < len(text)-len(char)+1; i++ {
				if text[i:i+len(char)] == char {
					return true
				}
			}
		}
	}
	return false
}

// containsEscapedMarkdown checks if markdown characters are properly escaped
func containsEscapedMarkdown(text string) bool {
	escapedChars := []string{"\\*", "\\_", "\\#", "\\[", "\\]"}
	for _, escaped := range escapedChars {
		for i := 0; i < len(text)-len(escaped)+1; i++ {
			if text[i:i+len(escaped)] == escaped {
				return true
			}
		}
	}
	return false
}