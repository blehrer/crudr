package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	forms "crudr/cmd/internal/forms"

	tea "github.com/charmbracelet/bubbletea"
	huh "github.com/charmbracelet/huh"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/spf13/cobra"
)

// Application state for CRUDR TUI
type CrudrState struct {
	Model       *libopenapi.DocumentModel[v3.Document] `json:"model"` // Stores updates that can be written to the API spec file
	Route       *string                                `json:"route"` // Route to current section
	Screen      *tea.Model                             `json:"screen"`
	CurrentForm *huh.Form
	NextForm    *huh.Form
	Cursor      int `json:"cursor"` // Position of the cursor. Nil is no cursor.
}

func initalCrudrState(filepath string) CrudrState {
	model, _ := openSpec(filepath)
	var route string
	return CrudrState{
		Model:       &model,
		Route:       &route,
		CurrentForm: forms.V3PointOfEntry(route),
	}
}

func (m CrudrState) Init() tea.Cmd {
	return m.CurrentForm.Init()
}

func (m CrudrState) View() string {
	s := strings.Join([]string{
		`
                    __       
.----.----.--.--.--|  |.----.
|  __|   _|  |  |  _  ||   _|
|____|__| |_____|_____||__| CRUDR - openapi tooling for the [c]rudderless
                             
		`,
		"%s", // placeholder for m.Screen?
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
			return m.CurrentForm.Update(msg)
		}
	}

	switch *m.Route {
	case "paths":
		m.CurrentForm = forms.EndpointsForm(m.Model)
		return m.CurrentForm.Update(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func Crudr() {
	file := flag.String("f", "./openapi.yaml", "Path to the OpenAPI spec file to work on.")
	flag.Parse()
	p := tea.NewProgram(initalCrudrState(*file))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "crudr",
	Short: "A CLI for CRUDR",
	Long:  `A suite of tools for maintaining OpenAPI specs, servers, and clients.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from cobra!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
