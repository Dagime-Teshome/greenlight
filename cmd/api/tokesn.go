package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Dagime-Teshome/greenlight/internal/data"
	"github.com/Dagime-Teshome/greenlight/internal/validator"
)

// 1.the client sends a JSON request to a new POST/v1/tokens/authentication endpoint
// containing their credentials (email and password).
// 2. We look up the user record based on the email, and check if the password provided is
// the correct one for the user. If it’s not, then we send an error response.
// 3. If the password is correct, we use our app.models.Tokens.New() method to generate a
// token with an expiry time of 24 hours and the scope "authentication".
// 4. We send this authentication token back to the client in a JSON response body.

func (app *app) authenticationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.validationError(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}
	validFor := 24 * time.Hour
	token, err := app.models.Tokens.New(user.ID, validFor, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusCreated, Envelope{"token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
