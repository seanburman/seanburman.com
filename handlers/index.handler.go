package handlers

import (
	"context"
	"net/http"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/models"
	"github.com/seanburman/seanburman.com/template"
	"github.com/seanburman/seanburman.com/types"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	c := template.Index()
	c.Render(r.Context(), w)
}

func HandleIndexFn(ctx context.Context) fncmp.FnComponent {
	user, err := fncmp.UseCache[models.User](ctx, types.UserKey)
	if err != nil || user.Value().Username == "" {
		return fncmp.RedirectURL(ctx, "/login")
	}

	return fncmp.NewFn(ctx, fncmp.HTML("<h1>Hello, "+user.Value().Username+"</h1>"))
}
