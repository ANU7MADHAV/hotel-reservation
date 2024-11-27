package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Applications) routes() *httprouter.Router {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.HealthCheck)

	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id", app.GetHotel)
	router.HandlerFunc(http.MethodPost, "/v1/hotels", app.CreateHotels)

	return router
}
