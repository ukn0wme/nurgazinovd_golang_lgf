package main

import (
	"fmt"
	"net/http"
	"nurgazinovd_golang_lg/internal/data"
	"nurgazinovd_golang_lg/internal/validator"
	"time"
)

func (app *application) createSongHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title    string        `json:"title"`
		Year     int32         `json:"year"`
		Duration data.Duration `json:"duration"`
		Genres   []string      `json:"genres"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	song := &data.Song{
		Title:    input.Title,
		Year:     input.Year,
		Duration: input.Duration,
		Genres:   input.Genres,
	}
	v := validator.New()
	if data.ValidateSong(v, song); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Songs.Insert(song)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/songs/%d", song.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"song": song}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showSongHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	song := data.Song{
		ID:       id,
		AddedAt:  time.Now(),
		Title:    "Oxxxymiron - Лига Опасного Интернета",
		Duration: 151,
		Genres:   []string{"rap", "hip-hop"},
		Version:  1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"song": song}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
