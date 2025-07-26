package util

import (
	"fmt"
	"os"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func ReadSpecFile(filepath string) (*libopenapi.DocumentModel[v3.Document], []error) {
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
	return model, errors
}

func WriteSpecFile(model libopenapi.DocumentModel[v3.Document], filepath string) {
	rendered, err := model.Model.Render()
	if err != nil {
		fmt.Printf("Error rendering model: %s\n", err)
		return
	}
	err = os.WriteFile(filepath, rendered, 0644)
	if err != nil {
		fmt.Printf("Error writing to %s: %s\n", filepath, err)
		return
	}
}
