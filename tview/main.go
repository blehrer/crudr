package main

import (
	"fmt"

	"github.com/rivo/tview"
)

const pageCount = endCredits

//go:generate stringer -type workflow
type workflow int

const (
	mainMenu workflow = iota + 1
	editSections
	endCredits
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()
	for page := range endCredits {
		func(page workflow) {
			pages.AddPage(fmt.Sprintf("page-%d", page),
				tview.NewModal().
					SetText(fmt.Sprintf("This is page %s. Choose where to go next.", page.String())).
					AddButtons([]string{"Next", "Quit"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonIndex == 0 {
							pages.SwitchToPage(fmt.Sprintf("page-%d", (page+1)%pageCount))
						} else {
							app.Stop()
						}
					}),
				false,
				page == 0)
		}(page)
	}
	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
