package main

import (
	"fmt"
	"time"

	"sketchdb.cozycole.net/internal/models"
)

func convertFormToSketch(form *sketchForm) (models.Sketch, error) {
	uploadDate, err := time.Parse(time.DateOnly, form.UploadDate)
	if err != nil {
		return models.Sketch{}, fmt.Errorf("unable to parse date")
	}

	creator := &models.Creator{
		ID: &form.CreatorID,
	}

	return models.Sketch{
		Title:         &form.Title,
		URL:           &form.URL,
		Slug:          &form.Slug,
		ThumbnailFile: form.Thumbnail,
		Rating:        &form.Rating,
		UploadDate:    &uploadDate,
		Creator:       creator,
	}, nil
}

func convertFormtoCastMember(form *castForm) models.CastMember {
	actor := models.Person{ID: &form.PersonID}
	character := models.Character{}
	if form.CharacterID != 0 {
		character.ID = &form.CharacterID
	}
	return models.CastMember{
		Actor:         &actor,
		Character:     &character,
		CharacterName: &form.CharacterName,
		ThumbnailFile: form.CharacterThumbnail,
		ProfileFile:   form.CharacterProfile,
	}
}

func convertFormtoCategory(form *categoryForm) models.Category {
	return models.Category{
		Name: &form.Name,
	}
}

func convertFormtoTag(form *tagForm) models.Tag {
	var categoryId int
	if form.CategoryId != 0 {
		categoryId = form.CategoryId
	}
	return models.Tag{
		Name:     &form.Name,
		Category: &models.Category{ID: &categoryId},
	}
}

func (app *application) convertFormtoSketchTags(form *sketchTagsForm) ([]*models.Tag, error) {
	var tags []*models.Tag
	for _, tagId := range form.TagIds {
		tag, err := app.tags.Get(tagId)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (app *application) convertFormtoShow(form *showForm) models.Show {
	return models.Show{
		Name: &form.Name,
		Slug: &form.Slug,
	}

}

func (app *application) convertFormtoEpisode(form *episodeForm) models.Episode {
	airDate, _ := time.Parse(time.DateOnly, form.AirDate)
	return models.Episode{
		ID:       &form.ID,
		Title:    &form.Title,
		Number:   &form.Number,
		AirDate:  &airDate,
		SeasonId: &form.SeasonId,
	}
}
