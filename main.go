package main

import (
	rw "crudr/io"
	"crudr/ui"
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// TODO: make CrudrState.Scene into a SceneContext.... this will require a bit of separation so as
// to not create a stack-overflow.... or maybe just use a reference? This will allow us to update
// CrudrState.Scene in the Update() method of any scene's model
type SceneContext struct {
	CrudrState CrudrState
	SceneModel tea.Model
}

// Application state for CRUDR TUI
type CrudrState struct {
	Model  *libopenapi.DocumentModel[v3.Document] // Stores updates that can be written to the API spec file
	Route  *string                                // Route to current section
	Scene  tea.Model                              // form holder
	Cursor int                                    // Position of the cursor. Nil is no cursor.
}

func initalCrudrState(filepath string) CrudrState {
	model, _ := rw.ReadSpecFile(filepath)
	scene := ui.NewRootScreenModel()
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

func main() {
	file := flag.String("f", "./openapi.yaml", "Path to the OpenAPI spec file to work on.")
	flag.Parse()
	p := tea.NewProgram(initalCrudrState(*file))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
