package main

import (
	"fmt"
	"net/http"
	"nurgazinovd_golang_lg/internal/data"
	"time"
)

// Add a createMovieHandler for the "POST /v1/movies" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createSongHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new song")
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showSongHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	song := data.Song{
		ID:       id,
		AddedAt:  time.Now(),
		Title:    "Oxxxymiron - Лига Опасного Интернета",
		Duration: 151,
		Genres:   []string{"rap", "hip-hop"},
		Version:  1,
	} // Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"song": song}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
