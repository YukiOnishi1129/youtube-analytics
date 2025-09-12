package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

type keywordUseCase struct {
	keywordRepo gateway.KeywordRepository
}

func NewKeywordUseCase(
	keywordRepo gateway.KeywordRepository,
) input.KeywordInputPort {
	return &keywordUseCase{
		keywordRepo: keywordRepo,
	}
}

func (u *keywordUseCase) ListKeywords(ctx context.Context) ([]*domain.Keyword, error) {
	// Get all keywords (not just enabled ones)
	return u.keywordRepo.FindAll(ctx, false)
}

func (u *keywordUseCase) CreateKeyword(ctx context.Context, input *input.CreateKeywordInput) (*domain.Keyword, error) {
	// Convert filter type
	ft := valueobject.FilterType(input.FilterType)
	if !ft.IsValid() {
		return nil, domain.ErrInvalidFilterType
	}

	// Create keyword using the constructor
	keyword, err := domain.NewKeyword(
		valueobject.UUID(uuid.New().String()),
		input.Name,
		ft,
		input.Pattern,
		input.Description,
	)
	if err != nil {
		return nil, err
	}

	// Save keyword
	if err := u.keywordRepo.Save(ctx, keyword); err != nil {
		return nil, err
	}

	return keyword, nil
}

func (u *keywordUseCase) GetKeyword(ctx context.Context, keywordID uuid.UUID) (*domain.Keyword, error) {
	keyword, err := u.keywordRepo.FindByID(ctx, valueobject.UUID(keywordID.String()))
	if err != nil {
		return nil, err
	}
	return keyword, nil
}

func (u *keywordUseCase) UpdateKeyword(ctx context.Context, input *input.UpdateKeywordInput) (*domain.Keyword, error) {
	// Get existing keyword
	keyword, err := u.keywordRepo.FindByID(ctx, valueobject.UUID(input.KeywordID.String()))
	if err != nil {
		return nil, err
	}

	// Convert filter type
	ft := valueobject.FilterType(input.FilterType)
	if !ft.IsValid() {
		return nil, domain.ErrInvalidFilterType
	}

	// Update using domain method
	if err := keyword.Update(&input.Name, &ft, &input.Pattern, input.Description); err != nil {
		return nil, err
	}

	// Update enabled status
	if input.Enabled && !keyword.Enabled {
		keyword.Enable()
	} else if !input.Enabled && keyword.Enabled {
		keyword.Disable()
	}

	// Save updated keyword
	if err := u.keywordRepo.Save(ctx, keyword); err != nil {
		return nil, err
	}

	return keyword, nil
}

func (u *keywordUseCase) EnableKeyword(ctx context.Context, keywordID uuid.UUID) (*domain.Keyword, error) {
	// Get keyword
	keyword, err := u.keywordRepo.FindByID(ctx, valueobject.UUID(keywordID.String()))
	if err != nil {
		return nil, err
	}

	// Enable it
	keyword.Enable()

	// Save
	if err := u.keywordRepo.Save(ctx, keyword); err != nil {
		return nil, err
	}

	return keyword, nil
}

func (u *keywordUseCase) DisableKeyword(ctx context.Context, keywordID uuid.UUID) (*domain.Keyword, error) {
	// Get keyword
	keyword, err := u.keywordRepo.FindByID(ctx, valueobject.UUID(keywordID.String()))
	if err != nil {
		return nil, err
	}

	// Disable it
	keyword.Disable()

	// Save
	if err := u.keywordRepo.Save(ctx, keyword); err != nil {
		return nil, err
	}

	return keyword, nil
}

func (u *keywordUseCase) DeleteKeyword(ctx context.Context, keywordID uuid.UUID) error {
	return u.keywordRepo.SoftDelete(ctx, valueobject.UUID(keywordID.String()))
}