package main

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

var maxFileNameLength = 50

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	videos, err := app.videos.GetAll(8)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Videos = videos

	app.render(w, http.StatusOK, "home.tmpl.html", "base", data)
}

func (app *application) searchPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "search.tmpl.html", "base", data)
}

func (app *application) videoView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	video, err := app.videos.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Video = video

	app.render(w, http.StatusOK, "view-video.tmpl.html", "base", data)
}

func (app *application) creatorView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	creator, err := app.creators.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	videos, err := app.videos.GetByCreator(creator.ID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Creator = creator
	data.Videos = videos

	app.render(w, http.StatusOK, "view-creator.tmpl.html", "base", data)
}

func (app *application) creatorAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = addCreatorForm{}
	app.render(w, http.StatusOK, "add-creator.tmpl.html", "base", data)
}

func (app *application) creatorAddPost(w http.ResponseWriter, r *http.Request) {
	var form addCreatorForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateAddCreatorForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "add-creator.tmpl.html", "base", data)
		return
	}

	date, _ := time.Parse(time.DateOnly, form.EstablishedDate)
	imgName := models.CreateSlugName(form.Name, maxFileNameLength)

	file, err := form.ProfileImage.Open()
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer file.Close()

	mimeType, err := utils.GetMultipartFileMime(file)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// the insert returns the fullImgName which is {fileName}-{id}.{ext}
	_, slug, fullImgName, err := app.creators.
		Insert(
			form.Name, form.URL, imgName,
			mimeToExt[mimeType], date,
		)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.fileStorage.SaveFile(path.Join("creator", fullImgName), file)
	if err != nil {
		// TODO: We gotta remove the db record on this error
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/creator/%s", slug), http.StatusSeeOther)
}

func (app *application) personView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	person, err := app.people.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	videos, err := app.videos.GetByPerson(*person.ID)

	data := app.newTemplateData(r)
	data.Person = person
	data.Videos = videos

	app.render(w, http.StatusOK, "view-person.tmpl.html", "base", data)
}

func (app *application) personAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = addPersonForm{}
	app.render(w, http.StatusOK, "add-person.tmpl.html", "base", data)
}

func (app *application) personAddPost(w http.ResponseWriter, r *http.Request) {
	var form addPersonForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateAddPersonForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "add-person.tmpl.html", "base", data)
		return
	}

	date, _ := time.Parse(time.DateOnly, form.BirthDate)
	imgName := models.CreateSlugName(form.First+" "+form.Last, maxFileNameLength)

	file, err := form.ProfileImage.Open()
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer file.Close()

	mimeType, err := utils.GetMultipartFileMime(file)
	if err != nil {
		app.serverError(w, err)
		return
	}

	_, slug, fullImgName, err := app.people.
		Insert(
			form.First, form.Last, imgName,
			mimeToExt[mimeType], date,
		)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.fileStorage.SaveFile(path.Join("person", fullImgName), file)
	if err != nil {
		// TODO: We gotta remove the db record on this error
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/person/%s", slug), http.StatusSeeOther)
}

func (app *application) videoAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Need to initialize form data since the template needs it to
	// render. It's a good place to put default values for the fields
	data.Form = addVideoForm{}
	app.render(w, http.StatusOK, "add-video.tmpl.html", "base", data)
}

func (app *application) videoAddPost(w http.ResponseWriter, r *http.Request) {
	var form addVideoForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateAddVideoForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "add-video.tmpl.html", "base", data)
		return
	}

	video, err := convertFormToVideo(&form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = addVideoImageNames(&video)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// NOTE: This mutates the video struct by adding the newly created db serial id
	// to the id field
	err = app.videos.Insert(&video)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.saveVideoImages(&video)
	if err != nil {
		app.serverError(w, err)
		// TODO: delete video entry now
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/video/%d/%s", video.ID, video.Slug), http.StatusSeeOther)
}

type dropdownSearchResults struct {
	Results      []result
	Redirect     string // e.g. /person/add
	RedirectText string // e.g. "Add Person +"
}

type result struct {
	ID   int
	Text string
	Img  string
}

func (app *application) personSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	redirLink := "/person/add"
	redirText := "Add Person +"
	results := dropdownSearchResults{
		Redirect:     redirLink,
		RedirectText: redirText,
	}

	if q != "" {
		q = strings.Replace(q, " ", "", -1)
		dbResults, err := app.people.Search(q)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				app.serverError(w, err)
			}
			return
		}

		if dbResults != nil {
			res := []result{}
			for _, row := range dbResults {
				r := result{}
				r.Text = *row.First + " " + *row.Last
				r.ID = *row.ID
				res = append(res, r)
			}

			results.Results = res
		}
	}

	data := app.newTemplateData(r)
	data.DropdownResults = results

	app.render(w, http.StatusOK, "dropdown.tmpl.html", "", data)
}

func (app *application) characterSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	redirLink := "/character/add"
	redirText := "Add Character +"
	results := dropdownSearchResults{
		Redirect:     redirLink,
		RedirectText: redirText,
	}

	if q != "" {
		q = strings.Replace(q, " ", "", -1)
		dbResults, err := app.characters.Search(q)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				app.serverError(w, err)
			}
			return
		}

		if dbResults != nil {
			res := []result{}
			for _, row := range dbResults {
				r := result{}
				r.Text = *row.Name
				r.ID = *row.ID
				res = append(res, r)
			}

			results.Results = res
		}
	}

	data := app.newTemplateData(r)
	data.DropdownResults = results

	app.render(w, http.StatusOK, "dropdown.tmpl.html", "", data)
}

func (app *application) searchPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	q := r.Form.Get("query")
	results, err := app.search.Search(q)
	if err != nil {
		if !errors.Is(err, models.ErrNoRecord) {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.SearchResults = results

	app.render(w, http.StatusOK, "search-result.tmpl.html", "", data)
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}
