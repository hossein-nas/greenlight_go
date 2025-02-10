package main

import (
	"fmt"
	"net/http"

	"greenlight.hosseinnasiri.ir/internal/data"
	"greenlight.hosseinnasiri.ir/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err, nil)
		return
	}

	v := validator.New()

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)

	if err != nil {
		app.writeError(w, err.Error(), http.StatusInternalServerError, nil)
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, movie, headers)

	if err != nil {
		app.writeError(w, err.Error(), http.StatusInternalServerError, nil)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil || id < 1 {
		http.NotFound(w, r)

		return
	}

	movie, err := app.models.Movies.Get(id)

	if err != nil {
		app.writeError(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	err = app.writeResponse(w, movie, nil)

	if err != nil {
		errMsg := fmt.Errorf("the server encountered a problem and could not process your request. %s", err)
		app.writeError(w, errMsg.Error(), http.StatusInternalServerError, nil)
	}
}
