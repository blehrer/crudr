package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type RootScreenModel struct {
	Form  *huh.Form
	Value *string
}

func NewRootScreenModel() RootScreenModel {
	var value string
	return RootScreenModel{
		Form:  V3PointOfEntry(value),
		Value: &value,
	}
}

func (m RootScreenModel) Init() tea.Cmd {
	return m.Form.Init()
}

func (m RootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.Form.Update(msg)
}

func (m RootScreenModel) View() string {
	return m.Form.View()
}

func V3PointOfEntry(selection string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Main").
				Key("root").
				Options(huh.NewOptions(sections...)...).
				Value(&selection),
		))
}

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
