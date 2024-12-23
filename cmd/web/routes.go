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

	router.HandleFunc("/ping", ping)

	router.Get("/search", app.searchPage)
	router.Post("/search", app.searchPost)

	router.Get("/video/{slug}", app.videoView)
	router.Get("/video/add", app.videoAdd)
	router.Post("/video/add", app.videoAddPost)

	router.Get("/creator/{slug}", app.creatorView)
	router.Get("/creator/add", app.creatorAdd)
	router.Post("/creator/add", app.creatorAddPost)

	router.Get("/person/{slug}", app.personView)
	router.Get("/person/add", app.personAdd)
	router.Post("/person/add", app.personAddPost)
	router.Get("/person/search", app.personSearch)

	router.Get("/character/search", app.characterSearch)

	return router
}
