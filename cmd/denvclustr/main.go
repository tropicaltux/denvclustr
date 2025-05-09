package main

import (
	"fmt"
	"os"

	"github.com/tropicaltux/denvclustr/internal/denvclustr"
)

func main() {
	if err := denvclustr.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
