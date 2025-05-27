package denvclustr

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/terraform-exec/tfexec"

	_ "github.com/tropicaltux/denvclustr/internal/logger"
	"github.com/tropicaltux/denvclustr/pkg/schema"
)

func showPlan(inputFile, workDirPath string) error {
	// Check if Terraform is installed
	if err := checkTerraformInstalled(); err != nil {
		return fmt.Errorf("terraform is required for plan: %w", err)
	}

	slog.Info("Showing deployment plan", "input", inputFile)

	// Process the input file
	// We don't need the root configuration for showPlan, but we need the HCL file
	_, hclFile, err := processInputFile(inputFile)
	if err != nil {
		return err
	}

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create output directory in current directory if it doesn't exist
	workingDir := filepath.Join(currentDir, workDirPath)
	if err := os.MkdirAll(workingDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	slog.Info("Using working directory for Terraform operations", "path", workingDir)

	// Create the main.tf file in the working directory
	tfFilePath := filepath.Join(workingDir, "main.tf")
	if err := os.WriteFile(tfFilePath, hclFile.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write terraform file: %w", err)
	}

	slog.Info("Created Terraform configuration for plan", "path", tfFilePath)

	// Initialize Terraform
	tf, err := tfexec.NewTerraform(workingDir, "terraform")
	if err != nil {
		return fmt.Errorf("failed to initialize terraform: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Run terraform init
	fmt.Println("Initializing Terraform...")
	err = tf.Init(ctx, tfexec.Upgrade(true))
	if err != nil {
		return fmt.Errorf("failed to run terraform init: %w", err)
	}

	// Run terraform plan and show output
	fmt.Println("Generating detailed plan...")
	planFilePath := filepath.Join(workingDir, "tfplan")

	// First save the plan to a file
	hasChanges, err := tf.Plan(ctx, tfexec.Out(planFilePath))
	if err != nil {
		return fmt.Errorf("failed to run terraform plan: %w", err)
	}

	// Display plan results
	if !hasChanges {
		fmt.Println("No changes. Infrastructure is up-to-date.")
	} else {
		fmt.Println("Detailed plan generated. The plan includes the following changes:")

		// Display plan details
		plan, err := tf.ShowPlanFile(ctx, planFilePath)
		if err != nil {
			fmt.Println("Failed to show detailed plan, but it would result in changes.")
			slog.Error("Failed to show plan details", "error", err)
		} else {
			displayResourceChanges(plan)
		}
	}

	fmt.Printf("\nTo apply this plan, run: denvclustr deploy %s -w %s\n", inputFile, workDirPath)
	fmt.Printf("Terraform files are preserved in: %s\n", workingDir)

	return nil
}

func deployDevcontainers(inputFile, workDirPath string) error {
	// Check if Terraform is installed
	if err := checkTerraformInstalled(); err != nil {
		return fmt.Errorf("terraform is required for deployment: %w", err)
	}

	slog.Info("Deploying devcontainers", "input", inputFile)

	// Process the input file
	root, hclFile, err := processInputFile(inputFile)
	if err != nil {
		return err
	}

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create output directory in current directory if it doesn't exist
	workingDir := filepath.Join(currentDir, workDirPath)
	if err := os.MkdirAll(workingDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	slog.Info("Using working directory for Terraform operations", "path", workingDir)

	// Create the main.tf file in the working directory
	tfFilePath := filepath.Join(workingDir, "main.tf")
	if err := os.WriteFile(tfFilePath, hclFile.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write terraform file: %w", err)
	}

	slog.Info("Created Terraform configuration", "path", tfFilePath)

	// Initialize Terraform
	tf, err := tfexec.NewTerraform(workingDir, "terraform")
	if err != nil {
		return fmt.Errorf("failed to initialize terraform: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Run terraform init
	fmt.Println("Initializing Terraform...")
	err = tf.Init(ctx, tfexec.Upgrade(true))
	if err != nil {
		return fmt.Errorf("failed to run terraform init: %w", err)
	}

	// First show the plan
	fmt.Println("\nGenerating deployment plan...")
	planFilePath := filepath.Join(workingDir, "tfplan")

	hasChanges, err := tf.Plan(ctx, tfexec.Out(planFilePath))
	if err != nil {
		return fmt.Errorf("failed to run terraform plan: %w", err)
	}

	if !hasChanges {
		fmt.Println("No changes to apply. Infrastructure is up-to-date.")
		return nil
	}

	// Display plan details
	fmt.Println("\nDeployment Plan:")
	fmt.Println("----------------")

	plan, err := tf.ShowPlanFile(ctx, planFilePath)
	if err != nil {
		fmt.Println("Could not display detailed plan. Continuing with deployment.")
		slog.Error("Failed to show plan details", "error", err)
	} else {
		displayResourceChanges(plan)
	}

	fmt.Println("\nProceeding with deployment...")

	// Run terraform apply
	err = tf.Apply(ctx)
	if err != nil {
		// Try to recover by destroying what was created
		fmt.Println("\nERROR: Deployment failed with error:", err)
		fmt.Println("Attempting to clean up any created resources...")

		// Try to destroy any resources that were created
		destroyErr := tf.Destroy(ctx)
		if destroyErr != nil {
			fmt.Println("WARNING: Cleanup also failed:", destroyErr)
			fmt.Println("You may need to manually clean up resources.")
			return fmt.Errorf("deployment failed and cleanup also failed: %w, cleanup error: %w", err, destroyErr)
		} else {
			fmt.Println("Cleanup successful. All created resources have been destroyed.")
			return fmt.Errorf("deployment failed but resources were cleaned up: %w", err)
		}
	}

	fmt.Println("\nDeployment completed successfully!")

	// Get and display the outputs
	outputs, err := tf.Output(ctx)
	if err != nil {
		return fmt.Errorf("failed to get outputs: %w", err)
	}

	if err := displayDeploymentOutputs(outputs, root); err != nil {
		return fmt.Errorf("failed to display outputs: %w", err)
	}

	fmt.Printf("\nTerraform files are preserved in: %s\n", workingDir)
	fmt.Printf("To destroy these resources, run: denvclustr destroy %s -w %s\n", inputFile, workDirPath)

	return nil
}

// displayDeploymentOutputs formats and displays the deployment outputs in a user-friendly way
func displayDeploymentOutputs(outputs map[string]tfexec.OutputMeta, root *schema.DenvclustrRoot) error {
	if len(outputs) == 0 {
		fmt.Println("\nNo outputs available.")
		return nil
	}

	fmt.Println("\n🚀 Available Access Methods:")
	fmt.Println("==========================")

	// Find the region from the infrastructure configuration
	var awsRegion string
	if root != nil && len(root.Infrastructure) > 0 {
		for _, infra := range root.Infrastructure {
			if infra.Provider == schema.ProviderAws {
				awsRegion = string(infra.Region)
				slog.Info("Using region from infrastructure configuration", "region", awsRegion)
				break
			}
		}
	}

	if awsRegion == "" {
		slog.Warn("AWS region not found in configuration. Cannot retrieve tokens from SSM.")
	}

	for _, output := range outputs {
		var outputData map[string]any
		if err := json.Unmarshal([]byte(output.Value), &outputData); err != nil {
			return fmt.Errorf("failed to unmarshal outputs: %w", err)
		}

		module := outputData["module"].(map[string]any)
		devcontainers := module["devcontainers"].([]any)

		for _, devcontainer := range devcontainers {
			devcontainerMap := devcontainer.(map[string]any)
			id := devcontainerMap["id"].(string)
			fmt.Printf("\n📦 Devcontainer %s:\n", id)
			remote_access := devcontainerMap["remote_access"].(map[string]any)

			if remote_access["openvscode_server"] != nil {
				openvscode_server := remote_access["openvscode_server"].(map[string]any)
				tokenSSMParameter := openvscode_server["token_ssm_parameter"].(string)
				urlTemplate := openvscode_server["url"].(string)

				if awsRegion == "" {
					// No region available, don't try to get token
					fmt.Printf("  🌐 VS Code Server URL template: %s\n", urlTemplate)
					fmt.Printf("  🔑 Get token from AWS SSM parameter: %s\n", tokenSSMParameter)
					fmt.Printf("  ℹ️  Region not specified in configuration. Cannot retrieve token automatically.\n")
					fmt.Printf("  ℹ️  AWS CLI command (specify your region): aws ssm get-parameter --region YOUR_REGION --name %s --with-decryption --query \"Parameter.Value\" --output text\n", tokenSSMParameter)
				} else {
					// Region is available, try to get token
					token, err := getTokenFromSSM(tokenSSMParameter, awsRegion)
					if err != nil {
						slog.Error("Failed to get OpenVSCode token", "error", err, "parameter", tokenSSMParameter)
						fmt.Printf("  🌐 VS Code Server URL template: %s\n", urlTemplate)
						fmt.Printf("  🔑 Get token from AWS SSM parameter: %s\n", tokenSSMParameter)
						fmt.Printf("  ℹ️  Could not retrieve token: %s\n", err)
						fmt.Printf("  ℹ️  AWS CLI command: aws ssm get-parameter --region %s --name %s --with-decryption --query \"Parameter.Value\" --output text\n",
							awsRegion, tokenSSMParameter)
					} else {
						// Replace {token} placeholder with the actual token
						url := strings.Replace(urlTemplate, "{token}", token, 1)
						fmt.Printf("  🌐 VS Code Server: %s\n", url)
					}
				}
			}

			if remote_access["ssh"] != nil {
				ssh := remote_access["ssh"].(map[string]any)
				fmt.Printf("  🔑 SSH Access: %s\n", ssh["command"])
			}
		}
	}

	return nil
}
