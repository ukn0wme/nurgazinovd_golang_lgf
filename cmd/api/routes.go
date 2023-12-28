package main

import (
	"expvar"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	// Use the requirePermission() middleware on each of the /v1/songs** endpoints,
	// passing in the required permission code as the first parameter.
	router.HandlerFunc(http.MethodGet, "/v1/songs", app.requirePermission("songs:read", app.listSongsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/songs", app.requirePermission("songs:write", app.createSongHandler))
	router.HandlerFunc(http.MethodGet, "/v1/songs/:id", app.requirePermission("songs:read", app.showSongHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/songs/:id", app.requirePermission("songs:write", app.updateSongHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/songs/:id", app.requirePermission("songs:write", app.deleteSongHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
