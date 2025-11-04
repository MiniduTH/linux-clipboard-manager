#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Detect Linux distribution
detect_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        DISTRO=$ID
        VERSION=$VERSION_ID
    elif [ -f /etc/redhat-release ]; then
        DISTRO="rhel"
    elif [ -f /etc/debian_version ]; then
        DISTRO="debian"
    else
        DISTRO="unknown"
    fi
    
    log_info "Detected distribution: $DISTRO"
}

# Check if running as root
check_root() {
    if [ "$EUID" -eq 0 ]; then
        log_warning "Running as root. This is not recommended for normal usage."
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# Install dependencies based on distribution
install_dependencies() {
    log_info "Installing clipboard utilities and dependencies..."
    
    case $DISTRO in
        ubuntu|debian|linuxmint|pop|elementary)
            log_info "Installing packages for Debian/Ubuntu-based system..."
            sudo apt update
            sudo apt install -y xclip xsel wl-clipboard
            ;;
        fedora|rhel|centos|rocky|almalinux)
            log_info "Installing packages for Red Hat-based system..."
            if command -v dnf >/dev/null 2>&1; then
                sudo dnf install -y xclip xsel wl-clipboard
            else
                sudo yum install -y xclip xsel wl-clipboard
            fi
            ;;
        arch|manjaro|endeavouros)
            log_info "Installing packages for Arch-based system..."
            sudo pacman -S --noconfirm xclip xsel wl-clipboard
            ;;
        opensuse*|sles)
            log_info "Installing packages for openSUSE-based system..."
            sudo zypper install -y xclip xsel wl-clipboard
            ;;
        alpine)
            log_info "Installing packages for Alpine Linux..."
            sudo apk add xclip xsel wl-clipboard
            ;;
        void)
            log_info "Installing packages for Void Linux..."
            sudo xbps-install -S xclip xsel wl-clipboard
            ;;
        gentoo)
            log_info "Installing packages for Gentoo..."
            sudo emerge x11-misc/xclip x11-misc/xsel gui-apps/wl-clipboard
            ;;
        *)
            log_warning "Unknown distribution: $DISTRO"
            log_info "Please install clipboard utilities manually:"
            echo "  • xclip"
            echo "  • xsel" 
            echo "  • wl-clipboard (for Wayland)"
            read -p "Continue with installation anyway? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                exit 1
            fi
            ;;
    esac
}

# Check if dependencies are available
check_dependencies() {
    log_info "Checking clipboard utilities..."
    
    local missing_deps=()
    local available_deps=()
    
    # Check for clipboard utilities
    for tool in xclip xsel wl-clipboard; do
        if command -v $tool >/dev/null 2>&1; then
            available_deps+=($tool)
            log_success "$tool is available"
        else
            missing_deps+=($tool)
        fi
    done
    
    if [ ${#available_deps[@]} -eq 0 ]; then
        log_error "No clipboard utilities found!"
        log_info "Attempting to install dependencies..."
        install_dependencies
        
        # Re-check after installation
        for tool in xclip xsel wl-clipboard; do
            if command -v $tool >/dev/null 2>&1; then
                available_deps+=($tool)
                log_success "$tool installed successfully"
            fi
        done
        
        if [ ${#available_deps[@]} -eq 0 ]; then
            log_error "Failed to install clipboard utilities. Please install manually."
            exit 1
        fi
    else
        log_success "Found ${#available_deps[@]} clipboard utilities: ${available_deps[*]}"
    fi
}

# Check display environment
check_display() {
    log_info "Checking display environment..."
    
    if [ -n "$DISPLAY" ]; then
        log_success "X11 display detected: $DISPLAY"
    elif [ -n "$WAYLAND_DISPLAY" ]; then
        log_success "Wayland display detected: $WAYLAND_DISPLAY"
    else
        log_warning "No graphical display detected"
        log_info "This might be a headless system or SSH session"
        log_info "The clipboard manager will still work in daemon-passive mode"
    fi
}

# Detect architecture and select binary
select_binary() {
    log_info "Detecting system architecture..."
    
    ARCH=$(uname -m)
    BINARY=""
    
    case $ARCH in
        x86_64)
            if [ -f "clipboard-manager-linux-amd64" ]; then
                BINARY="clipboard-manager-linux-amd64"
                log_success "Selected x86_64 binary"
            fi ;;
        aarch64|arm64)
            if [ -f "clipboard-manager-linux-arm64" ]; then
                BINARY="clipboard-manager-linux-arm64"
                log_success "Selected ARM64 binary"
            elif [ -f "clipboard-manager-linux-amd64" ]; then
                log_warning "ARM64 binary not available, using x86_64 (may not work)"
                BINARY="clipboard-manager-linux-amd64"
            fi ;;
        i386|i686)
            if [ -f "clipboard-manager-linux-386" ]; then
                BINARY="clipboard-manager-linux-386"
                log_success "Selected 32-bit binary"
            elif [ -f "clipboard-manager-linux-amd64" ]; then
                log_warning "32-bit binary not available, using x86_64"
                BINARY="clipboard-manager-linux-amd64"
            fi ;;
        *)
            log_warning "Unknown architecture: $ARCH, trying x86_64 binary"
            if [ -f "clipboard-manager-linux-amd64" ]; then
                BINARY="clipboard-manager-linux-amd64"
            fi ;;
    esac
    
    if [ -z "$BINARY" ]; then
        log_error "No compatible binary found for architecture: $ARCH"
        echo "Available binaries:"
        ls -1 clipboard-manager-linux-* 2>/dev/null || echo "  None found"
        exit 1
    fi
}

