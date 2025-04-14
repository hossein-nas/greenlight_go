package main

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"greenlight.hosseinnasiri.ir/internal/jsonlog"
	"greenlight.hosseinnasiri.ir/internal/utils"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, jsonlog.LoggerProperties{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) writeError(w http.ResponseWriter, _err interface{}, status int, header http.Header) error {
	var data interface{} = ""

	if utils.CheckError(_err) {
		data = _err.(error).Error()
	} else if reflect.TypeOf(_err).Kind() == reflect.String {
		data = _err.(string)
		app.logger.PrintError(errors.New(data.(string)), jsonlog.LoggerProperties{})
	} else {
		data = _err
	}

	response := map[string]interface{}{
		"status": "ERROR",
		"data":   data,
	}

	err := app.writeJSON(w, status, response, header)

	if err != nil {
		fmt.Println(" ERROR: " + "err")
		return err
	}

	return nil
}

func (a *application) badRequestResponse(w http.ResponseWriter, err error, headers http.Header) {
	a.writeError(w, err.Error(), http.StatusBadRequest, headers)
}

func (a *application) notFoundResponse(w http.ResponseWriter, headers http.Header) {
	a.writeError(w, errors.New("the resource can not be found").Error(), http.StatusNotFound, headers)
}

func (a *application) serverErrorResponse(w http.ResponseWriter, msg interface{}, headers http.Header) {
	if msg == nil {
		msg = errors.New("server can not respond right now. please try again later").Error()
	}

	a.writeError(w, msg, http.StatusNotFound, headers)
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.writeError(w, errors, http.StatusUnprocessableEntity, nil)
}

func (app *application) editConflictResponse(w http.ResponseWriter, headers http.Header) {
	message := "unable to update the record due to an edit conflict, please try again"

	app.writeError(w, message, http.StatusConflict, headers)
}

func (app *application) invalidCredentialsResponse(w http.ResponseWriter, headers http.Header) {
	message := "invalid authentication credentials"

	app.writeError(w, message, http.StatusUnauthorized, headers)
}

func (app *application) invalidAuthenticationResponse(w http.ResponseWriter, headers http.Header) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"

	app.writeError(w, message, http.StatusUnauthorized, headers)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {

	message := "rate limit exceeded"
	app.writeError(w, message, http.StatusTooManyRequests, nil)
}
