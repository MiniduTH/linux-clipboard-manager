# Release Notes - Version 3.7.0

## 🎨 Enhanced Edit Experience

This release focuses on improving the edit dialog usability and adding developer tools.

### ✨ New Features

#### Improved Edit Dialog
- **Larger Edit Area**: Increased from 10 to 15 visible rows for better content visibility
- **Bigger Dialog Window**: Expanded to 700x550 pixels (from 600x400)
- **Enhanced Scroll Area**: Explicit sizing (650x450) for smoother editing experience
- **Better Multi-line Support**: More comfortable editing of long text content

#### Developer Tools
- **Update Script**: New `scripts/update.sh` for easy rebuilding and reinstalling
- **Smart Process Management**: Handles "Text file busy" errors automatically
- **Graceful Shutdown**: Stops running instances before rebuilding
- **Auto-restart**: Automatically starts the manager after installation

### 🔧 Technical Improvements

#### Edit Dialog Enhancements
- Increased `SetMinRowsVisible` from 10 to 15 rows
- Set explicit minimum size for entry scroll container (650x450)
- Enlarged overall dialog dimensions (700x550)
- Better visual hierarchy and content visibility

#### Update Script Features
- Graceful process termination with fallback to force kill
- Uses existing `make build` workflow
- Waits between operations to ensure clean state
- Verifies successful restart
- Helpful error messages with emoji indicators

### 📖 Documentation

**Updated Files:**
- `ui.go` - Enhanced edit dialog sizing and layout
- `scripts/update.sh` - New developer update script

### 🎯 Use Cases

**Edit Dialog Improvements:**
- Editing longer text snippets more comfortably
- Better visibility of multi-line content
- Reduced scrolling when editing
- More natural editing experience

**Update Script:**
- Quick development iteration
- Safe rebuilding without manual process management
- Automated testing of changes
- Simplified deployment workflow

### 🚀 How to Use

#### Enhanced Edit Dialog
1. Open clipboard history with `Ctrl+Shift+V`
2. Click the edit button on any text item
3. Enjoy the larger, more comfortable editing area
4. Save your changes

#### Update Script
```bash
# After making code changes
./scripts/update.sh

# The script will:
# 1. Stop the running clipboard manager
# 2. Rebuild the binary
# 3. Install system-wide
# 4. Restart the service
```

### 📋 What's Changed

**Modified Files:**
- `ui.go` - Enhanced edit dialog dimensions and layout

**New Files:**
- `scripts/update.sh` - Developer update automation script

**Moved Files:**
- `update.sh` → `scripts/update.sh` - Better project organization

### 🔄 Upgrade Notes

This release is fully backward compatible. The enhanced edit dialog provides a better user experience with no breaking changes.

**Installation:**
```bash
# Download the release
wget https://github.com/MiniduTH/linux-clipboard-manager/releases/download/v3.7.0/clipboard-manager-v3.7.0-linux-amd64.tar.gz

# Extract
tar -xzf clipboard-manager-v3.7.0-linux-amd64.tar.gz

# Install
chmod +x install.sh
./install.sh
```

### 🎨 UI/UX Improvements

- 50% more visible content in edit dialog (15 vs 10 rows)
- 37% larger dialog window area
- Better proportions for comfortable editing
- Reduced need for scrolling during edits

### 📊 Performance

- No performance impact
- Same efficient database operations
- Minimal memory overhead from larger UI elements

### 🙏 Acknowledgments

Thanks to users who provided feedback about the edit dialog usability!

### 📝 Full Changelog

**Features:**
- Enhanced edit dialog with larger dimensions and better visibility
- Added developer update script for streamlined workflow

**Improvements:**
- Increased edit dialog visible rows from 10 to 15
- Expanded dialog window size to 700x550 pixels
- Set explicit scroll area dimensions (650x450)
- Automated rebuild and restart process

**Developer Experience:**
- Smart process management in update script
- Handles "Text file busy" errors gracefully
- Automated testing workflow
- Better project organization with scripts folder

---

## Previous Releases

### v3.6.0 - Edit Clipboard Items Feature
- Edit text items directly in clipboard history
- Multi-line editor with validation
- Modal dialog interface

### v3.5.0 - Auto-Startup on Installation
- Automatic startup application setup
- Enhanced installation script

### v3.4.0 - Smart Installer
- Automatic dependency management
- Improved error handling

---

**Download:** [clipboard-manager-v3.7.0-linux-amd64.tar.gz](https://github.com/MiniduTH/linux-clipboard-manager/releases/download/v3.7.0/clipboard-manager-v3.7.0-linux-amd64.tar.gz)

**Full Changelog:** https://github.com/MiniduTH/linux-clipboard-manager/compare/v3.6.0...v3.7.0
