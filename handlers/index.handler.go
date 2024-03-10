package handlers

import (
	"context"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/template"
)

func HandleIndexFn(ctx context.Context) fncmp.FnComponent {
	return fncmp.NewFn(ctx, template.Index())
}
