package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// {{{ V3 Point of entry

func V3PointOfEntry(selection string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Main").
				Key("root").
				Filtering(true).
				Options(huh.NewOptions(Constants.Sections...)...).
				Value(&selection),
		))
}

var V3PointOfEntry_ = huh.NewForm(
	huh.NewGroup(
		huh.NewInput().
			Title("OpenAPI spec version").
			Key("openapi").
			Description("The version number of the OpenAPI Specification that the OpenAPI document uses.").
			Suggestions([]string{"3.1.1"}).
			Validate(func(s string) error {
				_, err := regexp.Match("\\d+\\.\\d+\\.\\d+", []byte(s))
				return err
			})),
	huh.NewGroup(
		huh.NewInput().
			Title("info").
			Key("info").
			Description("Provides metadata about the API.")),
	huh.NewGroup(
		huh.NewInput().
			Title("jsonSchemaDialect").
			Key("jsonSchemaDialect").
			Description("The default value for the $schema keyword within Schema Objects contained within this OAS document.")),
	huh.NewGroup(
		huh.NewInput().
			Title("servers").
			Key("servers").
			Description("An array of Server Objects, which provide connectivity information to a target server.")),
	huh.NewGroup(
		huh.NewInput().
			Title("paths").
			Key("paths").
			Description("The available paths and operations for the API.")),
	huh.NewGroup(
		huh.NewInput().
			Title("webhooks").
			Key("webhooks").
			Description("The incoming webhooks that MAY be received as part of this API and that the API consumer MAY choose to implement.")),
	huh.NewGroup(
		huh.NewInput().
			Title("components").
			Key("components").
			Description("An element to hold various schemas for the document.")),
	huh.NewGroup(
		huh.NewInput().
			Title("security").
			Key("security").
			Description("A declaration of which security mechanisms can be used across the API.")),
	huh.NewGroup(
		huh.NewInput().
			Title("tags").
			Key("tags").
			Description("A list of tags used by the document with additional metadata.")),
	huh.NewGroup(
		huh.NewInput().
			Title("externalDocs").
			Key("externalDocs").
			Description("Additional external documentation.")))

// }}}

// {{{ Endpoints

func EndpointsForm(model *libopenapi.DocumentModel[v3.Document]) *huh.Form {
	var endpoint string
	paths := make([]huh.Option[string], 0, model.Model.Paths.PathItems.Len()+1)
	paths = append(paths, huh.NewOption("+ New Endpoint", Constants.NewEndpoint))
	for path := range model.Model.Paths.PathItems.OrderedMap.KeysFromNewest() {
		paths = append(paths, huh.NewOption(path, path))
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Path").
				Value(&endpoint),
		))
}

func SelectEndpoint(model *libopenapi.DocumentModel[v3.Document]) (string, *libopenapi.DocumentModel[v3.Document], error) {
	var endpoint string
	err := EndpointsForm(model).Run()
	if endpoint == Constants.NewEndpoint {
		return NewEndpoint(model)
	}
	return endpoint, model, err
}

func NewEndpoint(model *libopenapi.DocumentModel[v3.Document]) (string, *libopenapi.DocumentModel[v3.Document], error) {
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

func DeleteEndpoints(model libopenapi.DocumentModel[v3.Document]) error {
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

// }}}

// {{{Http Methods
func HttpMethods(endpointPath string) (string, error) {
	var method string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Path").
				Value(&endpointPath),
			huh.NewSelect[string]().
				Title("Method").
				Options(huh.NewOptions(Constants.HttpMethods...)...).
				Value(&method)))
	err := form.Run()
	return method, err
}

// }}}

type constants struct {
	Sections    []string
	NewEndpoint string
	HttpMethods []string
}

var Constants = constants{
	Sections: []string{
		"OpenAPI spec version",
		"info",
		"jsonSchemaDialect",
		"servers",
		"paths", "webhooks", "components", "security", "tags", "externalDocs"},
	NewEndpoint: "+Endpoint",

	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Methods
	HttpMethods: []string{"GET", "HEAD", "OPTIONS", "TRACE", "PUT", "DELETE", "POST", "PATCH", "CONNECT"},
}

// vim: foldmethod=marker
