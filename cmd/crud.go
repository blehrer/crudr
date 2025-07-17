package cmd

import (
	forms "crudr/cmd/internal/forms"
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/spf13/cobra"
)


func createEndpoint(document libopenapi.Document, model libopenapi.DocumentModel[v3.Document], filepath string) {
	var path string
	var newSummary, newOperationId, newDescription string

	httpMethod, err := forms.HttpMethods(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
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
	switch httpMethod {
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

	write(document, filepath)
	fmt.Println("Endpoint created successfully!")
}

func write(document libopenapi.Document, filepath string) {
	// Save the updated model to openapi.yaml
	rendered, err := document.Render()
	if err != nil {
		fmt.Printf("Error rendering document: %s\n", err)
		return
	}
	err = os.WriteFile(filepath, rendered, 0644)
	if err != nil {
		fmt.Printf("Error writing to %s: %s\n", filepath, err)
		return
	}
}

func openSpec(filepath string) (libopenapi.Document, libopenapi.DocumentModel[v3.Document], []error) {
	var errors []error = make([]error, 0, 3)
	openapiBytes, err := os.ReadFile(filepath)
	if err != nil {
		errors = append(errors, err)
	}

	document, err := libopenapi.NewDocument(openapiBytes)
	if err != nil {
		errors = append(errors, err)
	}

	model, errs := document.BuildV3Model()
	errors = append(errors, errs...)
	return document, *model, errors
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

var crudCmd = &cobra.Command{
	Use:       `crud`,
	Short:     `CRUD an OpenAPI spec`,
	Long:      `This command opens a wizard to perform CRUD operations on an OpenAPI spec.`,
	Args:      cobra.RangeArgs(0, 1),
	Example:   " crudr crud\n crudr crud path/to/spec.yaml",
	ValidArgs: []cobra.Completion{cobra.CompletionWithDesc(`filename`, `defaults to ./openapi.yaml`)},
	Run: func(cmd *cobra.Command, args []string) {
		// Parse args for path to openapi yaml spec
		var filepath = "openapi.yaml"
		if len(args) == 1 {
			filepath = args[0]
		}

		// Construct the document and model from the spec
		document, model, errs := openSpec(filepath)
		if len(errs) > 0 {
			fmt.Printf("Could not open the spec '%s'. Errors: %v\n", filepath, errs)
			return
		}

		// Form: Select endpoint to edit, setup
		paths := make([]string, 0, (model.Model.Paths.PathItems.Len()))
		for k := range model.Model.Paths.PathItems.OrderedMap.KeysFromNewest() {
			paths = append(paths, k)
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

		// Form: Select endpoint to edit, render
		var err = form.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Form: Select endpoint to edit, handle
		if selectedEndpoint == "+ New Endpoint" {
			createEndpoint(document, model, filepath)
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
			write(document, filepath)

			fmt.Println("Endpoint updated successfully!")
		}
	},
}

func init() {
	rootCmd.AddCommand(crudCmd)
}
