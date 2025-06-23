package kubesql

import (
	"testing"
)

func TestParseSelectClause(t *testing.T) {
	parser := NewParser("")

	testCases := []struct {
		input    string
		expected []SelectField
	}{
		{
			"name, namespace",
			[]SelectField{
				{Field: "name", Alias: ""},
				{Field: "namespace", Alias: ""},
			},
		},
		{
			"name AS pod_name, namespace AS ns",
			[]SelectField{
				{Field: "name", Alias: "pod_name"},
				{Field: "namespace", Alias: "ns"},
			},
		},
		{
			"*",
			[]SelectField{
				{Field: "*", Alias: ""},
			},
		},
		{
			"COUNT(*) AS total",
			[]SelectField{
				{Field: "COUNT(*)", Alias: "total"},
			},
		},
	}

	for _, tc := range testCases {
		result, err := parser.parseSelectClause(tc.input)
		if err != nil {
			t.Errorf("For input '%s', unexpected error: %v", tc.input, err)
			continue
		}

		if len(result) != len(tc.expected) {
			t.Errorf("For input '%s', expected %d fields, got %d",
				tc.input, len(tc.expected), len(result))
			continue
		}

		for i, field := range result {
			if field.Field != tc.expected[i].Field || field.Alias != tc.expected[i].Alias {
				t.Errorf("For input '%s', field %d: expected %+v, got %+v",
					tc.input, i, tc.expected[i], field)
			}
		}
	}
}

func TestParseOrderByClause(t *testing.T) {
	parser := NewParser("")

	testCases := []struct {
		input    string
		expected []OrderByField
	}{
		{
			"name ASC",
			[]OrderByField{
				{Field: "name", Direction: "ASC"},
			},
		},
		{
			"name DESC, namespace",
			[]OrderByField{
				{Field: "name", Direction: "DESC"},
				{Field: "namespace", Direction: "ASC"},
			},
		},
		{
			"created_at",
			[]OrderByField{
				{Field: "created_at", Direction: "ASC"},
			},
		},
	}

	for _, tc := range testCases {
		result, err := parser.parseOrderByClause(tc.input)
		if err != nil {
			t.Errorf("For input '%s', unexpected error: %v", tc.input, err)
			continue
		}

		if len(result) != len(tc.expected) {
			t.Errorf("For input '%s', expected %d fields, got %d",
				tc.input, len(tc.expected), len(result))
			continue
		}

		for i, field := range result {
			if field.Field != tc.expected[i].Field || field.Direction != tc.expected[i].Direction {
				t.Errorf("For input '%s', field %d: expected %+v, got %+v",
					tc.input, i, tc.expected[i], field)
			}
		}
	}
}

func TestParseOrderByClauseInvalidDirection(t *testing.T) {
	parser := NewParser("")

	invalidDirectionCases := []string{
		"name INVALID",
		"name UP",
		"name DOWN",
		"name asc desc",         // Multiple directions
		"name ASC DESC",         // Multiple directions (uppercase)
		"name DESC ASC",         // Multiple directions (reverse order)
		"name ASC INVALID",      // Valid direction followed by invalid
		"name INVALID ASC",      // Invalid direction followed by valid
		"name ASC DESC INVALID", // Three tokens
		"name ASCENDING",        // Wrong keyword
		"name DESCENDING",       // Wrong keyword
		"name xyz",
	}

	for _, input := range invalidDirectionCases {
		_, err := parser.parseOrderByClause(input)
		if err == nil {
			t.Errorf("For input '%s', expected error but got none", input)
		}
	}
}

func TestParseLimitClause(t *testing.T) {
	parser := NewParser("")

	testCases := []struct {
		input    string
		expected int
		hasError bool
	}{
		{"10", 10, false},
		{"0", 0, false},
		{"999", 999, false},
		{"abc", 0, true},
		{"-5", 0, true},
		{"", 0, true},
	}

	for _, tc := range testCases {
		result, err := parser.parseLimitClause(tc.input)

		if tc.hasError {
			if err == nil {
				t.Errorf("For input '%s', expected error but got none", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("For input '%s', unexpected error: %v", tc.input, err)
			}
			if result != tc.expected {
				t.Errorf("For input '%s', expected %d, got %d", tc.input, tc.expected, result)
			}
		}
	}
}

func TestParseOrderByClauseCaseInsensitive(t *testing.T) {
	parser := NewParser("")

	testCases := []struct {
		input    string
		expected []OrderByField
	}{
		{
			"name asc", // lowercase
			[]OrderByField{
				{Field: "name", Direction: "ASC"},
			},
		},
		{
			"name desc", // lowercase
			[]OrderByField{
				{Field: "name", Direction: "DESC"},
			},
		},
		{
			"name Asc", // mixed case
			[]OrderByField{
				{Field: "name", Direction: "ASC"},
			},
		},
		{
			"name DescENDing", // should fail
			[]OrderByField{},
		},
	}

	for i, tc := range testCases {
		result, err := parser.parseOrderByClause(tc.input)

		// Special case for the invalid direction test
		if i == 3 {
			if err == nil {
				t.Errorf("For input '%s', expected error but got none", tc.input)
			}
			continue
		}

		if err != nil {
			t.Errorf("For input '%s', unexpected error: %v", tc.input, err)
			continue
		}

		if len(result) != len(tc.expected) {
			t.Errorf("For input '%s', expected %d fields, got %d",
				tc.input, len(tc.expected), len(result))
			continue
		}

		for j, field := range result {
			if field.Field != tc.expected[j].Field || field.Direction != tc.expected[j].Direction {
				t.Errorf("For input '%s', field %d: expected %+v, got %+v",
					tc.input, j, tc.expected[j], field)
			}
		}
	}
}
