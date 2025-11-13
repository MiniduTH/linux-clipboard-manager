#!/bin/bash

# Update script for clipboard manager
# This script stops the running instance, rebuilds, and reinstalls the application

set -e  # Exit on error

echo "🛑 Stopping clipboard manager..."
pkill -f clipboard-manager || true
sleep 1  # Give it time to fully stop

# Check if still running
if pgrep -f clipboard-manager > /dev/null; then
    echo "⚠️  Force killing clipboard manager..."
    pkill -9 -f clipboard-manager || true
    sleep 1
fi

echo "🔨 Building clipboard manager..."
make build

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "📥 Installing system-wide..."
sudo cp clipboard-manager /usr/local/bin/

if [ $? -ne 0 ]; then
    echo "❌ Installation failed!"
    echo "💡 Tip: Make sure the process is fully stopped"
    exit 1
fi

echo "🚀 Starting clipboard manager..."
nohup clipboard-manager > /dev/null 2>&1 &
sleep 1

if pgrep -f clipboard-manager > /dev/null; then
    echo "✅ Clipboard manager updated and restarted successfully!"
    echo "📍 PID: $(pgrep -f clipboard-manager)"
else
    echo "⚠️  Clipboard manager may not have started. Check manually."
fi
