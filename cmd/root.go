package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gocrudr",
	Short: "A CLI for CRUDR",
	Long:  `A suite of tools for maintaining OpenAPI specs, servers, and clients.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from cobra!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
