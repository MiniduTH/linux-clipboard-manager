# Linux Clipboard Manager - AI Agent Instructions

## Project Overview
A Linux clipboard history manager written in Go that monitors clipboard content (text and images) and provides GUI/terminal access via global hotkey (Ctrl+Shift+V). Supports both X11 and Wayland, stores history in SQLite, and runs as a background daemon with system tray integration.

## Architecture & Core Components

### Multi-Mode Daemon Architecture
The application operates in multiple modes controlled via CLI arguments in `main.go`:
- **daemon**: Background monitoring with hotkeys
- **daemon-text-only**: Text-only monitoring (no images)
- **daemon-minimal**: Ultra-conservative polling for stability
- **daemon-passive**: No auto-monitoring (manual capture only)
- **show**: GUI popup (auto-starts daemon if needed)
- **tray**: System tray integration

**Key Pattern**: `ensureDaemonRunning()` checks for existing daemons before starting new instances to prevent duplicates. PID tracking via `~/.local/share/clipboard-manager/daemon.pid`.

### Data Flow
1. **Clipboard Monitoring** (`watchClipboard()` in `main.go`): Polls clipboard every 300ms
2. **Duplicate Detection**: Content hashing prevents duplicate entries
3. **Storage**: SQLite database at `~/.local/share/clipboard-manager/history.db` (max 50 items)
4. **UI Access**: Fyne-based GUI or terminal mode, triggered by hotkey or CLI

### File Responsibilities
- `main.go`: CLI routing, environment checks, daemon lifecycle, clipboard polling loops
- `daemon.go`: Process management (PID tracking, status checks, start/stop)
- `database.go`: SQLite operations (CRUD, migrations, duplicate handling)
- `history.go`: In-memory history management with mutex protection (`historyMu`)
- `ui.go`: Fyne GUI with scrollable list, edit/delete actions
- `history_list_item.go`: Custom Fyne widget for history items with hover states
- `hotkey.go`: GNOME/KDE hotkey setup via gsettings/kwriteconfig5
- `image_clipboard.go`: Image detection/restoration using xclip/wl-paste
- `tray.go`: System tray menu (start/stop/show/quit)

## Development Workflows

### Building & Testing
```bash
make build           # Standard build (creates ./clipboard-manager)
make install         # Build + system install to /usr/local/bin
make release         # Multi-arch binaries (amd64, arm64, 386) in build/release/
make test            # Run all tests via go test
make clean           # Remove build artifacts
```

**Testing Pattern**: Tests use `setupTestDB()` and `teardownTestDB()` helpers (see `history_test.go`) to create isolated SQLite databases. Always wrap history operations with `historyMu.Lock()` for thread safety.

### Installation Flow
The `Makefile` generates a comprehensive `install.sh` that:
1. Detects architecture and selects appropriate binary
2. Installs to `/usr/local/bin/clipboard-manager`
3. Creates `.desktop` files for launcher and autostart
4. Configures GNOME/KDE hotkeys automatically
5. Sets up `~/.config/autostart/clipboard-manager.desktop` for auto-launch

## Critical Conventions

### Thread Safety
- All `history` array access MUST use `historyMu.Lock()/Unlock()` or `RLock()/RUnlock()`
- Database operations don't need separate locking (SQLite handles concurrency)
- Example: See `addToHistory()`, `removeHistoryItem()`, `editHistoryItem()` in `history.go`

### Clipboard Content Filtering
System noise detection in `isSystemNoise()` (main.go) filters:
- Single characters
- Mouse selection artifacts ("1", "0")
- Whitespace-only content
- Invalid UTF-8 strings

### Database Schema
```sql
CREATE TABLE clipboard_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL CHECK(type IN ('text', 'image')),
    content TEXT NOT NULL,  -- Plain text or base64-encoded image
    timestamp DATETIME NOT NULL,
    image_format TEXT,      -- png, jpeg, etc.
    image_width INTEGER,
    image_height INTEGER,
    image_size INTEGER
);
```

