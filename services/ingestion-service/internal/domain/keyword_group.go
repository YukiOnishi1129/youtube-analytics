package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/google/uuid"
)

var (
	ErrEmptyGroupName    = errors.New("keyword group name cannot be empty")
	ErrNoKeywordItems    = errors.New("keyword group must have at least one keyword item")
	ErrDuplicateKeyword  = errors.New("duplicate keyword in group")
)

// KeywordGroup represents a group of related keywords
// This is the aggregate root for keyword management
type KeywordGroup struct {
	ID          valueobject.UUID
	GenreID     valueobject.UUID
	Name        string
	FilterType  valueobject.FilterType
	TargetField string
	Enabled     bool
	Description *string
	Items       []KeywordItem // Aggregate includes items
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

// KeywordItem represents an individual keyword within a group
type KeywordItem struct {
	ID        valueobject.UUID
	GroupID   valueobject.UUID
	Keyword   string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// NewKeywordGroup creates a new keyword group
func NewKeywordGroup(
	id valueobject.UUID,
	genreID valueobject.UUID,
	name string,
	filterType valueobject.FilterType,
	targetField string,
	description *string,
	keywords []string,
) (*KeywordGroup, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyGroupName
	}

	if len(keywords) == 0 {
		return nil, ErrNoKeywordItems
	}

	if !filterType.IsValid() {
		return nil, ErrInvalidFilterType
	}

	if targetField == "" {
		targetField = "title" // Default to title
	}

	// Create keyword items
	items, err := createKeywordItems(id, keywords)
	if err != nil {
		return nil, err
	}

	return &KeywordGroup{
		ID:          id,
		GenreID:     genreID,
		Name:        name,
		FilterType:  filterType,
		TargetField: targetField,
		Enabled:     true,
		Description: description,
		Items:       items,
		CreatedAt:   time.Now(),
	}, nil
}

// createKeywordItems creates keyword items from a list of keywords
func createKeywordItems(groupID valueobject.UUID, keywords []string) ([]KeywordItem, error) {
	seen := make(map[string]bool)
	items := make([]KeywordItem, 0, len(keywords))

	for _, keyword := range keywords {
		trimmed := strings.TrimSpace(keyword)
		if trimmed == "" {
			continue
		}

		// Check for duplicates
		if seen[trimmed] {
			return nil, ErrDuplicateKeyword
		}
		seen[trimmed] = true

		itemID := valueobject.UUID(uuid.New().String())

		items = append(items, KeywordItem{
			ID:        itemID,
			GroupID:   groupID,
			Keyword:   trimmed,
			CreatedAt: time.Now(),
		})
	}

	if len(items) == 0 {
		return nil, ErrNoKeywordItems
	}

	return items, nil
}

// Update updates the keyword group
func (kg *KeywordGroup) Update(name *string, filterType *valueobject.FilterType, description *string) error {
	if name != nil {
		if strings.TrimSpace(*name) == "" {
			return ErrEmptyGroupName
		}
		kg.Name = *name
	}

	if filterType != nil {
		if !filterType.IsValid() {
			return ErrInvalidFilterType
		}
		kg.FilterType = *filterType
	}

	if description != nil {
		kg.Description = description
	}

	now := time.Now()
	kg.UpdatedAt = &now

	return nil
}

// UpdateKeywords updates the keywords in the group
func (kg *KeywordGroup) UpdateKeywords(keywords []string) error {
	items, err := createKeywordItems(kg.ID, keywords)
	if err != nil {
		return err
	}

	kg.Items = items
	now := time.Now()
	kg.UpdatedAt = &now

	return nil
}

// AddKeyword adds a new keyword to the group
func (kg *KeywordGroup) AddKeyword(keyword string) error {
	trimmed := strings.TrimSpace(keyword)
	if trimmed == "" {
		return errors.New("keyword cannot be empty")
	}

	// Check for duplicates
	for _, item := range kg.Items {
		if item.Keyword == trimmed {
			return ErrDuplicateKeyword
		}
	}

	itemID := valueobject.UUID(uuid.New().String())

	kg.Items = append(kg.Items, KeywordItem{
		ID:        itemID,
		GroupID:   kg.ID,
		Keyword:   trimmed,
		CreatedAt: time.Now(),
	})

	now := time.Now()
	kg.UpdatedAt = &now

	return nil
}

// RemoveKeyword removes a keyword from the group
func (kg *KeywordGroup) RemoveKeyword(keyword string) error {
	trimmed := strings.TrimSpace(keyword)
	newItems := make([]KeywordItem, 0, len(kg.Items)-1)

	found := false
	for _, item := range kg.Items {
		if item.Keyword == trimmed {
			found = true
			continue
		}
		newItems = append(newItems, item)
	}

	if !found {
		return errors.New("keyword not found in group")
	}

	if len(newItems) == 0 {
		return ErrNoKeywordItems
	}

	kg.Items = newItems
	now := time.Now()
	kg.UpdatedAt = &now

	return nil
}

// GetKeywords returns all keywords as a slice of strings
func (kg *KeywordGroup) GetKeywords() []string {
	keywords := make([]string, 0, len(kg.Items))
	for _, item := range kg.Items {
		keywords = append(keywords, item.Keyword)
	}
	return keywords
}

// Enable enables the keyword group
func (kg *KeywordGroup) Enable() {
	kg.Enabled = true
	now := time.Now()
	kg.UpdatedAt = &now
}

// Disable disables the keyword group
func (kg *KeywordGroup) Disable() {
	kg.Enabled = false
	now := time.Now()
	kg.UpdatedAt = &now
}

// Delete performs soft delete
func (kg *KeywordGroup) Delete() {
	now := time.Now()
	kg.DeletedAt = &now
	kg.UpdatedAt = &now
}

// IsDeleted checks if the keyword group is deleted
func (kg *KeywordGroup) IsDeleted() bool {
	return kg.DeletedAt != nil
}