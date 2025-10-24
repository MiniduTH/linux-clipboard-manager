# Build Directory

This directory contains build artifacts and temporary files generated during development and testing.

## Contents

After running `make build` or `make build-all`, this directory contains:

- **clipboard-manager**: Standard release binary
- **clipboard-manager-debug**: Debug build with symbols (larger, for debugging)
- **clipboard-manager-optimized**: Optimized build (smaller, stripped symbols)
- **Test binaries**: Test-specific builds (when present)
- **Temporary files**: Files created during build process

## Cleanup

Build artifacts are automatically cleaned by:

```bash
# Clean all build artifacts
make clean

# Or manually
rm -rf build/*
```

## Gitignore

This directory and its contents are typically ignored by git to keep the repository clean. Only the README.md file is tracked.

## Build Process

Multiple build options are available:

```bash
# Standard build (creates binary in root + copy in build/)
make build

# Build all variants (standard, debug, optimized)
make build-all

# Build directly to build directory
go build -o build/clipboard-manager

# Build specific variants manually
go build -gcflags="all=-N -l" -o build/clipboard-manager-debug      # Debug
go build -ldflags="-s -w" -o build/clipboard-manager-optimized      # Optimized
```

## Binary Variants

- **Standard**: Regular build with default Go settings
- **Debug**: Includes debugging symbols, larger size, easier to debug
- **Optimized**: Stripped symbols, smaller size, production-ready