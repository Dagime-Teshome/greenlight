package main

import (
	"errors"
	"net/http"

	"github.com/Dagime-Teshome/greenlight/internal/data"
	"github.com/Dagime-Teshome/greenlight/internal/validator"
)

func (app *app) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	validator := validator.New()

	if data.ValidateUser(validator, user); !validator.Valid() {
		app.validationError(w, r, validator.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {

		case errors.Is(err, data.ErrDuplicateEmail):
			validator.AddError("email", "user with this email already exists")
			app.validationError(w, r, validator.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusCreated, Envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *app) getUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *app) updateUserHandler(w http.ResponseWriter, r *http.Request) {

}
func (app *app) listUsersHandler(w http.ResponseWriter, r *http.Request) {

}
func (app *app) deleteUserHandler(w http.ResponseWriter, r *http.Request) {

}
