package main

import (
	"fmt"
	"mime/multipart"

	"sketchdb.cozycole.net/internal/utils"
	"sketchdb.cozycole.net/internal/validator"
)

type Forms struct {
	Cast       *castForm
	Category   *categoryForm
	Creator    *creatorForm
	Episode    *episodeForm
	Login      *userLoginForm
	Person     *personForm
	Show       *showForm
	Signup     *userSignupForm
	Tag        *tagForm
	Sketch     *sketchForm
	SketchTags *sketchTagsForm
}

// Changes to the form fields must be updated in their respective
// validate functions
type creatorForm struct {
	ID                  int                   `form:"id"`
	Name                string                `form:"name"`
	URL                 string                `form:"url"`
	EstablishedDate     string                `form:"establishedDate"`
	ProfileImage        *multipart.FileHeader `img:"profileImg"`
	Action              string                `form:"-"`
	ImageUrl            string                `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) validateCreatorForm(form *creatorForm) {
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.URL), "url", "This field cannot be blank")
	if form.EstablishedDate != "" {
		form.CheckField(validator.ValidDate(form.EstablishedDate), "establishedDate", "Date not of correct format YYYY-MM-DD")
	}

	if form.ID == 0 {
		form.CheckField(form.ProfileImage != nil, "profileImg", "Please upload an image")
	}

	if form.ProfileImage == nil {
		return
	}

	profileImg, err := form.ProfileImage.Open()
	if err != nil {
		form.AddFieldError("profileImg", "Unable to open file, ensure it is a jpg")
		return
	}
	defer profileImg.Close()

	form.CheckField(validator.IsMime(profileImg, "image/jpeg"), "profileImg", "Uploaded file must be jpg")
}

type personForm struct {
	ID                  int                   `form:"id"`
	First               string                `form:"first"`
	Last                string                `form:"last"`
	BirthDate           string                `form:"birthDate"`
	Professions         string                `form:"professions"`
	ProfileImage        *multipart.FileHeader `img:"profileImg"`
	Action              string                `form:"-"`
	ImageUrl            string                `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) validatePersonForm(form *personForm) {
	form.CheckField(validator.NotBlank(form.First), "first", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Last), "last", "This field cannot be blank")
	if form.BirthDate != "" {
		form.CheckField(validator.ValidDate(form.BirthDate), "birthDate", "Date not of correct format YYYY-MM-DD")
	}

	// if it's an add instead of update
	if form.ID == 0 {
		form.CheckField(form.ProfileImage != nil, "profileImg", "Please upload an image")
	}

	if form.ProfileImage == nil {
		return
	}

	profileImg, err := form.ProfileImage.Open()
	if err != nil {
		form.AddFieldError("profileImg", "Unable to open file, ensure it is a jpg")
		return
	}
	defer profileImg.Close()

	form.CheckField(validator.IsMime(profileImg, "image/jpeg"), "profileImg", "Uploaded file must be jpg or png")
}

