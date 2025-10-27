package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	outputFile string
)

var generateCmd = &cobra.Command{
	Use:   "generate [path]",
	Short: "Generate IAM policy from Terraform files",
	Long: `Generate analyzes Terraform files and outputs a least-privilege
AWS IAM policy JSON that includes all required permissions.

Example:
  tf-iamgen generate ./terraform
  tf-iamgen generate . --output policy.json`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// For now, output a basic policy structure
		// Phase 2 will implement actual policy generation from mappings
		policy := map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Effect":   "Allow",
					"Action":   []string{},
					"Resource": "*",
				},
			},
		}

		// Marshal to JSON
		jsonData, err := json.MarshalIndent(policy, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal policy: %w", err)
		}

		// Output to file or stdout
		if outputFile != "" {
			if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
				return fmt.Errorf("failed to write policy file: %w", err)
			}
			fmt.Printf("Policy generated and saved to: %s\n", outputFile)
		} else {
			fmt.Println(string(jsonData))
		}

		return nil
	},
}

func init() {
	generateCmd.Flags().StringVar(&outputFile, "output", "", "Output file path (default: stdout)")
	generateCmd.Flags().StringP("format", "f", "json", "Output format (json)")
}
