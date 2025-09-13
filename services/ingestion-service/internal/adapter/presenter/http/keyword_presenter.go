package http

import (
	"net/http"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/http/generated"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/presenter"
)

// keywordPresenter implements presenter.KeywordPresenter for HTTP
type keywordPresenter struct{}

// NewKeywordPresenter creates a new HTTP keyword presenter
func NewKeywordPresenter() presenter.KeywordPresenter {
	return &keywordPresenter{}
}

// PresentKeyword presents a single keyword
// Note: This is for HTTP REST API, not used in current implementation
// as keyword operations are exposed via gRPC only
func (p *keywordPresenter) PresentKeyword(keyword *domain.Keyword) interface{} {
	// Not used in HTTP API - keywords are accessed via gRPC
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: "Keyword operations are available via gRPC API",
		},
	}
}

// PresentKeywords presents multiple keywords
// Note: This is for HTTP REST API, not used in current implementation
// as keyword operations are exposed via gRPC only
func (p *keywordPresenter) PresentKeywords(keywords []*domain.Keyword) interface{} {
	// Not used in HTTP API - keywords are accessed via gRPC
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: "Keyword operations are available via gRPC API",
		},
	}
}

// PresentDeleted presents a successful deletion
func (p *keywordPresenter) PresentDeleted() interface{} {
	return &HTTPResponse{
		StatusCode: http.StatusNoContent,
		Body:       nil,
	}
}

// PresentError presents an error
func (p *keywordPresenter) PresentError(err error) interface{} {
	switch err {
	case domain.ErrKeywordNotFound:
		return &HTTPResponse{
			StatusCode: http.StatusNotFound,
			Body: generated.Error{
				Code:    "KEYWORD_NOT_FOUND",
				Message: err.Error(),
			},
		}
	case domain.ErrKeywordDuplicate:
		return &HTTPResponse{
			StatusCode: http.StatusConflict,
			Body: generated.Error{
				Code:    "KEYWORD_DUPLICATE",
				Message: err.Error(),
			},
		}
	default:
		return &HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body: generated.Error{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		}
	}
}