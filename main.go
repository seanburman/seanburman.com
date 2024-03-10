package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/kitkitchen/fncmp"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	c := Index()
	c.Render(r.Context(), w)
}

func handleSquare(ctx context.Context) fncmp.FnComponent {
	event, err := fncmp.EventData[fncmp.EventTarget](ctx)
	if err != nil {
		return fncmp.FnErr(ctx, err)
	}
	fmt.Println(event)
	return Square(ctx)
}

func Square(ctx context.Context) fncmp.FnComponent {
	square := fncmp.NewFn(ctx, fncmp.HTML("<div name='f' class='square'></div>"))
	return square.WithEvents(handleSquare, fncmp.OnClick)
}

func handleRegister(ctx context.Context) fncmp.FnComponent {
	return fncmp.NewFn(ctx, fncmp.HTML("<div>Register</div>"))
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

func handlIndexFn(ctx context.Context) fncmp.FnComponent {
	register := fncmp.NewFn(ctx, Button("Register")).WithEvents(handleRegister, fncmp.OnClick)
	return fncmp.NewFn(ctx, LoginForm(register)).WithEvents(handleLogin, fncmp.OnSubmit)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", fncmp.MiddleWareFn(handleIndex, handlIndexFn))
	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	// http.ListenAndServe(":8080", mux)
}
