#!/bin/bash

# Clipboard Manager Uninstall Script for Linux
# This script completely removes the clipboard manager from your system

set -e

echo "🗑️  Clipboard Manager Uninstall Script"
echo "======================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Confirmation prompt
echo "This will completely remove the Clipboard Manager from your system including:"
echo "  • System-wide installation (/usr/local/bin/clipboard-manager)"
echo "  • Desktop entries and autostart files"
echo "  • All clipboard history data"
echo "  • Custom keyboard shortcuts"
echo "  • Running processes"
echo ""
read -p "Are you sure you want to continue? (y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Uninstall cancelled."
    exit 0
fi

echo ""
print_status "Starting uninstall process..."
echo ""

# Step 1: Kill any running processes
print_status "Stopping running clipboard manager processes..."
if pgrep -f "clipboard-manager" > /dev/null; then
    pkill -f "clipboard-manager" && print_success "Stopped running processes" || print_warning "Could not stop some processes"
else
    print_success "No running processes found"
fi

# Step 2: Remove system-wide installation
print_status "Removing system-wide installation..."
if [ -f "/usr/local/bin/clipboard-manager" ]; then
    if sudo rm -f /usr/local/bin/clipboard-manager; then
        print_success "Removed /usr/local/bin/clipboard-manager"
    else
        print_error "Failed to remove system-wide installation (permission denied?)"
    fi
else
    print_success "System-wide installation not found"
fi

# Step 3: Remove autostart entries
print_status "Removing autostart entries..."
AUTOSTART_FILE="$HOME/.config/autostart/clipboard-manager.desktop"
if [ -f "$AUTOSTART_FILE" ]; then
    rm -f "$AUTOSTART_FILE" && print_success "Removed autostart entry" || print_warning "Could not remove autostart entry"
else
    print_success "No autostart entry found"
fi

# Step 4: Remove desktop application entries
print_status "Removing desktop application entries..."
DESKTOP_FILE="$HOME/.local/share/applications/clipboard-manager.desktop"
if [ -f "$DESKTOP_FILE" ]; then
    rm -f "$DESKTOP_FILE" && print_success "Removed desktop entry" || print_warning "Could not remove desktop entry"
else
    print_success "No desktop entry found"
fi

# Step 5: Remove clipboard history data
print_status "Removing clipboard history data..."
DATA_DIR="$HOME/.local/share/clipboard-manager"
if [ -d "$DATA_DIR" ]; then
    rm -rf "$DATA_DIR" && print_success "Removed clipboard history data" || print_warning "Could not remove clipboard data"
else
    print_success "No clipboard data found"
fi

# Step 6: Remove custom keyboard shortcuts (GNOME)
print_status "Removing custom keyboard shortcuts..."
if command_exists gsettings; then
    # Check if the custom keybinding exists
    CUSTOM_KEYBINDINGS=$(gsettings get org.gnome.settings-daemon.plugins.media-keys custom-keybindings 2>/dev/null || echo "[]")
    
    if echo "$CUSTOM_KEYBINDINGS" | grep -q "clipboard-manager"; then
        # Reset the specific keybinding
        gsettings reset-recursively org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/clipboard-manager/ 2>/dev/null || true
        
        # Remove from custom keybindings list
        gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings "[]" 2>/dev/null || true
        
        print_success "Removed custom keyboard shortcuts"
    else
        print_success "No custom keyboard shortcuts found"
    fi
else
    print_warning "gsettings not found - skipping keyboard shortcut removal"
fi

# Step 7: Clean build artifacts (if in project directory)
print_status "Cleaning build artifacts..."
if [ -f "Makefile" ] && [ -f "go.mod" ]; then
    if make clean >/dev/null 2>&1; then
        print_success "Cleaned build artifacts"
    else
        print_warning "Could not clean build artifacts (not in project directory?)"
    fi
else
    print_success "Not in project directory - skipping build cleanup"
fi

# Step 8: Update desktop database
print_status "Updating desktop database..."
if command_exists update-desktop-database; then
    update-desktop-database "$HOME/.local/share/applications" 2>/dev/null && print_success "Updated desktop database" || print_warning "Could not update desktop database"
else
    print_success "update-desktop-database not found - skipping"
fi

echo ""
print_success "✅ Clipboard Manager has been completely uninstalled!"
echo ""
echo "The following items have been removed:"
echo "  ✓ System-wide binary"
echo "  ✓ Desktop entries and autostart files"
echo "  ✓ Clipboard history data"
echo "  ✓ Custom keyboard shortcuts"
echo "  ✓ Running processes"
echo "  ✓ Build artifacts (if applicable)"
echo ""

# Optional: Remove project directory
if [ -f "go.mod" ] && [ -f "main.go" ]; then
    echo "You are currently in the project directory."
    read -p "Do you want to remove the entire project directory? (y/N): " -n 1 -r
    echo ""
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        PROJECT_DIR=$(pwd)
        cd ..
        rm -rf "$PROJECT_DIR" && print_success "Removed project directory: $PROJECT_DIR" || print_error "Could not remove project directory"
    else
        echo "Project directory kept."
    fi
fi

echo ""
print_success "Uninstall completed successfully! 🎉"
echo ""
echo "Thank you for using Clipboard Manager!"