# Install the binary
install_binary() {
    log_info "Installing clipboard manager binary..."
    
    # Create directories if they don't exist
    sudo mkdir -p /usr/local/bin
    
    # Install the binary
    sudo cp "$BINARY" /usr/local/bin/clipboard-manager
    sudo chmod +x /usr/local/bin/clipboard-manager
    
    log_success "Binary installed to /usr/local/bin/clipboard-manager"
}

# Setup desktop integration
setup_desktop_integration() {
    log_info "Setting up desktop integration..."
    
    # Create applications directory
    mkdir -p ~/.local/share/applications
    
    # Create desktop entry
    cat > ~/.local/share/applications/clipboard-manager.desktop << EOF
[Desktop Entry]
Name=Clipboard Manager
Comment=Clipboard history manager for Linux
Exec=/usr/local/bin/clipboard-manager show
Icon=edit-copy
Terminal=false
Type=Application
Categories=Utility;
Keywords=clipboard;history;copy;paste;
EOF
    
    log_success "Desktop entry created"
}

# Setup autostart
setup_autostart() {
    log_info "Setting up autostart..."
    
    # Create autostart directory
    mkdir -p ~/.config/autostart
    
    # Create autostart entry
    cat > ~/.config/autostart/clipboard-manager.desktop << EOF
[Desktop Entry]
Name=Clipboard Manager
GenericName=Clipboard History Manager
Comment=Clipboard history manager with Ctrl+Shift+V hotkey
Exec=/usr/local/bin/clipboard-manager daemon
Icon=edit-copy
Terminal=false
Type=Application
Categories=Utility;System;Accessibility;
Keywords=clipboard;history;copy;paste;hotkey;
X-GNOME-Autostart-enabled=true
X-KDE-autostart-after=panel
X-MATE-Autostart-enabled=true
X-XFCE-Autostart-enabled=true
Hidden=false
NoDisplay=false
StartupNotify=false
X-GNOME-Autostart-Delay=3
X-KDE-StartupNotify=false
OnlyShowIn=GNOME;KDE;XFCE;MATE;Unity;Cinnamon;Pantheon;LXQt;LXDE;
EOF
    
    log_success "Autostart entry created"
}

# Test installation
test_installation() {
    log_info "Testing installation..."
    
    # Test if binary works
    if /usr/local/bin/clipboard-manager help >/dev/null 2>&1; then
        log_success "Binary is working correctly"
    else
        log_error "Binary test failed"
        log_info "Running diagnostics..."
        /usr/local/bin/clipboard-manager diagnose
        return 1
    fi
    
    # Test clipboard access
    if /usr/local/bin/clipboard-manager diagnose | grep -q "✓ Clipboard access working"; then
        log_success "Clipboard access is working"
    else
        log_warning "Clipboard access may have issues"
        log_info "Running full diagnostics..."
        /usr/local/bin/clipboard-manager diagnose
    fi
}

# Main installation function
main() {
    echo "🚀 Clipboard Manager Installer"
    echo "=============================="
    echo
    
    # Perform checks and installation
    check_root
    detect_distro
    check_display
    check_dependencies
    select_binary
    install_binary
    setup_desktop_integration
    setup_autostart
    
    echo
    log_info "Testing installation..."
    if test_installation; then
        echo
        log_success "Installation completed successfully!"
        echo
        echo "📋 Usage:"
        echo "   • Press Ctrl+Shift+V from anywhere to open clipboard history"
        echo "   • Clipboard manager will start automatically on login"
        echo "   • Run 'clipboard-manager help' for more options"
        echo "   • Run 'clipboard-manager diagnose' if you encounter issues"
        echo
        echo "🚀 The clipboard manager is now ready to use!"
        echo "   It will start automatically on your next login."
        echo "   To start it now, run: clipboard-manager daemon &"
    else
        echo
        log_error "Installation completed but there may be issues"
        echo "   Run 'clipboard-manager diagnose' for troubleshooting"
        exit 1
    fi
}

# Run main function
main "$@"