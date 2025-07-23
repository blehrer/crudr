package main

import (
	// rw "crudr/io"
	"crudr/ui"
	"crudr/ui/constants"
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	setupLogging()
	startTea()
	defer os.Exit(0)
}

func setupLogging() {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
}

func startTea() error {
	file := flag.String("f", "./openapi.yaml", "Path to the OpenAPI spec file to work on.")
	flag.Parse()
	m := ui.NewTuiModel(*file)
	constants.P = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := constants.P.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return nil
}
