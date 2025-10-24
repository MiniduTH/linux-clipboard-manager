#!/bin/bash

# Test Runner for Clipboard Manager
# This script runs all tests from the tests directory

echo "🧪 Running Clipboard Manager Tests..."
echo "======================================"

# Change to the parent directory (where the main Go files are)
cd "$(dirname "$0")/.."

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go to run tests."
    exit 1
fi

# Run tests with verbose output
echo "Running Go tests..."
if go test -v; then
    echo ""
    echo "✅ All tests passed!"
else
    echo ""
    echo "❌ Some tests failed!"
    exit 1
fi

echo ""
echo "🎯 Test Summary:"
echo "   • History management tests"
echo "   • Image clipboard tests" 
echo "   • UI integration tests"
echo "   • History list item tests"
echo ""
echo "✅ Test run completed!"