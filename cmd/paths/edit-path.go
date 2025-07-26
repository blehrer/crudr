package cmd

import (
	"crudr/util"
	"errors"
	"fmt"
	"time"

	// "github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/spf13/cobra"
)

type patheditorModel struct {
	filename    string
	oaModel     *libopenapi.DocumentModel[v3.Document]
	selectedKey string
	renderNext  string
}

type pathItem struct {
	id    string
	entry *v3.PathItem
}

func (pi pathItem) String() string {
	return pi.id
}

func pathItemKey(item *v3.PathItem) string {
	return item.GoLow().KeyNode.Value
}

func (m pathItem) FilterValue() string {
	return m.id
}

// getPathItems returns the keyset and values, in matching order, separated into slices.
func getPathItems(oam *libopenapi.DocumentModel[v3.Document]) ([]string, []*v3.PathItem) {
	entries := oam.Model.Paths.PathItems
	keys := make([]string, entries.Len())
	values := make([]*v3.PathItem, entries.Len())
	for k, v := range entries.FromNewest() {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

var model patheditorModel

var EditPathCmd = &cobra.Command{
	Use: "Edit a path in your api",
	ValidArgs: []cobra.Completion{
		cobra.CompletionWithDesc("filename", "The root spec file (openapi.yaml)"),
		cobra.CompletionWithDesc("path", "The API path (/path/to/endpoint)"),
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var errs []error
		model.filename = args[0]
		model.oaModel, errs = util.ReadSpecFile(model.filename)
		if len(errs) == 0 {
			if model.oaModel == nil {
				fmt.Print("oaModel is nil")
			}
			// keys, _ := getPathItems(model.oaModel)
			// form := huh.NewForm(
			// 	huh.NewGroup(
			// 		huh.NewSelect[string]().
			// 			Title("Paths").
			// 			Key("paths").
			// 			Options(huh.NewOptions(keys...)...).
			// 			Value(&model.selectedKey)))
			err := spinner.New().Title(fmt.Sprintf("picked %s", "ppoooooopp")).Action(func() {
				time.Sleep(time.Duration(4 * time.Second))
			}).Run()
			if err != nil {
				fmt.Printf("%v", err)
			}
			// p := tea.NewProgram(
			// 	form,
			// 	tea.WithFilter(func(model tea.Model, msg tea.Msg) tea.Msg {
			// 		switch msg := msg.(type) {
			// 		case tea.KeyMsg:
			// 			switch msg.String() {
			// 			case "ctrl+c":
			// 				return tea.Quit()
			// 			}
			// 		}
			// 		return msg
			// 	}))
			// _, err := p.Run()
			// if err != nil {
			// 	errs = append(errs, err)
			// }
			// spinner.New().Title(fmt.Sprintf("picked %s", form.GetString("paths"))).Run()
		}
		return errors.Join(errs...)
	},
}

func init() {

}
