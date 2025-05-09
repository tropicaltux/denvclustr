package denvclustr

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "denvclustr",
	Short: "A tool for deploying devcontainers.",
	Long: `denvclustr is a CLI tool that helps you deploy devcontainers
or generate Terraform HCL files based on denvclustr configuration files.`,
}

var deployCmd = &cobra.Command{
	Use:   "deploy [file]",
	Short: "Deploy devcontainers from a denvclustr file",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := "denvclustr.json"
		if len(args) > 0 {
			inputFile = args[0]
		}

		if planOnly {
			return showPlan(inputFile, workingDir)
		}
		return deployDevcontainers(inputFile, workingDir)
	},
}

var destroyCmd = &cobra.Command{
	Use:   "destroy [file]",
	Short: "Destroy devcontainers deployed from a denvclustr file",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := "denvclustr.json"
		if len(args) > 0 {
			inputFile = args[0]
		}

		if planOnly {
			return showDestroyPlan(inputFile, workingDir)
		}
		return destroyDevcontainers(inputFile, workingDir)
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate [file]",
	Short: "Generate Terraform HCL from a denvclustr file",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := "denvclustr.json"
		if len(args) > 0 {
			inputFile = args[0]
		}

		// If output file not specified, use current directory + denvclustr.tf
		if outputFile == "" {
			currentDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			outputFile = filepath.Join(currentDir, "denvclustr.tf")
		}

		return generateHcl(inputFile, outputFile)
	},
}
var (
	outputFile string
	planOnly   bool
	workingDir string
)

func init() {
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output Terraform file (default: ./denvclustr.tf)")

	deployCmd.Flags().BoolVarP(&planOnly, "plan", "p", false, "Show deployment plan without applying changes")
	deployCmd.Flags().StringVarP(&workingDir, "working-dir", "w", "output", "Working directory for Terraform operations")

	destroyCmd.Flags().BoolVarP(&planOnly, "plan", "p", false, "Show destroy plan without applying changes")
	destroyCmd.Flags().StringVarP(&workingDir, "working-dir", "w", "output", "Working directory for Terraform operations")

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(destroyCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
