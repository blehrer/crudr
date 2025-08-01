package cmd

import (
	"fmt"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/spf13/cobra"
	"log"
	"reflect"
)

var completions = func() []cobra.Completion {
	t := reflect.TypeFor[v3.PathItem]()
	cs := make([]cobra.Completion, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		cs = append(
			cs,
			cobra.CompletionWithDesc(
				f.Name,
				fmt.Sprintf("%s (%s)", f.Name, f.Type.Name())))
	}
	return cs
}()

var editPathCmd = &cobra.Command{
	Use:   "edit-path",
	Short: "edit the HTTP methods, tags, etc of a path",
	Long:  `edit the HTTP methods, tags, etc of a path`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var editWhat string

func init() {
	editPathCmd.Flags().StringVarP(&editWhat, "attr", "a", "", "what to edit")
	t := reflect.TypeFor[v3.PathItem]()
	completions := make([]cobra.Completion, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		completions = append(
			completions,
			cobra.CompletionWithDesc(
				f.Name,
				fmt.Sprintf("%s (%s)", f.Name, f.Type.Name())))
	}
	err := editPathCmd.RegisterFlagCompletionFunc(
		"attr",
		cobra.FixedCompletions(completions, cobra.ShellCompDirectiveNoFileComp))
	if err != nil {
		log.Fatal(err.Error())
	}
	rootCmd.AddCommand(editPathCmd)
}
