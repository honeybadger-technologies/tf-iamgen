package cmd

import (
	"fmt"
	"os"

	"github.com/honeybadger/tf-iamgen/internal/mapping"
	"github.com/honeybadger/tf-iamgen/internal/parser"
	"github.com/honeybadger/tf-iamgen/internal/policy"
	"github.com/spf13/cobra"
)

var showCoverage bool

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze Terraform files and identify AWS resources",
	Long: `Analyze scans Terraform files in the specified directory
and identifies all AWS resources that will be created.

The analyzer provides:
1. Summary of parsed files and resources
2. Detailed resource listing
3. IAM mapping coverage statistics
4. Preview of required IAM actions

Example:
  tf-iamgen analyze ./terraform
  tf-iamgen analyze . --coverage`,
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

		// Show coverage if requested
		if showCoverage && len(result.Resources) > 0 {
			separator := "=========================================================="
			fmt.Println("\n" + separator)
			fmt.Println("IAM Mapping Coverage Analysis")
			fmt.Println(separator)

			// Load mappings
			db := mapping.NewMappingDatabase()
			if err := db.LoadMappings("mappings"); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Could not load mappings: %v\n", err)
			}
			mappingService := mapping.NewMappingService(db)

			// Create generator for analysis
			opts := policy.PolicyGenerationOptions{
				GroupBy:              "service",
				UseWildcardResources: true,
				IncludeSids:          true,
			}
			generator := policy.NewGenerator(mappingService, opts)

			// Get coverage
			coverage, err := generator.GetPolicyCoverage(result)
			if err == nil {
				if cov, ok := coverage["resources_by_type"].(map[string]int); ok {
					fmt.Printf("\nResource Type Coverage:\n")
					mapped := 0
					unmapped := 0
					for resType := range cov {
						if db.HasMapping(resType) {
							mapped++
							fmt.Printf("  ✓ %s\n", resType)
						} else {
							unmapped++
							fmt.Printf("  ✗ %s (no mapping)\n", resType)
						}
					}
					fmt.Printf("\nCoverage: %d/%d resource types mapped\n", mapped, mapped+unmapped)
				}
			}

			// Show action preview
			pol, metadata, err := generator.GeneratePolicy(result)
			if err == nil && len(pol.Statement) > 0 {
				fmt.Printf("\nGenerated Policy Preview:\n")
				fmt.Printf("  Total Actions: %d\n", metadata.ActionCount)
				fmt.Printf("  Services: %v\n", metadata.Services)
				fmt.Printf("  Statements: %d\n", len(pol.Statement))
			}
		}

		return nil
	},
}

func init() {
	analyzeCmd.Flags().BoolVar(&showCoverage, "coverage", false, "Show IAM mapping coverage analysis")
}
