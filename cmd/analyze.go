package cmd

import (
	"fmt"

	"github.com/honeybadger/tf-iamgen/internal/parser"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze Terraform files and identify AWS resources",
	Long: `Analyze scans Terraform files in the specified directory
and identifies all AWS resources that will be created.

Example:
  tf-iamgen analyze ./terraform
  tf-iamgen analyze .`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dirPath := args[0]

		// Create parser
		p := parser.NewTerraformParser()

		// Parse directory
		result, err := p.ParseDirectory(dirPath)
		if err != nil {
			return fmt.Errorf("failed to parse directory: %w", err)
		}

		// Print summary
		fmt.Println(result.Summary())

		// Print discovered resources
		if len(result.Resources) > 0 {
			fmt.Println("\nDiscovered Resources:")
			for _, res := range result.Resources {
				fmt.Printf("  - %s.%s (%s:%d)\n", res.Type, res.Name, res.FilePath, res.LineNumber)
			}
		}

		// Print errors/warnings
		if len(result.Errors) > 0 {
			fmt.Println("\nWarnings:")
			for _, err := range result.Errors {
				fmt.Printf("  - %s\n", err.Error())
			}
		}

		return nil
	},
}
