package components

import (
	"context"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/types"
)

func DIV(cfg types.ComponentConfig, children ...fncmp.Component) fncmp.Component {
	div := fncmp.HTML("<div style='" + cfg.Style + "' class='" + cfg.Class + "'>")
	for _, child := range children {
		child.Render(context.Background(), &div)
	}
	for _, child := range cfg.Children {
		child.Render(context.Background(), &div)
	}
	div.Write([]byte("</div>"))
	return div
}
