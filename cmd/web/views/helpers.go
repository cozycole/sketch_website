package views

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"sketchdb.cozycole.net/internal/models"
)

func printCast(cast []*models.CastMember) string {
	castList := ""
	var personIds []int
	for i, cm := range cast {
		if cm.Actor == nil ||
			cm.Actor.ID == nil ||
			intSliceContains(personIds, *cm.Actor.ID) {
			continue
		}

		name := printPersonName(cm.Actor)
		if name != "" {
			if i != 0 {
				name = ", " + name
			}
			castList += name
		}

		personIds = append(personIds, *cm.Actor.ID)
	}

	return castList
}

func printPersonName(a *models.Person) string {
	if a == nil || a.First == nil {
		return ""
	}

	name := *a.First
	if a.Last == nil {
		return name
	}
	return name + " " + *a.Last
}

func humanDate(t *time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("Jan 2, 2006")
}

func createEpisodeTitle(episode *models.Episode) string {
	title := safeDerefString(episode.Title)
	if title == "" {
		if episode.Number == nil {
			title = "Episode ?"
		} else {
			title = fmt.Sprintf("Episode %d", *episode.Number)
		}
	}
	return title
}

func seasonEpisodeInfo(episode *models.Episode) string {
	var info string

	if episode.SeasonNumber != nil &&
		episode.Number != nil {

		info = fmt.Sprintf(
			"S%d · E%d · %s",
			*episode.SeasonNumber,
			*episode.Number,
			sketchCountLabel(len(episode.Sketches)),
		)
	}

	return info
}

func seasonEpisodeUrl(episode *models.Episode) string {
	var url string
	if episode.ShowID != nil &&
		episode.ShowSlug != nil &&
		episode.SeasonNumber != nil &&
		episode.Number != nil {

		url = fmt.Sprintf(
			"/show/%d/%s/season/%d/episode/%d",
			*episode.ShowID,
			*episode.ShowSlug,
			*episode.SeasonNumber,
			*episode.Number,
		)
	}

	return url
}

func getAge(birthDate *time.Time) int {
	today := time.Now()
	age := today.Year() - birthDate.Year()

	if today.YearDay() < birthDate.YearDay() {
		age--
	}

	return age
}

func BuildURL(baseURL string, currentPage int, filter *models.Filter) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	params := filter.Params()
	params.Add("page", strconv.Itoa(currentPage))
	if filter.Query != "" {
		params.Set("query", filter.Query)
	}

	u.RawQuery = params.Encode()

	return u.String(), nil
}

func safeDerefString(str *string) string {
	if str != nil {
		return *str
	}
	return ""
}

func sketchCountLabel(count int) string {
	labelString := "%d Sketch"
	if count != 1 {
		labelString += "es"
	}

	return fmt.Sprintf(labelString, count)
}

func intSliceContains(list []int, value int) bool {
	for _, n := range list {
		if n == value {
			return true
		}
	}
	return false
}

func getSeasonSketchCount(season *models.Season) int {
	var count int
	for _, ep := range season.Episodes {
		count += len(ep.Sketches)
	}

	return count
}

func getShowEpisodeCount(show *models.Show) int {
	var count int
	for _, season := range show.Seasons {
		count += len(season.Episodes)
	}

	return count
}

func getShowSketchCount(show *models.Show) int {
	var count int
	for _, season := range show.Seasons {
		for _, ep := range season.Episodes {
			count += len(ep.Sketches)
		}
	}

	return count
}
