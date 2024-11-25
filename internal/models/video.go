package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CastMember struct {
	Position  *int
	Actor     *Person
	Character *Character
	// TODO: add attributes about the role
}

type Video struct {
	ID          int
	Title       string
	URL         string
	Slug        string
	Thumbnail   string
	Rating      string
	Description *string
	UploadDate  *time.Time
	Creator     *Creator
	Cast        []*CastMember
}

type VideoModelInterface interface {
	Search(search string, offset int) ([]*Video, error)
	GetAll(offset int) ([]*Video, error)
	Get(id int) (*Video, error)
	GetBySlug(slug string) (*Video, error)
	GetByCreator(id int) ([]*Video, error)
	Insert(title, video_url, rating, imgName, imgExt string, upload_date time.Time) (int, string, string, error)
	InsertVideoCreatorRelation(vidId, creatorId int) error
	InsertVideoPersonRelation(vidId, personId, position int, characterId *int, imgName string) error
}

type VideoModel struct {
	DB         *pgxpool.Pool
	ResultSize int
}

func (m *VideoModel) Search(search string, offset int) ([]*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.thumbnail_name, v.upload_date,
		c.id, c.name, c.profile_img
		FROM video AS v
		LEFT JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		LEFT JOIN creator as c
		ON vcr.creator_id = c.id
		WHERE search_vector @@ to_tsquery('english', $1)
		LIMIT $2
		OFFSET $3;
	`
	rows, err := m.DB.Query(context.Background(), stmt, search, m.ResultSize, m.ResultSize*offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	videos := []*Video{}
	for rows.Next() {
		v := &Video{}
		err := rows.Scan(&v.ID, &v.Title, &v.URL, &v.Thumbnail)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return videos, nil
}

func (m *VideoModel) GetIdBySlug(slug string) (int, error) {
	stmt := `SELECT v.id FROM video as v WHERE v.slug = $1`
	id_row := m.DB.QueryRow(context.Background(), stmt, slug)

	var id int
	err := id_row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *VideoModel) GetBySlug(slug string) (*Video, error) {
	id, err := m.GetIdBySlug(slug)
	if err != nil {
		return nil, err
	}

	return m.Get(id)
}

func (m *VideoModel) GetByCreator(id int) ([]*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.slug, v.thumbnail_name, v.upload_date,
		c.id, c.name, c.page_url, c.slug, c.profile_img
		FROM video AS v
		JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		JOIN creator as c
		ON vcr.creator_id = c.id
		WHERE c.id = $1
	`
	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	videos := []*Video{}
	for rows.Next() {
		v := &Video{}
		c := &Creator{}
		err := rows.Scan(
			&v.ID, &v.Title, &v.URL, &v.Slug, &v.Thumbnail, &v.UploadDate,
			&c.ID, &c.Name, &c.URL, &c.Slug, &c.ProfileImage,
		)

		if err != nil {
			return nil, err
		}
		v.Creator = c
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return videos, nil
}

func (m *VideoModel) Get(id int) (*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.thumbnail_name, v.upload_date, v.description, v.pg_rating,
			c.id, c.name, c.profile_img
		FROM video AS v
		LEFT JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		LEFT JOIN creator as c
		ON vcr.creator_id = c.id
		WHERE v.id = $1
	`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	v := &Video{}
	c := &Creator{}
	err := row.Scan(
		&v.ID, &v.Title, &v.URL, &v.Thumbnail, &v.UploadDate, &v.Description, &v.Rating,
		&c.ID, &c.Name, &c.ProfileImage,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	v.Creator = c

	members, err := m.GetCastMembers(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	v.Cast = members

	return v, nil
}

func (m *VideoModel) GetCastMembers(video_id int) ([]*CastMember, error) {
	stmt := `
		SELECT p.id, p.slug, p.first, p.last, p.birthdate,
			p.description, p.profile_img, vpr.position,
			ch.id as char_id, ch.name as char_name, ch.img_name as char_img
		FROM video AS v
		LEFT JOIN video_person_rel as vpr
		ON v.id = vpr.video_id
		LEFT JOIN person as p
		ON vpr.person_id = p.id
		LEFT JOIN character as ch
		ON vpr.character_id = ch.id
		WHERE v.id = $1
	`
	rows, err := m.DB.Query(context.Background(), stmt, video_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	members := []*CastMember{}
	for rows.Next() {
		p := &Person{}
		ch := &Character{}
		cm := &CastMember{}
		err := rows.Scan(
			&p.ID, &p.Slug, &p.First, &p.Last, &p.BirthDate, &p.Description, &p.ProfileImg,
			&cm.Position, &ch.ID, &ch.Name, &ch.Image,
		)
		if err != nil {
			return nil, err
		}

		cm.Actor = p
		cm.Character = ch
		members = append(members, cm)
	}

	return members, nil
}

// TODO: make DRY later
func (m *VideoModel) GetAll(offset int) ([]*Video, error) {
	stmt := `
		SELECT v.id, v.title, v.video_url, v.slug, v.thumbnail_name, v.upload_date,
		c.id, c.name, c.page_url, c.slug, c.profile_img
		FROM video AS v
		JOIN video_creator_rel as vcr
		ON v.id = vcr.video_id
		JOIN creator as c
		ON vcr.creator_id = c.id
		LIMIT $1
		OFFSET $2;
	`
	rows, err := m.DB.Query(context.Background(), stmt, m.ResultSize, m.ResultSize*offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	videos := []*Video{}
	for rows.Next() {
		v := &Video{}
		c := &Creator{}
		err := rows.Scan(
			&v.ID, &v.Title, &v.URL, &v.Slug, &v.Thumbnail, &v.UploadDate,
			&c.ID, &c.Name, &c.URL, &c.Slug, &c.ProfileImage,
		)
		if err != nil {
			return nil, err
		}
		v.Creator = c
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return videos, nil
}

func (m *VideoModel) Insert(title, video_url, rating, slug, imgExt string, upload_date time.Time) (int, string, string, error) {
	stmt := `
	INSERT INTO video (title, video_url, upload_date, pg_rating, slug, thumbnail_name)
	VALUES ($1,$2,$3,$4,
		CONCAT($5::text, '-', currval(pg_get_serial_sequence('video', 'id'))), 
		CONCAT($5::text, '-', currval(pg_get_serial_sequence('video', 'id')), $6::text))
	RETURNING id, slug, thumbnail_name;`

	result := m.DB.QueryRow(
		context.Background(), stmt, title,
		video_url, upload_date, rating,
		slug, imgExt,
	)

	var id int
	var imgName string
	err := result.Scan(&id, &slug, &imgName)
	if err != nil {
		return 0, "", "", err
	}
	return id, slug, imgName, nil
}

func (m *VideoModel) InsertVideoCreatorRelation(vidId, creatorId int) error {
	stmt := `INSERT INTO video_creator_rel (video_id, creator_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(context.Background(), stmt, vidId, creatorId)
	return err
}

func (m *VideoModel) InsertVideoPersonRelation(vidId, personId, position int, characterId *int, imgName string) error {
	stmt := `INSERT INTO video_person_rel (video_id, person_id, characterId, position, img_name) VALUES ($1, $2, $3, $4, $5)`
	_, err := m.DB.Exec(context.Background(), stmt, vidId, personId, characterId, position, imgName)
	return err
}
