# Admin Management Use Cases

## Overview

Administrative functions for managing the YouTube Analytics system configuration and operations.

## YouTube Category Management

### View YouTube Categories

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: None
- **Processing**:
  1. Retrieve all YouTube categories from the system
  2. Display with current settings (name, assignable status)
- **Output**: List of categories with their properties
- **UI Elements**: 
  - Table with sortable columns
  - Assignable status toggle
  - Edit button per row

### Update YouTube Category

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: Category ID, updated name, assignable flag
- **Processing**:
  1. Validate the category exists
  2. Update category properties
  3. Log the change for audit
- **Output**: Success/error message
- **Notes**: 
  - Category IDs are fixed from YouTube API
  - Only name and assignable status can be modified

## Genre Management

### List Genres

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: Optional filters (enabled/disabled, region, language)
- **Processing**:
  1. Retrieve genres based on filters
  2. Include associated category names
  3. Show enabled/disabled status
- **Output**: List of genres with full details
- **UI Elements**:
  - Filter controls (region, language, status)
  - Table with genre details
  - Enable/disable toggle
  - Edit and delete buttons

### Create Genre

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: 
  - Code (e.g., "engineering_jp")
  - Name (e.g., "Engineering (JP)")
  - Language code (e.g., "ja")
  - Region code (e.g., "JP")
  - Category IDs (multiple selection from assignable categories)
- **Processing**:
  1. Validate unique code
  2. Validate region code exists in YouTube
  3. Validate all category IDs are assignable
  4. Create genre record
  5. Set as enabled by default
- **Output**: Created genre details or validation errors
- **UI Elements**:
  - Form with text inputs
  - Multi-select for categories (only shows assignable)
  - Region dropdown (from YouTube supported regions)
  - Language dropdown

### Update Genre

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: Genre ID, updated fields (name, enabled status, category IDs)
- **Processing**:
  1. Validate genre exists
  2. Validate category IDs if changed
  3. Update genre properties
  4. If disabling, warn about impact on collection
- **Output**: Updated genre or errors
- **Notes**: Code, language, and region cannot be changed after creation

### Delete Genre

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: Genre ID
- **Processing**:
  1. Check for associated videos
  2. If videos exist, require confirmation
  3. Soft delete or hard delete based on policy
  4. Remove associated keywords
- **Output**: Success/error message
- **UI Elements**:
  - Confirmation dialog
  - Warning about data loss

### Enable/Disable Genre

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: Genre ID, enabled flag
- **Processing**:
  1. Update genre enabled status
  2. Log the change
- **Output**: Updated status
- **Impact**: 
  - Disabled genres are skipped in batch collection
  - Existing videos remain accessible

## Keyword Management

### List Keyword Groups by Genre

- **Actor**: Administrator
- **Precondition**: Logged in with admin role, genre selected
- **Input**: Genre ID, optional filters (enabled/disabled, filter type)
- **Processing**:
  1. Retrieve keyword groups for the genre
  2. Load associated keyword items for each group
  3. Sort by filter type (include/exclude)
  4. Show enabled/disabled status
- **Output**: List of keyword groups with their individual keywords
- **UI Elements**:
  - Genre selector
  - Filter type tabs (All/Include/Exclude)
  - Keyword group cards showing individual keywords
  - Pattern preview (dynamically generated)

### Create Keyword Group

- **Actor**: Administrator
- **Precondition**: Logged in with admin role, genre selected
- **Input**:
  - Genre ID
  - Group name (e.g., "JavaScript/JS")
  - Keywords array (e.g., ["JavaScript", "JS", "React", "Vue"])
  - Filter type (include/exclude)
  - Target field (title/description)
  - Description (optional)
- **Processing**:
  1. Validate group name uniqueness in genre
  2. Validate at least one keyword provided
  3. Create keyword group with items
  4. Generate pattern dynamically for preview
  5. Set as enabled by default
- **Output**: Created keyword group or validation errors
- **UI Elements**:
  - Form within genre context
  - Filter type radio buttons
  - Keywords input (tag-style UI)
  - Pattern preview (auto-generated)

### Update Keyword Group

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: Group ID, updated fields (name, filter type, description)
- **Processing**:
  1. Update group properties (excluding keywords)
  2. Log the change
