package models

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Character struct {
	ID          *int
	Slug        *string
	Name        *string
	Image       *string
	Description *string
	Portrayal   *Person
}

type CharacterModelInterface interface {
	Search(search string) ([]*Character, error)
	Exists(id int) (bool, error)
}

type CharacterModel struct {
	DB *pgxpool.Pool
}

func (m *CharacterModel) Search(query string) ([]*Character, error) {
	query = query + "%"
	stmt := `SELECT c.id, c.slug, c.name, c.img_name
			FROM character as c
			WHERE LOWER(name) LIKE LOWER($1)
			ORDER BY name`

	rows, err := m.DB.Query(context.Background(), stmt, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	characters := []*Character{}
	for rows.Next() {
		c := &Character{}
		err := rows.Scan(
			&c.ID, &c.Slug, &c.Name, &c.Image,
		)
		if err != nil {
			return nil, err
		}
		characters = append(characters, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return characters, nil
}

func (m *CharacterModel) Exists(id int) (bool, error) {
	stmt := `SELECT id FROM character WHERE id = $1`
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