type characterForm struct {
	ID                  int                   `form:"id"`
	Name                string                `form:"name"`
	Type                string                `form:"type"`
	ProfileImage        *multipart.FileHeader `img:"profileImg"`
	PersonID            int                   `form:"personId"`
	PersonInput         string                `form:"personInput"`
	Action              string                `form:"-"`
	ImageUrl            string                `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) validateCharacterForm(form *characterForm) {
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Type), "type", "This field cannot be blank")

	if form.ID == 0 {
		form.CheckField(form.ProfileImage != nil, "profileImg", "Please upload an image")
	}

	form.CheckField(
		validator.PermittedValue(
			form.Type, "original", "impression", "fictional_impression", "generic",
		),
		"type", "Character must be either Original, Impression or Generic")
	form.CheckField(
		form.Type != "impression" || form.PersonID > 0,
		"personId", "Person must be selected")

	if form.ProfileImage == nil {
		return
	}

	profileImg, err := form.ProfileImage.Open()
	if err != nil {
		form.AddFieldError("profileImg", "Unable to open file, ensure it is a jpg")
		return
	}
	defer profileImg.Close()

	form.CheckField(validator.IsMime(profileImg, "image/jpeg"), "profileImg", "Uploaded file must be jpg")
}

type sketchForm struct {
	ID                  int                   `form:"id"`
	Title               string                `form:"title"`
	URL                 string                `form:"url"`
	Slug                string                `form:"slug"`
	Rating              string                `form:"rating"`
	UploadDate          string                `form:"uploadDate"`
	Number              int                   `form:"number"`
	Thumbnail           *multipart.FileHeader `img:"thumbnail"`
	CreatorID           int                   `form:"creatorId"`
	CreatorInput        string                `form:"creatorInput"`
	EpisodeID           int                   `form:"episodeId"`
	EpisodeInput        string                `form:"episodeInput"`
	EpisodeStart        int                   `form:"episodeStart"`
	SeriesID            int                   `form:"seriesId"`
	SeriesInput         string                `form:"seriesInput"`
	SeriesPart          int                   `form:"seriesPart"`
	Action              string                `form:"-"`
	ImageUrl            string                `form:"-"`
	validator.Validator `form:"-"`
}

// We need this function to have access to the apps state
// to validate based on database queries
func (app *application) validateSketchForm(form *sketchForm) {
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(form.CreatorID != 0 || form.EpisodeID != 0, "creatorId", "A creator or episode must be defined")
	form.CheckField(validator.NotBlank(form.UploadDate), "uploadDate", "This field cannot be blank")
	form.CheckField(validator.ValidDate(form.UploadDate), "uploadDate", "Date must be of the format YYYY-MM-DD")

	if form.CreatorID != 0 {
		form.CheckField(
			validator.BoolWithError(app.creators.Exists(form.CreatorID)),
			"creator",
			"Unable to find creator, please add them",
		)
	}

	if form.ID == 0 {
		form.CheckField(form.Thumbnail != nil, "thumbnail", "Please upload an image")
	}

	if form.Thumbnail == nil {
		return
	}

	thumbnail, err := form.Thumbnail.Open()
	if err != nil {
		form.AddFieldError("thumbnail", "Unable to open file, ensure it is a valid jpg")
		return
	}
	defer thumbnail.Close()

	form.CheckField(validator.IsMime(thumbnail, "image/jpeg", "image/png"), "thumbnail", "Uploaded file must be jpg or png")
	width, height, err := utils.GetImageDimensions(thumbnail)
	if err != nil {
		form.AddFieldError("thumbnail", "Unable to determine image dimensions")
		return
	}

	form.CheckField(width >= LargeThumbnailWidth && height >= LargeThumbnailHeight, "thumbnail", "Thumbnail dimensions must be at least 480x360")
}

type castForm struct {
	ID                  int                   `form:"id"`
	PersonID            int                   `form:"personId"`
	PersonInput         string                `form:"personInput"`
	CharacterName       string                `form:"characterName"`
	CastRole            string                `form:"castRole"`
	MinorRole           bool                  `form:"minorRole"`
	CharacterID         int                   `form:"characterId"`
	CharacterInput      string                `form:"characterInput"`
	ThumbnailName       string                `form:"-"`
	CharacterThumbnail  *multipart.FileHeader `img:"characterThumbnail"`
	ProfileImage        string                `form:"-"`
	CharacterProfile    *multipart.FileHeader `img:"characterProfile"`
	Action              string                `form:"-"`
	ImageUrl            string                `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) validateCastForm(form *castForm) {
	pid := form.PersonID
	form.CheckField(pid != 0, "personId", "This field cannot be blank. Select a person from dropdown.")
	if pid != 0 {
		form.CheckField(
			validator.BoolWithError(app.people.Exists(pid)),
			"personId",
			"Person does not exist. Please add it, then resubmit sketch!",
		)
	}

	form.CheckField(validator.NotBlank(form.CharacterName), "characterName", "This field cannot be blank.")

	cid := form.CharacterID
	if cid != 0 {
		form.CheckField(
			validator.BoolWithError(app.characters.Exists(cid)),
			"characterId",
			"Character does not exist. Please add it, then resubmit sketch!",
		)
	}

	thumb := form.CharacterThumbnail
	// cast members don't need a thumbnail or a profile image,
	// it will default to the sketch thumbnail and the person's profile image
	if thumb != nil {
		thumbnail, err := thumb.Open()
		if err != nil {
			form.AddFieldError("characterThumbnail", "File error, ensure it is a jpg")
			return
		}

		form.CheckField(validator.IsMime(thumbnail, "image/jpeg"),
			"characterThumbnail", "Uploaded file must be jpg or png")

		width, height, err := utils.GetImageDimensions(thumbnail)
		if err != nil {
			app.errorLog.Print(err)
			form.AddFieldError("characterThumbnail", "Unable to determine image dimensions")
			return
		}

		form.CheckField(width >= StandardThumbnailWidth && height >= StandardThumbnailHeight,
			"characterThumbnail",
			fmt.Sprintf(
				"Thumbnail dimensions must be at least %dx%d",
				StandardThumbnailWidth, StandardThumbnailHeight,
			),
		)
	}

	profile := form.CharacterProfile

	if profile == nil {
		return
	}

	thumbnail, err := profile.Open()
	if err != nil {
		form.AddFieldError("characterProfile", "File error, ensure it is a jpg")
		return
	}
	defer thumbnail.Close()

	form.CheckField(validator.IsMime(thumbnail, "image/jpeg"),
		"characterProfile", "Uploaded file must be jpg")

	width, height, err := utils.GetImageDimensions(thumbnail)
	form.CheckField(width >= MinProfileWidth && height >= MinProfileWidth,
		"characterProfile",
		fmt.Sprintf("Ensure photo width is larger than %dx%d", MinProfileWidth, MinProfileWidth))
}

