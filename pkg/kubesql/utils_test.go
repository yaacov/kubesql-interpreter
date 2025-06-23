package kubesql

import (
	"testing"
)

func TestSmartSplit(t *testing.T) {
	parser := NewParser("")

	testCases := []struct {
		input    string
		expected []string
	}{
		{"a, b, c", []string{"a", " b", " c"}},
		{"func(a, b), c", []string{"func(a, b)", " c"}},
		{"nested(func(a, b), c), d", []string{"nested(func(a, b), c)", " d"}},
		{"simple", []string{"simple"}},
		{"", []string{}}, // Empty string should return empty slice
		{"a,b,c", []string{"a", "b", "c"}},
		{"func(a,b,c),d", []string{"func(a,b,c)", "d"}},
	}

	for _, tc := range testCases {
		result := parser.smartSplit(tc.input, ',')
		if len(result) != len(tc.expected) {
			t.Errorf("For input '%s', expected %d parts, got %d",
				tc.input, len(tc.expected), len(result))
			continue
		}

		for i, part := range result {
			if part != tc.expected[i] {
				t.Errorf("For input '%s', part %d: expected '%s', got '%s'",
					tc.input, i, tc.expected[i], part)
			}
		}
	}
}

func TestNormalizeQuery(t *testing.T) {
	parser := NewParser("")

	testCases := []struct {
		input    string
		expected string
	}{
		{"  SELECT  name  FROM  pods  ", "SELECT name FROM pods"},
		{"SELECT\n\tname\nFROM\tpods", "SELECT name FROM pods"},
		{"select name from pods", "select name from pods"},
		{"", ""},
		{"   ", ""},
	}

	for _, tc := range testCases {
		result := parser.normalizeQuery(tc.input)
		if result != tc.expected {
			t.Errorf("For input '%s', expected '%s', got '%s'",
				tc.input, tc.expected, result)
		}
	}
}
