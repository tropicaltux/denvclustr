package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	denvclustr_specification "github.com/tropicaltux/denvclustr/pkg/denvclustr_specification"
)

func main() {
	var outputFile string
	flag.StringVar(&outputFile, "out", "", "Output JSON schema file (default: stdout).")
	flag.Parse()

	// Get the JSON schema
	schema := denvclustr_specification.GetDenvclustrSchema()

	// Marshal the schema to JSON with indentation
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling schema to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write to file or stdout
	if outputFile == "" {
		// Write to stdout
		fmt.Println(string(schemaJSON))
	} else {
		// Write to file
		err = os.WriteFile(outputFile, schemaJSON, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing schema to file: %v\n", err)
			os.Exit(1)
		}
	}
}
