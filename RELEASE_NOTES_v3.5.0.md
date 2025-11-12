# Linux Clipboard Manager v3.5.0 - Auto-Startup on Installation

## 🚀 What's New

### Automatic Startup Configuration
- **Auto-Enable on First Run**: The clipboard manager now automatically enables startup on first run
- **Seamless Experience**: No manual configuration needed - works immediately after installation
- **Smart Detection**: Only enables startup if not already configured
- **User Control**: Can still be disabled with `clipboard-manager startup-disable`

### Enhanced User Experience
- **Zero Configuration**: Install and forget - the app just works
- **Persistent Operation**: Automatically starts when you turn on your PC
- **Immediate Availability**: Ctrl+Shift+V hotkey ready after login
- **Non-Intrusive**: Silent startup configuration with one-time notification

## ✨ Key Features

### Automatic Startup Management
```bash
# Startup is automatically enabled on first run
./clipboard-manager

# Check startup status
clipboard-manager startup-status

# Disable if needed
clipboard-manager startup-disable

# Re-enable manually
clipboard-manager startup-enable
```

### How It Works
1. **First Run**: When you run the clipboard manager for the first time, it automatically:
   - Creates autostart desktop entry
   - Configures for all major desktop environments (GNOME, KDE, XFCE, MATE, etc.)
   - Shows a one-time notification about the setup
   
2. **Subsequent Runs**: The app checks if startup is already configured and skips setup if enabled

3. **User Control**: You can always disable/enable startup using the dedicated commands

## 📦 Installation

### Quick Installation (Recommended)
```bash
# Download and install
wget https://github.com/MiniduTH/linux-clipboard-manager/releases/download/v3.5.0/clipboard-manager-release.tar.gz
tar -xzf clipboard-manager-release.tar.gz
cd release/
./install.sh
```

After installation:
- ✅ Clipboard manager is installed
- ✅ Startup is automatically configured
- ✅ Will start on next login
- ✅ Ctrl+Shift+V hotkey is ready

### Build from Source
```bash
git clone https://github.com/MiniduTH/linux-clipboard-manager.git
cd linux-clipboard-manager
make install
```

## 🎯 Usage

### After Installation
1. **Immediate Use**: Run `clipboard-manager` or press Ctrl+Shift+V
2. **After Reboot**: The app starts automatically - just press Ctrl+Shift+V
3. **No Configuration**: Everything works out of the box

### Startup Management Commands
```bash
# Check if startup is enabled
clipboard-manager startup-status

# Disable automatic startup
clipboard-manager startup-disable

# Re-enable automatic startup
clipboard-manager startup-enable
```

## 🔧 Technical Details

### Startup Configuration
- **Location**: `~/.config/autostart/clipboard-manager.desktop`
- **Compatibility**: Works with GNOME, KDE, XFCE, MATE, Unity, Cinnamon, Pantheon, LXQt, LXDE
- **Delay**: 3-second startup delay to avoid system load issues
- **Visibility**: Appears in System Settings > Startup Applications

### Desktop Entry Features
- **Multi-DE Support**: Compatible with all major Linux desktop environments
- **Proper Integration**: Uses standard XDG autostart specification
- **User Manageable**: Can be managed through system settings GUI
- **Non-Hidden**: Visible in startup applications list for transparency

## 🐛 Bug Fixes & Improvements

- **Startup Reliability**: Improved startup configuration logic
- **First-Run Experience**: Better user feedback on initial setup
- **Silent Operation**: No unnecessary messages on subsequent runs
- **Error Handling**: Graceful handling of permission issues

## 📋 System Requirements

- Linux with X11 or Wayland
- Desktop environment (GNOME, KDE, XFCE, MATE, etc.)
- Clipboard utilities (xclip, xsel, or wl-clipboard)

## 🔄 Upgrade Notes

### From v3.4.0 or Earlier
- Existing installations will automatically enable startup on next run
- No manual migration needed
- Existing autostart configurations are preserved
- Run `clipboard-manager startup-status` to verify

### Clean Installation
- Fresh installations automatically configure startup
- No additional steps required
- Works immediately after installation

## 💡 User Benefits

### Before v3.5.0
```bash
# Manual steps required
./clipboard-manager
# Then manually add to startup applications
# Or run clipboard-manager startup-enable
```

### With v3.5.0
```bash
# Just run once - everything is automatic
./clipboard-manager
# ✅ Startup configured automatically
# ✅ Ready to use after reboot
```

## 🎉 What This Means for Users

1. **Install Once, Use Forever**: No need to remember to start the app
2. **Seamless Experience**: Works like a native system component
3. **Always Available**: Clipboard history is always being tracked
4. **No Configuration**: Zero-config experience for most users
5. **User Control**: Can still disable if preferred

## 🚀 Future Enhancements

- Configuration file for startup preferences
- Custom startup delay settings
- Startup mode selection (daemon, tray, etc.)
- Per-desktop environment optimization

---

**Full Changelog**: https://github.com/MiniduTH/linux-clipboard-manager/compare/v3.4.0...v3.5.0

## 🙏 Thank You

This release focuses on making the clipboard manager truly "install and forget" software. We believe that good software should work seamlessly without requiring users to think about configuration.

If you encounter any issues with the automatic startup feature, please use:
```bash
clipboard-manager diagnose
clipboard-manager startup-status
```

And report any problems on our GitHub issues page.
