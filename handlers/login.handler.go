package handlers

import (
	"context"
	"fmt"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/template"
)

func HandleLoginFn(ctx context.Context) fncmp.FnComponent {
	register := fncmp.NewFn(ctx, template.Button("Register")).WithEvents(handleRegister, fncmp.OnClick)
	return fncmp.NewFn(ctx, template.LoginForm(register)).WithEvents(handleLogin, fncmp.OnSubmit)
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func handleLogin(ctx context.Context) fncmp.FnComponent {
	login, err := fncmp.EventData[Login](ctx)
	if err != nil {
		return fncmp.FnErr(ctx, err)
	}
	fmt.Println(login)
	return fncmp.NewFn(ctx, fncmp.HTML("<h2>Hello, "+login.Email+"</h2>"))
}

func handleRegister(ctx context.Context) fncmp.FnComponent {
	return fncmp.NewFn(ctx, fncmp.HTML("<div>Register</div>"))
}
