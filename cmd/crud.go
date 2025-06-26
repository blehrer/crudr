package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
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

		paths := make([]string, 0, (model.Model.Paths.PathItems.Len()))
		for pair := model.Model.Paths.PathItems.First(); pair != nil; pair = pair.Next() {
			paths = append(paths, pair.Key())
		}
		sort.Strings(paths)

		options := make([]huh.Option[string], len(paths)+1)
		options[0] = huh.NewOption("+ New Endpoint", "+ New Endpoint")
		for i, path := range paths {
			options[i+1] = huh.NewOption(path, path)
		}

		var selectedEndpoint string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select an endpoint to edit").
					Options(options...).Value(&selectedEndpoint)),
		)

		err = form.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if selectedEndpoint == "+ New Endpoint" {
			var path, method string
			var newSummary, newOperationId, newDescription string

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Path").
						Value(&path),
					huh.NewSelect[string]().
						Title("Method").
						Options(
							huh.NewOption("GET", "GET"),
							huh.NewOption("POST", "POST"),
							huh.NewOption("PUT", "PUT"),
							huh.NewOption("DELETE", "DELETE"),
							huh.NewOption("PATCH", "PATCH"),
							huh.NewOption("HEAD", "HEAD"),
							huh.NewOption("OPTIONS", "OPTIONS"),
						).Value(&method)),
			)

			err = form.Run()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			editForm := createEditForm(&newSummary, &newOperationId, &newDescription)
			err = editForm.Run()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			// Create a new path item if it doesn't exist
			pathItem := model.Model.Paths.PathItems.GetOrZero(path)
			if pathItem == nil {
				pathItem = &v3.PathItem{}
				model.Model.Paths.PathItems.Set(path, pathItem)
			}

			// Create a new operation
			operation := &v3.Operation{
				Summary:     newSummary,
				OperationId: newOperationId,
				Description: newDescription,
			}

			switch method {
			case "GET":
				pathItem.Get = operation
			case "POST":
				pathItem.Post = operation
			case "PUT":
				pathItem.Put = operation
			case "DELETE":
				pathItem.Delete = operation
			case "PATCH":
				pathItem.Patch = operation
			case "HEAD":
				pathItem.Head = operation
			case "OPTIONS":
				pathItem.Options = operation
			}

			// Save the updated model to openapi.yaml
			rendered, err := document.Render()
			if err != nil {
				fmt.Printf("Error rendering document: %s\n", err)
				return
			}

			err = os.WriteFile("openapi.yaml", rendered, 0644)
			if err != nil {
				fmt.Printf("Error writing to openapi.yaml: %s\n", err)
				return
			}

			fmt.Println("Endpoint created successfully!")
		} else {
			pathItem := model.Model.Paths.PathItems.GetOrZero(selectedEndpoint)

			operations := make(map[string]*v3.Operation)
			if pathItem.Get != nil {
				operations["GET"] = pathItem.Get
			}
			if pathItem.Post != nil {
				operations["POST"] = pathItem.Post
			}
			if pathItem.Put != nil {
				operations["PUT"] = pathItem.Put
			}
			if pathItem.Delete != nil {
				operations["DELETE"] = pathItem.Delete
			}
			if pathItem.Patch != nil {
				operations["PATCH"] = pathItem.Patch
			}
			if pathItem.Head != nil {
				operations["HEAD"] = pathItem.Head
			}
			if pathItem.Options != nil {
				operations["OPTIONS"] = pathItem.Options
			}

			if len(operations) == 0 {
				fmt.Println("No operations found for this endpoint.")
				return
			}

			methodOptions := make([]huh.Option[string], 0, len(operations))
			for method := range operations {
				methodOptions = append(methodOptions, huh.NewOption(method, method))
			}

			var selectedMethod string
			methodSelectForm := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Select an HTTP method to edit").
						Options(methodOptions...).Value(&selectedMethod),
				),
			)

			err = methodSelectForm.Run()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			operation := operations[selectedMethod]
			// Initialize local variables with existing operation values
			currentSummary := operation.Summary
			currentOperationId := operation.OperationId
			currentDescription := operation.Description

			editForm := createEditForm(&currentSummary, &currentOperationId, &currentDescription)

			err = editForm.Run()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			// Update the operation object with the (potentially modified) values from the form
			operation.Summary = currentSummary
			operation.OperationId = currentOperationId
			operation.Description = currentDescription

			// Save the updated model to openapi.yaml
			rendered, err := document.Render()
			if err != nil {
				fmt.Printf("Error rendering document: %s\n", err)
				return
			}

			err = os.WriteFile("openapi.yaml", rendered, 0644)
			if err != nil {
				fmt.Printf("Error writing to openapi.yaml: %s\n", err)
				return
			}

			fmt.Println("Endpoint updated successfully!")
		}
	},
}

func init() {
	rootCmd.AddCommand(crudCmd)
}

func createEditForm(summary, operationId, description *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Summary").
				Value(summary),
			huh.NewInput().
				Title("Operation ID").
				Value(operationId),
			huh.NewText().
				Title("Description").
				Value(description),
		),
	)
}
