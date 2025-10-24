# Download and Install (No Go Required)

## Quick Installation

For users who don't have Go installed, you can use pre-built binaries:

### Method 1: Automatic Installation

```bash
# Download the release archive
wget https://github.com/MiniduTH/linux-clipboard-manager/releases/latest/download/clipboard-manager-release.tar.gz

# Extract and install
tar -xzf clipboard-manager-release.tar.gz
cd release/
chmod +x install.sh
./install.sh
```

### Method 2: Manual Installation

1. **Download the appropriate binary for your system:**
   - **Most users (64-bit Intel/AMD)**: `clipboard-manager-linux-amd64`
   - **Raspberry Pi 4, ARM servers**: `clipboard-manager-linux-arm64`
   - **Older 32-bit systems**: `clipboard-manager-linux-386`

2. **Install manually:**
   ```bash
   # Make it executable
   chmod +x clipboard-manager-linux-amd64
   
   # Install system-wide
   sudo cp clipboard-manager-linux-amd64 /usr/local/bin/clipboard-manager
   
   # Test it works
   clipboard-manager help
   ```

## Usage

After installation:

- **Press Super+V** (Windows key + V) to open clipboard history
- **Terminal interface**: `clipboard-manager list`
- **Help**: `clipboard-manager help`
- **Start daemon**: `clipboard-manager daemon`

## System Requirements

- **Linux** (any distribution)
- **X11 or Wayland** desktop environment
- **xclip** or **wl-clipboard** (usually pre-installed)

### Install clipboard utilities if needed:

```bash
# Ubuntu/Debian
sudo apt install xclip

# Fedora
sudo dnf install xclip

# Arch Linux
sudo pacman -S xclip
```

## Supported Architectures

- ✅ **x86_64** (Intel/AMD 64-bit) - Most common
- ✅ **ARM64** (aarch64) - Raspberry Pi 4, ARM servers
- ✅ **i386** (32-bit Intel/AMD) - Older systems

## File Sizes

- **x86_64 binary**: ~24MB
- **Complete archive**: ~12MB (compressed)

## Uninstalling

```bash
# Remove the binary
sudo rm /usr/local/bin/clipboard-manager

# Remove desktop integration
rm ~/.local/share/applications/clipboard-manager.desktop
rm ~/.config/autostart/clipboard-manager.desktop

# Remove data (optional)
rm -rf ~/.local/share/clipboard-manager
```

---

**No compilation required!** These are statically-linked binaries that work on any Linux system.