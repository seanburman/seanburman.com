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
		WithEvents(handleClickRegister, fncmp.OnClick)
	// Login form that returns data on submit
	return fncmp.NewFn(ctx, components.LoginForm(register)).
		WithEvents(handleLoginEvent, fncmp.OnSubmit)
}

func handleClickRegister(ctx context.Context) fncmp.FnComponent {
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
	// msg := fncmp.HTML("<h2 style='margin-top: 10px;'>Logging in...</h2>")
	// fncmp.NewFn(ctx, components.LogoSpinner(msg)).
	// 	SwapTagInner(template.FooterTag).Dispatch()
	fncmp.AddClasses(ctx, "logo", "animate-spin")

	// Get user from database
	creds, err := models.NewUser(db.Instance, login.UserName, login.Password, "")
	if err != nil {
		alertErr.Dispatch()
		return fncmp.FnErr(ctx, err)
	}
	user, err := creds.Get(db.Instance)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return components.AlertMessage(ctx, err.Error()).
				SwapTagInner(template.FooterTag)
		}
		return fncmp.FnErr(ctx, err)
	}

	// Check password
	ok := user.Authenticate(login.Password)
	if !ok {
		return components.AlertMessage(ctx, models.ErrWrongPassword.Error()).
			SwapTagInner(template.FooterTag)
	}

	// Store user in cache
	userCache, err := fncmp.UseCache[models.User](ctx, types.UserKey)
	if err != nil {
		return fncmp.FnErr(ctx, err)
	}
	fncmp.OnCacheChange(userCache, func() {
		fmt.Println("Cache updated")
	})

	userCache.Set(*user, 24*time.Hour)

	fmt.Println(login)
	return fncmp.NewFn(ctx, nil).WithRedirect("/")
}

func handleClickLogout(ctx context.Context) fncmp.FnComponent {
	cache, err := fncmp.UseCache[models.User](ctx, types.UserKey)
	if err != nil {
		return fncmp.FnErr(ctx, err)
	}
	cache.Delete()
	return fncmp.NewFn(ctx, nil).WithRedirect("/")
}
