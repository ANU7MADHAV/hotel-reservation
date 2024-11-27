package main

import (
	"net/http"
)

func (app *Applications) HealthCheck(w http.ResponseWriter, r *http.Request) {
	env := Envelope{
		"status": "available",
		"system info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJson(w, http.StatusOK, env)
	if err != nil {

		app.serverErrorResponse(w, r, err)
	}
}
