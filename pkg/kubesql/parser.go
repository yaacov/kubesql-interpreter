// package kubesql provides a parser for KubeSQL queries designed for kubectl operations.
// It enables SQL syntax for querying Kubernetes resources with support for
// SELECT, FROM, WHERE, ORDER BY, and LIMIT clauses.
package kubesql

import (
	"fmt"
	"strings"
)

// NewParser creates a new parser instance for the given KubeSQL query string.
func NewParser(query string) *Parser {
	return &Parser{query: strings.TrimSpace(query)}
}

// Parse parses the KubeSQL query into structured components.
func (p *Parser) Parse() (*Query, error) {
	result := &Query{
		Limit: DefaultLimit, // -1 indicates no limit
	}

	// Normalize the query by removing extra whitespace
	normalizedQuery := p.normalizeQuery(p.query)

	// Split into sections using regex
	sections, err := p.splitIntoSections(normalizedQuery)
	if err != nil {
		return nil, err
	}

	// Parse each section with detailed error handling
	if selectClause, exists := sections[SelectKeyword]; exists {
		result.Select, err = p.parseSelectClause(selectClause)
		if err != nil {
			return nil, fmt.Errorf("error parsing SELECT clause: %w", err)
		}
	}

	if fromClause, exists := sections[FromKeyword]; exists {
		result.From, err = p.parseFromClause(fromClause)
		if err != nil {
			return nil, fmt.Errorf("error parsing FROM clause: %w", err)
		}
	} else {
		return nil, fmt.Errorf("FROM clause is mandatory")
	}

	if whereClause, exists := sections[WhereKeyword]; exists {
		result.Where = TSLQuery(p.parseWhereClause(whereClause))
	}

	if orderByClause, exists := sections[OrderByKeyword]; exists {
		result.OrderBy, err = p.parseOrderByClause(orderByClause)
		if err != nil {
			return nil, fmt.Errorf("error parsing ORDER BY clause: %w", err)
		}
	}

	if limitClause, exists := sections[LimitKeyword]; exists {
		result.Limit, err = p.parseLimitClause(limitClause)
		if err != nil {
			return nil, fmt.Errorf("error parsing LIMIT clause: %w", err)
		}
	}

	return result, nil
}

// String returns a string representation of the parsed query.
// It reconstructs the KubeSQL syntax from the parsed components.
func (q *Query) String() string {
	var parts []string

	// Build SELECT clause
	if len(q.Select) > 0 {
		var selectParts []string
		for _, field := range q.Select {
			if field.Alias != "" {
				selectParts = append(selectParts, fmt.Sprintf("%s AS %s", field.Field, field.Alias))
			} else {
				selectParts = append(selectParts, string(field.Field))
			}
		}
		parts = append(parts, fmt.Sprintf("SELECT %s", strings.Join(selectParts, ", ")))
	}

	// Add FROM clause (required)
	if q.From != "" {
		parts = append(parts, fmt.Sprintf("FROM %s", q.From))
	}

	// Add WHERE clause if present
	if q.Where != "" {
		parts = append(parts, fmt.Sprintf("WHERE %s", q.Where))
	}

	// Add ORDER BY clause if present
	if len(q.OrderBy) > 0 {
		var orderParts []string
		for _, field := range q.OrderBy {
			orderParts = append(orderParts, fmt.Sprintf("%s %s", field.Field, field.Direction))
		}
		parts = append(parts, fmt.Sprintf("ORDER BY %s", strings.Join(orderParts, ", ")))
	}

	// Add LIMIT clause if specified
	if q.Limit >= 0 {
		parts = append(parts, fmt.Sprintf("LIMIT %d", q.Limit))
	}

	return strings.Join(parts, " ")
}
