package ui

import (
	rw "crudr/io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// Application state for CRUDR TUI
type CrudrState struct {
	Model  *libopenapi.DocumentModel[v3.Document] // Stores updates that can be written to the API spec file
	Route  *string                                // Route to current section
	Scene  tea.Model                              // form holder
	Cursor int                                    // Position of the cursor. Nil is no cursor.
}

func NewCrudrState(filepath string) CrudrState {
	model, _ := rw.ReadSpecFile(filepath)
	scene := NewRootSceneModel()
	route := "root"
	return CrudrState{
		Model: &model,
		Route: &route,
		Scene: scene,
	}
}

func (m CrudrState) Init() tea.Cmd {
	return m.Scene.Init()
}

func (m CrudrState) View() string {
	s := strings.Join([]string{
		`
                    __       
.----.----.--.--.--|  |.----.
|  __|   _|  |  |  _  ||   _|
|____|__| |_____|_____||__| CRUDR - openapi tooling for the [c]rudderless
                             
		`,
		m.Scene.View(),
		"\tctrl+c - close",
	}, "\n")

	return s
}

func (m CrudrState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c":
			return m, tea.Quit

		default:
			return m.Scene.Update(msg)
		}
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

type SceneContext interface {
	NewScene(CrudrState) SceneContext
	FromSceneContext(SceneContext) SceneContext
	GetCrudrState() CrudrState
	GetModel() tea.Model
}
