#!/bin/bash

# Clipboard Manager Control Script
# This script helps manage the clipboard manager to avoid conflicts

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLIPBOARD_MANAGER="$SCRIPT_DIR/../clipboard-manager"

# Make sure the binary exists and is executable
if [ ! -f "$CLIPBOARD_MANAGER" ]; then
    echo "Error: clipboard-manager binary not found at $CLIPBOARD_MANAGER"
    echo "Please run 'make build' first."
    exit 1
fi

chmod +x "$CLIPBOARD_MANAGER"

show_help() {
    echo "Clipboard Manager Control Script"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  start-daemon    - Start clipboard daemon (background monitoring only)"
    echo "  start-with-gui  - Start with hotkey support (Super+Z opens GUI)"
    echo "  start-tray      - Start with system tray integration"
    echo "  show-gui        - Show GUI once (starts daemon if needed)"
    echo "  stop            - Stop all clipboard manager processes"
    echo "  status          - Show current status"
    echo "  clean           - Remove autostart files and stop processes"
    echo "  help            - Show this help"
    echo ""
    echo "Recommended usage:"
    echo "  For background monitoring: $0 start-daemon"
    echo "  For GUI access: $0 show-gui (when needed)"
    echo "  For system tray: $0 start-tray"
}

stop_all() {
    echo "Stopping all clipboard manager processes..."
    pkill -f clipboard-manager 2>/dev/null || true
    
    # Remove PID file
    rm -f ~/.local/share/clipboard-manager/daemon.pid 2>/dev/null || true
    
    echo "✓ All clipboard manager processes stopped"
}

clean_autostart() {
    echo "Cleaning up autostart files..."
    rm -f ~/.config/autostart/clipboard-manager.desktop 2>/dev/null || true
    echo "✓ Autostart files removed"
}

check_status() {
    if pgrep -f "clipboard-manager" > /dev/null; then
        echo "✓ Clipboard manager is running"
        echo "Running processes:"
        ps aux | grep clipboard-manager | grep -v grep | while read line; do
            echo "  $line"
        done
    else
        echo "✗ Clipboard manager is not running"
    fi
    
    if [ -f ~/.config/autostart/clipboard-manager.desktop ]; then
        echo "⚠ Autostart file exists (may cause conflicts)"
    else
        echo "✓ No autostart file found"
    fi
}

case "$1" in
    "start-daemon")
        echo "Starting clipboard daemon (background monitoring only)..."
        stop_all
        "$CLIPBOARD_MANAGER" daemon-only &
        echo "✓ Daemon started. Use '$0 show-gui' to open GUI when needed."
        echo "  No hotkeys or automatic GUI windows will appear."
        ;;
    "start-with-gui")
        echo "Starting clipboard manager with hotkey support..."
        stop_all
        "$CLIPBOARD_MANAGER" &
        echo "✓ Started with hotkey support. Press Super+Z to open GUI."
        ;;
    "start-tray")
        echo "Starting clipboard manager with system tray..."
        stop_all
        "$CLIPBOARD_MANAGER" tray &
        echo "✓ System tray started. Right-click tray icon for options."
        ;;
    "show-gui")
        echo "Opening clipboard GUI..."
        "$CLIPBOARD_MANAGER" show
        ;;
    "stop")
        stop_all
        ;;
    "status")
        check_status
        ;;
    "clean")
        stop_all
        clean_autostart
        echo "✓ Cleanup complete"
        ;;
    "help"|"")
        show_help
        ;;
    *)
        echo "Unknown command: $1"
        echo "Use '$0 help' for available commands."
        exit 1
        ;;
esac