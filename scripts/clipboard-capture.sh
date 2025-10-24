#!/bin/bash

# Simple script to capture clipboard content manually
# Useful when running in passive mode

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLIPBOARD_MANAGER="$SCRIPT_DIR/../clipboard-manager"

if [ ! -f "$CLIPBOARD_MANAGER" ]; then
    echo "Error: clipboard-manager not found at $CLIPBOARD_MANAGER"
    echo "Please build the application first: go build -o clipboard-manager"
    exit 1
fi

echo "Capturing current clipboard content..."
"$CLIPBOARD_MANAGER" capture

echo ""
echo "Recent clipboard history:"
"$CLIPBOARD_MANAGER" list | head -10