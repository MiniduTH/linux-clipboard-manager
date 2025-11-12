# Edit Clipboard Items Feature

## Overview
Added the ability to edit text content in clipboard history items directly from the GUI.

## What's New

### User-Facing Changes
- **Edit Button**: Text items now have a pencil icon button next to the delete button
- **Edit Dialog**: Click the edit button to open a dialog with a multi-line text editor
- **Save/Cancel**: Edit dialog has Save and Cancel buttons
- **Validation**: Empty content is rejected with an error message
- **Auto-Refresh**: UI automatically refreshes after successful edit

### Technical Implementation

#### New Functions
1. **database.go**
   - `updateClipboardItem()`: Updates content and timestamp in SQLite database

2. **history.go**
   - `editHistoryItem()`: Validates and updates text items in history
   - Includes UTF-8 validation and empty content checks

3. **ui.go**
   - `showEditDialog()`: Modal dialog for editing text content
   - `onEdit` handler: Converts UI index to history index and shows dialog

4. **history_list_item.go**
   - Added `editButton` widget for text items
   - Added `editHovered` state tracking
   - Updated hover detection to handle both edit and delete buttons
   - Edit button only appears for text items (not images)

## Usage

### For Users
1. Open clipboard history with `clipboard-manager show` or press Super+V
2. Find a text item you want to edit
3. Click the pencil icon (edit button) on the right side
4. Modify the text in the dialog
5. Click "Save" to apply changes or "Cancel" to discard

### For Developers
```go
// Edit an item by history index
err := editHistoryItem(index, newContent)
if err != nil {
    // Handle error (invalid index, empty content, etc.)
}
```

## Limitations
- Only text items can be edited (images cannot be edited)
- Content cannot be empty after editing
- Content must be valid UTF-8

## Testing
Build and run the application:
```bash
go build -o clipboard-manager
./clipboard-manager show
```

Then:
1. Copy some text to create history items
2. Open the GUI
3. Click the edit button on any text item
4. Modify the content and save
5. Verify the item is updated in the list

## Files Modified
- `database.go` - Added `updateClipboardItem()` function
- `history.go` - Added `editHistoryItem()` function
- `ui.go` - Added `showEditDialog()` and `onEdit` handler
- `history_list_item.go` - Added edit button and hover states
- `README.md` - Updated feature list and usage documentation
