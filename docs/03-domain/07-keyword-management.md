# Keyword Management Design

## Overview

The keyword management system manages keywords used for filtering YouTube videos. Keywords are organized into groups, and each group can contain multiple keywords.

## Domain Model

### Aggregate Design

**KeywordGroup (Aggregate Root)**
```
KeywordGroup (Aggregate Root)
├── KeywordGroup Entity
└── KeywordItem Value Objects[]
```

### Entity: KeywordGroup

KeywordGroup is the aggregate root that manages related keywords together.

**Properties:**
- `ID`: UUID - Unique identifier for the group
- `GenreID`: UUID - ID of the related genre
- `Name`: string - Group name (e.g., JavaScript/JS, Python)
- `FilterType`: FilterType - Filter type (include/exclude)
- `TargetField`: string - Search target field (title/description)
- `Enabled`: bool - Enabled/disabled flag
- `Description`: *string - Group description (optional)
- `Items`: []KeywordItem - Keywords belonging to the group
- `CreatedAt`: time.Time
- `UpdatedAt`: *time.Time
- `DeletedAt`: *time.Time

**Business Rules:**
1. Group name is required and cannot be empty
2. At least one keyword item is required
3. No duplicate keywords allowed within the same group
4. Default TargetField is "title"

### Value Object: KeywordItem

Represents individual keywords as value objects.

**Properties:**
- `ID`: UUID - Unique identifier for the item
- `GroupID`: UUID - ID of the parent group
- `Keyword`: string - The actual keyword
- `CreatedAt`: time.Time
- `UpdatedAt`: *time.Time

## Database Schema

### keyword_groups table
```sql
CREATE TABLE keyword_groups (
  id           uuid PRIMARY KEY,
  genre_id     uuid NOT NULL REFERENCES genres(id),
  name         text NOT NULL,
  filter_type  text NOT NULL CHECK (filter_type IN ('include', 'exclude')),
  target_field varchar(20) NOT NULL DEFAULT 'title',
  enabled      boolean DEFAULT true,
  description  text,
  created_at   timestamptz DEFAULT now(),
  updated_at   timestamptz,
  deleted_at   timestamptz
);
```

### keyword_items table
```sql
CREATE TABLE keyword_items (
  id                uuid PRIMARY KEY,
  keyword_group_id  uuid NOT NULL REFERENCES keyword_groups(id) ON DELETE CASCADE,
  keyword           text NOT NULL,
  created_at        timestamptz DEFAULT now(),
  updated_at        timestamptz
);
```

## Pattern Generation

Regular expression patterns are not stored in the database but generated dynamically in the application layer.

### Pattern Generation Flow

1. Retrieve keyword_items from KeywordGroup
2. Generate pattern using KeywordPatternGenerator service
   - Automatically detect English and Japanese
   - Generate appropriate variations for each language
   - Combine as regular expression

### Example

Input Keywords:
```
["Next.js", "React", "ネクスト"]
```

Generated Pattern:
```
(?i)(Next\.js|Nextjs|React|ネクスト)
```

## Repository Interface

```go
type KeywordGroupRepository interface {
    Create(ctx context.Context, group *domain.KeywordGroup) error
    Update(ctx context.Context, group *domain.KeywordGroup) error
    UpdateWithItems(ctx context.Context, group *domain.KeywordGroup) error
    Delete(ctx context.Context, id valueobject.UUID) error
    FindByID(ctx context.Context, id valueobject.UUID) (*domain.KeywordGroup, error)
    FindByGenreID(ctx context.Context, genreID valueobject.UUID) ([]*domain.KeywordGroup, error)
    List(ctx context.Context, limit, offset int) ([]*domain.KeywordGroup, error)
    ListByEnabled(ctx context.Context, enabled bool, limit, offset int) ([]*domain.KeywordGroup, error)
}
```

## Use Cases

### 1. Create Keyword Group
```go
type CreateKeywordGroupInput struct {
    GenreID     uuid.UUID
    Name        string
    Keywords    []string
    FilterType  valueobject.FilterType
    TargetField string
    Description *string
}
```

### 2. Update Keywords in Group
```go
func UpdateKeywords(ctx context.Context, groupID uuid.UUID, keywords []string) (*domain.KeywordGroup, error)
```

### 3. Generate Pattern for Video Filtering
```go
// Used during video filtering
patterns, err := useCase.GeneratePatternsForGenre(ctx, genreID)
// Returns map[FilterType][]string with patterns for include/exclude
```

## Migration from Old Structure

Migration from old structure (keywords table) to new structure (keyword_groups + keyword_items):

1. Rename keywords table to keyword_groups
2. Create keyword_items table
3. Extract individual keywords from pattern field and insert into keyword_items
4. Remove pattern field from keyword_groups

## Benefits

1. **Data Integrity**: Adding/removing keywords immediately reflects in patterns
2. **Simplified Management**: Easy to manage individual keywords
3. **Extensibility**: Easy to modify pattern generation logic
4. **Visibility**: Display and edit individual keywords in admin UI

## Integration Points

### Video Filtering
When filtering videos by genre:
1. Load all enabled KeywordGroups for the genre
2. Generate patterns dynamically for each group
3. Apply include/exclude filters based on FilterType

### Admin Interface
The admin interface can:
1. Display keyword groups with their individual keywords
2. Add/remove keywords without dealing with regex
3. Preview generated patterns
4. Enable/disable groups

## Future Enhancements

1. **Keyword Suggestions**: Suggest related keywords based on existing ones
2. **Pattern Caching**: Cache generated patterns for performance
3. **Analytics**: Track which keywords match most videos
4. **Bulk Operations**: Import/export keywords in CSV format