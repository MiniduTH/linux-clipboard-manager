# Clipboard Manager for Linux

A clipboard history manager for Linux systems (Ubuntu, Fedora, etc.) that tracks your clipboard history and allows you to restore previous clipboard entries.

## Features

- 📋 Automatic clipboard monitoring
- 🖥️ GUI interface using Fyne
- 💻 Terminal fallback interface
- 💾 Persistent history storage
- 🔄 Duplicate detection and removal
- 🧹 History management (clear, limit)
- 🐧 Full Linux support (X11 and Wayland)

## Requirements

- Go 1.21 or later
- Linux with X11 or Wayland
- Clipboard utility: `xclip`, `xsel`, or `wl-clipboard`
- For GUI: GTK development libraries

## Installation

### Quick Install

```bash
chmod +x install.sh
./install.sh
```

### Manual Installation

1. **Install dependencies:**

   **Ubuntu/Debian:**
   ```bash
   sudo apt update
   sudo apt install -y xclip libgtk-3-dev libayatana-appindicator3-dev
   ```

   **Fedora:**
   ```bash
   sudo dnf install -y xclip gtk3-devel libayatana-appindicator-gtk3-devel
   ```

2. **Build the application:**
   ```bash
   go mod tidy
   go build -o clipboard-manager
   ```

## Usage

### Start Clipboard Watcher
```bash
./clipboard-manager
```
This starts the background clipboard monitoring service.

### Show GUI History
```bash
./clipboard-manager show
```
Opens a graphical window showing clipboard history. Click any item to restore it to clipboard.

### Show Terminal History
```bash
./clipboard-manager list
```
Displays clipboard history in the terminal.

### Run in Background
```bash
nohup ./clipboard-manager > /dev/null 2>&1 &
```

## Configuration

- History is stored in `~/.local/share/clipboard-manager/history.json`
- Maximum history items: 50 (configurable in code)
- Auto-save interval: 10 seconds

## Troubleshooting

### GUI Not Working
If the GUI doesn't work, the application will automatically fall back to terminal mode.

**Common issues:**
- No display server: Ensure you're running in X11 or Wayland session
- Missing GTK libraries: Install development packages for your distribution
- Permission issues: Ensure proper access to clipboard

### Clipboard Not Working
- Install clipboard utilities: `sudo apt install xclip` (Ubuntu) or `sudo dnf install xclip` (Fedora)
- Check if running in proper graphical session
- Verify `$DISPLAY` or `$WAYLAND_DISPLAY` environment variables are set

### Build Issues
- Ensure Go 1.21+ is installed
- Run `go mod tidy` to resolve dependencies
- Install CGO dependencies for Fyne GUI

## Development

The application consists of three main components:

- `main.go` - Main application logic and clipboard monitoring
- `ui.go` - GUI interface using Fyne
- `history.go` - History management and persistence

## License

MIT License - feel free to modify and distribute.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test on both Ubuntu and Fedora
5. Submit a pull request

## System Integration

To start automatically with your desktop session, create a desktop entry:

```bash
cat > ~/.config/autostart/clipboard-manager.desktop << EOF
[Desktop Entry]
Type=Application
Name=Clipboard Manager
Exec=/path/to/clipboard-manager
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
EOF
```