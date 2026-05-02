package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *app) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	// movies end points
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requireActivatedUser(app.requirePermissionCode("movies:read", app.listMovieHandler)))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requireActivatedUser(app.requirePermissionCode("movies:write", app.createMovieHandler)))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requireActivatedUser(app.requirePermissionCode("movie:read", app.showMovieHandler)))
	router.HandlerFunc(http.MethodPost, "/v1/movies/", app.requireActivatedUser(app.requirePermissionCode("movie:write", app.createMovieHandler)))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requireActivatedUser(app.requirePermissionCode("movies:write", app.updateMoveHandler)))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requireActivatedUser(app.requirePermissionCode("movies:write", app.deleteMovieHandler)))
	// users end points
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	// Authentication
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.authenticationHandler)
	return app.recoverPanic(app.rateLimiter(app.authenticate(router)))
}
