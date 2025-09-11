package service

import (
	"regexp"
	"strings"
)

// PatternBuilder is a domain service for building regex patterns from keywords
type PatternBuilder interface {
	BuildPattern(name string) (string, error)
}

type patternBuilder struct{}

// NewPatternBuilder creates a new pattern builder
func NewPatternBuilder() PatternBuilder {
	return &patternBuilder{}
}

// BuildPattern builds a regex pattern from a keyword name
// Example: "Next.js", "Next JS", "NextJS" → (?i)next\.?js
func (pb *patternBuilder) BuildPattern(name string) (string, error) {
	// Normalize the name
	normalized := strings.TrimSpace(name)
	if normalized == "" {
		return "", nil
	}

	// Escape special regex characters
	escaped := regexp.QuoteMeta(normalized)
	
	// Replace spaces with flexible pattern (\s* for optional spaces)
	pattern := strings.ReplaceAll(escaped, " ", `\s*`)
	
	// Replace dots with optional dots (\.? for optional dots)
	pattern = strings.ReplaceAll(pattern, `\.`, `\.?`)
	
	// Make case insensitive
	pattern = "(?i)" + pattern
	
	// Validate the pattern
	if _, err := regexp.Compile(pattern); err != nil {
		return "", err
	}
	
	return pattern, nil
}