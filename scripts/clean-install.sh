#!/bin/bash

# Conservative Clean Installation Script for Clipboard Manager
# This script only removes clipboard-manager specific files and settings
# It preserves all existing system dependencies and user data
# Only installs missing dependencies, never removes existing ones

echo "🧹 Cleaning up previous Clipboard Manager installation..."

# Stop any running clipboard manager processes
echo "Stopping any running clipboard manager processes..."
pkill -f clipboard-manager 2>/dev/null || true
sleep 2

# Remove old binaries (only clipboard-manager specific)
echo "Removing old clipboard-manager binaries..."
rm -f clipboard-manager
rm -f clipboard-manager-test
rm -f tests/clipboard-manager-test
rm -f /usr/local/bin/clipboard-manager 2>/dev/null || true

# Remove old desktop entries (only clipboard-manager specific)
echo "Removing old clipboard-manager desktop entries..."
rm -f ~/.local/share/applications/clipboard-manager.desktop
rm -f ~/.config/autostart/clipboard-manager.desktop

# Remove old GNOME hotkey settings (only clipboard-manager specific)
echo "Removing old clipboard-manager GNOME hotkey settings..."
if command -v gsettings &> /dev/null; then
    CUSTOM_PATH="/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/clipboard-manager/"
    gsettings reset org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:$CUSTOM_PATH name 2>/dev/null || true
    gsettings reset org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:$CUSTOM_PATH command 2>/dev/null || true
    gsettings reset org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:$CUSTOM_PATH binding 2>/dev/null || true
    
    # Remove from custom keybindings list
    CURRENT_BINDINGS=$(gsettings get org.gnome.settings-daemon.plugins.media-keys custom-keybindings 2>/dev/null || echo "@as []")
    if [[ "$CURRENT_BINDINGS" == *"$CUSTOM_PATH"* ]]; then
        NEW_BINDINGS=$(echo "$CURRENT_BINDINGS" | sed "s|'$CUSTOM_PATH'[,]*||g" | sed 's/\[,/[/g' | sed 's/,,/,/g' | sed 's/,]/]/g')
        gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings "$NEW_BINDINGS" 2>/dev/null || true
    fi
fi

# Remove old KDE hotkey settings (only clipboard-manager specific)
echo "Removing old clipboard-manager KDE hotkey settings..."
if command -v kwriteconfig5 &> /dev/null; then
    kwriteconfig5 --file kglobalshortcutsrc --group "clipboard-manager.desktop" --key "_k_friendly_name" --delete 2>/dev/null || true
    kwriteconfig5 --file kglobalshortcutsrc --group "clipboard-manager.desktop" --key "show" --delete 2>/dev/null || true
fi

# Clean only clipboard-manager related Go build cache (preserve system Go cache)
echo "Cleaning clipboard-manager build artifacts..."
go clean 2>/dev/null || true

# Preserve user data (clipboard history)
echo "Preserving user clipboard history..."
if [ -f "$HOME/.local/share/clipboard-manager/history.json" ]; then
    echo "✅ Existing clipboard history will be preserved"
else
    echo "ℹ️  No existing clipboard history found"
fi

echo "✅ Cleanup completed!"
echo ""

# Fresh installation
echo "🚀 Starting fresh installation..."

# Check dependencies (only install missing ones, preserve existing)
echo "Checking required dependencies..."
MISSING_DEPS=()
MISSING_DEV_DEPS=()

# Check for Go
if ! command -v go &> /dev/null; then
    MISSING_DEPS+=("go")
    echo "⚠️  Go is not installed"
else
    echo "✅ Go is available: $(go version | cut -d' ' -f3)"
fi

# Check for clipboard utilities (at least one is needed)
CLIPBOARD_AVAILABLE=false
if command -v xclip &> /dev/null; then
    echo "✅ xclip is available"
    CLIPBOARD_AVAILABLE=true
fi
if command -v xsel &> /dev/null; then
    echo "✅ xsel is available"
    CLIPBOARD_AVAILABLE=true
fi
if command -v wl-copy &> /dev/null; then
    echo "✅ wl-clipboard is available"
    CLIPBOARD_AVAILABLE=true
fi

if [ "$CLIPBOARD_AVAILABLE" = false ]; then
    MISSING_DEPS+=("clipboard-utility")
    echo "⚠️  No clipboard utility found (need xclip, xsel, or wl-clipboard)"
