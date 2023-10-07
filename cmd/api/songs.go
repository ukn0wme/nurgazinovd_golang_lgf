package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nurgazinovd_golang_lg/internal/data"
	"time"
)

func (app *application) createSongHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title    string   `json:"title"`
		Year     int32    `json:"year"`
		Duration int32    `json:"duration"`
		Genres   []string `json:"genres"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
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
	} // Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"song": song}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
