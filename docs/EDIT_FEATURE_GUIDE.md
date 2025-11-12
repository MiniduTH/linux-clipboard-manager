# Edit Feature User Guide

## How to Edit Clipboard Items

### Step 1: Open Clipboard History
Press **Super+V** (Windows key + V) or run:
```bash
clipboard-manager show
```

### Step 2: Locate the Item
Find the text item you want to edit in the list. Each text item will have two buttons on the right:
- 📝 **Edit button** (pencil icon)
- ❌ **Delete button** (X icon)

### Step 3: Click Edit
Click the pencil icon to open the edit dialog.

### Step 4: Modify Content
- A dialog window will appear with the current text
- The text is in a multi-line editor with word wrapping
- Modify the text as needed
- You can add or remove lines, change words, etc.

### Step 5: Save or Cancel
- Click **Save** to apply your changes
- Click **Cancel** to discard changes
- If you try to save empty content, you'll get an error

### Step 6: Verify Changes
After saving, the clipboard history will refresh automatically and show your edited content.

## Features

### What You Can Edit
✅ Any text item in your clipboard history
✅ Multi-line text with proper formatting
✅ Text of any length

### What You Cannot Edit
❌ Image items (they don't have an edit button)
❌ Cannot make content empty (validation prevents this)

### Visual Indicators
- **Hover Effect**: Edit button highlights when you hover over it
- **Button Importance**: Edit button becomes more prominent on hover
- **Dialog Size**: 600x400 pixels with scrollable content area

## Tips

1. **Quick Access**: The edit button only appears for text items
2. **Keyboard Navigation**: Use Tab to navigate between Save and Cancel buttons
3. **Multi-line Editing**: The editor supports multiple lines with proper wrapping
4. **Validation**: Empty content is rejected to prevent data loss
5. **Auto-Refresh**: No need to close and reopen - the list updates automatically

## Troubleshooting

### Edit Button Not Visible
- **Cause**: Item is an image, not text
- **Solution**: Only text items can be edited

### Cannot Save Empty Content
- **Cause**: Validation prevents empty clipboard items
- **Solution**: Add some text content before saving

### Changes Not Appearing
- **Cause**: Database update failed
- **Solution**: Check console output for error messages

## Technical Details

### Database Updates
- Edited items get a new timestamp
- Original content is replaced with new content
- Changes are persisted to SQLite database immediately

### UI Behavior
- Edit dialog is modal (blocks interaction with main window)
- Main window refreshes after successful edit
- Error dialogs appear for validation failures

### Performance
- Editing is instant for small to medium text
- Large text (>10KB) may take a moment to load in editor
- Database updates are atomic and safe
