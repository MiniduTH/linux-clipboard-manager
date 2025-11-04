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
    echo "⚠️  Warning: No clipboard utility found."
    echo "⚙️  Installing clipboard utilities automatically..."
    echo
    
    INSTALL_SUCCESS=false
    
    # Detect distribution and install appropriate package
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        
        case "$ID" in
            ubuntu|debian|linuxmint)
                echo "Detected: $ID"
                if [ -n "$WAYLAND_DISPLAY" ]; then
                    echo "Installing wl-clipboard for Wayland..."
                    sudo apt-get update && sudo apt-get install -y wl-clipboard && INSTALL_SUCCESS=true
                else
                    echo "Installing xclip for X11..."
                    sudo apt-get update && sudo apt-get install -y xclip && INSTALL_SUCCESS=true
                fi
                ;;
            fedora|rhel|centos)
                echo "Detected: $ID"
                if [ -n "$WAYLAND_DISPLAY" ]; then
                    echo "Installing wl-clipboard for Wayland..."
                    sudo dnf install -y wl-clipboard && INSTALL_SUCCESS=true
                else
                    echo "Installing xclip for X11..."
                    sudo dnf install -y xclip && INSTALL_SUCCESS=true
                fi
                ;;
            arch|manjaro)
                echo "Detected: $ID"
                if [ -n "$WAYLAND_DISPLAY" ]; then
                    echo "Installing wl-clipboard for Wayland..."
                    sudo pacman -S --noconfirm wl-clipboard && INSTALL_SUCCESS=true
                else
                    echo "Installing xclip for X11..."
                    sudo pacman -S --noconfirm xclip && INSTALL_SUCCESS=true
                fi
                ;;
            opensuse*|sles)
                echo "Detected: $ID"
                if [ -n "$WAYLAND_DISPLAY" ]; then
                    echo "Installing wl-clipboard for Wayland..."
                    sudo zypper install -y wl-clipboard && INSTALL_SUCCESS=true
                else
                    echo "Installing xclip for X11..."
                    sudo zypper install -y xclip && INSTALL_SUCCESS=true
                fi
                ;;
            *)
                echo "⚠️  Unknown distribution: $ID"
                ;;
        esac
    fi
    
    if [ "$INSTALL_SUCCESS" = true ]; then
        echo "✓ Clipboard utilities installed successfully!"
        
        # Re-check for clipboard utilities
        if command -v xclip &> /dev/null; then
            CLIPBOARD_TOOL="xclip"
        elif command -v xsel &> /dev/null; then
            CLIPBOARD_TOOL="xsel"
        elif command -v wl-copy &> /dev/null; then
            CLIPBOARD_TOOL="wl-clipboard"
        fi
        
        if [ -n "$CLIPBOARD_TOOL" ]; then
            echo "✓ Found clipboard utility: $CLIPBOARD_TOOL"
        else
            echo "❌ Error: Clipboard utilities still not available after installation"
            exit 1
        fi
    else
        echo
        echo "❌ Error: Could not automatically install clipboard utilities"
        echo "   Please install manually:"
        echo "   Ubuntu/Debian: sudo apt install xclip"
        echo "   Fedora: sudo dnf install xclip"
        echo "   Arch: sudo pacman -S xclip"
        echo "   For Wayland: install wl-clipboard instead of xclip"
        exit 1
    fi
else
    echo "✓ Found clipboard utility: $CLIPBOARD_TOOL"
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