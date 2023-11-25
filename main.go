package main

import (
	"context"
	"log"
	"strings"

	"github.com/angel-cdo/dotui/render"

	"github.com/angel-cdo/dotui/config"

	"github.com/derailed/tcell/v2"
	"github.com/derailed/tview"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

var (
	global_context *config.GlobalContext
	app            *tview.Application
	flex           *tview.Flex
	token                            = config.GetDoToken("cornix-prod")
	current        render.BaseRender = &render.Registry{}
	words                            = []string{"database", "droplet", "kubernetes", "context", "registry"}
)

func main() {
	authToken := &oauth2.Token{AccessToken: token}
	oauthClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(authToken))
	global_context = &config.GlobalContext{
		Client:  godo.NewClient(oauthClient),
		Table:   tview.NewTable(),
		Context: "",
	}
	global_context.Table.SetFixed(1, 0).SetSelectable(true, false)
	app = tview.NewApplication()
	flex = tview.NewFlex().SetDirection(tview.FlexRow)
	pages := tview.NewPages()
	modal := tview.NewModal()
	flex.AddItem(global_context.Table, 0, 1, true)
	refreshTable()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == ':' {
			openSerachBox(flex)
			return nil
		} else if event.Rune() == 'q' {
			modal.ClearButtons()
			modal.
				SetText("Do you want to quit the application?").
				AddButtons([]string{"Quit", "Cancel"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Quit" {
						app.Stop()
					} else if buttonLabel == "Cancel" {
						pages.RemovePage("modal")
					}
				})
			pages.AddPage("modal", modal, true, true)
		}
		return event
	})

	pages.AddPage("background", flex, true, true)

	if err := app.SetRoot(pages, true).Run(); err != nil {
		log.Fatal(err)
	}
}

func openSerachBox(flex *tview.Flex) {
	flex.Clear()
	inputField := tview.NewInputField()
	inputField.SetDoneFunc(func(key tcell.Key) {
		text := inputField.GetText()
		if text == "droplet" {
			current = &render.Droplet{}
			// } else if text == "database" {
			// 	current = &render.Database{}
			// } else if text == "kubernetes" {
			// 	current = &render.Kubernetes{}
			// } else if text == "context" {
			// 	current = &render.Context{}
		} else if text == "registry" {
			current = &render.Registry{}
		} else {
			flex.Clear()
			flex.AddItem(global_context.Table, 0, 1, true)
			app.SetFocus(global_context.Table)
			return
		}
		flex.Clear()
		flex.AddItem(global_context.Table, 0, 1, true)
		app.SetFocus(global_context.Table)
		refreshTable()
	})
	inputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}
		for _, word := range words {
			if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
				entries = append(entries, word)
				break
			}
		}
		return
	})
	flex.AddItem(inputField, 3, 1, true)
	flex.AddItem(global_context.Table, 0, 1, false)
	app.SetFocus(inputField)

}

func refreshTable() {
	global_context.Table.Clear()
	global_context.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})
	current.Render(global_context, &current)
}
