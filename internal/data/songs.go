package data

import (
	"context"
	"database/sql"
	"errors"
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&song.ID, &song.AddedAt, &song.Version)
}

func (m SongModel) Get(id int64) (*Song, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, added_at, title, year, duration, genres, version
FROM musics
WHERE id = $1`
	var song Song
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&song.ID,
		&song.AddedAt,
		&song.Title,
		&song.Year,
		&song.Duration,
		pq.Array(&song.Genres),
		&song.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &song, nil
}
func (m SongModel) Update(song *Song) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
UPDATE musics
SET title = $1, year = $2, duration = $3, genres = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		song.Title,
		song.Year,
		song.Duration,
		pq.Array(song.Genres),
		song.ID,
		song.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&song.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m SongModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the song ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
DELETE FROM songs
WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use ExecContext() and pass the context as the first argument.
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
