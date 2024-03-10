package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/component"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleLoginFn(ctx context.Context) fncmp.FnComponent {
	// Register button with it's own event
	register := fncmp.NewFn(ctx, component.BlackButton("Register")).
		WithEvents(HandleRegisterFn, fncmp.OnClick)
	// Login form that returns data on submit
	return fncmp.NewFn(ctx, component.LoginForm(register)).
		WithEvents(handleLoginEvent, fncmp.OnSubmit)
}

func handleLoginEvent(ctx context.Context) fncmp.FnComponent {
	// Get login data from form on submit
	login, err := fncmp.EventData[Login](ctx)
	if err != nil {
		// Log to console
		return fncmp.FnErr(ctx, err)
	}
	// Check db
	msg := fncmp.HTML("<h2 style='margin-top: 10px;'>Logging in...</h2>")
	fncmp.NewFn(ctx, component.LoadingSpinner(msg)).
		SwapTagInner("main").
		Dispatch()
	time.Sleep(2 * time.Second)

	fmt.Println(login)
	return fncmp.NewFn(ctx, fncmp.HTML("<h2>Hello, "+login.Email+"</h2>"))
}
