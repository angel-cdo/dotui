package render

import (
	"context"
	"fmt"
	"log"

	"github.com/angel-cdo/dotui/config"

	"github.com/digitalocean/godo"
)

type RepositoryImage struct {
}

func ByteCountIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func (r *RepositoryImage) Render(global_context *config.GlobalContext, base_render *BaseRender) {
	drawTable(global_context, *base_render)
}

func (r *RepositoryImage) GetHeaders() []string {
	headers := []string{"TAG", "DATE", "SIZE"}
	return headers
}

func (r *RepositoryImage) GetData(global_context *config.GlobalContext) [][]string {
	page := 1
	list_value := godo.ListOptions{Page: page, PerPage: 100, WithProjects: false}
	registry, _, _ := global_context.Client.Registry.Get(context.Background())
	images_list := []*godo.RepositoryTag{}
	images, _, err := global_context.Client.Registry.ListRepositoryTags(context.Background(), registry.Name, global_context.Context, &list_value)

	for len(images) > 0 {
		if err != nil {
			log.Fatal(err)
		}
		images_list = append(images_list, images...)

		page += 1
		list_value.Page = page
		images, _, err = global_context.Client.Registry.ListRepositoryTags(context.Background(), registry.Name, global_context.Context, &list_value)
	}
	ret := make([][]string, len(images_list))
	for i, image := range images_list {
		ret[i] = make([]string, 3)
		ret[i][0] = image.Tag
		ret[i][1] = image.UpdatedAt.String()
		ret[i][2] = ByteCountIEC(image.SizeBytes)
	}
	return ret
}

func (r *RepositoryImage) GetName() string {
	return "images"
}