**Migration Path**: Old JSON format auto-migrates to SQLite on first run (see `loadHistory()` in main.go).

### UI Refresh Pattern
When modifying history from GUI:
1. Update database via `updateClipboardItem()` or delete
2. Call `refreshHistoryFromDB()` to sync in-memory state
3. Call `refreshUI(w)` to rebuild Fyne widget tree
Example: `showEditDialog()` in `ui.go` demonstrates this flow.

### Image Handling
- **Storage**: Base64-encoded in database `content` field
- **Detection**: Tries xclip TARGETS, then wl-paste --list-types
- **Restoration**: Uses `restoreImageToSystemClipboard()` which writes temp files
- **Formats**: PNG, JPEG/JPG, GIF, BMP (priority: PNG > JPEG)

## Environment Requirements

### Runtime Dependencies
- X11 (`DISPLAY` env) or Wayland (`WAYLAND_DISPLAY` env)
- Clipboard utilities: xclip (X11) or wl-clipboard (Wayland)
- GTK3 libraries for Fyne GUI

**Environment Check**: `checkEnvironment()` validates display and clipboard access. Use `clipboard-manager diagnose` for troubleshooting.

### Build Dependencies
- Go 1.24+ (due to toolchain directive in `go.mod`)
- CGO_ENABLED=1 for SQLite and Fyne
- GTK3 dev headers: `libgtk-3-dev`, `libayatana-appindicator3-dev`
- For cross-compilation: `gcc-aarch64-linux-gnu`, `gcc-multilib`

## Common Patterns

### Adding New Clipboard Item Types
1. Add type constant to `ClipboardItemType` in `history.go`
2. Extend database schema in `createTables()` (database.go)
3. Update `saveClipboardItem()` for serialization
4. Modify `watchClipboard()` polling logic to detect new type
5. Update `NewHistoryListItem()` widget for display

### Adding CLI Commands
1. Add argument check in `main.go` (e.g., `if os.Args[1] == "newcmd"`)
2. Implement handler function
3. Update help text in the help command block
4. Add to `scripts/manage-clipboard.sh` if needed for automation

### Debugging Clipboard Issues
- Use `daemon-text-only` mode to isolate image monitoring problems
- Check `/proc/<pid>/cmdline` to identify daemon processes (see `isDaemonRunning()`)
- Test clipboard access with `clipboard-manager capture`
- Use `diagnose` command for environment validation

## Testing Guidelines

### Test Structure
- Use `testing.T` for assertions
- Setup: `setupTestDB(t)` creates isolated DB
- Teardown: `defer teardownTestDB(t)` cleans up
- Helpers: `addTestItems()`, `clearTestHistory()`, `getTestHistoryItem()`

### Running Tests
```bash
go test -v                    # All tests with verbose output
go test -run TestEditItem     # Specific test
./tests/run_tests.sh          # Via test script
```

### UI Testing
`ui_integration_test.go` uses Fyne test utilities. Mock clipboard operations with `test_helpers.go` functions.

## Key Files for Features

- **Edit feature**: `ui.go` (showEditDialog), `history.go` (editHistoryItem), `database.go` (updateClipboardItem)
- **Hotkey setup**: `hotkey.go` (setupGnomeHotkey, setupKDEHotkey)
- **Autostart**: `main.go` (ensureStartupEnabled, createAutostartEntry)
- **System tray**: `tray.go` (runWithSystemTray, system menu creation)
- **Image support**: `image_clipboard.go` (detectImageInClipboard, tryXclipImage)

## Version & Release Management
- Version tracked in `VERSION` file (currently v3.7.0)
- Release notes in `docs/releases/RELEASE_NOTES_v*.md`
- Use `make release` to generate multi-arch binaries with install script
- Package as `clipboard-manager-release.tar.gz` for GitHub releases
