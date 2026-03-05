package main

import (
	"fmt"
	"net/http"
)

func (app *app) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "created a movie")

}
func (app *app) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Printf("showing movie details of %d", id)
}
