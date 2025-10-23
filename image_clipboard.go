package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// detectImageInClipboard checks if there's an image in the clipboard
// Returns image data and format if found, nil otherwise
func detectImageInClipboard() ([]byte, string, error) {
	// Try different clipboard utilities based on what's available
	
	// First try xclip (X11)
	if imageData, format, err := tryXclipImage(); err == nil {
		return imageData, format, nil
	}
	
	// Then try wl-clipboard (Wayland)
	if imageData, format, err := tryWlClipboardImage(); err == nil {
		return imageData, format, nil
	}
	
	return nil, "", fmt.Errorf("no image found in clipboard")
}

// tryXclipImage attempts to get image data using xclip
func tryXclipImage() ([]byte, string, error) {
	// Check if xclip is available
	if _, err := exec.LookPath("xclip"); err != nil {
		return nil, "", fmt.Errorf("xclip not available")
	}
	
	// List available clipboard targets to see if there are images
	cmd := exec.Command("xclip", "-selection", "clipboard", "-t", "TARGETS", "-o")
	output, err := cmd.Output()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get clipboard targets: %v", err)
	}
	
	targets := strings.Split(string(output), "\n")
	
	// Look for image formats in the targets
	var imageFormat string
	for _, target := range targets {
		target = strings.TrimSpace(target)
		switch target {
		case "image/png":
			imageFormat = "png"
		case "image/jpeg":
			imageFormat = "jpeg"
		case "image/jpg":
			imageFormat = "jpg"
		case "image/gif":
			imageFormat = "gif"
		case "image/bmp":
			imageFormat = "bmp"
		}
		if imageFormat != "" {
			break
		}
	}
	
	if imageFormat == "" {
		return nil, "", fmt.Errorf("no supported image format found in clipboard")
	}
	
	// Get the image data
	mimeType := fmt.Sprintf("image/%s", imageFormat)
	cmd = exec.Command("xclip", "-selection", "clipboard", "-t", mimeType, "-o")
	imageData, err := cmd.Output()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get image data: %v", err)
	}
	
	if len(imageData) == 0 {
		return nil, "", fmt.Errorf("empty image data")
	}
	
	return imageData, imageFormat, nil
}

// tryWlClipboardImage attempts to get image data using wl-clipboard
func tryWlClipboardImage() ([]byte, string, error) {
	// Check if wl-paste is available
	if _, err := exec.LookPath("wl-paste"); err != nil {
		return nil, "", fmt.Errorf("wl-paste not available")
	}
	
	// List available clipboard types
	cmd := exec.Command("wl-paste", "--list-types")
	output, err := cmd.Output()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get clipboard types: %v", err)
	}
	
	types := strings.Split(string(output), "\n")
	
	// Look for image formats in the types
	var imageFormat string
	var mimeType string
	for _, clipType := range types {
		clipType = strings.TrimSpace(clipType)
		switch clipType {
		case "image/png":
			imageFormat = "png"
			mimeType = clipType
		case "image/jpeg":
			imageFormat = "jpeg"
			mimeType = clipType
		case "image/jpg":
			imageFormat = "jpg"
			mimeType = clipType
		case "image/gif":
			imageFormat = "gif"
			mimeType = clipType
		case "image/bmp":
			imageFormat = "bmp"
			mimeType = clipType
		}
		if imageFormat != "" {
			break
		}
	}
	
	if imageFormat == "" {
		return nil, "", fmt.Errorf("no supported image format found in clipboard")
	}
	
	// Get the image data
	cmd = exec.Command("wl-paste", "--type", mimeType)
	imageData, err := cmd.Output()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get image data: %v", err)
	}
	
	if len(imageData) == 0 {
		return nil, "", fmt.Errorf("empty image data")
	}
	
	return imageData, imageFormat, nil
}

// restoreImageToSystemClipboard restores image data to the system clipboard
func restoreImageToSystemClipboard(imageData []byte, format string) error {
	// Try xclip first (X11)
	if err := tryRestoreWithXclip(imageData, format); err == nil {
		return nil
	}
	
	// Then try wl-clipboard (Wayland)
	if err := tryRestoreWithWlClipboard(imageData, format); err == nil {
		return nil
	}
	
	return fmt.Errorf("failed to restore image to clipboard")
}

// tryRestoreWithXclip attempts to restore image using xclip
func tryRestoreWithXclip(imageData []byte, format string) error {
	if _, err := exec.LookPath("xclip"); err != nil {
		return fmt.Errorf("xclip not available")
	}
	
	mimeType := fmt.Sprintf("image/%s", format)
	cmd := exec.Command("xclip", "-selection", "clipboard", "-t", mimeType)
	cmd.Stdin = bytes.NewReader(imageData)
	
	return cmd.Run()
}

// tryRestoreWithWlClipboard attempts to restore image using wl-clipboard
func tryRestoreWithWlClipboard(imageData []byte, format string) error {
	if _, err := exec.LookPath("wl-copy"); err != nil {
		return fmt.Errorf("wl-copy not available")
	}
	
	mimeType := fmt.Sprintf("image/%s", format)
	cmd := exec.Command("wl-copy", "--type", mimeType)
	cmd.Stdin = bytes.NewReader(imageData)
	
	return cmd.Run()
}