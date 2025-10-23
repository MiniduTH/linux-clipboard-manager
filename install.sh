#!/bin/bash

# Clipboard Manager Installation Script for Linux

echo "Installing Clipboard Manager for Linux..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first."
    echo "Visit: https://golang.org/doc/install"
    exit 1
fi

# Check for clipboard utilities
echo "Checking for clipboard utilities..."
CLIPBOARD_TOOL=""

if command -v xclip &> /dev/null; then
    CLIPBOARD_TOOL="xclip"
elif command -v xsel &> /dev/null; then
    CLIPBOARD_TOOL="xsel"
elif command -v wl-copy &> /dev/null; then
    CLIPBOARD_TOOL="wl-clipboard"
fi

if [ -z "$CLIPBOARD_TOOL" ]; then
    echo "Warning: No clipboard utility found."
    echo "Installing clipboard utilities..."
    
    # Detect distribution and install appropriate package
    if command -v apt &> /dev/null; then
        # Ubuntu/Debian
        echo "Detected Ubuntu/Debian system"
        sudo apt update
        sudo apt install -y xclip
    elif command -v dnf &> /dev/null; then
        # Fedora
        echo "Detected Fedora system"
        sudo dnf install -y xclip
    elif command -v yum &> /dev/null; then
        # CentOS/RHEL
        echo "Detected CentOS/RHEL system"
        sudo yum install -y xclip
    else
        echo "Please install xclip, xsel, or wl-clipboard manually:"
        echo "  Ubuntu/Debian: sudo apt install xclip"
        echo "  Fedora: sudo dnf install xclip"
        echo "  Arch: sudo pacman -S xclip"
        exit 1
    fi
else
    echo "Found clipboard utility: $CLIPBOARD_TOOL"
fi

# Build the application
echo "Building clipboard manager..."
go mod tidy
go build -o clipboard-manager

if [ $? -eq 0 ]; then
    echo "Build successful!"
    
    # Make it executable
    chmod +x clipboard-manager
    
    # Optionally install to system
    read -p "Install to /usr/local/bin? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        sudo cp clipboard-manager /usr/local/bin/
        echo "Installed to /usr/local/bin/clipboard-manager"
    fi
    
    echo ""
    echo "Installation complete!"
    echo ""
    echo "Usage:"
    echo "  ./clipboard-manager        - Start clipboard watcher"
    echo "  ./clipboard-manager show   - Show GUI history"
    echo "  ./clipboard-manager list   - Show terminal history"
    echo ""
    echo "To start the clipboard watcher in background:"
    echo "  nohup ./clipboard-manager > /dev/null 2>&1 &"
    
else
    echo "Build failed!"
    exit 1
fi