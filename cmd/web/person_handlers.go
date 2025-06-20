package main

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

func (app *application) personView(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	persondId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	person, err := app.people.GetById(persondId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	popular, err := app.sketches.Get(
		&models.Filter{
			Limit:  16,
			Offset: 0,
			SortBy: "az",
			People: []*models.Person{
				&models.Person{ID: person.ID},
			},
		},
	)

	stats, err := app.people.GetPersonStats(persondId)

	data := app.newTemplateData(r)
	page, err := views.PersonPageView(person, stats, popular, app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
	}

	data.Page = page
	app.render(r, w, http.StatusOK, "view-person.gohtml", "base", data)
}

func (app *application) personAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Forms.Person = &personForm{}
	app.render(r, w, http.StatusOK, "add-person.gohtml", "base", data)
}

func (app *application) personAddPost(w http.ResponseWriter, r *http.Request) {
	var form personForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateAddPersonForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Person = &form
		app.render(r, w, http.StatusUnprocessableEntity, "add-person.gohtml", "base", data)
		return
	}

	date, _ := time.Parse(time.DateOnly, form.BirthDate)
	imgName := models.CreateSlugName(form.First + " " + form.Last)

	file, err := form.ProfileImage.Open()
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	defer file.Close()

	mimeType, err := utils.GetMultipartFileMime(file)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	_, slug, fullImgName, err := app.people.
		Insert(
			form.First, form.Last, imgName,
			mimeToExt[mimeType], date,
		)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.fileStorage.SaveFile(path.Join("person", fullImgName), file)
	if err != nil {
		// TODO: We gotta remove the db record on this error
		app.serverError(r, w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/person/%s", slug), http.StatusSeeOther)
}
