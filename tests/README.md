# Clipboard Manager Tests

This directory contains test utilities and scripts for the Clipboard Manager project.

## Test Files

The actual Go test files (`*_test.go`) are located in the root directory alongside the source code, following Go conventions.

## Test Scripts

- `run_tests.sh` - Comprehensive test runner script
- `test.sh` - Legacy test script (moved here for organization)

## Test Binaries

- `clipboard-manager-test` - Test binary (if present)

## Running Tests

### Option 1: Use the test runner script
```bash
./tests/run_tests.sh
```

### Option 2: Run Go tests directly from root
```bash
go test -v
```

### Option 3: Run specific test files
```bash
go test -v -run TestImageClipboard
go test -v -run TestHistory
```

## Test Coverage

The test suite covers:
- ✅ Clipboard history management
- ✅ Image clipboard functionality  
- ✅ UI component rendering
- ✅ History list item interactions
- ✅ Base64 image encoding/decoding
- ✅ Duplicate detection
- ✅ Mixed text and image content

## Test Environment

Tests are designed to work in headless environments and don't require actual clipboard utilities to be functional during testing.