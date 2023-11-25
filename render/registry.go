package render

import (
	"context"
	"log"

	"github.com/angel-cdo/dotui/config"

	"github.com/derailed/tcell/v2"
	"github.com/digitalocean/godo"
)

type Registry struct {
}

func (r *Registry) Render(global_context *config.GlobalContext, base_render *BaseRender) {
	drawTable(global_context, *base_render)
	global_context.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, column := global_context.Table.GetSelection()
		if event.Key() == tcell.KeyEnter {
			cell := global_context.Table.GetCell(row, column)
			context := cell.Text
			global_context.Context = context
			*base_render = &RepositoryImage{}
			global_context.Table.Clear()
			global_context.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				return event
			})
			(*base_render).Render(global_context, base_render)
		}
		return event
	})
}

func (r *Registry) GetHeaders() []string {
	headers := []string{"NAME", "REGION"}
	return headers
}

func (r *Registry) GetData(global_context *config.GlobalContext) [][]string {
	page := 1
	list_value := godo.ListOptions{Page: page, PerPage: 100, WithProjects: false}
	registry, _, _ := global_context.Client.Registry.Get(context.Background())
	repositories_list := []*godo.Repository{}
	repositories, _, err := global_context.Client.Registry.ListRepositories(context.Background(), registry.Name, &list_value)

	for len(repositories) > 0 {
		if err != nil {
			log.Fatal(err)
		}
		repositories_list = append(repositories_list, repositories...)
		page += 1
		list_value.Page = page
		repositories, _, err = global_context.Client.Registry.ListRepositories(context.Background(), registry.Name, &list_value)
	}
	ret := make([][]string, len(repositories_list))
	for i, registery := range repositories_list {
		ret[i] = make([]string, 2)
		ret[i][0] = registery.Name
		ret[i][1] = registery.RegistryName
	}
	return ret
}

func (r *Registry) GetName() string {
	return "registry"
}
