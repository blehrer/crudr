package cmd

import (
	"fmt"
	"os"

	"github.com/pb33f/libopenapi"
	"github.com/spf13/cobra"
)

var crudCmd = &cobra.Command{
	Use:   "crud",
	Short: "CRUD an OpenAPI spec",
	Long:  `This command opens a wizard to perform CRUD operations on an OpenAPI spec.`,
	Run: func(cmd *cobra.Command, args []string) {
		openapiBytes, err := os.ReadFile("openapi.yaml")
		if err != nil {
			fmt.Println("No openapi.yaml found. Please run 'gocrudr init' first.")
			return
		}

		document, err := libopenapi.NewDocument(openapiBytes)
		if err != nil {
			fmt.Printf("Error parsing openapi.yaml: %s\n", err)

			return
		}

		model, _ := document.BuildV3Model()

		fmt.Println("Available Endpoints:")
		for _, pathItem := range model.Model.Paths.PathItems.FromNewest() {
			if pathItem.Get != nil {
				fmt.Printf("- GET %s ", pathItem.Get.Summary)
			}
			if pathItem.Post != nil {
				fmt.Printf("- POST %s ", pathItem.Get.Summary)
			}
			if pathItem.Put != nil {
				fmt.Printf("- PUT %s ", pathItem.Get.Summary)
			}
			if pathItem.Delete != nil {
				fmt.Printf("- DELETE %s ", pathItem.Get.Summary)
			}
			if pathItem.Patch != nil {
				fmt.Printf("- PATCH %s ", pathItem.Get.Summary)
			}
			if pathItem.Head != nil {
				fmt.Printf("- HEAD %s ", pathItem.Get.Summary)
			}
			if pathItem.Options != nil {
				fmt.Printf("- OPTIONS %s ", pathItem.Get.Summary)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(crudCmd)
}
