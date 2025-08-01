package TUI

import (
	"crudr/util"
	"fmt"
	"io/fs"
	"log"
	"os"
	"reflect"

	"github.com/charmbracelet/huh"
)

func ListPaths(filePath string) {
	// Read the OpenAPI specification file
	if filePath == "" {
		f := os.DirFS(".")
		result, err := fs.Glob(f, "openapi.[json|yaml|yml]")
		if err != nil {
			log.Fatal(err)
		}
		filePath = result[0]
	}
	v3Model := util.ReadSpecFile(filePath)

	// Get all the paths from the specification
	paths := v3Model.Model.Paths.PathItems

	// Create a slice of huh.Option for each path
	var huhPaths []huh.Option[string]
	if paths != nil {
		for pair := paths.First(); pair != nil; pair = pair.Next() {
			path := pair.Key()
			huhPaths = append(huhPaths, huh.NewOption(path, path))
		}
	}

	var selection string
	// Create a new huh.Form with a select field for the paths
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a path").
				Options(huhPaths...).
				Value(&selection),
		),
	)

	// Run the form
	err := form.Run()
	if err != nil {
		log.Fatalf("Failed to run form: %v", err)
	}
	if selection == "" {
		log.Fatalf("Failed to find any path")
	} else {
		fmt.Println("you selected", selection)
	}

	pathItem := paths.Value(selection).Get

	var pathItemFields []huh.Option[string]
	el := reflect.TypeOf(*pathItem)
	for i := 0; i < el.NumField(); i++ {
		f := el.Field(i)
		pathItemFields = append(pathItemFields, huh.NewOption(f.Name, f.Type.String()))
	}

	var editWhat string
	pathEditForm := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title(fmt.Sprintf("What do you want to do with the path %s?", selection)).
			Options(pathItemFields...).Value(&editWhat)))

	err = pathEditForm.Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("chose to edit", editWhat)

	//var fieldValue string
	//reflect.ValueOf(pathItem).FieldByName(editWhat).Set(reflect.ValueOf(&fieldValue))
}
