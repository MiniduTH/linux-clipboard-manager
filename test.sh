#!/bin/bash

echo "Testing Clipboard Manager..."

# Test 1: Help command
echo "1. Testing help command:"
./clipboard-manager help
echo ""

# Test 2: List command (should show existing history)
echo "2. Testing list command:"
./clipboard-manager list
echo ""

# Test 3: Add some test content
echo "3. Adding test content to clipboard:"
echo "Test content 1" | xclip -selection clipboard
sleep 1
echo "Test content 2" | xclip -selection clipboard
sleep 1

# Test 4: Run watcher briefly
echo "4. Running watcher for 3 seconds:"
timeout 3s ./clipboard-manager &
WATCHER_PID=$!
sleep 1
echo "Adding more content while watcher is running..."
echo "Test content 3" | xclip -selection clipboard
wait $WATCHER_PID 2>/dev/null

# Test 5: Check updated history
echo ""
echo "5. Final history check:"
./clipboard-manager list

echo ""
echo "All tests completed!"