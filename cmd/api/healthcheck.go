package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	location, err := time.LoadLocation("Asia/Tehran")

	if err != nil {
		http.Error(w, "Loading timezone errored.", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
		"time":        now.UTC(),
		"time(IR)":    now.In(location),
		"byte":        []byte{'h', 'e', 'l', 'l', 'o'},
	}

	err = app.writeResponse(w, data, nil)

	if err != nil {
		errMsg := fmt.Errorf("the server encountered a problem and could not process your request. %s", err)
		app.writeError(w, errMsg.Error(), http.StatusInternalServerError, nil)
		return
	}

}
