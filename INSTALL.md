# Installation Guide

## Quick Install (Recommended)

```bash
git clone https://github.com/MiniduTH/linux-clipboard-manager.git
cd linux-clipboard-manager
make install
```

That's it! Press **Super+V** (Windows key + V) from anywhere to access your clipboard history.

## Build System

This project uses a comprehensive Makefile for all operations:

### Essential Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the application |
| `make install` | Install system-wide with hotkey setup |
| `make test` | Run all tests |
| `make clean` | Clean build artifacts |
| `make uninstall` | Remove installation |
| `make help` | Show all available commands |

### What `make install` Does

1. **Builds** the clipboard manager binary
2. **Installs** it to `/usr/local/bin/clipboard-manager`
3. **Creates** desktop entry for application launcher
4. **Configures** Super+V hotkey automatically
5. **Starts** background daemon for clipboard monitoring

### Manual Installation Steps

If you prefer manual control:

```bash
# 1. Build
make build

# 2. Install binary
sudo cp clipboard-manager /usr/local/bin/

# 3. Set up hotkey (optional)
./scripts/setup-hotkey.sh

# 4. Start daemon
clipboard-manager daemon &
```

## Usage After Installation

- **Super+V**: Open clipboard history GUI
- **Terminal**: `clipboard-manager list` for text interface
- **Help**: `clipboard-manager help` for all options

## Uninstalling

```bash
# Remove everything
make uninstall-complete

# Or just system installation
make uninstall
```

## Troubleshooting

**Build fails?**
- Ensure Go 1.21+ is installed
- Run `make deps` to install dependencies

**Hotkey doesn't work?**
- Run `make install` again
- Check you're in a graphical session
- Set up manually in system settings

**GUI doesn't open?**
- Install GUI dependencies (see README.md)
- Use `clipboard-manager list` for terminal interface