# Release Notes - Version 3.6.0

## 🎉 Edit Clipboard Items Feature

This release introduces the ability to edit text content directly in your clipboard history!

### ✨ New Features

#### Edit Text Items
- **Edit Button**: Text items now display a pencil icon button for quick editing
- **Multi-line Editor**: Full-featured text editor with word wrapping
- **Modal Dialog**: Clean, focused editing experience with Save/Cancel options
- **Real-time Updates**: Changes are immediately reflected in your history
- **Smart Validation**: Prevents empty content and ensures UTF-8 encoding

#### Enhanced UI
- **Dual Action Buttons**: Edit and Delete buttons side-by-side for text items
- **Improved Hover States**: Better visual feedback when hovering over buttons
- **Responsive Design**: Edit dialog adapts to content size with scrolling
- **Error Handling**: Clear error messages for validation failures

### 🔧 Technical Improvements

#### Database Layer
- New `updateClipboardItem()` function for atomic updates
- Timestamp updates on edit to maintain chronological order
- Efficient SQLite queries for content updates

#### History Management
- `editHistoryItem()` function with comprehensive validation
- UTF-8 string validation
- Empty content prevention
- Automatic history refresh after edits

#### User Interface
- Modal edit dialog with 600x400 pixel default size
- Scrollable content area for large text
- Keyboard navigation support
- Auto-refresh on successful save

### 📖 Documentation

#### New Documentation
- **EDIT_FEATURE.md**: Technical overview of the edit feature
- **docs/EDIT_FEATURE_GUIDE.md**: Comprehensive user guide
- **Updated README.md**: Feature list and usage instructions

### 🎯 Use Cases

Perfect for:
- Fixing typos in copied text
- Adjusting formatting before pasting
- Cleaning up copied content
- Modifying text without re-copying
- Quick text transformations

### 🚀 How to Use

1. Open clipboard history with `Super+V` or `clipboard-manager show`
2. Find a text item you want to edit
3. Click the pencil icon (edit button)
4. Modify the text in the dialog
5. Click "Save" to apply or "Cancel" to discard

### 📋 What's Changed

**Modified Files:**
- `database.go` - Added update functionality
- `history.go` - Added edit validation and logic
- `ui.go` - Added edit dialog and handlers
- `history_list_item.go` - Added edit button UI
- `README.md` - Updated documentation

**New Files:**
- `EDIT_FEATURE.md` - Technical documentation
- `docs/EDIT_FEATURE_GUIDE.md` - User guide

### ⚠️ Limitations

- Only text items can be edited (images are view-only)
- Content cannot be empty after editing
- Content must be valid UTF-8 text

### 🔄 Upgrade Notes

This release is fully backward compatible. Your existing clipboard history will work seamlessly with the new edit feature.

**Installation:**
```bash
# Download the release
wget https://github.com/MiniduTH/linux-clipboard-manager/releases/download/v3.6.0/clipboard-manager-v3.6.0-linux-amd64.tar.gz

# Extract
tar -xzf clipboard-manager-v3.6.0-linux-amd64.tar.gz

# Install
chmod +x install.sh
./install.sh
```

### 🐛 Bug Fixes

- Improved button hover detection with better padding
- Fixed button importance states during hover
- Enhanced error dialog handling

### 🎨 UI/UX Improvements

- Better visual hierarchy with edit/delete buttons
- Consistent button sizing (28x28 pixels)
- Improved accessibility with larger touch targets
- Smoother hover transitions

### 📊 Performance

- Minimal overhead for edit functionality
- Efficient database updates
- Fast UI refresh after edits
- No impact on clipboard monitoring

### 🔐 Security

- Input validation prevents malformed data
- UTF-8 encoding validation
- Safe database transactions
- No injection vulnerabilities

### 🙏 Acknowledgments

Thanks to all users who requested this feature! Your feedback helps make this clipboard manager better.

### 📝 Full Changelog

**Features:**
- Add edit button for text items in clipboard history
- Implement modal edit dialog with multi-line editor
- Add database update function for content modification
- Include validation for empty content and UTF-8 encoding

**Improvements:**
- Enhanced UI with dual action buttons
- Better hover states and visual feedback
- Improved error handling and user messages
- Auto-refresh after successful edits

**Documentation:**
- Added comprehensive edit feature guide
- Updated README with edit instructions
- Created technical documentation

**Technical:**
- New `updateClipboardItem()` database function
- New `editHistoryItem()` validation function
- New `showEditDialog()` UI function
- Enhanced `HistoryListItem` widget with edit support

---

## Previous Releases

### v3.5.0 - Auto-Startup on Installation
- Automatic startup application setup
- Enhanced installation script
- Improved system integration

### v3.4.0 - Smart Installer
- Automatic dependency management
- Improved error handling
- Better user experience

### v3.3.0 - Enhanced Installation
- Clipboard utility detection
- Improved error messages
- Better compatibility

---

**Download:** [clipboard-manager-v3.6.0-linux-amd64.tar.gz](https://github.com/MiniduTH/linux-clipboard-manager/releases/download/v3.6.0/clipboard-manager-v3.6.0-linux-amd64.tar.gz)

**Full Changelog:** https://github.com/MiniduTH/linux-clipboard-manager/compare/v3.5.0...v3.6.0
