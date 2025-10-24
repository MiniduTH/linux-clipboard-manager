package main

import (
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"
)

// createTestImage creates a simple test image for testing
func createTestImage() ([]byte, error) {
	// Create a simple 100x100 red square
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	
	// Fill with red color
	red := color.RGBA{255, 0, 0, 255}
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, red)
		}
	}
	
	// Encode as PNG
	var buf strings.Builder
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	
	return []byte(buf.String()), nil
}

func TestAddImageToHistory(t *testing.T) {
	// Setup test database
	setupTestDB(t)
	defer teardownTestDB(t)
	
	// Clear history for test
	clearTestHistory(t)
	
	// Create test image
	imageData, err := createTestImage()
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	
	// Add image to history
	addImageToHistory(imageData, "png")
	
	// Verify image was added
	actualLen := getTestHistoryLength()
	if actualLen != 1 {
		t.Fatalf("Expected 1 item in history, got %d", actualLen)
	}
	
	item := getTestHistoryItem(0)
	if item.Type != ItemTypeImage {
		t.Errorf("Expected item type %s, got %s", ItemTypeImage, item.Type)
	}
	
	// Verify base64 content
	_, err = base64.StdEncoding.DecodeString(item.Content)
	if err != nil {
		t.Errorf("Failed to decode base64 content: %v", err)
	}
	
	// Verify metadata
	if item.ImageMeta == nil {
		t.Error("Expected image metadata, got nil")
	} else {
		if item.ImageMeta.Format != "png" {
			t.Errorf("Expected format 'png', got '%s'", item.ImageMeta.Format)
		}
		if item.ImageMeta.Width != 100 {
			t.Errorf("Expected width 100, got %d", item.ImageMeta.Width)
		}
		if item.ImageMeta.Height != 100 {
			t.Errorf("Expected height 100, got %d", item.ImageMeta.Height)
		}
	}
}

func TestImageDuplicateDetection(t *testing.T) {
	// Setup test database
	setupTestDB(t)
	defer teardownTestDB(t)
	
	// Clear history for test
	clearTestHistory(t)
	
	// Create test image
	imageData, err := createTestImage()
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	
	// Add same image twice
	addImageToHistory(imageData, "png")
	addImageToHistory(imageData, "png")
	
	// Should only have one item (duplicate detection)
	actualLen := getTestHistoryLength()
	if actualLen != 1 {
		t.Errorf("Expected 1 item in history (duplicate detection), got %d", actualLen)
	}
}

func TestMixedTextAndImageHistory(t *testing.T) {
	// Setup test database
	setupTestDB(t)
	defer teardownTestDB(t)
	
	// Clear history for test
	clearTestHistory(t)
	
	// Add text item
	addToHistory("Test text content")
	
	// Create and add image
	imageData, err := createTestImage()
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	addImageToHistory(imageData, "png")
	
	// Add another text item
	addToHistory("Another text item")
	
	// Verify mixed history
	actualLen := getTestHistoryLength()
	if actualLen != 3 {
		t.Fatalf("Expected 3 items in history, got %d", actualLen)
	}
	
	// Check order and types
	item0 := getTestHistoryItem(0)
	if item0.Type != ItemTypeText || item0.Content != "Test text content" {
		t.Error("First item should be text")
	}
	
	item1 := getTestHistoryItem(1)
	if item1.Type != ItemTypeImage {
		t.Error("Second item should be image")
	}
	
	item2 := getTestHistoryItem(2)
	if item2.Type != ItemTypeText || item2.Content != "Another text item" {
		t.Error("Third item should be text")
	}
}

func TestRestoreImageToClipboard(t *testing.T) {
	// Create test image
	imageData, err := createTestImage()
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	
	// Encode as base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)
	
	// Test restoration (this will likely fail in test environment without clipboard utilities)
	err = restoreImageToClipboard(base64Data, "png")
	
	// We expect this to fail in test environment, but it should not panic
	// The function should handle the error gracefully
	if err != nil {
		t.Logf("Image restoration failed as expected in test environment: %v", err)
	}
}

func TestImageClipboardDetection(t *testing.T) {
	// Test image detection (this will likely fail in test environment)
	imageData, format, err := detectImageInClipboard()
	
	// We expect this to fail in test environment without actual clipboard content
	if err != nil {
		t.Logf("Image detection failed as expected in test environment: %v", err)
		return
	}
	
	// If it somehow succeeds, verify the data
	if len(imageData) == 0 {
		t.Error("Expected non-empty image data")
	}
	
	if format == "" {
		t.Error("Expected non-empty format")
	}
	
	t.Logf("Detected image: %d bytes, format: %s", len(imageData), format)
}