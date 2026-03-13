package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dagime-Teshome/greenlight/internal/data"
	"github.com/Dagime-Teshome/greenlight/internal/validator"
)

func (app *app) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year,omitempty"`
		Runtime data.Runtime `json:"runtime,omitempty"`
		Genres  []string     `json:"genres,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	data.ValidateMovie(v, movie)

	if !v.Valid() {
		app.validationError(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)

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
