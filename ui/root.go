package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Implements SceneContext
type RootScene struct {
	CrudrState CrudrState
	Model      RootSceneModel
}

// Implements tea.Model
type RootSceneModel struct {
	Form  *huh.Form
	Value *string
}

func (s RootScene) NewScene(cs CrudrState) RootScene {
	return RootScene{
		CrudrState: cs,
		Model:      NewRootSceneModel(),
	}
}

func (s RootScene) FromSceneContext(sc SceneContext) RootScene {
	return RootScene{
		CrudrState: sc.GetCrudrState(),
		Model:      NewRootSceneModel(),
	}
}

func NewRootSceneModel() RootSceneModel {
	var value string
	return RootSceneModel{
		Form:  V3PointOfEntry(value),
		Value: &value,
	}
}

func (s RootScene) Init() tea.Cmd {
	return s.Model.Form.Init()
}

func (s RootScene) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if s.Model.Form.State == huh.StateCompleted {
		// update routes, spawn scenes, etc
	}
	return s.Model.Form.Update(msg)
}

func (s RootScene) View() string {
	return s.Model.Form.View()
}

func (s RootScene) GetCrudrState() CrudrState {
	return s.CrudrState
}

func (s RootScene) GetSceneModel() RootSceneModel {
	return s.Model
}

func V3PointOfEntry(selection string) *huh.Form {
	var sections = []string{
		"openapi",
		"info",
		"jsonSchemaDialect",
		"servers",
		"paths",
		"webhooks",
		"components",
		"security",
		"tags",
		"externalDocs",
	}
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Main").
				Key("selection").
				Options(huh.NewOptions(sections...)...).
				Value(&selection),
		))
}
