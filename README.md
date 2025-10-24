# Clipboard Manager for Linux

A powerful clipboard history manager for Linux systems (Ubuntu, Fedora, etc.) that tracks your clipboard history and allows you to restore previous clipboard entries with a global hotkey.

## ✨ Features

- 📋 **Smart clipboard monitoring** with automatic filtering
- ⌨️ **Global hotkey support (Super+Z)** for instant access from anywhere
- 🖥️ **GUI interface** using Fyne with automatic terminal fallback
- 🔧 **System tray integration** with right-click menu
- 💾 **Persistent history storage** (up to 50 items)
- 🔄 **Intelligent duplicate detection** and removal
- 🧹 **Advanced history management** (clear, limit, validation)
- 🐧 **Full Linux support** (X11 and Wayland)
- 🚀 **Multiple run modes** (daemon, tray, GUI, terminal)
- ⚡ **Automatic desktop environment detection** (GNOME, KDE)
- 🔗 **Desktop integration** with .desktop files and autostart

## 🚀 Quick Start

### 1. Install Dependencies

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install -y xclip libgtk-3-dev libayatana-appindicator3-dev
```

**Fedora:**
```bash
sudo dnf install -y xclip gtk3-devel libayatana-appindicator-gtk3-devel
```

### 2. Build and Setup

```bash
# Clone the repository
git clone https://github.com/MiniduTH/linux-clipboard-manager.git
cd linux-clipboard-manager

# Build the application
go build -o clipboard-manager

# Quick setup with hotkeys (recommended)
./setup-hotkey.sh
```

That's it! Press **Super+Z** (Windows key + Z) from anywhere to access your clipboard history.

## 📖 Usage

### Global Hotkey (Recommended)
After running `./setup-hotkey.sh`, press **Super+Z** from anywhere to open clipboard history.

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
Opens a graphical window showing clipboard history. Click any item to restore it to clipboard.

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

- **History location**: `~/.local/share/clipboard-manager/history.json`
- **Maximum history items**: 50 (configurable in code)
- **Auto-save interval**: 10 seconds
- **Desktop entries**: `~/.local/share/applications/`
- **Autostart**: `~/.config/autostart/` (optional)

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
   - **Command**: `/path/to/clipboard-manager show`
   - **Shortcut**: Super+Z

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
- **Issue**: Super+Z doesn't work
- **Solutions**:
  1. Run `./setup-hotkey.sh` again
  2. Check if running in proper graphical session
  3. Verify `$DISPLAY` or `$WAYLAND_DISPLAY` environment variables
  4. Set up manually in system settings

### Build Issues
- **Issue**: Compilation fails
- **Solution**: 
  1. Ensure Go 1.21+ is installed
  2. Run `go mod tidy`
  3. Install CGO dependencies for Fyne GUI

## 🛠️ Development

### Build Commands
```bash
# Build the application
make build

# Clean build artifacts
make clean

# Install dependencies
make deps

# Run the application
make run

# Run in daemon mode
make daemon
```

### Available Make Targets
- `make build` - Build the clipboard manager binary
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report
- `make clean` - Clean build artifacts
- `make install` - Install system-wide (requires sudo)
- `make help` - Show all available targets

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

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
| Global Hotkey | ✅ Super+Z | ❌ Usually not |
| GUI + Terminal | ✅ Both | ❌ Usually one |
| Smart Filtering | ✅ Yes | ❌ No |
| System Tray | ✅ Yes | ❌ Rare |
| Auto Setup | ✅ Yes | ❌ Manual |
| Cross-DE Support | ✅ GNOME/KDE | ❌ Limited |

## 🚀 Roadmap

- [ ] Wayland-native clipboard support
- [ ] Plugin system for custom filters
- [ ] Cloud sync support
- [ ] Encrypted history storage
- [ ] Custom hotkey configuration
- [ ] Clipboard search functionality
- [ ] Image clipboard support

---

**Made with ❤️ for the Linux community**

If you find this useful, please ⭐ star the repository!