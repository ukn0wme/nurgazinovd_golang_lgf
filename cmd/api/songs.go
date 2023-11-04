package main

import (
	"errors"
	"fmt"
	"net/http"
	"nurgazinovd_golang_lg/internal/data"
	"nurgazinovd_golang_lg/internal/validator"
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
	song, err := app.models.Songs.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"song": song}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) updateSongHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the song ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the existing song record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	song, err := app.models.Songs.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Declare an input struct to hold the expected data from the client.
	var input struct {
		Title    string        `json:"title"`
		Year     int32         `json:"year"`
		Duration data.Duration `json:"duration"`
		Genres   []string      `json:"genres"`
	}
	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the request body to the appropriate fields of the song
	// record.
	song.Title = input.Title
	song.Year = input.Year
	song.Duration = input.Duration
	song.Genres = input.Genres
	// Validate the updated song record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	if data.ValidateSong(v, song); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Pass the updated song record to our new Update() method.
	err = app.models.Songs.Update(song)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Write the updated song record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"song": song}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteSongHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the song ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the song from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = app.models.Songs.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "song successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
