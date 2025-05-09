package denvclustr

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func showDestroyPlan(inputFile, workDirPath string) error {
	// Check if Terraform is installed
	if err := checkTerraformInstalled(); err != nil {
		return fmt.Errorf("terraform is required for destroy plan: %w", err)
	}

	slog.Info("Showing destroy plan", "input", inputFile, "working-dir", workDirPath)

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Get full working directory path
	workingDir := filepath.Join(currentDir, workDirPath)

	// Check if working directory exists
	if _, err := os.Stat(workingDir); os.IsNotExist(err) {
		return fmt.Errorf("working directory not found: %s", workingDir)
	}

	// Check if Terraform state exists
	statePath := filepath.Join(workingDir, "terraform.tfstate")
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return fmt.Errorf("terraform state not found in %s - nothing to destroy", workingDir)
	}

	slog.Info("Using working directory for Terraform operations", "path", workingDir)

	// Initialize Terraform
	tf, err := tfexec.NewTerraform(workingDir, "terraform")
	if err != nil {
		return fmt.Errorf("failed to initialize terraform: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Run terraform init (required even for destroy plan)
	slog.Info("Running terraform init...")
	err = tf.Init(ctx)
	if err != nil {
		return fmt.Errorf("failed to run terraform init: %w", err)
	}

	// Run terraform plan with destroy flag and show output
	fmt.Println("Generating destroy plan...")
	planFilePath := filepath.Join(workingDir, "tfplan-destroy")

	// Save the destroy plan to a file
	hasChanges, err := tf.Plan(ctx, tfexec.Out(planFilePath), tfexec.Destroy(true))
	if err != nil {
		return fmt.Errorf("failed to run terraform destroy plan: %w", err)
	}

	// Display plan results
	if !hasChanges {
		fmt.Println("No resources to destroy. Infrastructure is empty.")
	} else {
		fmt.Println("Destroy plan generated. The following resources will be destroyed:")

		// Display plan details
		plan, err := tf.ShowPlanFile(ctx, planFilePath)
		if err != nil {
			fmt.Println("Failed to show detailed plan, but resources would be destroyed.")
			slog.Error("Failed to show destroy plan details", "error", err)
		} else {
			displayResourceChanges(plan)
		}
	}

	fmt.Printf("\nTo execute this destroy operation, run: denvclustr destroy %s -w %s\n", inputFile, workDirPath)

	return nil
}

func destroyDevcontainers(inputFile, workDirPath string) error {
	// Check if Terraform is installed
	if err := checkTerraformInstalled(); err != nil {
		return fmt.Errorf("terraform is required for destroy operation: %w", err)
	}

	slog.Info("Destroying devcontainers", "input", inputFile, "working-dir", workDirPath)

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Get full working directory path
	workingDir := filepath.Join(currentDir, workDirPath)

	// Check if working directory exists
	if _, err := os.Stat(workingDir); os.IsNotExist(err) {
		return fmt.Errorf("working directory not found: %s", workingDir)
	}

	// Check if Terraform state exists
	statePath := filepath.Join(workingDir, "terraform.tfstate")
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return fmt.Errorf("terraform state not found in %s - nothing to destroy", workingDir)
	}

	slog.Info("Using working directory for Terraform operations", "path", workingDir)

	// Initialize Terraform
	tf, err := tfexec.NewTerraform(workingDir, "terraform")
	if err != nil {
		return fmt.Errorf("failed to initialize terraform: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Run terraform init (required for destroy)
	fmt.Println("Initializing Terraform...")
	err = tf.Init(ctx)
	if err != nil {
		return fmt.Errorf("failed to run terraform init: %w", err)
	}

	// First show the destroy plan
	fmt.Println("\nGenerating destroy plan...")
	planFilePath := filepath.Join(workingDir, "tfplan-destroy")

	hasChanges, err := tf.Plan(ctx, tfexec.Out(planFilePath), tfexec.Destroy(true))
	if err != nil {
		return fmt.Errorf("failed to run terraform destroy plan: %w", err)
	}

	if !hasChanges {
		fmt.Println("No resources to destroy. Infrastructure is empty.")
		return nil
	}

	// Display plan details
	fmt.Println("\nDestroy Plan:")
	fmt.Println("-------------")

	plan, err := tf.ShowPlanFile(ctx, planFilePath)
	if err != nil {
		fmt.Println("Could not display detailed plan. Resources will still be destroyed if you proceed.")
		slog.Error("Failed to show destroy plan details", "error", err)
	} else {
		displayResourceChanges(plan)
	}

	// Confirm before destroying
	fmt.Println("\nWARNING: This will destroy all resources shown above.")
	fmt.Println("You cannot recover from this operation.")
	fmt.Print("Do you want to proceed? (yes/no): ")

	var response string
	fmt.Scanln(&response)
	if response != "yes" {
		fmt.Println("Destroy operation cancelled.")
		return nil
	}

	fmt.Println("\nDestroying resources...")

	err = tf.Destroy(ctx)
	if err != nil {
		return fmt.Errorf("failed to run terraform destroy: %w", err)
	}

	fmt.Printf("\nAll resources have been successfully destroyed!\n")
	fmt.Printf("Terraform files are still preserved in: %s\n", workingDir)

	return nil
}
