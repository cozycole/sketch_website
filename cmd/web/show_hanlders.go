package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	// "path"
	// "time"
	"sketchdb.cozycole.net/internal/models"
	// "sketchdb.cozycole.net/internal/utils"
)

func (app *application) viewShow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	showId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	show, err := app.shows.GetById(showId)
	if err != nil {
		app.serverError(r, w, err)
		return
	} else if show.ID == nil {
		app.notFound(w)
		return
	}

	filter := &models.Filter{
		Limit:  8,
		SortBy: "latest",
		Shows:  []*models.Show{show},
	}

	videos, err := app.videos.Get(filter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	cast, err := app.shows.GetShowCast(*show.ID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Show = show
	data.People = cast
	data.Videos = videos

	app.render(r, w, http.StatusOK, "view-show.tmpl.html", "base", data)
}

func (app *application) addShowPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Show = &models.Show{}
	data.Forms.Show = &showForm{}
	app.render(r, w, http.StatusOK, "add-show.tmpl.html", "base", data)
}

func (app *application) addShow(w http.ResponseWriter, r *http.Request) {
	var form showForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateShowForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Show = &form
		data.Show = &models.Show{}
		app.render(r, w, http.StatusUnprocessableEntity, "add-show.tmpl.html", "base", data)
		return
	}

	show := app.convertFormtoShow(&form)
	slug := models.CreateSlugName(*show.Name)
	show.Slug = &slug

	thumbName, err := generateThumbnailName(form.ProfileImg)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	show.ProfileImg = &thumbName

	id, err := app.shows.Insert(&show)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	show.ID = &id

	err = app.saveProfileImage(*show.ProfileImg, "show", form.ProfileImg)
	if err != nil {
		app.shows.Delete(&show)
		app.serverError(r, w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/show/%d/%s", *show.ID, *show.Slug), http.StatusSeeOther)
}

func (app *application) updateShowPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	showId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	show, err := app.shows.GetById(showId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Show = show
	for _, season := range data.Show.Seasons {
		for _, ep := range season.Episodes {
			app.infoLog.Printf("Episode %d\n", *ep.Number)
			for _, vid := range ep.Videos {
				app.infoLog.Printf("Video %d\n", *vid.ID)
			}
		}
	}
	data.Episode = &models.Episode{}
	data.Forms.Show = &showForm{}
	data.Forms.Episode = &episodeForm{}
	for _, ep := range data.Show.Seasons[0].Episodes {
		if ep.Title != nil {
			app.infoLog.Println(*ep.Title)
		}

	}
	app.render(r, w, http.StatusOK, "update-show.tmpl.html", "base", data)
}

func (app *application) updateShow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	showId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	var form showForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	oldShow, err := app.shows.GetById(showId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	app.validateUpdateShowForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Show = &form
		data.Show = oldShow
		app.render(r, w, http.StatusUnprocessableEntity, "show-form.tmpl.html", "show-form", data)
		return
	}

	newShow := app.convertFormtoShow(&form)
	newShow.ID = &showId

	var profileImg string
	if oldShow.ProfileImg != nil {
		profileImg = *oldShow.ProfileImg
	}

	if form.ProfileImg != nil {
		var err error
		profileImg, err = generateThumbnailName(form.ProfileImg)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
		err = app.saveProfileImage(profileImg, "show", form.ProfileImg)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	newShow.ProfileImg = &profileImg
	err = app.shows.Update(&newShow)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if form.ProfileImg != nil && oldShow.ProfileImg != nil {
		err = app.deleteImage("show", *oldShow.ProfileImg)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	data := app.newTemplateData(r)
	data.Forms.Show = &form
	data.Show = &newShow
	data.Flash = flashMessage{
		Level:   "success",
		Message: "Show successfully updated!",
	}
	app.render(r, w, http.StatusOK, "show-form.tmpl.html", "show-form", data)
}

func (app *application) addSeason(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	showIdParam := r.PathValue("id")
	showId, err := strconv.Atoi(showIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	seasonId, err := app.shows.AddSeason(showId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	season, err := app.shows.GetSeason(seasonId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Season = season
	data.Episode = &models.Episode{}
	data.Forms.Episode = &episodeForm{}
	app.render(r, w, http.StatusOK, "season-form.tmpl.html", "season-form", data)
}

func (app *application) addEpisode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var form episodeForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.infoLog.Printf("EP FORM: %+v\n", form)

	app.validateEpisodeForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Episode = &form
		data.Episode = &models.Episode{}
		app.render(r, w, http.StatusUnprocessableEntity, "episode-form.tmpl.html", "episode-form", data)
		return
	}

	episode := app.convertFormtoEpisode(&form)

	thumbnailName, err := generateThumbnailName(form.Thumbnail)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	episode.Thumbnail = &thumbnailName

	id, err := app.shows.InsertEpisode(&episode)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	episode.ID = &id

	app.saveThumbnail(thumbnailName, "episode", form.Thumbnail)

	season, err := app.shows.GetSeason(*episode.SeasonId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.infoLog.Printf("Season: %+v\n", season)
	data := app.newTemplateData(r)
	data.Season = season
	data.Episode = &models.Episode{}
	data.Forms.Episode = &episodeForm{}
	data.Flash = flashMessage{Level: "success", Message: "Episode added!"}
	app.render(r, w, http.StatusOK, "season-form.tmpl.html", "season-form", data)
}

func (app *application) updateEpisode(w http.ResponseWriter, r *http.Request) {
	epIdParam := r.PathValue("id")
	epId, err := strconv.Atoi(epIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	var form episodeForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	oldEpisode, err := app.shows.GetEpisode(epId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	app.validateUpdateEpisodeForm(&form)
	app.infoLog.Printf("%+v\n", form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Episode = &form
		data.Episode = oldEpisode
		app.render(r, w, http.StatusUnprocessableEntity, "episode-form.tmpl.html", "episode-form", data)
		return
	}

	episode := app.convertFormtoEpisode(&form)

	var thumbnailName string
	if oldEpisode.Thumbnail != nil {
		thumbnailName = *oldEpisode.Thumbnail
	} else {
		thumbnailName = ""
	}

	if form.Thumbnail != nil {
		var err error
		thumbnailName, err = generateThumbnailName(form.Thumbnail)
		if err != nil {
			app.serverError(r, w, err)
			return
		}

		err = app.saveThumbnail(thumbnailName, "episode", form.Thumbnail)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	episode.ID = &epId
	episode.Thumbnail = &thumbnailName
	app.infoLog.Printf("%+v\n", episode)
	err = app.shows.UpdateEpisode(&episode)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if form.Thumbnail != nil && oldEpisode.Thumbnail != nil {
		err = app.deleteImage("episode", *oldEpisode.Thumbnail)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	data := app.newTemplateData(r)
	data.Forms.Episode = &form
	data.Episode = &episode
	data.Flash = flashMessage{
		Level:   "success",
		Message: "Episode updated successfully!",
	}

	app.render(r, w, http.StatusOK, "episode-form.tmpl.html", "episode-form", data)
}

func (app *application) deleteEpisode(w http.ResponseWriter, r *http.Request) {
	epIdParam := r.PathValue("id")
	epId, err := strconv.Atoi(epIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	data := app.newTemplateData(r)
	episode, err := app.shows.GetEpisode(epId)
	if err != nil {
		var status int
		if errors.Is(err, models.ErrNoRecord) {
			data.Flash = flashMessage{
				Level:   "error",
				Message: "404 Episode does not exist",
			}
			status = http.StatusNotFound
		} else {
			data.Flash = flashMessage{
				Level:   "error",
				Message: "500 Internal Server Error",
			}
			status = http.StatusInternalServerError
		}
		app.render(r, w, status, "flash-message.tmpl.html", "flash-message", data)
		return
	}

	if len(episode.Videos) != 0 {
		data.Flash = flashMessage{
			Level:   "error",
			Message: "400 Cannot delete episode with sketches",
		}
		app.render(r, w, http.StatusBadRequest, "flash-message.tmpl.html", "flash-message", data)
		return
	}

	err = app.shows.DeleteEpisode(*episode.ID)
	if err != nil {
		data.Flash = flashMessage{
			Level:   "error",
			Message: "500 Internal Server Error",
		}
		app.render(r, w, http.StatusInternalServerError, "flash-message.tmpl.html", "flash-message", data)
		return
	}

	data.Flash = flashMessage{
		Level:   "success",
		Message: "Episode deleted successfully!",
	}
	app.render(r, w, http.StatusOK, "flash-message.tmpl.html", "flash-message", data)
}
