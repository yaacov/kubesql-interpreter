package kubesql

import (
	"testing"
)

func TestNewParser(t *testing.T) {
	query := "  SELECT name FROM pods  "
	parser := NewParser(query)

	if parser.query != "SELECT name FROM pods" {
		t.Errorf("Expected trimmed query, got: %s", parser.query)
	}
}

func TestParseBasicQuery(t *testing.T) {
	query := "SELECT name, namespace FROM pods WHERE status='Running' ORDER BY name ASC LIMIT 10"
	parser := NewParser(query)

	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse query: %v", err)
	}

	// Test integration - all parts working together
	if result.From != "pods" {
		t.Errorf("Expected FROM 'pods', got: %s", result.From)
	}

	if len(result.Select) != 2 {
		t.Errorf("Expected 2 SELECT fields, got: %d", len(result.Select))
	}

	if result.Where != "status='Running'" {
		t.Errorf("Expected WHERE 'status='Running'', got: %s", result.Where)
	}

	if len(result.OrderBy) != 1 {
		t.Errorf("Expected 1 ORDER BY field, got: %d", len(result.OrderBy))
	}

	if result.Limit != 10 {
		t.Errorf("Expected LIMIT 10, got: %d", result.Limit)
	}
}

func TestParseWithAlias(t *testing.T) {
	query := "SELECT name AS pod_name, namespace AS ns FROM pods"
	parser := NewParser(query)

	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse query: %v", err)
	}

	if len(result.Select) != 2 {
		t.Errorf("Expected 2 SELECT fields, got: %d", len(result.Select))
	}

	if result.Select[0].Field != "name" || result.Select[0].Alias != "pod_name" {
		t.Errorf("Expected 'name AS pod_name', got: %s AS %s",
			result.Select[0].Field, result.Select[0].Alias)
	}

	if result.Select[1].Field != "namespace" || result.Select[1].Alias != "ns" {
		t.Errorf("Expected 'namespace AS ns', got: %s AS %s",
			result.Select[1].Field, result.Select[1].Alias)
	}
}

func TestParseMissingFromClause(t *testing.T) {
	query := "SELECT name, namespace"
	parser := NewParser(query)

	_, err := parser.Parse()
	if err == nil {
		t.Error("Expected error for missing FROM clause")
	}
}

func TestQueryString(t *testing.T) {
	query := "SELECT name AS pod_name FROM pods WHERE status='Running' ORDER BY name DESC LIMIT 5"
	parser := NewParser(query)

	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse query: %v", err)
	}

	reconstructed := result.String()
	expected := "SELECT name AS pod_name FROM pods WHERE status='Running' ORDER BY name DESC LIMIT 5"

	if reconstructed != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, reconstructed)
	}
}
