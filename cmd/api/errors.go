package main

import "net/http"

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.writeError(w, errors, http.StatusUnprocessableEntity, nil)
}
