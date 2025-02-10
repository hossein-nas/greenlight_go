package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.hosseinnasiri.ir/internal/data"
	"greenlight.hosseinnasiri.ir/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  data.Genres  `json:"genres"`
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

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil || id < 1 {
		http.NotFound(w, r)

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

	err = app.writeResponse(w, movie, nil)

	if err != nil {

		errMsg := fmt.Errorf("the server encountered a problem and could not process your request. %s", err)
		app.writeError(w, errMsg.Error(), http.StatusInternalServerError, nil)
	}
}
