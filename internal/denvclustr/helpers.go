package denvclustr

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/hashicorp/hcl/v2/hclwrite"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/tropicaltux/denvclustr/pkg/dc2tf"
	"github.com/tropicaltux/denvclustr/pkg/schema"
)

// processInputFile reads, parses, and converts a denvclustr JSON file to HCL.
// It returns the parsed configuration and HCL content, or an error if any step fails.
func processInputFile(inputFile string) (*schema.DenvclustrRoot, *hclwrite.File, error) {
	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("input file not found: %s", inputFile)
	}

	// Read the input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read input file: %w", err)
	}

	// Parse the denvclustr file
	root, err := schema.Parse(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse denvclustr file: %w", err)
	}

	// Convert to Terraform HCL
	hclFile, err := dc2tf.Convert(root)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert to Terraform: %w", err)
	}

	return root, hclFile, nil
}

// checkTerraformInstalled verifies that the Terraform CLI is available
func checkTerraformInstalled() error {
	// Execute a simple command to check if terraform is available
	path, err := exec.LookPath("terraform")
	if err != nil {
		return fmt.Errorf("terraform CLI not found: %w", err)
	}

	slog.Info("Terraform found", "path", path)
	return nil
}

// displayResourceChanges formats and displays the changes from a Terraform plan
func displayResourceChanges(plan *tfjson.Plan) {
	if plan == nil || len(plan.ResourceChanges) == 0 {
		fmt.Println("No resource changes detected.")
		return
	}

	fmt.Printf("\nResources to be created/modified: %d\n", len(plan.ResourceChanges))

	for _, change := range plan.ResourceChanges {
		if change.Change != nil && len(change.Change.Actions) > 0 {
			fmt.Printf("  %s: %s\n", change.Change.Actions[0], change.Address)
		}
	}
}
