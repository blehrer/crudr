package ui

import (
	"crudr/ui/constants"
	"crypto/rand"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type pathsModel struct {
	printNext string
	oaModel   *libopenapi.DocumentModel[v3.Document]
	pathsList list.Model
	selected  pathItem
}

type pathItem struct {
	key   string
	value *v3.PathItem
}

func (pi pathItem) FilterValue() string {
	return pi.key
}

func newPathsModel(oaModel *libopenapi.DocumentModel[v3.Document]) pathsModel {
	items := []list.Item{}
	for k, v := range oaModel.Model.Paths.PathItems.FromNewest() {
		items = append(items, pathItem{key: k, value: v})
	}
	ls := list.New(
		items,
		list.NewDefaultDelegate(),
		constants.WindowSize.Width,
		constants.WindowSize.Height)
	return pathsModel{
		printNext: "",
		oaModel:   oaModel,
		pathsList: ls,
	}
}

func (m pathsModel) Init() tea.Cmd {
	return nil
}

func responseCmd(pathName string, pathItem *v3.PathItem) func() tea.Msg {
	return func() tea.Msg {
		return sceneChangeMsg{
			sceneState: scene_pathItem,
			data: map[string]any{
				"pathName": pathName,
				"pathItem": &pathItem,
			},
		}
	}
}

func (m pathsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.printNext = rand.Text()
	// switch msg := msg.(type) {
	//
	// case tea.KeyMsg:
	// 	if slices.Contains(constants.Keymap.Create.Keys(), msg.String()) {
	// 		pathName := m.pathsList.SelectedItem().FilterValue()
	// 		oa3pathItem, ok := m.oaModel.Model.Components.PathItems.Get(pathName)
	// 		if ok {
	// 			m.selected = pathItem{
	// 				key:   pathName,
	// 				value: oa3pathItem,
	// 			}
	// 			return m, responseCmd(pathName, oa3pathItem)
	// 		}
	// 	}
	// }

	return m, nil
}

func (m pathsModel) View() string {
	s := strings.Join([]string{
		m.pathsList.View(),
		"printNext: " + m.printNext,
	}, "\n")

	return s
}
