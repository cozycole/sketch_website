package models

import (
	"errors"
)

var (
	ErrNoRecord               = errors.New("models: no matching record found")
	ErrInvalidCredentials     = errors.New("models: invalid crednetials")
	ErrDuplicateEmail         = errors.New("models: duplicate email")
	ErrDuplicateUsername      = errors.New("models: duplicate username")
	ErrDuplicateVidCreatorRel = errors.New("models: duplicate sketch creator relation")
	ErrDuplicateVidPersonRel  = errors.New("models: duplicate sketch person relation")
)
