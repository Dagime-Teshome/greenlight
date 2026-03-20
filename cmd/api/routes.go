package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *app) routes() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	// movies end points
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies/", app.createMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMoveHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)
	return router
}
