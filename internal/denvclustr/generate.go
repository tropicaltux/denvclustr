package denvclustr

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func generateHcl(inputFile, outputFile string) error {
	slog.Info("Generating Terraform HCL", "input", inputFile, "output", outputFile)

	// Process the input file
	_, hclFile, err := processInputFile(inputFile)
	if err != nil {
		return err
	}

	// Create output directory if it doesn't exist
	outDir := filepath.Dir(outputFile)
	if outDir != "." && outDir != "" {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Write the HCL file (will overwrite if it exists)
	if err := os.WriteFile(outputFile, hclFile.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	slog.Info("Successfully generated Terraform HCL", "output", outputFile)
	fmt.Printf("Successfully generated Terraform HCL: %s\n", outputFile)
	return nil
}
