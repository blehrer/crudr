package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new CRUDR project",
	Long:  `This command initializes a new CRUDR project. It checks for an existing OpenAPI specification file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("openapi.yaml"); err == nil {
			fmt.Println("Found openapi.yaml")
			return
		}
		if _, err := os.Stat("openapi.json"); err == nil {
			fmt.Println("Found openapi.json")
			return
		}
		fmt.Println("No OpenAPI spec found. Run the setup wizard here.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
