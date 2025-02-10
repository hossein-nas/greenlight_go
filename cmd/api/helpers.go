package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("invalid :id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, header http.Header) error {
	js, err := json.Marshal(data)

	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range header {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) writeResponse(w http.ResponseWriter, data interface{}, header http.Header) error {
	response := map[string]interface{}{
		"status": "OK",
		"data":   data,
	}

	return app.writeJSON(w, http.StatusOK, response, header)
}

func (app *application) writeError(w http.ResponseWriter, error interface{}, status int, header http.Header) error {

	response := map[string]interface{}{
		"status": "ERROR",
		"data":   error,
	}

	err := app.writeJSON(w, status, response, header)

	if err != nil {
		return err
	}
	http.Error(w, "", status)

	return nil
}

func (a *application) badRequestResponse(w http.ResponseWriter, err error, headers http.Header) {
	a.writeError(w, err.Error(), http.StatusBadRequest, headers)
}

func (a *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.DisallowUnknownFields()

	err := bodyDecoder.Decode(dst)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %s)", syntaxError)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("bobdy must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	var _dst interface{}
	err = bodyDecoder.Decode(&_dst)

	if !errors.Is(err, io.EOF) {
		return errors.New("json structure is malformed")
	}

	return nil
}
