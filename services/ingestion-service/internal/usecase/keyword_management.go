package usecase

import (
	"context"
	"fmt"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/service"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/repository"
	"github.com/google/uuid"
)

// KeywordGroupManagementUseCase handles keyword group operations
type KeywordGroupManagementUseCase interface {
	CreateKeywordGroup(ctx context.Context, input CreateKeywordGroupInput) (*domain.KeywordGroup, error)
	UpdateKeywordGroup(ctx context.Context, groupID uuid.UUID, input UpdateKeywordGroupInput) (*domain.KeywordGroup, error)
	UpdateKeywords(ctx context.Context, groupID uuid.UUID, keywords []string) (*domain.KeywordGroup, error)
	DeleteKeywordGroup(ctx context.Context, groupID uuid.UUID) error
	GetKeywordGroup(ctx context.Context, groupID uuid.UUID) (*domain.KeywordGroup, error)
	ListKeywordGroupsByGenre(ctx context.Context, genreID uuid.UUID) ([]*domain.KeywordGroup, error)
	GeneratePatternForGroup(ctx context.Context, groupID uuid.UUID) (string, error)
}

type CreateKeywordGroupInput struct {
	GenreID     uuid.UUID
	Name        string
	Keywords    []string
	FilterType  valueobject.FilterType
	TargetField string
	Description *string
}

type UpdateKeywordGroupInput struct {
	Name        *string
	FilterType  *valueobject.FilterType
	Description *string
}

type keywordGroupManagementUseCase struct {
	groupRepo        repository.KeywordGroupRepository
	patternGenerator *service.KeywordPatternGenerator
}

// NewKeywordGroupManagementUseCase creates a new keyword group management use case
func NewKeywordGroupManagementUseCase(
	groupRepo repository.KeywordGroupRepository,
) KeywordGroupManagementUseCase {
	return &keywordGroupManagementUseCase{
		groupRepo:        groupRepo,
		patternGenerator: service.NewKeywordPatternGenerator(),
	}
}

// CreateKeywordGroup creates a new keyword group with keywords
func (u *keywordGroupManagementUseCase) CreateKeywordGroup(
	ctx context.Context,
	input CreateKeywordGroupInput,
) (*domain.KeywordGroup, error) {
	// Generate group ID
	groupID := valueobject.UUID(uuid.New().String())

	// Create keyword group domain object
	group, err := domain.NewKeywordGroup(
		groupID,
		valueobject.UUID(input.GenreID.String()),
		input.Name,
		input.FilterType,
		input.TargetField,
		input.Description,
		input.Keywords,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create keyword group: %w", err)
	}

	// Save to repository
	if err := u.groupRepo.Create(ctx, group); err != nil {
		return nil, fmt.Errorf("failed to save keyword group: %w", err)
	}

	return group, nil
}

// UpdateKeywordGroup updates a keyword group (excluding keywords)
func (u *keywordGroupManagementUseCase) UpdateKeywordGroup(
	ctx context.Context,
	groupID uuid.UUID,
	input UpdateKeywordGroupInput,
) (*domain.KeywordGroup, error) {
	// Find existing group
	group, err := u.groupRepo.FindByID(ctx, valueobject.UUID(groupID.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to find keyword group: %w", err)
	}

	// Update group
	if err := group.Update(input.Name, input.FilterType, input.Description); err != nil {
		return nil, fmt.Errorf("failed to update keyword group: %w", err)
	}

	// Save to repository
	if err := u.groupRepo.Update(ctx, group); err != nil {
		return nil, fmt.Errorf("failed to save updated keyword group: %w", err)
	}

	return group, nil
}

// UpdateKeywords updates the keywords in a group
func (u *keywordGroupManagementUseCase) UpdateKeywords(
	ctx context.Context,
	groupID uuid.UUID,
	keywords []string,
) (*domain.KeywordGroup, error) {
	// Find existing group
	group, err := u.groupRepo.FindByID(ctx, valueobject.UUID(groupID.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to find keyword group: %w", err)
	}

	// Update keywords
	if err := group.UpdateKeywords(keywords); err != nil {
		return nil, fmt.Errorf("failed to update keywords: %w", err)
	}

	// Save to repository with items
	if err := u.groupRepo.UpdateWithItems(ctx, group); err != nil {
		return nil, fmt.Errorf("failed to save updated keywords: %w", err)
	}

	return group, nil
}

// DeleteKeywordGroup soft deletes a keyword group
func (u *keywordGroupManagementUseCase) DeleteKeywordGroup(
	ctx context.Context,
	groupID uuid.UUID,
) error {
	return u.groupRepo.Delete(ctx, valueobject.UUID(groupID.String()))
}

// GetKeywordGroup retrieves a keyword group by ID
func (u *keywordGroupManagementUseCase) GetKeywordGroup(
	ctx context.Context,
	groupID uuid.UUID,
) (*domain.KeywordGroup, error) {
	return u.groupRepo.FindByID(ctx, valueobject.UUID(groupID.String()))
}

// ListKeywordGroupsByGenre lists all keyword groups for a genre
func (u *keywordGroupManagementUseCase) ListKeywordGroupsByGenre(
	ctx context.Context,
	genreID uuid.UUID,
) ([]*domain.KeywordGroup, error) {
	return u.groupRepo.FindByGenreID(ctx, valueobject.UUID(genreID.String()))
}

// GeneratePatternForGroup generates a regex pattern for a keyword group
func (u *keywordGroupManagementUseCase) GeneratePatternForGroup(
	ctx context.Context,
	groupID uuid.UUID,
) (string, error) {
	// Find group
	group, err := u.groupRepo.FindByID(ctx, valueobject.UUID(groupID.String()))
	if err != nil {
		return "", fmt.Errorf("failed to find keyword group: %w", err)
	}

	// Generate pattern from keywords
	keywords := group.GetKeywords()
	pattern := u.patternGenerator.GeneratePattern(keywords)

	return pattern, nil
}

// GeneratePatternsForGenre generates patterns for all enabled keyword groups in a genre
func (u *keywordGroupManagementUseCase) GeneratePatternsForGenre(
	ctx context.Context,
	genreID uuid.UUID,
) (map[valueobject.FilterType][]string, error) {
	// Get all groups for genre
	groups, err := u.groupRepo.FindByGenreID(ctx, valueobject.UUID(genreID.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to find keyword groups: %w", err)
	}

	// Generate patterns by filter type
	patterns := make(map[valueobject.FilterType][]string)
	for _, group := range groups {
		if !group.Enabled {
			continue
		}

		keywords := group.GetKeywords()
		pattern := u.patternGenerator.GeneratePattern(keywords)
		
		if pattern != "" {
			patterns[group.FilterType] = append(patterns[group.FilterType], pattern)
		}
	}

	return patterns, nil
}