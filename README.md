# kubesql-interpreter

[![GoDoc](https://godoc.org/github.com/yaacov/kubesql-interpreter?status.svg)](https://godoc.org/github.com/yaacov/kubesql-interpreter)
[![Go Report Card](https://goreportcard.com/badge/github.com/yaacov/kubesql-interpreter)](https://goreportcard.com/report/github.com/yaacov/kubesql-interpreter)

A Go library for parsing KubeSQL queries designed for kubectl sql operations. This parser enables KubeSQL syntax for querying Kubernetes resources, making it easier to filter, sort, and limit results when working with kubectl.

## Usage

### Command Line Tool

The project includes a command-line tool that can parse KubeSQL queries and output the results in JSON or YAML format.

#### Building

```bash
# Build the command-line tool
make build
```

#### Usage Examples

```bash
# Parse a simple query
./bin/kubesql "SELECT name FROM pods"

# Complex query with filtering and sorting
./bin/kubesql "SELECT name, namespace, status FROM pods WHERE status='Running' ORDER BY name ASC LIMIT 5"

# Output as YAML
./bin/kubesql -format yaml "SELECT * FROM services WHERE namespace='default'"
```

### Library Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yaacov/kubesql-interpreter/pkg/kubesql"
)

func main() {
    query := "SELECT name, namespace, status AS pod_status FROM pods WHERE status='Running' ORDER BY name ASC LIMIT 10"
    
    parser := kubesql.NewParser(query)
    result, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Parsed query: %s\n", result.String())
    fmt.Printf("Resource: %s\n", result.From)
    fmt.Printf("Fields: %v\n", result.Select)
}
```

### Supported Syntax

#### SELECT Clause

```sql
SELECT name, namespace
SELECT name AS pod_name, namespace AS ns
SELECT *
SELECT metadata.name, status.phase
```

#### FROM Clause

```sql
FROM pods
FROM services
FROM deployments
```

#### WHERE Clause

```sql
WHERE status='Running'
WHERE namespace='default'
WHERE metadata.labels.app='nginx'
```

#### ORDER BY Clause

```sql
ORDER BY name ASC
ORDER BY creationTimestamp DESC
ORDER BY name ASC, namespace DESC
```

#### LIMIT Clause

```sql
LIMIT 10
LIMIT 100
```

#### Examples

```sql
-- Get all running pods with their names and status
SELECT name, status FROM pods WHERE status='Running'

-- Get services in default namespace, ordered by name
SELECT name, type, clusterIP FROM services WHERE namespace='default' ORDER BY name

-- Get top 5 deployments by creation time
SELECT name, replicas, creationTimestamp FROM deployments ORDER BY creationTimestamp DESC LIMIT 5

-- Complex query with aliases
SELECT metadata.name AS name, status.phase AS status, spec.nodeName AS node 
FROM pods 
WHERE status.phase='Running' 
ORDER BY metadata.creationTimestamp DESC 
LIMIT 20
```

## API Reference

### Types

#### `TSLQuery`

Represents a raw TSL (Time Series Logic) query string.

```go
type TSLQuery string
```

#### `Query`

Represents a parsed KubeSQL query with all its components.

```go
type Query struct {
    Select  []SelectField  // Fields to select from the resource
    From    string         // Kubernetes resource type (e.g., "pods", "mynamespace/services")
    Where   TSLQuery       // Filter conditions (stored as raw TSL string)
    OrderBy []OrderByField // Sorting specifications
    Limit   int            // Maximum number of results (-1 means no limit)
}
```

#### `SelectField`

Represents a field in the SELECT clause with optional alias support.

```go
type SelectField struct {
    Field TSLQuery // The field expression (e.g., "metadata.name", "status.phase")
    Alias string   // Optional alias for the field (empty if no alias)
}
```

#### `OrderByField`

Represents a field in the ORDER BY clause with sort direction.

```go
type OrderByField struct {
    Field     TSLQuery // Field expression to sort by
    Direction string   // Sort direction: "ASC" or "DESC"
}
```

#### `Parser`

Handles the parsing of KubeSQL queries into structured components.

```go
type Parser struct {
    query string // The original SQL-like query string
}
```

### Methods

#### `NewParser(query string) *Parser`

Creates a new parser instance for the given query string.

#### `Parse() (*Query, error)`

Parses the query and returns a structured representation or an error.

#### `String() string`

Returns a string representation of the parsed query (on `Query`).

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
