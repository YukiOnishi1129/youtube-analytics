package service

import (
	"fmt"
	"strings"
	"unicode"
)

// KeywordPatternGenerator generates regex patterns from keywords
type KeywordPatternGenerator struct{}

// NewKeywordPatternGenerator creates a new pattern generator
func NewKeywordPatternGenerator() *KeywordPatternGenerator {
	return &KeywordPatternGenerator{}
}

// GeneratePattern generates a regex pattern from keywords
// It creates variations within the same language only
func (g *KeywordPatternGenerator) GeneratePattern(keywords []string) string {
	if len(keywords) == 0 {
		return ""
	}

	var allVariations []string
	for _, keyword := range keywords {
		keyword = strings.TrimSpace(keyword)
		if keyword == "" {
			continue
		}

		// Detect if keyword is Japanese or English
		if isJapanese(keyword) {
			allVariations = append(allVariations, g.generateJapaneseVariations(keyword)...)
		} else {
			allVariations = append(allVariations, g.generateEnglishVariations(keyword)...)
		}
	}

	if len(allVariations) == 0 {
		return ""
	}

	// Remove duplicates across all keywords
	patterns := uniqueStrings(allVariations)

	// Create case-insensitive regex pattern
	return fmt.Sprintf("(?i)(%s)", strings.Join(patterns, "|"))
}

// generateEnglishVariations generates variations for English keywords
func (g *KeywordPatternGenerator) generateEnglishVariations(keyword string) []string {
	// Start with the escaped version
	variations := []string{escapeRegex(keyword)}

	// Handle dots in technology names (e.g., Next.js)
	if strings.Contains(keyword, ".") {
		// Add variation without dots
		noDots := strings.ReplaceAll(keyword, ".", "")
		variations = append(variations, escapeRegex(noDots))
	}

	// Handle spaces and hyphens
	if strings.Contains(keyword, " ") {
		// Add variations with hyphens
		hyphenated := strings.ReplaceAll(keyword, " ", "-")
		variations = append(variations, escapeRegex(hyphenated))
		
		// Add variations without spaces
		noSpaces := strings.ReplaceAll(keyword, " ", "")
		variations = append(variations, escapeRegex(noSpaces))
	}
	
	if strings.Contains(keyword, "-") {
		// Add variations with spaces
		spaced := strings.ReplaceAll(keyword, "-", " ")
		variations = append(variations, escapeRegex(spaced))
		
		// Add variations without hyphens
		noHyphens := strings.ReplaceAll(keyword, "-", "")
		variations = append(variations, escapeRegex(noHyphens))
	}

	// Handle common abbreviations
	variations = append(variations, g.handleCommonAbbreviations(keyword)...)

	// Remove duplicates
	return uniqueStrings(variations)
}

// generateJapaneseVariations generates variations for Japanese keywords
func (g *KeywordPatternGenerator) generateJapaneseVariations(keyword string) []string {
	variations := []string{escapeRegex(keyword)}

	// Add katakana variations for hiragana
	if hasHiragana(keyword) {
		katakanaVersion := hiraganaToKatakana(keyword)
		if katakanaVersion != keyword {
			variations = append(variations, escapeRegex(katakanaVersion))
		}
	}

	// Add hiragana variations for katakana
	if hasKatakana(keyword) {
		hiraganaVersion := katakanaToHiragana(keyword)
		if hiraganaVersion != keyword {
			variations = append(variations, escapeRegex(hiraganaVersion))
		}
	}

	// Handle common Japanese abbreviations
	variations = append(variations, g.handleJapaneseAbbreviations(keyword)...)

	return uniqueStrings(variations)
}

// handleCommonAbbreviations handles common English abbreviations
func (g *KeywordPatternGenerator) handleCommonAbbreviations(keyword string) []string {
	var variations []string

	// Common technology abbreviations
	abbreviations := map[string][]string{
		"javascript": {"js"},
		"typescript": {"ts"},
		"kubernetes": {"k8s"},
		"internationalization": {"i18n"},
		"localization": {"l10n"},
		"continuous integration": {"ci"},
		"continuous deployment": {"cd"},
		"continuous delivery": {"cd"},
	}

	for full, abbrevs := range abbreviations {
		if strings.EqualFold(keyword, full) {
			for _, abbrev := range abbrevs {
				variations = append(variations, escapeRegex(abbrev))
			}
		}
		// Also check reverse (abbreviation to full form)
		for _, abbrev := range abbrevs {
			if strings.EqualFold(keyword, abbrev) {
				variations = append(variations, escapeRegex(full))
			}
		}
	}

	return variations
}

// handleJapaneseAbbreviations handles common Japanese abbreviations
func (g *KeywordPatternGenerator) handleJapaneseAbbreviations(keyword string) []string {
	var variations []string

	// Common Japanese tech abbreviations
	abbreviations := map[string][]string{
		"プログラミング": {"プログラム"},
		"アプリケーション": {"アプリ"},
		"インフラストラクチャー": {"インフラ"},
		"フロントエンド": {"フロント"},
		"バックエンド": {"バック"},
	}

	for full, abbrevs := range abbreviations {
		if keyword == full {
			variations = append(variations, abbrevs...)
		}
		for _, abbrev := range abbrevs {
			if keyword == abbrev {
				variations = append(variations, full)
			}
		}
	}

	return variations
}

// Helper functions

func isJapanese(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Hiragana, unicode.Katakana, unicode.Han) {
			return true
		}
	}
	return false
}

func hasHiragana(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Hiragana) {
			return true
		}
	}
	return false
}

func hasKatakana(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Katakana) {
			return true
		}
	}
	return false
}

func hiraganaToKatakana(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if r >= 0x3041 && r <= 0x3096 {
			runes[i] = r + 0x60
		}
	}
	return string(runes)
}

func katakanaToHiragana(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if r >= 0x30A1 && r <= 0x30F6 {
			runes[i] = r - 0x60
		}
	}
	return string(runes)
}

func escapeRegex(s string) string {
	specialChars := []string{"\\", ".", "*", "+", "?", "^", "$", "(", ")", "[", "]", "{", "}", "|"}
	result := s
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}
	return result
}

func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}