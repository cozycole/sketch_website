package main

import (
	"testing"

	"sketchdb.cozycole.net/internal/assert"
	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

func TestAddSketchImageNames(t *testing.T) {
	validThumbnail, err := utils.CreateMultipartFileHeader("./testdata/test-thumbnail.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
	validImg, err := utils.CreateMultipartFileHeader("./testdata/test-img.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
	validImg2, err := utils.CreateMultipartFileHeader("./testdata/test-img2.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}

	tests := []struct {
		name           string
		sketch         models.Sketch
		wantVidThumb   string
		wantCharThumb1 string
		wantCharThumb2 string
	}{
		{
			name: "Valid Entry",
			sketch: models.Sketch{
				ID:            1,
				Title:         "Test sketch",
				ThumbnailFile: validThumbnail,
				Creator:       &models.Creator{ID: 1},
				Cast: []*models.CastMember{
					{
						Position:      utils.GetIntPtr(0),
						Actor:         &models.Person{ID: utils.GetIntPtr(1)},
						Character:     &models.Character{ID: utils.GetIntPtr(1)},
						ThumbnailFile: validImg,
					},
					{
						Position:      utils.GetIntPtr(1),
						Actor:         &models.Person{ID: utils.GetIntPtr(2)},
						Character:     &models.Character{ID: utils.GetIntPtr(2)},
						ThumbnailFile: validImg2,
					},
				},
			},
			wantVidThumb:   "test-sketch-1.jpg",
			wantCharThumb1: "LyIaGOuGOANpVwsu0UfYtA.jpg",
			wantCharThumb2: "5hcw5nF7F4fPCdRJFP-5IA.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := addSketchImageNames(&tt.sketch)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.sketch.ThumbnailName, tt.wantVidThumb)
			assert.Equal(t, *tt.sketch.Cast[0].ThumbnailName, tt.wantCharThumb1)
			assert.Equal(t, *tt.sketch.Cast[1].ThumbnailName, tt.wantCharThumb2)
		})
	}
}

func TestSaveSketchImages(t *testing.T) {
	validThumbnail, err := utils.CreateMultipartFileHeader("./testdata/test-thumbnail.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
	validImg, err := utils.CreateMultipartFileHeader("./testdata/test-img.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
	validImg2, err := utils.CreateMultipartFileHeader("./testdata/test-img2.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}

	app := newTestApplication(t)

	wantCharThumb1 := "LyIaGOuGOANpVwsu0UfYtA.jpg"
	wantCharThumb2 := "5hcw5nF7F4fPCdRJFP-5IA.jpg"

	tests := []struct {
		name      string
		sketch    models.Sketch
		wantError error
	}{
		{
			name: "Valid Entry",
			sketch: models.Sketch{
				ID:            1,
				Title:         "Test sketch",
				ThumbnailFile: validThumbnail,
				ThumbnailName: "test-sketch-1.jpg",
				Creator:       &models.Creator{ID: 1},
				Cast: []*models.CastMember{
					{
						Position:      utils.GetIntPtr(0),
						Actor:         &models.Person{ID: utils.GetIntPtr(1)},
						Character:     &models.Character{ID: utils.GetIntPtr(1)},
						ThumbnailFile: validImg,
						ThumbnailName: &wantCharThumb1,
					},
					{
						Position:      utils.GetIntPtr(1),
						Actor:         &models.Person{ID: utils.GetIntPtr(2)},
						Character:     &models.Character{ID: utils.GetIntPtr(2)},
						ThumbnailFile: validImg2,
						ThumbnailName: &wantCharThumb2,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.saveSketchImages(&tt.sketch)
			assert.Equal(t, err, tt.wantError)
		})
	}

}
