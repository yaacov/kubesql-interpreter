package kubesql

import (
	"fmt"
	"regexp"
	"strings"
)

// clausePattern represents a SQL clause pattern for regex matching
type clausePattern struct {
	name    string // The clause name (SELECT, FROM, etc.)
	pattern string // The regex pattern to match the clause
}

// splitIntoSections splits the normalized query into its constituent SQL clauses.
func (p *Parser) splitIntoSections(query string) (map[string]string, error) {
	sections := make(map[string]string)

	// Define SQL clause patterns in the order they typically appear
	// Each pattern captures the content of the clause while looking ahead
	// for the next clause or end of string
	clausePatterns := []clausePattern{
		{SelectKeyword, `(?i)^select\s+(.+?)(?:\s+from\s+|\s*$)`},
		{FromKeyword, `(?i)\bfrom\s+([^\s]+)(?:\s+where\s+|\s+order\s+by\s+|\s+limit\s+|\s*$)`},
		{WhereKeyword, `(?i)\bwhere\s+(.+?)(?:\s+order\s+by\s+|\s+limit\s+|\s*$)`},
		{OrderByKeyword, `(?i)\border\s+by\s+(.+?)(?:\s+limit\s+|\s*$)`},
		{LimitKeyword, `(?i)\blimit\s+(\d+)\s*$`},
	}

	// Apply each pattern to extract clause content
	for _, clause := range clausePatterns {
		re := regexp.MustCompile(clause.pattern)
		matches := re.FindStringSubmatch(query)
		if len(matches) > 1 {
			sections[clause.name] = strings.TrimSpace(matches[1])
		}
	}

	return sections, nil
}

// parseSelectClause parses the SELECT clause into fields and optional aliases.
// It handles comma-separated field lists and AS keyword for aliases.
// Examples:
//   - "name, namespace" -> [{Expression: "name"}, {Expression: "namespace"}]
//   - "name AS pod_name" -> [{Expression: "name", Alias: "pod_name"}]
func (p *Parser) parseSelectClause(selectClause string) ([]SelectField, error) {
	if selectClause == "" {
		return nil, nil
	}

	var fields []SelectField

	// Split by comma, respecting parentheses for complex expressions
	parts := p.smartSplit(selectClause, ',')

	for _, part := range parts {
		field := SelectField{}
		part = strings.TrimSpace(part)

		// Check for alias using AS keyword (case insensitive)
		asRegex := regexp.MustCompile(`(?i)^(.+?)\s+as\s+([^\s]+)$`)
		matches := asRegex.FindStringSubmatch(part)

		if len(matches) == 3 {
			// Found an alias
			field.Field = TSLQuery(strings.TrimSpace(matches[1]))
			field.Alias = strings.TrimSpace(matches[2])
		} else {
			// No alias, just the expression
			field.Field = TSLQuery(part)
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// parseFromClause parses the FROM clause to extract the Kubernetes resource name.
func (p *Parser) parseFromClause(fromClause string) (string, error) {
	resource := strings.TrimSpace(fromClause)
	if resource == "" {
		return "", fmt.Errorf("FROM clause cannot be empty")
	}
	return resource, nil
}

// parseWhereClause returns the WHERE clause content.
func (p *Parser) parseWhereClause(whereClause string) string {
	return strings.TrimSpace(whereClause)
}

// parseOrderByClause parses the ORDER BY clause into fields and sort directions.
// It handles comma-separated field lists with optional ASC/DESC directions.
// Examples:
//   - "name" -> [{Field: "name", Direction: "ASC"}]
//   - "name DESC, namespace ASC" -> [{Field: "name", Direction: "DESC"}, {Field: "namespace", Direction: "ASC"}]
func (p *Parser) parseOrderByClause(orderByClause string) ([]OrderByField, error) {
	var fields []OrderByField

	// Split by comma, respecting parentheses
	parts := p.smartSplit(orderByClause, ',')

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Split by whitespace to separate field from direction
		tokens := strings.Fields(part)
		if len(tokens) == 0 {
			continue
		}

		field := OrderByField{
			Field:     TSLQuery(tokens[0]),
			Direction: DefaultSortDirection, // Default to ASC
		}

		// Check for explicit direction
		if len(tokens) > 1 {
			// Validate that there are exactly 2 tokens (field and direction)
			if len(tokens) > 2 {
				return nil, fmt.Errorf("invalid ORDER BY field '%s': expected 'field [ASC|DESC]' but found multiple directions", part)
			}

			direction := strings.ToUpper(tokens[1])
			if direction == "DESC" || direction == "ASC" {
				field.Direction = direction
			} else {
				return nil, fmt.Errorf("invalid sort direction '%s': must be 'ASC' or 'DESC'", tokens[1])
			}
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// parseLimitClause parses the LIMIT clause to extract the maximum number of results.
func (p *Parser) parseLimitClause(limitClause string) (int, error) {
	limitStr := strings.TrimSpace(limitClause)
	var limit int
	_, err := fmt.Sscanf(limitStr, "%d", &limit)
	if err != nil {
		return DefaultLimit, fmt.Errorf("invalid LIMIT value: %s", limitStr)
	}
	if limit < 0 {
		return DefaultLimit, fmt.Errorf("LIMIT value must be non-negative: %d", limit)
	}
	return limit, nil
}
