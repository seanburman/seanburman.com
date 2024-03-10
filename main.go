package main

import (
	"net/http"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/handlers"
	"github.com/seanburman/seanburman.com/template"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	c := template.Index()
	c.Render(r.Context(), w)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", fncmp.MiddleWareFn(handleIndex, handlers.HandleLoginFn))
	mux.HandleFunc("/login", fncmp.MiddleWareFn(handleIndex, handlers.HandleLoginFn))
	// http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	http.ListenAndServe(":8080", mux)
}
