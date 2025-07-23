package ui

import (
	rw "crudr/io"
	"crudr/ui/constants"
	"crypto/rand"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

//go:generate stringer -type sceneState -trimprefix=scene_
type sceneState int

const (
	scene_mainmenu sceneState = iota
	scene_root
	scene_paths
	scene_pathItem
)

type tuiModel struct {
	specfile    string
	oaModel     libopenapi.DocumentModel[v3.Document]
	mainMenu    list.Model
	sceneState  sceneState
	updateAtTop string
	scenes      map[sceneState]tea.Model
	quitting    bool
}

func NewTuiModel(file string) tuiModel {
	oaModel, errs := rw.ReadSpecFile(file)
	if len(errs) > 0 {
		log.Fatalf("%s could not be parsed. %v", file, errs)
		os.Exit(1)
	}
	return tuiModel{
		specfile:    file,
		oaModel:     oaModel,
		mainMenu:    mainMenu(),
		sceneState:  scene_paths,
		quitting:    false,
		updateAtTop: sceneState.String(scene_paths),
		scenes: map[sceneState]tea.Model{
			scene_paths: newPathsModel(&oaModel),
		},
	}
}

func (m tuiModel) Init() tea.Cmd {
	return nil
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.updateAtTop = rand.Text()
	scene := m.scenes[m.sceneState]
	scene.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.mainMenu.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	case sceneChangeMsg:
		m.sceneState = msg.sceneState
	}
	return m, nil
}

func (m tuiModel) View() string {
	if m.quitting {
		return ""
	}
	scene := m.scenes[m.sceneState]
	s := strings.Join([]string{
		"updateAtTop: " + m.updateAtTop,
		constants.ASCIITitle,
		scene.View(),
	}, "\n")
	return s
}

type mainMenuItem struct {
	name       string
	sceneState sceneState
}

func (i mainMenuItem) FilterValue() string {
	return i.name
}

func mainMenu() list.Model {
	menuItems := []list.Item{
		mainMenuItem{name: "paths", sceneState: scene_paths},
	}
	return list.New(menuItems, list.NewDefaultDelegate(), constants.WindowSize.Width, constants.WindowSize.Height)
}

type sceneChangeMsg struct {
	sceneState sceneState
	data       map[string]any
}

type printMsg struct {
	val string
}

func (pm printMsg) FilterValue() string {
	return pm.val
}

func printCmd(s string) func() printMsg {
	return func() printMsg {
		return printMsg{val: s}
	}
}
