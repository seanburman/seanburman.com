package handlers

import (
	"context"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/class"
	"github.com/seanburman/seanburman.com/component"
	"github.com/seanburman/seanburman.com/template"
	"github.com/seanburman/seanburman.com/types"
)

func HandleRegisterFn(ctx context.Context) fncmp.FnComponent {
	text := "Your email and password are stored as hashes<br/> and cannot be recovered."
	class := class.TextCenter.Add(class.MB4).String()

	// Register button with it's own event
	form := component.RegisterForm(ctx)
	p := template.Page(types.ComponentConfig{
		Class: "",
		Children: []fncmp.Component{
			fncmp.HTML("<p class=" + class + ">" + text + "</p>"),
		},
	}, form)
	return fncmp.NewFn(ctx, p)
}
