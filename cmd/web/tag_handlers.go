package main

import (
	"fmt"
	"net/http"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *application) tagAddPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Forms.Tag = &tagForm{}
	app.render(r, w, http.StatusOK, "add-tag.gohtml", "base", data)
}

func (app *application) tagAdd(w http.ResponseWriter, r *http.Request) {
	var form tagForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateTagForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Tag = &form
		app.render(r, w, http.StatusUnprocessableEntity, "add-tag.gohtml", "base", data)
		return
	}

	tag := convertFormtoTag(&form)
	tagSlug := *tag.Name
	if tag.Category != nil && tag.Category.ID != nil {
		category, err := app.categories.Get(*tag.Category.ID)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
		tagSlug = *category.Name + " " + tagSlug
	}
	slug := models.CreateSlugName(tagSlug)
	tag.Slug = &slug
	_, err = app.tags.Insert(&tag)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Tag added: %s", tag.Name))

	data := app.newTemplateData(r)
	data.Forms.Tag = &tagForm{}
	app.render(r, w, http.StatusOK, "add-tag.gohtml", "base", data)
}

func (app *application) tagRow(w http.ResponseWriter, r *http.Request) {
	app.render(r, w, http.StatusOK, "tag-table.gohtml", "tag-row", views.TagRow{})
}
