package main

import (
	"net/http"

	// "sketchdb.cozycole.net/ui"
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	dir := "./ui/static/"
	fs := http.FileServer(http.Dir(dir))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))
	router.HandleFunc("/", app.home)

	router.Get("/search", app.search)
	router.Get("/add/video", app.videoAdd)
	router.Post("/add/video", app.videoAddPost)

	return router
}