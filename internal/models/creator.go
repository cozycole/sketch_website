package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Creator struct {
	ID              int
	Name            string
	URL             string
	ProfileImage    string
	Slug            string
	EstablishedDate time.Time
}

type CreatorModelInterface interface {
	Insert(name, url, imgName, imgExt string, establishedDate time.Time) (int, string, string, error)
	Get(id int) (*Creator, error)
	Exists(id int) (bool, error)
	GetBySlug(slug string) (*Creator, error)
}

type CreatorModel struct {
	DB *pgxpool.Pool
}

func (m *CreatorModel) GetBySlug(slug string) (*Creator, error) {
	stmt := `SELECT id, name, slug, page_url, profile_img, date_established FROM creator
	WHERE slug = $1`

	row := m.DB.QueryRow(context.Background(), stmt, slug)

	c := &Creator{}

	err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.URL, &c.ProfileImage, &c.EstablishedDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return c, nil

}

func (m *CreatorModel) Insert(name, url, slug, imgExt string, establishedDate time.Time) (int, string, string, error) {
	stmt := `
	INSERT INTO creator (name, page_url, date_established, slug, profile_img)
	VALUES ($1,$2,$3,
		CONCAT($4::text, '-', currval(pg_get_serial_sequence('creator', 'id'))),
		CONCAT($4::text, '-', currval(pg_get_serial_sequence('creator', 'id')), $5::text))
	RETURNING id, slug, profile_img;`

	var id int
	var fullImgName string

	row := m.DB.QueryRow(context.Background(), stmt, name, url, establishedDate, slug, imgExt)
	err := row.Scan(&id, &slug, &fullImgName)
	if err != nil {
		return 0, "", "", err
	}
	return id, slug, fullImgName, err
}

func (m *CreatorModel) Get(id int) (*Creator, error) {
	stmt := `SELECT id, name, slug, page_url, profile_img, date_established FROM creator
	WHERE id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	c := &Creator{}

	err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.URL, &c.ProfileImage, &c.EstablishedDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return c, nil
}

func (m *CreatorModel) Exists(id int) (bool, error) {
	stmt := `SELECT id FROM creator WHERE id = $1`
	row := m.DB.QueryRow(context.Background(), stmt, id)

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
