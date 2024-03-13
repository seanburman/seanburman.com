package components

import (
	"context"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/types"
)

func AlertMessage(ctx context.Context, message string) fncmp.FnComponent {
	container := DIV(types.ComponentConfig{
		ID:    "alert-box",
		Class: "w-full flex justify-center items-center text-sm bg-red-500 text-white p-4 shadow-md position-fixed top-0",
		Style: "height: 40px",
		Children: []fncmp.Component{
			fncmp.HTML("<p>" + message + "</p>"),
		},
	})

	return fncmp.NewFn(ctx, container)
}
