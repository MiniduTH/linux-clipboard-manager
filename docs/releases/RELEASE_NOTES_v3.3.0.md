# Linux Clipboard Manager v3.3.0

## 🚀 What's New

### Build Optimizations
- **Fast Build Targets**: Added multiple build options for faster development
  - `make fast` - Quick build with minimal output
  - `make ultra-fast` - Fastest compilation (no optimizations)
  - `make dev` - Development build for iteration
  - `make parallel` - Uses all CPU cores for compilation
- **Improved Build Times**: Optimized Makefile reduces compilation time
- **Build Cache Support**: Better utilization of Go's build cache

### Release Package
- **Professional Installer**: Automated `install.sh` script for easy setup
- **Cross-Platform Support**: Optimized binaries for different architectures
- **Complete Documentation**: Comprehensive installation and usage guides
- **No Dependencies**: Users don't need Go installed to use the release

## 📦 Installation

### Quick Install (Recommended)
```bash
# Download and extract the release
wget https://github.com/MiniduTH/linux-clipboard-manager/releases/download/v3.3.0/clipboard-manager-release.tar.gz
tar -xzf clipboard-manager-release.tar.gz
cd release

# Run the installer
chmod +x install.sh
./install.sh
```

### Manual Installation
```bash
# Download the binary for your architecture
sudo cp clipboard-manager-linux-amd64 /usr/local/bin/clipboard-manager
sudo chmod +x /usr/local/bin/clipboard-manager
```

## ✨ Features

- **Clipboard History**: Stores text and image clipboard content
- **System Integration**: Ctrl+Shift+V hotkey for instant access
- **Cross-Desktop Support**: Works with GNOME, KDE, XFCE, MATE, and more
- **System Tray**: Optional system tray integration
- **Auto-Start**: Automatically starts with your desktop session
- **SQLite Storage**: Reliable database storage for clipboard history
- **Image Support**: Handles PNG, JPEG, and other image formats

## 🎯 Usage

After installation:
- Press **Ctrl+Shift+V** from anywhere to open clipboard history
- Click any item to restore it to clipboard
- Use `clipboard-manager help` for all available commands
- Access through system tray (if enabled)

## 🔧 Development

For developers:
```bash
# Fast development builds
make fast          # Quick build
make dev           # Development build
make parallel      # Multi-core compilation

# Create releases
make release       # Full release package
make dist          # Distributable archive
```

## 📋 System Requirements

- Linux with X11 or Wayland
- Clipboard utilities (xclip, xsel, or wl-clipboard)
- Desktop environment (GNOME, KDE, XFCE, etc.)

## 🐛 Bug Fixes & Improvements

- Optimized build process for faster development cycles
- Enhanced release packaging with better user experience
- Improved documentation and installation guides
- Better error handling in build scripts

---

**Full Changelog**: https://github.com/MiniduTH/linux-clipboard-manager/compare/v3.2.0...v3.3.0