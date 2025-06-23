package kubesql

import (
	"regexp"
	"strings"
)

// normalizeQuery removes extra whitespace and normalizes the query string.
func (p *Parser) normalizeQuery(query string) string {
	// Replace multiple whitespace with single space
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(strings.TrimSpace(query), " ")
}

// smartSplit splits a string by delimiter while respecting parentheses nesting.
func (p *Parser) smartSplit(s string, delimiter rune) []string {
	var result []string
	var current strings.Builder
	parenDepth := 0

	for _, char := range s {
		switch char {
		case '(':
			parenDepth++
			current.WriteRune(char)
		case ')':
			parenDepth--
			current.WriteRune(char)
		case delimiter:
			if parenDepth == 0 {
				// We're not inside parentheses, so this is a real separator
				result = append(result, current.String())
				current.Reset()
			} else {
				// We're inside parentheses, so treat as regular character
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	// Add the last part if there's any content
	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
