package main

import (
	"net/http"
)

func (app *app) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	hj := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJson(w, http.StatusOK, Envelope{"HealthCheck": hj}, nil)
	if err != nil {
		app.logger.Print(err)
		app.serverErrorResponse(w, r, err)
	}

}
