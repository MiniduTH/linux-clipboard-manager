# Clipboard Manager for Linux

A powerful clipboard history manager for Linux systems (Ubuntu, Fedora, etc.) that tracks your clipboard history and allows you to restore previous clipboard entries with a global hotkey.

**🚀 No Go installation required!** Download pre-built binaries from [Releases](https://github.com/MiniduTH/linux-clipboard-manager/releases/latest) or see [DOWNLOAD.md](DOWNLOAD.md) for quick installation.

**⚡ Auto-starts on login!** Installation automatically sets up the clipboard manager to start when you turn on your PC.

**✨ Latest Release: v3.7.0** - Enhanced edit dialog with better usability and developer tools.

## ✨ Features

- 📋 **Smart clipboard monitoring** with automatic filtering
- ⌨️ **Global hotkey support (Ctrl+Shift+V)** for instant access from anywhere
- 🖥️ **GUI interface** using Fyne with automatic terminal fallback
- ✏️ **Edit clipboard items** - modify text content directly in the history
- 🔧 **System tray integration** with right-click menu
- 💾 **SQLite database storage** with automatic JSON migration (up to 50 items)
- 🔄 **Intelligent duplicate detection** and removal
- 🧹 **Advanced history management** (clear, limit, validation)
- 🐧 **Full Linux support** (X11 and Wayland)
- 🚀 **Multiple run modes** (daemon, tray, GUI, terminal)
- ⚡ **Automatic desktop environment detection** (GNOME, KDE)
- 🔗 **Desktop integration** with .desktop files and automatic autostart
- 🖥️ **Auto-starts on login** - works immediately after PC restart

## 🚀 Quick Start

### Option 1: No Go Required (Recommended for Most Users)

**Download pre-built release:**
1. Go to [Releases](https://github.com/MiniduTH/linux-clipboard-manager/releases/latest)
2. Download `clipboard-manager-release.tar.gz`
3. Extract and install:
   ```bash
   wget https://github.com/MiniduTH/linux-clipboard-manager/releases/download/v3.7.0/clipboard-manager-release.tar.gz
   tar -xzf clipboard-manager-release.tar.gz
   cd release/
   chmod +x install.sh
   ./install.sh
   ```

That's it! Press **Ctrl+Shift+V** from anywhere to access your clipboard history.

### Option 2: Build from Source (For Developers)

```bash
# Clone and install in one go
git clone https://github.com/MiniduTH/linux-clipboard-manager.git
cd linux-clipboard-manager
make install
```

### Option 2: Manual Setup

#### 1. Install Dependencies

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install -y xclip libgtk-3-dev libayatana-appindicator3-dev \
    libxxf86vm-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libgl1-mesa-dev
```

**Fedora:**
```bash
sudo dnf install -y xclip gtk3-devel libayatana-appindicator-gtk3-devel \
    libXxf86vm-devel libXrandr-devel libXinerama-devel libXcursor-devel libXi-devel mesa-libGL-devel
```

**Arch Linux:**
```bash
sudo pacman -S xclip gtk3 libayatana-appindicator \
    libxxf86vm libxrandr libxinerama libxcursor libxi mesa
```

#### 2. Build and Install

```bash
# Clone the repository
git clone https://github.com/MiniduTH/linux-clipboard-manager.git
cd linux-clipboard-manager

# Build and install system-wide
make build
make install

# Or just build for development
make build
./clipboard-manager help
```

## 📖 Usage

### Global Hotkey (Recommended)
After running `make install`, press **Ctrl+Shift+V** from anywhere to open clipboard history.

### Manual Commands

#### System Integration Mode
```bash
./clipboard-manager
```
Starts clipboard monitoring and attempts to set up system hotkeys automatically.

#### Show GUI History
```bash
./clipboard-manager show
```
Opens a graphical window showing clipboard history. 
- **Click** any item to restore it to clipboard
- **Edit button** (pencil icon) to modify text content
- **Delete button** (X icon) to remove items

#### Show Terminal History
```bash
./clipboard-manager list
```
Displays clipboard history in the terminal - perfect for SSH sessions or minimal setups.

#### System Tray Mode
```bash
./clipboard-manager tray
```
Runs with a system tray icon. Right-click the tray icon for menu options.

#### Background Daemon Mode
```bash
./clipboard-manager daemon
```
Runs in background without GUI or hotkeys - ideal for servers or minimal setups.

#### Help
```bash
./clipboard-manager help
```
Shows all available commands and options.

## ⚙️ Configuration

- **History storage**: `~/.local/share/clipboard-manager/history.db` (SQLite database)
- **Maximum history items**: 50 (configurable in code)
- **Database features**: Automatic migration from JSON, duplicate detection, efficient queries
- **Desktop entries**: `~/.local/share/applications/`
- **Autostart**: `~/.config/autostart/` (optional)

### Database Migration

The application automatically migrates existing JSON history files to SQLite database format:
- **Old format**: `~/.local/share/clipboard-manager/history.json`
- **New format**: `~/.local/share/clipboard-manager/history.db`
- **Backup**: Original JSON file is backed up as `history.json.backup`

## 🧪 Testing

The project includes comprehensive tests for all functionality:

### Run Tests
```bash
# Using Makefile (recommended)
make test

# Using test runner script
./tests/run_tests.sh

# Direct Go test command
go test -v
```

### Test Coverage
```bash
# Generate coverage report
make test-coverage

# View coverage in browser
open coverage.html
```

### Test Organization
- **Test files**: `*_test.go` files in root directory (following Go conventions)
- **Test utilities**: `tests/` directory contains test scripts and documentation
- **Test coverage**: History management, image clipboard, UI components, and integrations

## 🔧 Advanced Setup

### Manual Hotkey Setup
If automatic setup doesn't work:

1. Open your system settings
2. Go to Keyboard Shortcuts
3. Add a custom shortcut:
   - **Name**: Clipboard Manager
   - **Command**: `/usr/local/bin/clipboard-manager show`
   - **Shortcut**: Ctrl+Shift+V

### Autostart Setup
To start automatically on login:
```bash
# The setup script can do this automatically, or manually:
cp clipboard-manager.desktop ~/.config/autostart/
```

### System Service (Advanced)
For system-wide installation:
```bash
sudo cp clipboard-manager /usr/local/bin/
sudo systemctl --user enable clipboard-manager.service
```

## 🐛 Troubleshooting

### GUI Not Working
- **Issue**: GUI doesn't open
- **Solution**: The app automatically falls back to terminal mode. Install GTK development libraries:
  ```bash
  # Ubuntu/Debian
  sudo apt install libgtk-3-dev libayatana-appindicator3-dev
  
  # Fedora
  sudo dnf install gtk3-devel libayatana-appindicator-gtk3-devel
  ```

### Clipboard Not Working
- **Issue**: Clipboard monitoring fails
- **Solution**: Install clipboard utilities:
  ```bash
  # Ubuntu/Debian
  sudo apt install xclip
  
  # Fedora
  sudo dnf install xclip
  ```

### Hotkey Not Working
- **Issue**: Ctrl+Shift+V doesn't work
- **Solutions**:
  1. Run `make install` again to reconfigure hotkeys
  2. Check if running in proper graphical session
  3. Verify `$DISPLAY` or `$WAYLAND_DISPLAY` environment variables
  4. Set up manually in system settings
  5. Check if another application is using the same hotkey

### Build Issues
- **Issue**: Compilation fails
- **Solution**: 
  1. Ensure Go 1.21+ is installed
  2. Run `go mod tidy`
  3. Install CGO dependencies for Fyne GUI

## 📁 Project Structure

```
clipboard-manager/
├── docs/              # Documentation and guides
│   ├── releases/      # Release notes
│   ├── CONTRIBUTING.md
│   ├── EDIT_FEATURE_GUIDE.md
│   └── README.md
├── scripts/           # Installation and setup scripts
│   ├── clean-install.sh
│   ├── install.sh
│   ├── setup-hotkey.sh
│   └── update.sh      # Developer update script
├── tests/             # Test utilities and documentation
│   ├── run_tests.sh
│   └── README.md
├── build/             # Build artifacts (gitignored)
├── .github/           # GitHub workflows
├── .kiro/             # Kiro IDE specifications
├── *.go               # Go source files (main package)
├── *_test.go          # Go test files
├── Makefile           # Build automation
└── README.md          # This file
```

## 🗑️ Uninstalling

### Complete Uninstall
To completely remove the clipboard manager from your system:

```bash
# Using the uninstall script (recommended)
./scripts/uninstall.sh

# Or using Makefile
make uninstall-complete
```

This will remove:
- System-wide installation
- Desktop entries and autostart files
- All clipboard history data
- Custom keyboard shortcuts
- Running processes
- Build artifacts

### System-wide Only
To remove only the system-wide installation:
```bash
make uninstall
```

## 🛠️ Development

### Build System

This project uses a comprehensive Makefile for all build and installation tasks:

```bash
# Build the application
make build

# Build all variants (standard, debug, optimized)
make build-all

# Install system-wide with hotkey setup
make install

# Run tests
make test

# Run tests with coverage
make test-coverage

# Clean build artifacts
make clean

# Install dependencies
make deps

# Create release package
make release
```

### Available Make Targets

**Build & Development:**
- `make build` - Build the clipboard manager binary
- `make build-all` - Build multiple variants (standard, debug, optimized)
- `make release` - Create release package with documentation
- `make deps` - Install/update Go dependencies
- `make clean` - Clean build artifacts

**Testing:**
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report

**Installation:**
- `make install` - Install system-wide with hotkey setup (requires sudo)
- `make uninstall` - Uninstall system-wide (requires sudo)
- `make uninstall-complete` - Complete uninstall (removes everything)

**Runtime:**
- `make run` - Run the application
- `make daemon` - Run in daemon mode
- `make show` - Show GUI

**Help:**
- `make help` - Show all available targets with descriptions

### Quick Development Workflow

```bash
# Initial setup
git clone <repository>
cd clipboard-manager
make deps

# Development cycle
make build    # Build
make test     # Test
make run      # Test run

# Quick update and restart (after code changes)
./scripts/update.sh  # Stops, rebuilds, installs, and restarts

# Installation
make install  # Install with hotkeys

# Cleanup
make clean    # Clean build files
make uninstall-complete  # Remove everything
```

### Developer Update Script

The `scripts/update.sh` script automates the rebuild and restart process:

```bash
./scripts/update.sh
```

This script will:
1. Stop the running clipboard manager
2. Rebuild the binary using `make build`
3. Install system-wide (requires sudo)
4. Restart the service automatically

Perfect for quick iteration during development!

### Creating Releases (For Maintainers)

```bash
# Create release binaries for all architectures
make release

# Create distributable archive
make dist

# The archive will be at: build/clipboard-manager-release.tar.gz
# Users can extract and run ./install.sh without Go installed
```

## 🔧 Troubleshooting

### Build Issues

**Error: `cannot find -lXxf86vm`**
```bash
# Ubuntu/Debian
sudo apt install libxxf86vm-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libgl1-mesa-dev

# Fedora
sudo dnf install libXxf86vm-devel libXrandr-devel libXinerama-devel libXcursor-devel libXi-devel mesa-libGL-devel

# Arch Linux
sudo pacman -S libxxf86vm libxrandr libxinerama libxcursor libxi mesa
```

**Error: GUI doesn't work**
- Ensure you're running in a graphical environment (X11 or Wayland)
- Check that `$DISPLAY` or `$WAYLAND_DISPLAY` environment variables are set
- Install required GUI libraries as shown in the dependencies section

**Error: Clipboard utilities not found**
- Install at least one clipboard utility: `xclip`, `xsel`, or `wl-clipboard`
- For Wayland: `sudo apt install wl-clipboard`
- For X11: `sudo apt install xclip` or `sudo apt install xsel`

## 🤝 Contributing

We welcome contributions! Please see [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for guidelines.

### Development Setup
```bash
git clone https://github.com/MiniduTH/linux-clipboard-manager.git
cd linux-clipboard-manager
go mod tidy
go build -o clipboard-manager
```

### Testing
```bash
# Run the test script
./test.sh

# Test different modes
./clipboard-manager daemon &
./clipboard-manager list
./clipboard-manager show
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🌟 Features Comparison

| Feature | This App | Other Tools |
|---------|----------|-------------|
| Global Hotkey | ✅ Ctrl+Shift+V | ❌ Usually not |
| Edit Items | ✅ Yes | ❌ Rare |
| GUI + Terminal | ✅ Both | ❌ Usually one |
| Image Support | ✅ Yes | ❌ Limited |
| Smart Filtering | ✅ Yes | ❌ No |
| System Tray | ✅ Yes | ❌ Rare |
| Auto Setup | ✅ Yes | ❌ Manual |
| Cross-DE Support | ✅ GNOME/KDE | ❌ Limited |

## 📝 Recent Updates

### Version 3.7.0 - Enhanced Edit Experience
- **Larger Edit Dialog**: 50% more visible content (15 vs 10 rows)
- **Bigger Window**: Expanded to 700x550 pixels for comfortable editing
- **Developer Tools**: New `scripts/update.sh` for easy rebuilding
- **Bug Fixes**: Fixed GitHub workflows and test compatibility

[View Full Changelog](https://github.com/MiniduTH/linux-clipboard-manager/releases/tag/v3.7.0)

## 🚀 Roadmap

- [x] Image clipboard support
- [x] Edit clipboard items
- [ ] Wayland-native clipboard support
- [ ] Plugin system for custom filters
- [ ] Cloud sync support
- [ ] Encrypted history storage
- [ ] Custom hotkey configuration
- [ ] Clipboard search functionality

---

**Made with ❤️ for the Linux community**

If you find this useful, please ⭐ star the repository!