package forms

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type V3Model libopenapi.DocumentModel[v3.Document]

func Endpoints(model V3Model) (string, V3Model, error) {
	var endpoint string
	paths := make([]huh.Option[string], 0, model.Model.Paths.PathItems.Len()+1)
	paths = append(paths, huh.NewOption("+ New Endpoint", __.NewEndpoint))
	for path := range model.Model.Paths.PathItems.OrderedMap.KeysFromNewest() {
		paths = append(paths, huh.NewOption(path, path))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Path").
				Value(&endpoint),
		))
	err := form.Run()
	if endpoint == __.NewEndpoint {
		return NewEndpoint(model)
	}
	return endpoint, model, err
}

func NewEndpoint(model V3Model) (string, V3Model, error) {
	var endpoint string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Endpoint path").
				Description("ex `/path/{param}`").
				Value(&endpoint)))
	err := form.Run()
	endpoint = sanitizeEndpoint(endpoint)
	model.Model.Paths.PathItems.Set(endpoint, &v3.PathItem{})
	return endpoint, model, err
}

func DeleteEndpoints(model V3Model) error {
	var paths []huh.Option[string]
	for path := range model.Model.Paths.PathItems.KeysFromNewest() {
		paths = append(paths, huh.NewOption(path, path))
	}
	var selectedEndpoints []string
	var confirmed bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				OptionsFunc(func() []huh.Option[string] { return paths }, nil).
				Value(&selectedEndpoints),
			huh.NewConfirm().
				DescriptionFunc(func() string {
					return fmt.Sprintf("Are you sure you want to delete the following endpoints?:\n %s",
						strings.Join(selectedEndpoints, "\n- "))
				}, nil).
				Value(&confirmed)))
	err := form.Run()
	// TODO: if err, it probably means that the user Ctrl+c'd ... in this case you should go back to state 0
	// TODO: if !confirmed, ask the user to update the multi-select. make sure to preserve the state from previous selection
	// TODO: do something with failed deletions... notify user somehow.
	var failed []huh.Option[string]
	// TODO: put deleted in an undo-history?
	var deleted []*v3.PathItem
	if confirmed {
		for index := range selectedEndpoints {
			deletedPath, ok := model.Model.Paths.PathItems.Delete(paths[index].Key)
			if ok {
				deleted = append(deleted, deletedPath)
			} else {
				failed = append(failed, paths[index])
			}
		}
	}
	return err
}

func sanitizeEndpoint(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	slash := "/"
	if !strings.HasPrefix(endpoint, slash) {
		endpoint = slash + endpoint
	}
	return endpoint
}

func HttpMethods(endpointPath string) (string, error) {
	var method string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Path").
				Value(&endpointPath),
			huh.NewSelect[string]().
				Title("Method").
				Options(huh.NewOptions(__.HttpMethods...)...).
				Value(&method)))
	err := form.Run()
	return method, err
}

type constants struct {
	NewEndpoint string
	HttpMethods []string
}

var __ = constants{
	NewEndpoint: "+Endpoint",

	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Methods
	HttpMethods: []string{"GET", "HEAD", "OPTIONS", "TRACE", "PUT", "DELETE", "POST", "PATCH", "CONNECT"},
}
