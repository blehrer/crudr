package cmd

import (
	"crudr/pkg/TUI"

	"github.com/spf13/cobra"
)

var listPathsCmd = &cobra.Command{
	Use:   "list-paths [file]",
	Short: "List all paths in an OpenAPI specification",
	Long:  `List all paths in an OpenAPI specification.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := "openapi.yaml"
		if len(args)>0{
			filePath = args[0]
		}
		TUI.ListPaths(filePath)
	},
}

func init() {
	rootCmd.AddCommand(listPathsCmd)
}
