package service

import (
	"regexp"
	"strings"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// FilterResult represents the result of filtering
type FilterResult int

const (
	// FilterResultNeutral means no filter matched
	FilterResultNeutral FilterResult = iota
	// FilterResultInclude means include filter matched
	FilterResultInclude
	// FilterResultExclude means exclude filter matched
	FilterResultExclude
)

// FilterService is a domain service for filtering videos by keywords
type FilterService interface {
	Filter(title string, keywords []*domain.Keyword) FilterResult
}

type filterService struct{}

// NewFilterService creates a new filter service
func NewFilterService() FilterService {
	return &filterService{}
}

// Filter applies keyword filters to a video title
func (fs *filterService) Filter(title string, keywords []*domain.Keyword) FilterResult {
	normalizedTitle := strings.ToLower(title)
	
	// Check exclude filters first (higher priority)
	for _, kw := range keywords {
		if !kw.Enabled || kw.IsDeleted() {
			continue
		}
		
		if kw.FilterType == valueobject.FilterTypeExclude {
			if matches, _ := regexp.MatchString(kw.Pattern, normalizedTitle); matches {
				return FilterResultExclude
			}
		}
	}
	
	// Check include filters
	for _, kw := range keywords {
		if !kw.Enabled || kw.IsDeleted() {
			continue
		}
		
		if kw.FilterType == valueobject.FilterTypeInclude {
			if matches, _ := regexp.MatchString(kw.Pattern, normalizedTitle); matches {
				return FilterResultInclude
			}
		}
	}
	
	return FilterResultNeutral
}