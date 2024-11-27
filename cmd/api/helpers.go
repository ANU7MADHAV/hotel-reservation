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

type Envelope map[string]interface{}

func (app *Applications) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)

	if err != nil || id < 1 {
		app.logger.Print(err)
	}
	return id, nil
}

func (app *Applications) writeJson(w http.ResponseWriter, status int, data interface{}) error {
	js, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *Applications) readJson(w http.ResponseWriter, r *http.Request, message interface{}) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&message)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalError):
			if unmarshalError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type of filed %q", unmarshalError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json : unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json:unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}

	}
	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must only contains a single JSON value")
	}
	return nil

}
