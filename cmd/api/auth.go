package main

import (
	"errors"
	"net/http"
	"time"

	"greenlight.hosseinnasiri.ir/internal/data"
	"greenlight.hosseinnasiri.ir/internal/validator"
)

func (app *application) createSignin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, err, nil)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, nil)
		default:
			app.serverErrorResponse(w, err, nil)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)

	if err != nil {
		app.serverErrorResponse(w, err, nil)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, nil)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)

	if err != nil {
		app.serverErrorResponse(w, err, nil)
		return
	}

	app.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"authentication_token": token,
	}, nil)
}
