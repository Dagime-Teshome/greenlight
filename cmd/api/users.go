package main

import (
	"errors"
	"net/http"
	"time"

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

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	err = app.writeJson(w, http.StatusAccepted, Envelope{"user": user}, nil)
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

func (app *app) activateUserHandler(w http.ResponseWriter, r *http.Request) {

	//"token": "{{.activationToken}}"

	var input struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	validator := validator.New()
	data.ValidateTokenPlaintext(validator, input.Token)
	if !validator.Valid() {
		app.validationError(w, r, validator.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.Token)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			validator.AddError("token", "invalid or expired activation token")
			app.validationError(w, r, validator.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, Envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
