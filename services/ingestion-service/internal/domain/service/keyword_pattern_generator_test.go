package service

import (
	"strings"
	"testing"
)

func TestKeywordPatternGenerator_GeneratePattern(t *testing.T) {
	generator := NewKeywordPatternGenerator()

	tests := []struct {
		name     string
		keywords []string
		want     []string // Expected patterns within the regex
		notWant  []string // Patterns that should NOT be in the regex
	}{
		{
			name:     "English with dots - Next.js",
			keywords: []string{"Next.js"},
			want:     []string{"Next\\.js", "Nextjs"},
			notWant:  []string{"ネクスト", "次"},
		},
		{
			name:     "English with spaces",
			keywords: []string{"machine learning"},
			want:     []string{"machine learning", "machine-learning", "machinelearning"},
			notWant:  []string{"機械学習"},
		},
		{
			name:     "English abbreviation - JavaScript",
			keywords: []string{"JavaScript"},
			want:     []string{"JavaScript", "js"},
			notWant:  []string{"ジャバスクリプト"},
		},
		{
			name:     "Japanese only - プログラミング",
			keywords: []string{"プログラミング"},
			want:     []string{"プログラミング", "プログラム"},
			notWant:  []string{"programming", "coding"},
		},
		{
			name:     "Japanese katakana/hiragana conversion",
			keywords: []string{"ぷろぐらみんぐ"},
			want:     []string{"ぷろぐらみんぐ", "プログラミング"},
			notWant:  []string{"programming"},
		},
		{
			name:     "Multiple keywords same language",
			keywords: []string{"React", "Vue", "Angular"},
			want:     []string{"React", "Vue", "Angular"},
			notWant:  []string{"リアクト", "ビュー", "アンギュラー"},
		},
		{
			name:     "Mixed English and Japanese - separate patterns",
			keywords: []string{"Python", "パイソン"},
			want:     []string{"Python", "パイソン"},
			notWant:  []string{}, // They should be separate, not cross-translated
		},
		{
			name:     "Empty keywords",
			keywords: []string{},
			want:     []string{},
			notWant:  []string{},
		},
		{
			name:     "Kubernetes abbreviation",
			keywords: []string{"Kubernetes"},
			want:     []string{"Kubernetes", "k8s"},
			notWant:  []string{"クバネティス"},
		},
		{
			name:     "Japanese abbreviation - アプリケーション",
			keywords: []string{"アプリケーション"},
			want:     []string{"アプリケーション", "アプリ"},
			notWant:  []string{"application", "app"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generator.GeneratePattern(tt.keywords)

			// Check if pattern starts with (?i)( and ends with )
			if len(tt.keywords) > 0 && (!strings.HasPrefix(got, "(?i)(") || !strings.HasSuffix(got, ")")) {
				t.Errorf("Pattern should be wrapped with (?i)(...), got: %s", got)
			}

			// Extract the pattern content
			if len(got) > 5 {
				patternContent := got[5 : len(got)-1] // Remove (?i)( and )

				// Check wanted patterns
				for _, want := range tt.want {
					if !strings.Contains(patternContent, want) {
						t.Errorf("Pattern should contain %q, got: %s", want, got)
					}
				}

				// Check unwanted patterns
				for _, notWant := range tt.notWant {
					if strings.Contains(patternContent, notWant) {
						t.Errorf("Pattern should NOT contain %q, got: %s", notWant, got)
					}
				}
			} else if len(tt.want) > 0 {
				t.Errorf("Expected pattern with content, got: %s", got)
			}
		})
	}
}

func TestKeywordPatternGenerator_UniquePatterns(t *testing.T) {
	generator := NewKeywordPatternGenerator()

	// Test that duplicate variations are removed
	keywords := []string{"next.js", "nextjs"} // Both should generate "nextjs" variation
	pattern := generator.GeneratePattern(keywords)

	// Count occurrences of "nextjs" in the pattern
	count := strings.Count(pattern, "nextjs")
	if count > 1 {
		t.Errorf("Pattern should not contain duplicate variations, found %d occurrences of 'nextjs'", count)
	}
}

func TestKeywordPatternGenerator_RegexEscaping(t *testing.T) {
	generator := NewKeywordPatternGenerator()

	// Test that special regex characters are escaped
	keywords := []string{"C++", "Node.js", "ASP.NET"}
	pattern := generator.GeneratePattern(keywords)

	// Check proper escaping
	expectedEscapes := map[string]string{
		"C++":     "C\\+\\+",
		"Node.js": "Node\\.js",
		"ASP.NET": "ASP\\.NET",
	}

	for original, escaped := range expectedEscapes {
		if !strings.Contains(pattern, escaped) {
			t.Errorf("Pattern should contain properly escaped %q as %q, got: %s", original, escaped, pattern)
		}
	}
}