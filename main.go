package main

import (
	"crudr/ui"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	file := flag.String("f", "./openapi.yaml", "Path to the OpenAPI spec file to work on.")
	flag.Parse()
	p := tea.NewProgram(ui.NewCrudrState(*file))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
