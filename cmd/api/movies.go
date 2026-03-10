package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dagime-Teshome/greenlight/internal/data"
)

func (app *app) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "created a movie")

}
func (app *app) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}
	err = app.writeJson(w, http.StatusOK, Envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Print(err)
		app.serverErrorResponse(w, r, err)
	}
}
