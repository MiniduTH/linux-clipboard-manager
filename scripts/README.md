# Scripts Directory

This directory contains installation and setup scripts for the Clipboard Manager.

## Installation Scripts

### `clean-install.sh`
Complete clean installation script that:
- Removes old clipboard-manager installations
- Preserves existing system dependencies
- Only installs missing dependencies
- Sets up hotkeys and desktop integration
- Tests the installation

**Usage:**
```bash
./scripts/clean-install.sh
```

### `install.sh`
Basic installation script for quick setup.

**Usage:**
```bash
./scripts/install.sh
```

### `setup-hotkey.sh`
Sets up system hotkeys and desktop integration:
- GNOME hotkey configuration (Super+Z)
- KDE hotkey configuration
- Desktop entry creation
- Autostart configuration

**Usage:**
```bash
./scripts/setup-hotkey.sh
```

## Running Scripts

All scripts should be run from the project root directory:

```bash
# From project root
./scripts/clean-install.sh
./scripts/setup-hotkey.sh
```

## Script Features

- **Conservative approach** - Only modifies clipboard-manager specific settings
- **Multi-distro support** - Works with Ubuntu, Fedora, Arch Linux
- **Dependency preservation** - Never removes existing system packages
- **Error handling** - Graceful fallbacks and clear error messages