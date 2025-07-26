package cmd

import (
	patheditor "crudr/cmd/edit-path"
	"fmt"
	"os"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/spf13/cobra"
)

type CrudrModel struct {
	spec    string
	oaModel *libopenapi.DocumentModel[v3.Document]
}

var cm CrudrModel

var rootCmd = &cobra.Command{
	Use:        "crudr",
	Short:      "CRUDr is a TUI for pb33f/libopenapi",
	Args:       cobra.NoArgs,
	ArgAliases: []string{"file", "spec"},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(patheditor.EditPathCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
