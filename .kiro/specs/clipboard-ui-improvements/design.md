# Design Document

## Overview

The clipboard UI improvements will enhance the existing Fyne-based interface by implementing proper text wrapping, fixing the clear history functionality, adding individual item deletion with X buttons, and improving hover states for better visibility and user experience.

## Architecture

The design maintains the existing single-window architecture using Fyne widgets but introduces:

- Custom list item widgets with embedded delete buttons
- Enhanced text rendering with wrapping capabilities  
- Improved event handling for individual item deletion
- Better styling and theming for hover states

## Components and Interfaces

### Enhanced List Item Widget

```go
type HistoryListItem struct {
    widget.BaseWidget
    text string
    index int
    onDelete func(int)
    onSelect func(int)
}
```

The custom list item will contain:
- Multi-line text label with wrapping
- Delete button (X icon) positioned on the right
- Click handlers for both selection and deletion
- Proper hover state styling

### Updated History Management

The existing `history.go` functions will be extended with:
- `removeHistoryItem(index int)` - Remove specific item by index
- Enhanced `clearHistory()` with proper UI callback
- Thread-safe operations for concurrent UI updates

### UI Layout Structure

```
Window
├── Header (title + item count)
├── Separator
├── ScrollContainer
│   └── VBox of HistoryListItems
│       ├── HBox per item
│       │   ├── MultiLine Text (flex: 1)
│       │   └── Delete Button (fixed width)
├── Separator  
└── Button Container (Clear All + Close)
```

## Data Models

### HistoryItem Display Model

```go
type DisplayItem struct {
    Index int
    Text string
    DisplayText string
    IsWrapped bool
}
```

This model separates the original clipboard text from the display representation, allowing for proper truncation and wrapping without losing the original content.

## Error Handling

### Delete Operations
- Validate index bounds before deletion
- Handle concurrent access during deletion
- Graceful fallback if deletion fails
- UI feedback for failed operations

### Clear History Operations  
- Confirm operation success before UI update
- Handle file system errors during save
- Maintain UI consistency on failure
- Provide user feedback for errors

### UI Rendering
- Handle empty states gracefully
- Fallback styling if custom themes fail
- Proper cleanup of event handlers
- Memory management for large history lists

## Testing Strategy

### Unit Tests
- History management functions (add, remove, clear)
- Text wrapping and truncation logic
- Index validation and bounds checking

### Integration Tests
- UI component interaction (click, hover, delete)
- File persistence after operations
- Concurrent access scenarios

### Manual Testing
- Visual verification of text wrapping
- Hover state visibility across different themes
- Delete button positioning and sizing
- Performance with maximum history items (50)

## Implementation Details

### Text Wrapping Implementation
- Use Fyne's `widget.RichText` for multi-line display
- Set maximum height to prevent excessive vertical space
- Implement smart truncation at word boundaries
- Preserve original formatting for clipboard operations

### Delete Button Styling
- Use Fyne's `widget.Button` with custom icon
- Position using `container.NewHBox` with proper spacing
- Apply consistent sizing (20x20 pixels)
- Implement hover effects with color changes

### Hover State Enhancement
- Override default Fyne hover colors
- Use higher contrast colors for better visibility
- Apply consistent styling across all interactive elements
- Ensure accessibility compliance for color contrast

### Performance Considerations
- Lazy loading for large history lists
- Efficient re-rendering on item deletion
- Minimal memory allocation during UI updates
- Proper cleanup of event handlers and resources