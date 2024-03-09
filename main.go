package main

import (
	"context"
	"net/http"

	"github.com/kitkitchen/fncmp"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	c := Index()
	c.Render(r.Context(), w)
}

func handlIndexFn(ctx context.Context) fncmp.FnComponent {
	return fncmp.NewFn(ctx, fncmp.HTML("<h1>Sean Burman</h1>"))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", fncmp.MiddleWareFn(handleIndex, handlIndexFn))
	http.ListenAndServe(":8080", mux)
}