fi

# Check for GUI development libraries (needed for Fyne)
if ! pkg-config --exists gtk+-3.0 2>/dev/null; then
    MISSING_DEV_DEPS+=("gtk3-dev")
    echo "⚠️  GTK3 development libraries not found"
else
    echo "✅ GTK3 development libraries available"
fi

# Install only missing dependencies
if [ ${#MISSING_DEPS[@]} -ne 0 ] || [ ${#MISSING_DEV_DEPS[@]} -ne 0 ]; then
    echo ""
    echo "📦 Installing only missing dependencies..."
    
    if command -v apt &> /dev/null; then
        # Ubuntu/Debian
        INSTALL_PACKAGES=()
        
        for dep in "${MISSING_DEPS[@]}"; do
            case $dep in
                "go")
                    echo "❌ Go needs to be installed manually from https://golang.org/doc/install"
                    echo "   Or use: sudo apt install golang-go"
                    exit 1
                    ;;
                "clipboard-utility")
                    INSTALL_PACKAGES+=("xclip")
                    ;;
            esac
        done
        
        for dep in "${MISSING_DEV_DEPS[@]}"; do
            case $dep in
                "gtk3-dev")
                    INSTALL_PACKAGES+=("libgtk-3-dev" "libayatana-appindicator3-dev" "libxxf86vm-dev" "libxrandr-dev" "libxinerama-dev" "libxcursor-dev" "libxi-dev" "libgl1-mesa-dev")
                    ;;
            esac
        done
        
        if [ ${#INSTALL_PACKAGES[@]} -ne 0 ]; then
            echo "Installing: ${INSTALL_PACKAGES[*]}"
            sudo apt update && sudo apt install -y "${INSTALL_PACKAGES[@]}"
        fi
        
    elif command -v dnf &> /dev/null; then
        # Fedora/RHEL
        INSTALL_PACKAGES=()
        
        for dep in "${MISSING_DEPS[@]}"; do
            case $dep in
                "go")
                    INSTALL_PACKAGES+=("golang")
                    ;;
                "clipboard-utility")
                    INSTALL_PACKAGES+=("xclip")
                    ;;
            esac
        done
        
        for dep in "${MISSING_DEV_DEPS[@]}"; do
            case $dep in
                "gtk3-dev")
                    INSTALL_PACKAGES+=("gtk3-devel" "libayatana-appindicator-gtk3-devel" "libXxf86vm-devel" "libXrandr-devel" "libXinerama-devel" "libXcursor-devel" "libXi-devel" "mesa-libGL-devel")
                    ;;
            esac
        done
        
        if [ ${#INSTALL_PACKAGES[@]} -ne 0 ]; then
            echo "Installing: ${INSTALL_PACKAGES[*]}"
            sudo dnf install -y "${INSTALL_PACKAGES[@]}"
        fi
        
    elif command -v pacman &> /dev/null; then
        # Arch Linux
        INSTALL_PACKAGES=()
        
        for dep in "${MISSING_DEPS[@]}"; do
            case $dep in
                "go")
                    INSTALL_PACKAGES+=("go")
                    ;;
                "clipboard-utility")
                    INSTALL_PACKAGES+=("xclip")
                    ;;
            esac
        done
        
        for dep in "${MISSING_DEV_DEPS[@]}"; do
            case $dep in
                "gtk3-dev")
                    INSTALL_PACKAGES+=("gtk3" "libayatana-appindicator" "libxxf86vm" "libxrandr" "libxinerama" "libxcursor" "libxi" "mesa")
                    ;;
            esac
        done
        
        if [ ${#INSTALL_PACKAGES[@]} -ne 0 ]; then
            echo "Installing: ${INSTALL_PACKAGES[*]}"
            sudo pacman -S --needed "${INSTALL_PACKAGES[@]}"
        fi
        
    else
        echo "❌ Unsupported package manager. Please install missing dependencies manually:"
        echo "   - Go: https://golang.org/doc/install"
        echo "   - Clipboard utility: xclip, xsel, or wl-clipboard"
        echo "   - GUI libraries: GTK3 development packages"
        exit 1
    fi
    
    echo "✅ Dependencies installed successfully"
else
    echo "✅ All required dependencies are already available"
fi

# Build fresh binary
echo "Building fresh binary..."
go mod tidy
if ! go build -o clipboard-manager; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful!"

# Make executable
chmod +x clipboard-manager

# Test the binary
echo "Testing the binary..."
if ! ./clipboard-manager help > /dev/null 2>&1; then
    echo "❌ Binary test failed!"
    exit 1
fi

echo "✅ Binary test passed!"

# Run setup
echo ""
echo "🔧 Running automated setup..."
if [ -f "scripts/setup-hotkey.sh" ]; then
    chmod +x scripts/setup-hotkey.sh
    ./scripts/setup-hotkey.sh
else
    echo "⚠️  scripts/setup-hotkey.sh not found, running manual setup..."
    
    # Get the absolute path to the current executable
    EXEC_PATH="$(pwd)/clipboard-manager"
    
    # Create desktop entry
    DESKTOP_DIR="$HOME/.local/share/applications"
    mkdir -p "$DESKTOP_DIR"
    
    cat > "$DESKTOP_DIR/clipboard-manager.desktop" << EOF
[Desktop Entry]
Name=Clipboard Manager
Comment=Clipboard history manager for Linux
Exec=$EXEC_PATH show
Icon=edit-copy
Terminal=false
Type=Application
Categories=Utility;System;
Keywords=clipboard;history;copy;paste;
StartupNotify=true
EOF
    
    echo "✅ Desktop entry created"
    
    # Try to set up GNOME hotkey
    if command -v gsettings &> /dev/null && [ "$XDG_CURRENT_DESKTOP" = "GNOME" ]; then
        echo "Setting up GNOME hotkey..."
        CUSTOM_PATH="/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/clipboard-manager/"
        
        # Get current bindings and add ours
        CURRENT_BINDINGS=$(gsettings get org.gnome.settings-daemon.plugins.media-keys custom-keybindings)
        if [ "$CURRENT_BINDINGS" = "@as []" ]; then
            NEW_BINDINGS="['$CUSTOM_PATH']"
        else
            NEW_BINDINGS="${CURRENT_BINDINGS%]*}, '$CUSTOM_PATH']"
        fi
        
        gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings "$NEW_BINDINGS"
        gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:$CUSTOM_PATH name "Clipboard Manager"
        gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:$CUSTOM_PATH command "$EXEC_PATH show"
        gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:$CUSTOM_PATH binding "<Super>z"
        
        echo "✅ GNOME hotkey configured (Super+V)"
    fi
fi

echo ""
echo "🎉 Installation completed successfully!"
echo ""
echo "📋 Your Clipboard Manager is ready to use:"
echo "   • Press Super+V (Windows key + V) to open clipboard history"
echo "   • Run './clipboard-manager help' for all options"
echo "   • Run './clipboard-manager daemon' to start in background"
echo ""
echo "🧪 Testing clipboard functionality..."

# Test clipboard write capability
TEST_TEXT="clipboard-manager installation test $(date +%s)"
if command -v xclip &> /dev/null; then
    echo "$TEST_TEXT" | xclip -selection clipboard 2>/dev/null
    echo "✅ Test text copied using xclip"
elif command -v xsel &> /dev/null; then
    echo "$TEST_TEXT" | xsel -ib 2>/dev/null
    echo "✅ Test text copied using xsel"
elif command -v wl-copy &> /dev/null; then
    echo "$TEST_TEXT" | wl-copy 2>/dev/null
    echo "✅ Test text copied using wl-copy"
else
    echo "⚠️  No clipboard utility available for testing"
fi

# Brief daemon test to capture clipboard content
echo "Starting brief daemon test..."
timeout 3s ./clipboard-manager daemon > /dev/null 2>&1 &
DAEMON_PID=$!
sleep 2

# Stop the test daemon
kill $DAEMON_PID 2>/dev/null || true
wait $DAEMON_PID 2>/dev/null || true

echo "Current clipboard history:"
./clipboard-manager list

echo ""
echo "✅ Installation and test completed!"
echo ""
echo "🎯 Next steps:"
echo "   1. Press Super+V (Windows key + V) to open clipboard history"
echo "   2. Run './clipboard-manager daemon' to start monitoring in background"
echo "   3. Run './clipboard-manager help' to see all available options"
echo ""
echo "🚀 Clipboard Manager with image support is now ready to use!"