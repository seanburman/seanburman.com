package main

import (
	"net/http"

	"github.com/kitkitchen/fncmp"
	"github.com/seanburman/seanburman.com/config"
	"github.com/seanburman/seanburman.com/db"
	"github.com/seanburman/seanburman.com/handlers"
)

func main() {
	// Connect to the database
	db.Instance = &db.Database{
		Postgres: db.NewPostgresDriver(),
	}

	// Create the server
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", fncmp.MiddleWareFn(handlers.HandleIndex, handlers.HandleIndexFn))
	mux.HandleFunc("/login", fncmp.MiddleWareFn(handlers.HandleIndex, handlers.HandleLoginFn))
	mux.HandleFunc("/register", fncmp.MiddleWareFn(handlers.HandleIndex, handlers.HandleRegisterFn))
	http.ListenAndServe(":"+config.Env().PORT, mux)
}