- **Output**: Updated group or errors
- **Notes**: To update keywords, use separate keyword management endpoints

### Manage Keywords in Group

- **Actor**: Administrator
- **Precondition**: Logged in with admin role, keyword group selected
- **Input**: Group ID, keywords array
- **Processing**:
  1. Replace all keywords in the group
  2. Validate no duplicates
  3. Ensure at least one keyword remains
- **Output**: Updated keyword list
- **UI Elements**:
  - Add/remove keywords UI
  - Bulk edit option
  - Real-time pattern preview

### Delete Keyword Group

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: Group ID
- **Processing**:
  1. Soft delete group (set deleted_at)
  2. Cascade delete to keyword items
  3. Maintain history for audit
- **Output**: Success message
- **Notes**: Soft delete preserves historical data

### Test Generated Pattern

- **Actor**: Administrator
- **Precondition**: Creating or editing keyword group
- **Input**: Keywords array, test strings
- **Processing**:
  1. Generate pattern from keywords
  2. Test pattern against provided strings
  3. Show matches
  3. Show matches/non-matches
- **Output**: Test results with highlighted matches
- **UI Elements**:
  - Pattern input
  - Test string textarea (multiple lines)
  - Results display with match highlighting

## Batch Operations Management

### View Collection Status

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: Date range, genre filter
- **Processing**:
  1. Retrieve collection job history
  2. Show per-genre statistics
  3. Display success/failure counts
- **Output**: Collection statistics and history
- **UI Elements**:
  - Date range picker
  - Genre filter
  - Statistics cards (collected, adopted, failed)
  - Job history table

### Trigger Manual Collection

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: Selected genres (optional, defaults to all enabled)
- **Processing**:
  1. Validate at least one genre selected
  2. Queue collection job
  3. Show job progress
- **Output**: Job ID and status
- **UI Elements**:
  - Genre multi-select
  - Run button
  - Progress indicator
  - Job log viewer

### View WebSub Subscription Status

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: Channel filter, subscription status filter
- **Processing**:
  1. Retrieve channel subscription status
  2. Show lease expiration times
  3. Highlight expiring soon
- **Output**: Channel subscription details
- **UI Elements**:
  - Status filter (active/expired/all)
  - Channel search
  - Subscription table with expiry countdown
  - Bulk renewal button

### Manage Channel Subscriptions

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: Channel ID, action (subscribe/unsubscribe/renew)
- **Processing**:
  1. For subscribe: Send subscription request to YouTube WebSub hub
  2. For unsubscribe: Send unsubscription request
  3. For renew: Refresh subscription before expiry
  4. Update channel subscription status
  5. Log the action
- **Output**: Updated subscription status
- **UI Elements**:
  - Subscribe/Unsubscribe toggle per channel
  - Renew button for active subscriptions
  - Subscription history log
- **Notes**: 
  - WebSub subscriptions typically last 10 days
  - Auto-renewal handled by batch job

## System Monitoring

### View System Health

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: None
- **Processing**:
  1. Check service health endpoints
  2. Verify database connectivity
  3. Check external API status
  4. Review job queue status
- **Output**: System health dashboard
- **UI Elements**:
  - Service status cards
  - API quota usage
  - Database connection pool stats
  - Recent error logs

### View Audit Logs

- **Actor**: Administrator
- **Precondition**: Logged in with admin role
- **Input**: Date range, action type, actor
- **Processing**:
  1. Retrieve audit logs
  2. Filter by criteria
  3. Sort by timestamp
- **Output**: Audit trail of admin actions
- **UI Elements**:
  - Date range picker
  - Action type filter
  - Actor (admin user) filter
  - Detailed log table

## Navigation Flow

```
Admin Dashboard
├── YouTube Categories
│   └── Edit Category
├── Genres
│   ├── Create Genre
│   ├── Edit Genre
│   └── Keywords
│       ├── Create Keyword
│       └── Edit Keyword
├── Batch Operations
│   ├── Collection Status
│   └── Manual Trigger
├── Subscriptions
│   └── Channel Status
└── System
    ├── Health Dashboard
    └── Audit Logs
```

## Security Considerations

1. All operations require admin role authentication
2. All changes are logged with timestamp and actor
3. Sensitive operations require confirmation
4. Rate limiting on manual triggers
5. Input validation on all forms
6. CSRF protection on state-changing operations