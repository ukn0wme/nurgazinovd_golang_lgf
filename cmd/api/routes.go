package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/songs", app.requireActivatedUser(app.listSongsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/songs", app.requireActivatedUser(app.createSongHandler))
	router.HandlerFunc(http.MethodGet, "/v1/songs/:id", app.requireActivatedUser(app.showSongHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/songs/:id", app.requireActivatedUser(app.updateSongHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/songs/:id", app.requireActivatedUser(app.deleteSongHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
