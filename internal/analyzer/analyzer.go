package analyzer

import (
	"regexp"
	"strings"

	"github.com/geodask/clipboard-manager/internal/domain"
)

type Analyzer interface {
	Analyze(entry *domain.ClipboardEntry) *domain.Analysis
}

type SimpleAnalyzer struct {
	passwordPattern *regexp.Regexp
	tokenPattern    *regexp.Regexp
	apiKeyPattern   *regexp.Regexp
}

func NewSimpleAnalyzer() *SimpleAnalyzer {
	return &SimpleAnalyzer{
		passwordPattern: regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[:=]\s*\S+`),
		tokenPattern:    regexp.MustCompile(`(?i)(token|bearer)\s*[:=]?\s*[A-Za-z0-9_-]{20,}`),
		apiKeyPattern:   regexp.MustCompile(`(?i)(api[_-]?key|secret[_-]?key)\s*[:=]\s*\S+`),
	}
}

func (a *SimpleAnalyzer) Analyze(entry *domain.ClipboardEntry) *domain.Analysis {
	content := entry.Content

	if a.passwordPattern.MatchString(content) {
		return &domain.Analysis{
			Type:        domain.ContentTypeText,
			IsSensitive: true,
			Reason:      "contains password",
		}
	}

	if a.tokenPattern.MatchString(content) {
		return &domain.Analysis{
			Type:        domain.ContentTypeText,
			IsSensitive: true,
			Reason:      "contains token",
		}
	}

	if a.apiKeyPattern.MatchString(content) {
		return &domain.Analysis{
			Type:        domain.ContentTypeText,
			IsSensitive: true,
			Reason:      "contains API key",
		}
	}

	contenType := a.detectType(content)

	return &domain.Analysis{
		Type:        contenType,
		IsSensitive: false,
		Reason:      "",
	}
}

func (a *SimpleAnalyzer) detectType(content string) domain.ContentType {
	if strings.HasPrefix(content, "http://") || strings.HasPrefix(content, "https://") {
		return domain.ContentTypeURL
	}

	if strings.HasPrefix(content, "/") || strings.Contains(content, ":\\") {
		return domain.ContentTypeFilePath
	}

	codeIndicators := []string{"func ", "def ", "class ", "import ", "const ", "let ", "var "}

	for _, indicator := range codeIndicators {
		if strings.Contains(content, indicator) {
			return domain.ContentTypeCode
		}
	}

	return domain.ContentTypeText
}
