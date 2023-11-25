package config

import (
	"github.com/derailed/tview"
	"github.com/digitalocean/godo"
)

type GlobalContext struct {
	Client  *godo.Client
	Table   *tview.Table
	Context string
}
