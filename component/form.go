package component

import (
	"context"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/class"
	"github.com/seanburman/seanburman.com/types"
)

func RegisterForm(ctx context.Context) fncmp.FnComponent {
	return fncmp.NewFn(ctx, Form(types.ComponentConfig{
		ID:    "register-form",
		Class: "flex flex-col bg-white p-8 rounded-lg shadow-md",
		Children: []fncmp.Component{
			Label(types.ComponentConfig{
				For:   "email",
				Class: class.Label.String(),
				Label: "Email",
			}),
			Input(types.ComponentConfig{
				ID:          "email",
				Class:       class.Input.Add(class.MB4).String(),
				Type:        "email",
				Name:        "email",
				Placeholder: "Email",
			}),
			Label(types.ComponentConfig{
				For:   "password",
				Class: class.Label.String(),
				Label: "Password",
			}),
			Input(types.ComponentConfig{
				ID:          "password",
				Class:       class.Input.Add(class.MB4).String(),
				Type:        "password",
				Name:        "password",
				Placeholder: "Password",
			}),
			Label(types.ComponentConfig{
				For:   "confirm-password",
				Class: class.Label.String(),
				Label: "Confirm Password",
			}),
			Input(types.ComponentConfig{
				ID:          "confirm-password",
				Class:       class.Input.Add(class.MB4).String(),
				Type:        "password",
				Name:        "confirm-password",
				Placeholder: "Confirm Password",
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
