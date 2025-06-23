package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yaacov/kubesql-interpreter/pkg/kubesql"
	"gopkg.in/yaml.v3"
)

var (
	outputFormat = flag.String("format", "json", "Output format: json or yaml")
	helpFlag     = flag.Bool("help", false, "Show help message")
)

func main() {
	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	// Get the KubeSQL query from command line arguments
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No SQL query provided\n")
		showUsage()
		os.Exit(1)
	}

	query := strings.Join(args, " ")

	// Parse the KubeSQL query
	parser := kubesql.NewParser(query)
	result, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing query: %v\n", err)
		os.Exit(1)
	}

	// Convert to the desired output format
	var output []byte
	switch strings.ToLower(*outputFormat) {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)
		err = encoder.Encode(result)
		if err != nil {
			log.Fatalf("Error marshaling to JSON: %v", err)
		}
		return
	case "yaml":
		output, err = yaml.Marshal(result)
		if err != nil {
			log.Fatalf("Error marshaling to YAML: %v", err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unsupported output format '%s'. Supported formats: json, yaml\n", *outputFormat)
		os.Exit(1)
	}

	fmt.Print(string(output))
}

func showHelp() {
	fmt.Printf(`KubeSQL Parser Command Line Tool

DESCRIPTION:
    Parse KubeSQL queries and output the result as JSON or YAML.
    
USAGE:
    sql [OPTIONS] <SQL_QUERY>

OPTIONS:
    -format string
            Output format: json or yaml (default "json")
    -help
            Show this help message
    -version
            Show version information

EXAMPLES:
    # Parse a simple SELECT query and output as JSON
    sql "SELECT name, namespace FROM pods WHERE status='Running'"
    
    # Parse a query with aliases and output as YAML
    sql -format yaml "SELECT metadata.name AS pod_name, status.phase AS status FROM pods"
    
    # Parse a complex query with ORDER BY and LIMIT
    sql "SELECT * FROM deployments ORDER BY creationTimestamp DESC LIMIT 5"
    
    # Query using quotes to handle special characters
    sql "SELECT name, type, clusterIP FROM services WHERE namespace='default'"

SUPPORTED SQL FEATURES:
    - SELECT with field selection and aliases
    - FROM with Kubernetes resource types
    - WHERE with filter conditions
    - ORDER BY with ASC/DESC sorting
    - LIMIT for result count restriction

`)
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: sql [OPTIONS] <KUBESQL_QUERY>\n")
	fmt.Fprintf(os.Stderr, "Run 'sql -help' for more information.\n")
}
