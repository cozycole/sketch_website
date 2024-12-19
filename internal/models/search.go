package models

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Logic for comprehensive search of entire all resources

type SearchModel struct {
	DB *pgxpool.Pool
}

type SearchResult struct {
	Type *string
	Id   *int
	Name *string
	Slug *string
	Img  *string
	Rank *float32
}

func (m *SearchModel) Search(query string) ([]*SearchResult, error) {
	stmt := `
		SELECT 'video' AS type, id, title AS name, slug, thumbnail_name as img, ts_rank(search_vector, plainto_tsquery('english', $1)) AS rank
		FROM video
		WHERE search_vector @@ plainto_tsquery('english', $1)
		UNION ALL
		SELECT 'person' AS type, id, CONCAT(first, ' ', last) as name, slug, profile_img as img, ts_rank(search_vector, plainto_tsquery('english', $1)) AS rank
		FROM person
		WHERE search_vector @@ plainto_tsquery('english', $1)
		UNION ALL
		SELECT 'character' AS type, id, name, slug, img_name as img, ts_rank(search_vector, plainto_tsquery('english', $1)) AS rank
		FROM character
		WHERE search_vector @@ plainto_tsquery('english', $1)
		UNION ALL
		SELECT 'creator' AS type, id, name, slug, profile_img as img, ts_rank(search_vector, plainto_tsquery('english', $1)) AS rank
		FROM creator
		WHERE search_vector @@ plainto_tsquery('english', $1)
		ORDER BY rank, name DESC
	`

	rows, err := m.DB.Query(context.Background(), stmt, query)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	results := []*SearchResult{}
	for rows.Next() {
		sr := &SearchResult{}
		err := rows.Scan(
			&sr.Type, &sr.Id, &sr.Name, &sr.Slug, &sr.Img, &sr.Rank,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, sr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
