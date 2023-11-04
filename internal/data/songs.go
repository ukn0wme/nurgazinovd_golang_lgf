package data

import (
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
	return m.DB.QueryRow(query, args...).Scan(&song.ID, &song.AddedAt, &song.Version)
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
	err := m.DB.QueryRow(query, id).Scan(
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
	// Otherwise, return a pointer to the song struct.
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
	err := m.DB.QueryRow(query, args...).Scan(&song.Version)
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
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the songs table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
