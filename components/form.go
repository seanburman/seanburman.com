package components

import (
	"context"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/class"
	"github.com/seanburman/seanburman.com/types"
)

type RegisterFormDetails struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm-password"`
	Email           string `json:"email"`
}

func RegisterForm(ctx context.Context) fncmp.FnComponent {
	return fncmp.NewFn(ctx, Form(types.ComponentConfig{
		ID:    "register-form",
		Class: "flex flex-col bg-white p-8 rounded-lg shadow-md w-96",
		Children: []fncmp.Component{
			Label(types.ComponentConfig{
				For:   "username",
				Class: class.Label.String(),
				Label: "Username",
			}),
			Input(types.ComponentConfig{
				ID:          "username",
				Name:        "username",
				Class:       class.Input.Add(class.MB4).String(),
				Type:        "text",
				Placeholder: "",
				Required:    "true",
			}),
			Label(types.ComponentConfig{
				For:   "email",
				Class: class.Label.String(),
				Label: "Email",
			}),
			Input(types.ComponentConfig{
				ID:          "email",
				Name:        "email",
				Class:       class.Input.Add(class.MB4).String(),
				Type:        "email",
				Placeholder: "",
				Required:    "true",
			}),
			Label(types.ComponentConfig{
				For:   "password",
				Class: class.Label.String(),
				Label: "Password",
			}),
			Input(types.ComponentConfig{
				ID:          "password",
				Name:        "password",
				Class:       class.Input.Add(class.MB4).String(),
				Type:        "password",
				Placeholder: "Password",
				Required:    "true",
			}),
			Label(types.ComponentConfig{
				For:   "confirm-password",
				Class: class.Label.String(),
				Label: "Confirm Password",
			}),
			Input(types.ComponentConfig{
				ID:          "confirm-password",
				Name:        "confirm-password",
				Class:       class.Input.Add(class.MB4).String(),
				Type:        "password",
				Placeholder: "Confirm Password",
				Required:    "true",
			}),
			Button(types.ComponentConfig{
				ID:    "register",
				Name:  "register",
				Class: class.ButtonBlack.String(),
				Type:  "submit",
				Value: "Register",
			}),
		},
	}))
}
