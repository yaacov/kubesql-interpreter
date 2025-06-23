package kubesql

// SQL clause names and default values
const (
	// DefaultLimit indicates no limit should be applied
	DefaultLimit = -1

	// Default sort direction
	DefaultSortDirection = "ASC"

	// SQL Keywords
	SelectKeyword  = "SELECT"
	FromKeyword    = "FROM"
	WhereKeyword   = "WHERE"
	OrderByKeyword = "ORDER BY"
	LimitKeyword   = "LIMIT"
)

type TSLQuery string // TSLQuery represents a raw TSL query string

// SelectField represents a field in the SELECT clause with optional alias support.
type SelectField struct {
	Field TSLQuery // The field expression (e.g., "metadata.name", "status.phase")
	Alias string   // Optional alias for the field (empty if no alias)
}

// OrderByField represents a field in the ORDER BY clause with sort direction.
type OrderByField struct {
	Field     TSLQuery // Field expression to sort by
	Direction string   // Sort direction: "ASC" or "DESC"
}

// Query represents a parsed KubeSQL query with all its components.
type Query struct {
	Select  []SelectField  // Fields to select from the resource
	From    string         // Kubernetes resource type (e.g., "pods", "mynamespace/services")
	Where   TSLQuery       // Filter conditions (stored as raw TSL string)
	OrderBy []OrderByField // Sorting specifications
	Limit   int            // Maximum number of results (-1 means no limit)
}

// Parser handles the parsing of KubeSQL queries into structured components.
type Parser struct {
	query string // The original SQL-like query string
}
