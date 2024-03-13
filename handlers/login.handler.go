package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/components"
	"github.com/seanburman/seanburman.com/db"
	"github.com/seanburman/seanburman.com/models"
	"github.com/seanburman/seanburman.com/template"
	"github.com/seanburman/seanburman.com/types"
)

type Login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func HandleLoginFn(ctx context.Context) fncmp.FnComponent {
	request, ok := ctx.Value(fncmp.RequestKey).(*http.Request)
	if !ok {
		return fncmp.FnErr(ctx, errors.New("could not get request"))
	}

	// TODO: Once we can get the original request, we can use this to redirect to the original page
	fmt.Println(request.URL.Path)

	// Register button with it's own event
	register := fncmp.NewFn(ctx, components.BlackButton("Register")).
		WithEvents(handleRegisterClick, fncmp.OnClick)
	// Login form that returns data on submit
	return fncmp.NewFn(ctx, components.LoginForm(register)).
		WithEvents(handleLoginEvent, fncmp.OnSubmit)
}

func handleRegisterClick(ctx context.Context) fncmp.FnComponent {
	return fncmp.NewFn(ctx, nil).WithRedirect("/register")
}

func handleLoginEvent(ctx context.Context) fncmp.FnComponent {
	alertErr := components.AlertMessage(ctx, "An error occured. Please try again.").
		SwapTagInner("header")

	// Get login data from form on submit
	login, err := fncmp.EventData[Login](ctx)
	if err != nil {
		alertErr.Dispatch()
		return fncmp.FnErr(ctx, err)
	}

	// Show loading spinner
	msg := fncmp.HTML("<h2 style='margin-top: 10px;'>Logging in...</h2>")
	fncmp.NewFn(ctx, components.LoadingSpinner(msg)).
		SwapTagInner(template.HeaderTag).Dispatch()

	// Get user from database
	creds, err := models.NewUser(db.Instance, login.UserName, login.Password, "")
	if err != nil {
		alertErr.Dispatch()
		return fncmp.FnErr(ctx, err)
	}
	user, err := creds.Get(db.Instance)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) || errors.Is(err, models.ErrWrongPassword) {
			components.AlertMessage(ctx, err.Error()).
				SwapTagInner(template.HeaderTag).Dispatch()
			return fncmp.FnErr(ctx, err)
		}
	}

	// Check password
	ok := user.Authenticate(login.Password)
	if !ok {
		components.AlertMessage(ctx, models.ErrWrongPassword.Error()).
			SwapTagInner(template.HeaderTag).Dispatch()
		return fncmp.FnErr(ctx, err)
	}

	// Store user in cache
	authenticatedUser, err := fncmp.UseCache[models.User](ctx, types.UserKey)
	if err != nil {
		return fncmp.FnErr(ctx, err)
	}
	authenticatedUser.Set(*user, 24*time.Hour)

	fmt.Println(login)
	return fncmp.NewFn(ctx, nil).WithRedirect("/")
}

func handleLogoutEvent(ctx context.Context) fncmp.FnComponent {
	return fncmp.NewFn(ctx, nil).WithRedirect("/")
}
