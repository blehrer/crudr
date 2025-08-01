package util

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"log"
	"os"
	"reflect"
)

func ReadSpecFile(filePath string) *libopenapi.DocumentModel[v3.Document] {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	config := datamodel.NewDocumentConfiguration()
	config.AllowFileReferences = true

	doc, err := libopenapi.NewDocumentWithConfiguration(fileContent, config)
	if err != nil {
		log.Fatalf("Failed to create new document: %v", err)
	}

	// Build a v3 model from the document.
	v3Model, errors := doc.BuildV3Model()
	if len(errors) > 0 {
		for i := range errors {
			log.Printf("error: %e\n", errors[i])
		}
		log.Fatalf("cannot create v3 model from document: %d errors reported", len(errors))
	}

	return v3Model
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

//	func KindToForm(kind reflect.Kind) func(s string) func(s string) (huh.Group, *string) {
//		var value string
//		input := switch kind {
//		case reflect.Bool:
//			input = huh.NewGroup(huh.NewSelect[string]().
//				Title(f.Name).
//				Value(&value).
//				Options([]huh.Option[string]{
//					huh.NewOption("true", "true"),
//					huh.NewOption("false", "false")}...))
//			break
//		}
//		return func(s string){return input, *value}
//	}
func TypeToForm(t reflect.Type) map[string]func(s string) huh.Group {
	submaps := []map[string]func(s string) huh.Group{}
	var rv = map[string]func(s string) huh.Group{}
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			rv[f.Name] = func(s string) func(s string) huh.Group {
				var value string
				var input *huh.Group
				typename := f.Type.String()
				switch typename {
				case "string":
				case "float":
				case "int":
					input = huh.NewGroup(huh.NewInput().
						Title(f.Name).
						Value(&value))
					break
				case "bool":
					input = huh.NewGroup(huh.NewSelect[string]().
						Title(f.Name).
						Value(&value).
						Options([]huh.Option[string]{
							huh.NewOption("true", "true"),
							huh.NewOption("false", "false")}...))
					break
				default:
					submaps = append(submaps, TypeToForm(f.Type))
				}
				return func(ss string) huh.Group { return *input }
			}(f.Type.String())
		}
	} else {
	}
	for _, submap := range submaps {
		for k, v := range submap {
			rv[k] = v
		}
	}
	return rv
}
func ToForm(datastructure any) map[string]func(s string) huh.Group {
	t := reflect.TypeOf(datastructure)
	return TypeToForm(t)
}
