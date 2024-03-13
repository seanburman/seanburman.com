package handlers

import (
	"context"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/components"
	"github.com/seanburman/seanburman.com/db"
	"github.com/seanburman/seanburman.com/models"
	"github.com/seanburman/seanburman.com/template"
	"github.com/seanburman/seanburman.com/types"
)

func handleRegisterEvent(ctx context.Context) fncmp.FnComponent {
	alertErr := components.AlertMessage(ctx, "An error occured. Please try again.").
		SwapTagInner(template.HeaderTag)

	// Get register data from form on submit
	register, err := fncmp.EventData[components.RegisterFormDetails](ctx)
	if err != nil {
		return alertErr
	}

	if register.Password != register.ConfirmPassword {
		return components.AlertMessage(ctx, models.ErrPasswordsDoNotMatch.Error()).
			SwapTagInner(template.FooterTag)
	}

	// Show loading spinner
	msg := fncmp.HTML("<h2 style='margin-top: 10px;'>Registering " + register.UserName + "...</h2>")
	fncmp.NewFn(ctx, components.LogoSpinner(msg)).
		SwapTagInner(template.FooterTag).
		Dispatch()

	// Create user in database
	newUser, err := models.NewUser(db.Instance, register.UserName, register.Password, register.Email)
	if err != nil {
		return components.AlertMessage(ctx, "An error occured. Please try again.").
			SwapTagInner(template.FooterTag)
	}
	_, err = newUser.Create(db.Instance)
	if err != nil {
		return components.AlertMessage(ctx, err.Error()).
			SwapTagInner(template.FooterTag)
	}

	// Set user in cache
	user, err := fncmp.UseCache[models.User](ctx, types.UserKey)
	if err != nil {
		return fncmp.FnErr(ctx, err)
	}
	user.Set(*newUser)

	// Redirect to home page
	return fncmp.NewFn(ctx, nil).WithRedirect("/")
}

func HandleRegisterFn(ctx context.Context) fncmp.FnComponent {
	// Register button with it's own event
	form := components.RegisterForm(ctx).WithEvents(handleRegisterEvent, fncmp.OnSubmit)
	p := template.Page(types.ComponentConfig{
		Class: "w-full flex justify-center items-center",
	}, form)
	return fncmp.NewFn(ctx, p)
}