type userSignupForm struct {
	Username            string `form:"username"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) validateUserSignupForm(form *userSignupForm) {
	form.CheckField(validator.NotBlank(form.Username), "username", "Field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "Field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegEx), "email", "Please enter a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "Field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be 8-20 chararacters")
	form.CheckField(validator.MaxChars(form.Password, 20), "password", "Password Must be 8-20 chararacters")
}

type userLoginForm struct {
	Username            string `form:"username"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) validateUserLoginForm(form *userLoginForm) {
	form.CheckField(validator.NotBlank(form.Username), "username", "Field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "Field cannot be blank")
}

type categoryForm struct {
	ID                  int    `form:"id"`
	Name                string `form:"categoryName"`
	Action              string `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) validateCategoryForm(form *categoryForm) {
	form.CheckField(validator.NotBlank(form.Name), "categoryName", "Field cannot be blank")
}

type tagForm struct {
	ID                  int    `form:"id"`
	Name                string `form:"tag"`
	CategoryID          int    `form:"categoryId"`
	CategoryInput       string `form:"categoryInput"`
	Action              string `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) validateTagForm(form *tagForm) {
	form.CheckField(validator.NotBlank(form.Name), "tag", "Field cannot be blank")
	if form.CategoryID != 0 {
		form.CheckField(
			validator.BoolWithError(app.categories.Exists(form.CategoryID)),
			"categoryId",
			"Category does not exist. Please add it, then resubmit tag!",
		)
	}
}

type sketchTagsForm struct {
	TagIds              []int    `form:"tagId"`
	TagInputs           []string `form:"tagName"`
	validator.Validator `form:"-"`
}

func (app *application) validateSketchTagsForm(form *sketchTagsForm) {
	for i, tagId := range form.TagIds {
		if tagId == 0 {
			form.AddMultiFieldError("tagId", i, "Please input tag")
		}

		if exists, _ := app.tags.Exists(tagId); !exists {
			form.AddMultiFieldError("tagId", i, "Tag does not exist")
		}
	}
}

type showForm struct {
	ID                  int                   `form:"id"`
	Name                string                `form:"name"`
	Slug                string                `form:"slug"`
	ProfileImg          *multipart.FileHeader `img:"profileImg"`
	ProfileImgUrl       string                `form:"-"`
	Action              string                `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) validateShowForm(form *showForm) {
	form.CheckField(validator.NotBlank(form.Name), "name", "Field cannot be blank")
	if form.ID != 0 {
		form.CheckField(validator.NotBlank(form.Slug), "slug", "Field cannot be blank")
	}

	if form.ID == 0 {
		form.CheckField(form.ProfileImg != nil, "profileImg", "Please upload an image")
	}

	if form.ProfileImg == nil {
		return
	}

	thumbnail, err := form.ProfileImg.Open()
	if err != nil {
		form.AddFieldError("profileImg", "Unable to open file, ensure it is a valid jpg")
		return
	}
	defer thumbnail.Close()

	form.CheckField(validator.IsMime(thumbnail, "image/jpeg"), "thumbnail", "Uploaded file must be jpg")
}

type episodeForm struct {
	ID                  int                   `form:"id"`
	Number              int                   `form:"number"`
	Title               string                `form:"title"`
	URL                 string                `form:"url"`
	AirDate             string                `form:"airDate"`
	Thumbnail           *multipart.FileHeader `img:"thumbnail"`
	ThumbnailName       string                `form:"-"`
	SeasonId            int                   `form:"seasonId"`
	ThumbnailUrl        string                `form:"-"`
	Action              string                `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) validateEpisodeForm(form *episodeForm) {
	form.CheckField(form.Number != 0, "number", "Please enter a valid number")

	// validate episode number
	if form.SeasonId == 0 {
		form.AddNonFieldError("Season ID not defined")
	}

	season, err := app.shows.GetSeason(form.SeasonId)
	if err != nil {
		app.errorLog.Printf("Error getting season for episode form validation: $s", err)
		form.AddNonFieldError("Error getting season")
	}

	for _, ep := range season.Episodes {
		if safeDeref(ep.Number) == form.Number && safeDeref(ep.ID) != form.ID {
			form.AddFieldError("number", fmt.Sprintf("Episode number %d already exists for Season %d", *ep.Number, *season.Number))
		}
	}

	if form.AirDate != "" {
		form.CheckField(validator.ValidDate(form.AirDate), "airDate", "Date must be of the format YYYY-MM-DD")
	}

	if form.ID == 0 {
		form.CheckField(form.Thumbnail != nil, "thumbnail", "Please upload an image")
	}

	if form.Thumbnail == nil {
		return
	}

	thumbnail, err := form.Thumbnail.Open()
	if err != nil {
		form.AddFieldError("thumbnail", "Unable to open file, ensure it is a valid jpg")
		return
	}
	defer thumbnail.Close()

	form.CheckField(validator.IsMime(thumbnail, "image/jpeg"), "thumbnail", "Uploaded file must be jpg")
	width, height, err := utils.GetImageDimensions(thumbnail)
	if err != nil {
		form.AddFieldError("thumbnail", "Unable to determine image dimensions")
		return
	}

	form.CheckField(
		width >= LargeThumbnailWidth && height >= LargeThumbnailHeight,
		"thumbnail",
		fmt.Sprintf("Thumbnail dimensions must be at least %dx%d", LargeThumbnailWidth, LargeThumbnailHeight),
	)
}

type seriesForm struct {
	ID                  int                   `form:"id"`
	Title               string                `form:"title"`
	Thumbnail           *multipart.FileHeader `img:"thumbnail"`
	ThumbnailName       string                `form:"-"`
	ImageUrl            string                `form:"-"`
	Action              string                `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) validateSeriesForm(form *seriesForm) {
	form.CheckField(validator.NotBlank(form.Title), "title", "Please enter a title")

	if form.ID == 0 {
		form.CheckField(form.Thumbnail != nil, "thumbnail", "Please upload an image")
	}

	if form.Thumbnail == nil {
		return
	}

	thumbnail, err := form.Thumbnail.Open()
	if err != nil {
		form.AddFieldError("thumbnail", "Unable to open file, ensure it is a valid jpg")
		return
	}
	defer thumbnail.Close()

	form.CheckField(validator.IsMime(thumbnail, "image/jpeg"), "thumbnail", "Uploaded file must be jpg")
	width, height, err := utils.GetImageDimensions(thumbnail)
	if err != nil {
		form.AddFieldError("thumbnail", "Unable to determine image dimensions")
		return
	}

	form.CheckField(
		width >= LargeThumbnailWidth && height >= LargeThumbnailHeight,
		"thumbnail",
		fmt.Sprintf("Thumbnail dimensions must be at least %dx%d", LargeThumbnailWidth, LargeThumbnailHeight),
	)
}
