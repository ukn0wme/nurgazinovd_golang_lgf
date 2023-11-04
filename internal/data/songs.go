package data

import (
	"database/sql"
	"github.com/lib/pq"
	"nurgazinovd_golang_lg/internal/validator"
	"time"
)

type Song struct {
	ID       int64     `json:"id"`
	AddedAt  time.Time `json:"-"`
	Title    string    `json:"title"`
	Year     int32     `json:"year,omitempty"`
	Duration Duration  `json:"duration,omitempty,string"`
	Genres   []string  `json:"genres,omitempty"`
	Version  int32     `json:"version"`
}

func ValidateSong(v *validator.Validator, song *Song) {
	v.Check(song.Title != "", "title", "must be provided")
	v.Check(len(song.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(song.Year != 0, "year", "must be provided")
	v.Check(song.Year >= 1888, "year", "must be greater than 1888")
	v.Check(song.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(song.Duration != 0, "duration", "must be provided")
	v.Check(song.Duration > 0, "duration", "must be a positive integer")

	v.Check(song.Genres != nil, "genres", "must be provided")
	v.Check(len(song.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(song.Genres) <= 3, "genres", "must not contain more than 3 genres")
	v.Check(validator.Unique(song.Genres), "genres", "must not contain duplicate values")
}

type SongModel struct {
	DB *sql.DB
}

func (m SongModel) Insert(song *Song) error {
	query := `
INSERT INTO songs (title, year, duration, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, added_at, version`
	args := []interface{}{song.Title, song.Year, song.Duration, pq.Array(song.Genres)}
	return m.DB.QueryRow(query, args...).Scan(&song.ID, &song.AddedAt, &song.Version)
}
func (m SongModel) Get(id int64) (*Song, error) {
	return nil, nil
}
func (m SongModel) Update(song *Song) error {
	return nil
}
func (m SongModel) Delete(id int64) error {
	return nil
}
