package render

import (
	"context"
	"fmt"
	"log"

	"github.con/angel-cdo/dotui/config"

	"github.com/digitalocean/godo"
)

type Droplet struct {
}

func (d *Droplet) Render(global_context *config.GlobalContext, base_render *BaseRender) {
	drawTable(global_context, d)
}

func (d *Droplet) GetHeaders() []string {
	headers := []string{"ID", "NAME", "REGION", "MEMORY", "CPUs", "STATUS"}
	return headers
}

func (d *Droplet) GetData(global_context *config.GlobalContext) [][]string {
	page := 1
	list_value := godo.ListOptions{Page: page, PerPage: 100, WithProjects: false}
	droplets_list := []godo.Droplet{}
	droplets, _, err := global_context.Client.Droplets.List(context.Background(), &list_value)

	for len(droplets) > 0 {
		if err != nil {
			log.Fatal(err)
		}
		droplets_list = append(droplets_list, droplets...)
		page += 1
		list_value.Page = page
		droplets, _, err = global_context.Client.Droplets.List(context.Background(), &list_value)
	}
	ret := make([][]string, len(droplets_list))
	for i, droplet := range droplets_list {
		ret[i] = make([]string, 6)
		ret[i][0] = fmt.Sprintf("%d", droplet.ID)
		ret[i][1] = droplet.Name
		ret[i][2] = droplet.Region.Slug
		ret[i][3] = fmt.Sprintf("%d GB", droplet.Memory)
		ret[i][4] = fmt.Sprintf("%d", droplet.Vcpus)
		ret[i][5] = droplet.Status
	}
	return ret
}

func (d *Droplet) GetName() string {
	return "droplet"
}
