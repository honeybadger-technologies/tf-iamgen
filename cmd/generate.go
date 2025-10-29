package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/honeybadger/tf-iamgen/internal/mapping"
	"github.com/honeybadger/tf-iamgen/internal/parser"
	"github.com/honeybadger/tf-iamgen/internal/policy"
)

var (
	outputFile   string
	outputFormat string
	groupBy      string
)

var generateCmd = &cobra.Command{
	Use:   "generate [path]",
	Short: "Generate IAM policy from Terraform files",
	Long: `Generate analyzes Terraform files and outputs a least-privilege
AWS IAM policy JSON that includes all required permissions.

The generator:
1. Parses Terraform files in the specified directory
2. Identifies all AWS resources
3. Maps resources to required IAM actions
4. Generates a least-privilege policy document

Example:
  tf-iamgen generate ./terraform
  tf-iamgen generate . --output policy.json
  tf-iamgen generate . --format json --group-by service`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		terraformPath := args[0]

		// Step 1: Parse Terraform files
		tfParser := parser.NewTerraformParser()
		parseResult, err := tfParser.ParseDirectory(terraformPath)
		if err != nil {
			return fmt.Errorf("failed to parse Terraform files: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Found %d resources in %s\n", len(parseResult.Resources), terraformPath)

		// Step 2: Load IAM mappings
		db := mapping.NewMappingDatabase()
		if err := db.LoadMappings("mappings"); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not load all mappings: %v\n", err)
		}
		mappingService := mapping.NewMappingService(db)

		// Step 3: Generate policy
		opts := policy.PolicyGenerationOptions{
			GroupBy:              groupBy,
			UseWildcardResources: true,
			IncludeSids:          true,
			Minimize:             false,
		}
		generator := policy.NewGenerator(mappingService, opts)

		pol, metadata, err := generator.GeneratePolicy(parseResult)
		if err != nil {
			return fmt.Errorf("failed to generate policy: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Generated policy for %d resources (%d unique actions)\n",
			metadata.ResourceCount, metadata.ActionCount)

		// Step 4: Validate policy
		warnings, _ := generator.ValidatePolicy(pol)
		for _, warning := range warnings {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
		}

		// Step 5: Format and output
		var policyOutput string
		switch outputFormat {
		case "json":
			policyOutput, err = pol.ToJSON()
		default:
			return fmt.Errorf("unsupported output format: %s", outputFormat)
		}

		if err != nil {
			return fmt.Errorf("failed to format policy: %w", err)
		}

		// Output to file or stdout
		if outputFile != "" {
			if err := os.WriteFile(outputFile, []byte(policyOutput), 0644); err != nil {
				return fmt.Errorf("failed to write policy file: %w", err)
			}
			fmt.Printf("Policy generated and saved to: %s\n", outputFile)
		} else {
			fmt.Println(policyOutput)
		}

		return nil
	},
}

func init() {
	generateCmd.Flags().StringVar(&outputFile, "output", "", "Output file path (default: stdout)")
	generateCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format: json (default: json)")
	generateCmd.Flags().StringVar(&groupBy, "group-by", "flat", "Group statements by: service, resource, or flat (default: flat)")
}
