package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the version and build information for tf-iamgen.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Version info is injected at build time from main.go
		fmt.Printf("tf-iamgen version %s\n", "dev")
		fmt.Printf("Built at: %s\n", "unknown")
		fmt.Printf("Git commit: %s\n", "unknown")
	},
}
