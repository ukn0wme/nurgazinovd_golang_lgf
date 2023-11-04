package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Songs SongModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Songs: SongModel{DB: db},
	}
}
