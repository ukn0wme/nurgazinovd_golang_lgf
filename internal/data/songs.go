package data

import (
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
