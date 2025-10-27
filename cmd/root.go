package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tf-iamgen",
	Short: "Automatically generate AWS IAM policies from Terraform",
	Long: `tf-iamgen is a tool that analyzes your Terraform infrastructure
and automatically generates the minimum required AWS IAM permissions.

It helps DevOps and security teams achieve least-privilege IAM configurations
without manual policy writing or guesswork.

Example:
  tf-iamgen analyze ./terraform
  tf-iamgen generate --output policy.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(versionCmd)
}
