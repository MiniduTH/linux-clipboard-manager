package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Setup system-level hotkey using gsettings (GNOME) or other methods
func setupLinuxHotkeys() {
	fmt.Println("Setting up system hotkey integration...")
	
	// Get the absolute path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Warning: Could not get executable path: %v\n", err)
		execPath = "./clipboard-manager"
	}
	
	// Try to set up GNOME hotkey
	if setupGnomeHotkey(execPath) {
		fmt.Println("✓ GNOME hotkey configured: Super+Z")
	} else if setupKDEHotkey(execPath) {
		fmt.Println("✓ KDE hotkey configured: Super+Z")
	} else {
		fmt.Println("⚠ Could not configure system hotkey automatically")
		fmt.Println("Manual setup instructions:")
		fmt.Println("1. Open your system settings")
		fmt.Println("2. Go to Keyboard Shortcuts")
		fmt.Println("3. Add a custom shortcut:")
		fmt.Printf("   Command: %s show\n", execPath)
		fmt.Println("   Shortcut: Super+Z")
	}
	
	// Keep the application running
	fmt.Println("Clipboard manager is running. Press Ctrl+C to stop.")
	select {} // Block forever
}

// Setup GNOME hotkey using gsettings
func setupGnomeHotkey(execPath string) bool {
	// Check if gsettings is available
	if _, err := exec.LookPath("gsettings"); err != nil {
		return false
	}
	
	// Set custom keybinding
	schemaPath := "org.gnome.settings-daemon.plugins.media-keys"
	customPath := "/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/clipboard-manager/"
	
	commands := [][]string{
		{"gsettings", "set", schemaPath, "custom-keybindings", "['" + customPath + "']"},
		{"gsettings", "set", schemaPath + ".custom-keybinding:" + customPath, "name", "Clipboard Manager"},
		{"gsettings", "set", schemaPath + ".custom-keybinding:" + customPath, "command", execPath + " show"},
		{"gsettings", "set", schemaPath + ".custom-keybinding:" + customPath, "binding", "<Super>z"},
	}
	
	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			fmt.Printf("Warning: Failed to set GNOME hotkey: %v\n", err)
			return false
		}
	}
	
	return true
}

// Setup KDE hotkey using kwriteconfig5
func setupKDEHotkey(execPath string) bool {
	// Check if kwriteconfig5 is available
	if _, err := exec.LookPath("kwriteconfig5"); err != nil {
		return false
	}
	
	commands := [][]string{
		{"kwriteconfig5", "--file", "kglobalshortcutsrc", "--group", "clipboard-manager", "--key", "_k_friendly_name", "Clipboard Manager"},
		{"kwriteconfig5", "--file", "kglobalshortcutsrc", "--group", "clipboard-manager", "--key", "show", "Meta+Z,none,Show Clipboard History"},
	}
	
	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			fmt.Printf("Warning: Failed to set KDE hotkey: %v\n", err)
			return false
		}
	}
	
	return true
}

// Create a desktop entry for better system integration
func createDesktopEntry(execPath string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	
	desktopDir := filepath.Join(homeDir, ".local", "share", "applications")
	os.MkdirAll(desktopDir, 0755)
	
	desktopFile := filepath.Join(desktopDir, "clipboard-manager.desktop")
	content := fmt.Sprintf(`[Desktop Entry]
Name=Clipboard Manager
Comment=Clipboard history manager for Linux
Exec=%s show
Icon=edit-copy
Terminal=false
Type=Application
Categories=Utility;
Keywords=clipboard;history;copy;paste;
`, execPath)
	
	os.WriteFile(desktopFile, []byte(content), 0644)
}