package main

import (
	"fmt"
	"hotel-reservation/internal/data"
	"net/http"
	"time"
)

func (app *Applications) GetHotel(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	hotel := &data.HotelsType{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "Stay",
		Address:   "Noida sector 62",
		Location:  "Noida",
	}

	err = app.writeJson(w, http.StatusOK, Envelope{"hotels": hotel})

	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Applications) CreateHotels(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		Location string `json:"location"`
	}

	err := app.readJson(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	hotel := &data.HotelsType{
		Name:     input.Name,
		Address:  input.Address,
		Location: input.Location,
	}

	err = app.data.Hotels.Insert(hotel)

	fmt.Println("err", err)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", hotel.ID))

	err = app.writeJson(w, http.StatusCreated, Envelope{"hotel": hotel})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
