package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/songs", app.createSongHandler)
	router.HandlerFunc(http.MethodGet, "/v1/songs/:id", app.showSongHandler)
	// Add the route for the PUT /v1/songs/:id endpoint.
	router.HandlerFunc(http.MethodPut, "/v1/songs/:id", app.updateSongHandler)
	return router
}
