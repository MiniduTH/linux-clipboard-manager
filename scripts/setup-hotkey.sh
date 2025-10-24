#!/bin/bash

# Setup script for clipboard manager hotkeys

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXEC_PATH="$SCRIPT_DIR/clipboard-manager"

echo "Setting up Clipboard Manager hotkeys..."

# Make sure the binary exists
if [ ! -f "$EXEC_PATH" ]; then
    echo "Error: clipboard-manager binary not found at $EXEC_PATH"
    echo "Please run 'go build -o clipboard-manager' first"
    exit 1
fi

# Make it executable
chmod +x "$EXEC_PATH"

# Detect desktop environment and set up hotkey
if [ "$XDG_CURRENT_DESKTOP" = "GNOME" ] || [ "$DESKTOP_SESSION" = "gnome" ]; then
    echo "Detected GNOME desktop environment"
    
    # Check if gsettings is available
    if command -v gsettings &> /dev/null; then
        echo "Setting up GNOME hotkey (Super+V)..."
        
        # Get current custom keybindings
        CURRENT_BINDINGS=$(gsettings get org.gnome.settings-daemon.plugins.media-keys custom-keybindings)
        CUSTOM_PATH="/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/clipboard-manager/"
        
        # Add our custom path if not already present
        if [[ "$CURRENT_BINDINGS" != *"$CUSTOM_PATH"* ]]; then
            if [ "$CURRENT_BINDINGS" = "@as []" ]; then
                NEW_BINDINGS="['$CUSTOM_PATH']"
            else
                # Remove the closing bracket and add our path
                NEW_BINDINGS="${CURRENT_BINDINGS%]*}, '$CUSTOM_PATH']"
            fi
            gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings "$NEW_BINDINGS"
        fi
        
        # Set the custom keybinding
        gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:$CUSTOM_PATH name "Clipboard Manager"
        gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:$CUSTOM_PATH command "$EXEC_PATH show"
        gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:$CUSTOM_PATH binding "<Super>v"
        
        echo "✓ GNOME hotkey configured successfully!"
        echo "Press Super+V (Windows key + V) to open clipboard history"
    else
        echo "gsettings not found. Please set up the hotkey manually in GNOME Settings."
    fi

elif [ "$XDG_CURRENT_DESKTOP" = "KDE" ] || [ "$DESKTOP_SESSION" = "plasma" ]; then
    echo "Detected KDE desktop environment"
    
    if command -v kwriteconfig5 &> /dev/null; then
        echo "Setting up KDE hotkey (Meta+V)..."
        
        # Create the shortcut
        kwriteconfig5 --file kglobalshortcutsrc --group "clipboard-manager.desktop" --key "_k_friendly_name" "Clipboard Manager"
        kwriteconfig5 --file kglobalshortcutsrc --group "clipboard-manager.desktop" --key "show" "Meta+V,none,Show Clipboard History"
        
        # Restart KDE shortcuts daemon
        if command -v kquitapp5 &> /dev/null && command -v kstart5 &> /dev/null; then
            kquitapp5 kglobalaccel && sleep 2 && kstart5 kglobalaccel
        fi
        
        echo "✓ KDE hotkey configured successfully!"
        echo "Press Meta+V (Windows key + V) to open clipboard history"
    else
        echo "kwriteconfig5 not found. Please set up the hotkey manually in KDE System Settings."
    fi

else
    echo "Desktop environment not detected or not supported for automatic setup."
    echo ""
    echo "Manual setup instructions:"
    echo "1. Open your system settings"
    echo "2. Go to Keyboard Shortcuts"
    echo "3. Add a custom shortcut:"
    echo "   Name: Clipboard Manager"
    echo "   Command: $EXEC_PATH show"
    echo "   Shortcut: Super+V (or Windows key + V)"
fi

# Create desktop entry
echo ""
echo "Creating desktop entry..."
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

echo "✓ Desktop entry created"

# Create autostart entry
echo ""
read -p "Do you want to start Clipboard Manager automatically on login? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    AUTOSTART_DIR="$HOME/.config/autostart"
    mkdir -p "$AUTOSTART_DIR"
    
    cat > "$AUTOSTART_DIR/clipboard-manager.desktop" << EOF
[Desktop Entry]
Name=Clipboard Manager
Comment=Clipboard history manager for Linux
Exec=$EXEC_PATH daemon
Icon=edit-copy
Terminal=false
Type=Application
Categories=Utility;System;
X-GNOME-Autostart-enabled=true
Hidden=false
NoDisplay=false
EOF
    
    echo "✓ Autostart entry created"
    echo "Clipboard Manager will start automatically on next login"
fi

echo ""
echo "Setup complete! You can now:"
echo "1. Press Super+V to open clipboard history"
echo "2. Run '$EXEC_PATH help' for more options"
echo "3. Start the daemon with '$EXEC_PATH daemon'"

# Test the setup
echo ""
read -p "Would you like to test the clipboard manager now? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Starting clipboard manager in daemon mode for 5 seconds..."
    echo "Copy some text to test it, then check with '$EXEC_PATH list'"
    timeout 5s "$EXEC_PATH" daemon &
    DAEMON_PID=$!
    
    echo "Daemon started (PID: $DAEMON_PID). Copy some text now..."
    sleep 5
    
    echo ""
    echo "Current clipboard history:"
    "$EXEC_PATH" list
fi