#!/bin/bash

echo "🧹 Cleaning up previous Clipboard Manager installation..."

# Stop any running clipboard manager processes
echo "Stopping any running clipboard manager processes..."
pkill -f clipboard-manager 2>/dev/null || true
sleep 2

# Remove old binaries
echo "Removing old binaries..."
rm -f clipboard-manager
rm -f /usr/local/bin/clipboard-manager 2>/dev/null || true

# Remove old desktop entries
echo "Removing old desktop entries..."
rm -f ~/.local/share/applications/clipboard-manager.desktop
rm -f ~/.config/autostart/clipboard-manager.desktop

# Remove old GNOME hotkey settings
echo "Removing old GNOME hotkey settings..."
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

# Remove old KDE hotkey settings
echo "Removing old KDE hotkey settings..."
if command -v kwriteconfig5 &> /dev/null; then
    kwriteconfig5 --file kglobalshortcutsrc --group "clipboard-manager.desktop" --key "_k_friendly_name" --delete 2>/dev/null || true
    kwriteconfig5 --file kglobalshortcutsrc --group "clipboard-manager.desktop" --key "show" --delete 2>/dev/null || true
fi

# Clean Go module cache (optional)
echo "Cleaning Go module cache..."
go clean -modcache 2>/dev/null || true

echo "✅ Cleanup completed!"
echo ""

# Fresh installation
echo "🚀 Starting fresh installation..."

# Check dependencies
echo "Checking dependencies..."
MISSING_DEPS=()

if ! command -v go &> /dev/null; then
    MISSING_DEPS+=("go")
fi

if ! command -v xclip &> /dev/null && ! command -v xsel &> /dev/null && ! command -v wl-copy &> /dev/null; then
    MISSING_DEPS+=("xclip")
fi

if [ ${#MISSING_DEPS[@]} -ne 0 ]; then
    echo "❌ Missing dependencies: ${MISSING_DEPS[*]}"
    echo "Installing missing dependencies..."
    
    if command -v apt &> /dev/null; then
        # Ubuntu/Debian
        sudo apt update
        for dep in "${MISSING_DEPS[@]}"; do
            case $dep in
                "go")
                    echo "Please install Go from https://golang.org/doc/install"
                    exit 1
                    ;;
                "xclip")
                    sudo apt install -y xclip libgtk-3-dev libayatana-appindicator3-dev
                    ;;
            esac
        done
    elif command -v dnf &> /dev/null; then
        # Fedora
        for dep in "${MISSING_DEPS[@]}"; do
            case $dep in
                "go")
                    sudo dnf install -y golang
                    ;;
                "xclip")
                    sudo dnf install -y xclip gtk3-devel libayatana-appindicator-gtk3-devel
                    ;;
            esac
        done
    else
        echo "Please install missing dependencies manually"
        exit 1
    fi
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
if [ -f "setup-hotkey.sh" ]; then
    chmod +x setup-hotkey.sh
    ./setup-hotkey.sh
else
    echo "⚠️  setup-hotkey.sh not found, running manual setup..."
    
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
        
        echo "✅ GNOME hotkey configured (Super+Z)"
    fi
fi

echo ""
echo "🎉 Installation completed successfully!"
echo ""
echo "📋 Your Clipboard Manager is ready to use:"
echo "   • Press Super+Z (Windows key + Z) to open clipboard history"
echo "   • Run './clipboard-manager help' for all options"
echo "   • Run './clipboard-manager daemon' to start in background"
echo ""
echo "🧪 Testing clipboard functionality..."
echo "test installation" | xclip -selection clipboard 2>/dev/null || echo "test installation" | xsel -ib 2>/dev/null || true
sleep 1

# Start daemon for a few seconds to capture the test
timeout 3s ./clipboard-manager daemon > /dev/null 2>&1 &
sleep 2

echo "Current clipboard history:"
./clipboard-manager list

echo ""
echo "✅ Installation and test completed!"
echo "🚀 Clipboard Manager is now ready to use!"