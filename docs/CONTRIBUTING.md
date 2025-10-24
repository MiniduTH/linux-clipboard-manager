# Contributing to Linux Clipboard Manager

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Development Setup

1. **Prerequisites:**
   - Go 1.21 or later
   - Linux with X11 or Wayland
   - Clipboard utilities: `xclip`, `xsel`, or `wl-clipboard`
   - GTK development libraries

2. **Clone and build:**
   ```bash
   git clone https://github.com/MiniduTH/linux-clipboard-manager.git
   cd linux-clipboard-manager
   go mod tidy
   go build -o clipboard-manager
   ```

## Testing

- Test on both Ubuntu and Fedora if possible
- Verify GUI and terminal modes work
- Test clipboard monitoring functionality
- Run the provided test script: `./test.sh`

## Code Style

- Follow standard Go formatting (`go fmt`)
- Add comments for exported functions
- Keep functions focused and small
- Handle errors appropriately

## Submitting Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Test thoroughly
5. Commit with descriptive messages
6. Push to your fork
7. Create a pull request

## Reporting Issues

When reporting bugs, please include:
- Linux distribution and version
- Go version
- Desktop environment (GNOME, KDE, etc.)
- Steps to reproduce
- Expected vs actual behavior
- Any error messages

## Feature Requests

We welcome feature requests! Please:
- Check existing issues first
- Describe the use case
- Explain why it would be valuable
- Consider implementation complexity

## Areas for Contribution

- Support for additional clipboard utilities
- Performance optimizations
- UI/UX improvements
- Documentation improvements
- Testing on different distributions
- Packaging for various package managers

Thank you for contributing!