# Linux Clipboard Manager v3.4.0 - Smart Installer

## 🚀 What's New

### Smart Installation System
- **Automatic Dependency Management**: The installer now automatically detects and installs missing clipboard utilities
- **Multi-Distribution Support**: Works with Ubuntu, Debian, Fedora, RHEL, Arch, openSUSE, Alpine, Void, Gentoo, and more
- **Environment Diagnostics**: Comprehensive environment checking with detailed error messages
- **Professional Experience**: Colored output, progress indicators, and clear instructions

### Enhanced Error Handling
- **Better Diagnostics**: New `clipboard-manager diagnose` command for troubleshooting
- **Actionable Solutions**: Error messages now include specific steps to resolve issues
- **Headless Support**: Improved support for server/SSH environments
- **Graceful Fallbacks**: Better handling of missing dependencies

## 📦 Installation

### One-Command Installation
```bash
# Download and install automatically
wget https://github.com/MiniduTH/linux-clipboard-manager/releases/download/v3.4.0/clipboard-manager-smart-release.tar.gz
tar -xzf clipboard-manager-smart-release.tar.gz
cd release
./install.sh
```

The smart installer will:
- ✅ Detect your Linux distribution automatically
- ✅ Install missing clipboard utilities (xclip, xsel, wl-clipboard)
- ✅ Validate your environment
- ✅ Set up desktop integration
- ✅ Configure autostart
- ✅ Test the installation

### Supported Distributions
- **Debian/Ubuntu**: `apt install xclip xsel wl-clipboard`
- **Fedora/RHEL**: `dnf install xclip xsel wl-clipboard`
- **Arch Linux**: `pacman -S xclip xsel wl-clipboard`
- **openSUSE**: `zypper install xclip xsel wl-clipboard`
- **Alpine**: `apk add xclip xsel wl-clipboard`
- **Void Linux**: `xbps-install xclip xsel wl-clipboard`
- **Gentoo**: `emerge xclip xsel wl-clipboard`

## ✨ Features

### Core Functionality
- **Clipboard History**: Stores text and image clipboard content with SQLite
- **System Integration**: Ctrl+Shift+V hotkey for instant access
- **Cross-Desktop Support**: GNOME, KDE, XFCE, MATE, Cinnamon, and more
- **System Tray**: Optional system tray integration
- **Auto-Start**: Automatically starts with your desktop session

### New Diagnostic Tools
```bash
# Check environment and dependencies
clipboard-manager diagnose

# Get help for any issues
clipboard-manager help

# Check if daemon is running
clipboard-manager status
```

## 🔧 Troubleshooting

If you encounter the error "This application requires a graphical environment", run:
```bash
clipboard-manager diagnose
```

This will show you:
- ✅ Display server status (X11/Wayland)
- ✅ Available clipboard utilities
- ✅ Clipboard access test results
- 💡 Specific solutions for your system

### Common Solutions
- **Missing clipboard utilities**: The installer now handles this automatically
- **SSH/Remote sessions**: Use `ssh -X` for X11 forwarding
- **Headless servers**: Use `clipboard-manager daemon-passive` mode
- **Permission issues**: The installer sets up proper permissions

## 🎯 Usage

After installation:
- Press **Ctrl+Shift+V** from anywhere to open clipboard history
- Click any item to restore it to clipboard
- Use `clipboard-manager help` for all available commands
- Run `clipboard-manager diagnose` if you encounter issues

## 🐛 Bug Fixes & Improvements

- **Environment Detection**: More robust checking for graphical environments
- **Error Messages**: Clear, actionable error messages with solutions
- **Installation Process**: Automated dependency management
- **Cross-Platform**: Better support for different Linux distributions
- **User Experience**: Professional installation with progress feedback

## 📋 System Requirements

- Linux with X11 or Wayland
- Desktop environment (automatically detected)
- Clipboard utilities (automatically installed by the smart installer)

---

**Full Changelog**: https://github.com/MiniduTH/linux-clipboard-manager/compare/v3.3.0...v3.4.0

## 🙏 Note for Users

This release significantly improves the installation experience. If you previously had issues with missing dependencies or environment setup, this version should resolve those problems automatically. The smart installer takes care of everything needed to get the clipboard manager working on your system.