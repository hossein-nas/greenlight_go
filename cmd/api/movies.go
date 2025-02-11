package main

import (
	"errors"
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

	err = app.writeResponseCreated201(w, movie, headers)

	if err != nil {
		app.writeError(w, err.Error(), http.StatusInternalServerError, nil)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil || id < 1 {
		app.notFoundResponse(w, nil)

		return
	}

	movie, err := app.models.Movies.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, nil)
		default:
			app.serverErrorResponse(w, nil, nil)
		}
		return
	}

	err = app.writeResponse(w, movie, nil)

	if err != nil {
		errMsg := fmt.Errorf("the server encountered a problem and could not process your request. %s", err)
		app.serverErrorResponse(w, errMsg, nil)
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {

	movies, err := app.models.Movies.List()

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, nil)
		default:
			app.serverErrorResponse(w, nil, nil)
		}
		return
	}

	err = app.writeResponse(w, movies, nil)

	if err != nil {
		errMsg := fmt.Errorf("the server encountered a problem and could not process your request. %s", err)
		app.serverErrorResponse(w, errMsg, nil)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil || id < 1 {
		app.notFoundResponse(w, nil)

		return
	}

	movieToBeChanged, err := app.models.Movies.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, nil)
		default:
			app.serverErrorResponse(w, nil, nil)
		}
		return
	}

	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, err, nil)
		return
	}

	if input.Title != nil {
		movieToBeChanged.Title = *input.Title
	}
	if input.Year != nil {
		movieToBeChanged.Year = *input.Year
	}
	if input.Runtime != nil {
		movieToBeChanged.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movieToBeChanged.Genres = input.Genres
	}

	v := validator.New()

	if data.ValidateMovie(v, movieToBeChanged); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movieToBeChanged)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, nil)
		default:
			app.serverErrorResponse(w, nil, nil)
		}
		return

	}

	err = app.writeResponse(w, movieToBeChanged, nil)

	if err != nil {
		errMsg := fmt.Errorf("the server encountered a problem and could not process your request. %s", err)
		app.serverErrorResponse(w, errMsg, nil)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil || id < 1 {
		app.notFoundResponse(w, nil)

		return
	}

	err = app.models.Movies.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, nil)
		default:
			app.serverErrorResponse(w, err, nil)
		}
		return
	}

	response := fmt.Sprintf("the movie with id of %d has removed", id)

	err = app.writeResponse(w, response, nil)

	if err != nil {
		errMsg := fmt.Errorf("the server encountered a problem and could not process your request. %s", err)
		app.serverErrorResponse(w, errMsg, nil)
	}
}
