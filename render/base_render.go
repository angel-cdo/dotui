package render

import (
	"fmt"

	"github.com/angel-cdo/dotui/config"

	"github.com/derailed/tcell/v2"
	"github.com/derailed/tview"
)

type BaseRender interface {
	Render(global_context *config.GlobalContext, base_render *BaseRender)
	GetHeaders() []string
	GetData(lobal_context *config.GlobalContext) [][]string
	GetName() string
}

func setTableCell(table *tview.Table, row, col int, text string) {
	cell := tview.NewTableCell(text).SetMaxWidth(1).SetExpansion(1)
	table.SetCell(row, col, cell)
}

func drawHeaders(render BaseRender, table *tview.Table) {
	headers := render.GetHeaders()
	for i, header := range headers {
		cell := tview.NewTableCell(header).SetTextColor(tcell.ColorYellow).SetSelectable(false)
		table.SetCell(0, i, cell)
	}
}

func drawTable(global_context *config.GlobalContext, render BaseRender) {
	drawHeaders(render, global_context.Table)
	data := render.GetData(global_context)
	for i, row := range data {
		for j, col := range row {
			setTableCell(global_context.Table, i+1, j, col)
		}
	}
	global_context.Table.SetBorder(true).SetTitle(fmt.Sprintf(" %s [%d]", render.GetName(), len(data)))
}
