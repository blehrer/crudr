package cmd

import (
	"crudr/util"

	"github.com/charmbracelet/huh"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type pathitemChange struct {
	target *v3.PathItem
	change v3.PathItem
}

type pathsModel struct {
	filename     string
	oaModel      *libopenapi.DocumentModel[v3.Document]
	selectedPath string
	description string
}

func newPathsModel(filename string) pathsModel {
	oaModel, errs := util.ReadSpecFile(filename)
	if len(errs) > 0 {
		panic(errs)
	}
	return pathsModel{
		filename: filename,
		oaModel:  oaModel,
	}
}

func (pm pathsModel) keys() []string {
	pathItemsSeq := pm.oaModel.Model.Paths.PathItems
	keyset := make([]string, pathItemsSeq.Len())
	for k := range pathItemsSeq.FromOldest() {
		keyset = append(keyset, k)
	}
	return keyset
}

var pm = newPathsModel("openapi.yaml")
var changeset []pathitemChange

var form = huh.NewForm(
	huh.NewGroup(
		huh.NewSelect[string]().
			Key("path").
			Value(&pm.selectedPath).
			Options(huh.NewOptions(pm.keys()...)...)).
		Description("Choose an API path to edit"),

	// Fields of a given path
	// - Route name
	huh.NewGroup(
		huh.NewInput().
			Description("Route").
			Placeholder("/").
			Validate(func(s string) error {
				if s[0] != '/' {
					return *new(error)
				}
				return nil
			}).
			Key("route-name").
		Value(pm.description),
			// - Route description
		huh.NewInput().
		Description("Description of the route").
		Key("route-description").
		Value(),
		).
		Description("route"),

	huh.NewGroup(
		Key
		).
		Description("Route"),

	huh.NewGroup().
		Description("Methods"),
)
