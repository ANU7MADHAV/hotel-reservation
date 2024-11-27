package main

import (
	"fmt"
	"net/http"
)

func (app *Applications) logError(r *http.Request, err error) {
	app.logger.Println(r, err)
}

func (app *Applications) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := Envelope{"message": message}

	err := app.writeJson(w, status, env)

	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}

}

func (app *Applications) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *Applications) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *Applications) methodNotAllowResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *Applications) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
