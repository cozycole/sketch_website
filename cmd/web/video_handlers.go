package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/internal/models"
)

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

	user, ok := r.Context().Value(userContextKey).(*models.User)
	if ok {
		hasLike, _ := app.videos.HasLike(video.ID, user.ID)
		video.Liked = hasLike
	}

	data := app.newTemplateData(r)
	if video.YoutubeID != nil && *video.YoutubeID != "" {
		videoUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", *video.YoutubeID)
		video.URL = &videoUrl
	}
	data.Video = video

	app.render(w, http.StatusOK, "view-video.tmpl.html", "base", data)
}

func (app *application) videoAddPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Need to initialize form data since the template needs it to
	// render. It's a good place to put default values for the fields
	data.Forms.Video = videoForm{}
	data.Video = &models.Video{}
	app.render(w, http.StatusOK, "add-video.tmpl.html", "base", data)
}

func (app *application) videoAdd(w http.ResponseWriter, r *http.Request) {
	var form videoForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}
	app.infoLog.Printf("%+v", form)

	app.validateAddVideoForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Video = form
		app.render(w, http.StatusUnprocessableEntity, "add-video.tmpl.html", "base", data)
		return
	}

	video, err := convertFormToVideo(&form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	id, err := app.videos.Insert(&video)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if video.Creator.ID != nil {
		err = app.videos.InsertVideoCreatorRelation(id, *video.Creator.ID)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	thumbnailHash := generateThumbnailHash(id)
	thumbnailExtension, err := getFileExtension(video.ThumbnailFile)
	if err != nil {
		// TODO: delete video entry here
		app.serverError(w, err)
		return
	}

	thumbnailName := thumbnailHash + thumbnailExtension
	err = app.videos.InsertThumbnailName(id, thumbnailName)
	if err != nil {
		app.serverError(w, err)
		return
	}

	video.ThumbnailName = thumbnailName
	err = app.saveVideoThumbnail(&video)
	if err != nil {
		app.serverError(w, err)
		// TODO: delete video entry here
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/video/%s", video.Slug), http.StatusSeeOther)
}

func (app *application) videoUpdatePage(w http.ResponseWriter, r *http.Request) {
	videoIdParam := r.PathValue("id")
	videoId, err := strconv.Atoi(videoIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	video, err := app.videos.Get(videoId)
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
	data.Forms.Video = videoForm{}
	data.Forms.VideoActor = videoActorForm{}
	app.render(w, http.StatusOK, "update-video.tmpl.html", "base", data)
}

func (app *application) videoUpdate(w http.ResponseWriter, r *http.Request) {
	videoIdParam := r.PathValue("id")
	videoId, err := strconv.Atoi(videoIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	var form videoForm

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	video, err := convertFormToVideo(&form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	video.ID = videoId
	err = app.videos.Update(&video)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) videoAddLike(w http.ResponseWriter, r *http.Request) {
	videoIdParam := r.PathValue("id")
	videoId, err := strconv.Atoi(videoIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok || nil == user {
		app.infoLog.Println("User not logged in!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = app.users.AddLike(user.ID, videoId)
	if err != nil {
		// check if problem with primary key constraint
		app.badRequest(w)
		return
	}
}

func (app *application) videoRemoveLike(w http.ResponseWriter, r *http.Request) {
	videoIdParam := r.PathValue("id")
	videoId, err := strconv.Atoi(videoIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	err = app.users.RemoveLike(user.ID, videoId)
	if err != nil {
		app.badRequest(w)
		return
	}
}
