package mocks

import (
	"time"

	"sketchdb.cozycole.net/internal/models"
)

var mockActor = &models.Actor{
	ID:         1,
	First:      "Brad",
	Last:       "Pitt",
	ProfileImg: "brad-pitt-1.jpg",
	BirthDate:  time.Now(),
}

type ActorModel struct{}

func (m *ActorModel) Get(id int) (*models.Actor, error) {
	return mockActor, nil
}

func (m *ActorModel) Exists(id int) (bool, error) {
	switch id {
	case 1, 2, 3:
		return true, nil
	default:
		return false, nil
	}
}

func (m *ActorModel) Insert(first, last, imgName, imgExt string, birthDate time.Time) (int, string, error) {
	return 1, "brad-pitt-1.jpg", nil
}
