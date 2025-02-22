package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/logs"
	"github.com/rs/zerolog"
)

// run background functions
func (app *application) Background(fn func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				err := fmt.Errorf("%s", err)
				app.log.FatalLog(err)
			}
		}()
		fn()
	}()
}

// extrach user params from the url
func getparams(r *http.Request) (int, error) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	uid, err := strconv.Atoi(id)
	if err != nil || uid <= 0 {
		return 0, err
	}

	return uid, nil
}

// read from json
func readFromJson(r *http.Request, toStruct interface{}) error {
	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(toStruct)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidMarshalError *json.InvalidUnmarshalError

		// check errors type
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed JSON at (%d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type character %d", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body cannot be empty")
		case errors.As(err, &invalidMarshalError):
			panic(err)
		default:
			return err
		}

	}
	return nil
}

// send json response to client
func (app *application) writeResponse(w http.ResponseWriter, statusCode int, message toJson) (int, error) {

	byt, err := writeToJson(message)

	if err != nil {
		return 0, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return w.Write(byt)
}

// write to json...
func writeToJson(value interface{}) ([]byte, error) {
	byte, err := json.Marshal(value)

	if err != nil {
		return nil, err
	}
	return byte, err
}

// log user request in application middleware
func (app *application) logRequest(w MiddleLogRequest, r *http.Request, timeS time.Time) {
	statusCode := w.statusCode
	message := http.StatusText(statusCode)
	logs := logs.Get(zerolog.InfoLevel)

	switch {
	case statusCode >= 200 && statusCode <= 400:
		logs.Info().Str("method", r.Method).
			Int("status", statusCode).
			Str("url", r.URL.RequestURI()).
			Dur("elapsed time", time.Since(timeS)).
			Str("message", message).Send()
	case statusCode >= 400 && statusCode <= 500:
		logs.Warn().Str("method", r.Method).
			Int("status", statusCode).
			Str("url", r.URL.RequestURI()).
			Dur("elapsed time", time.Since(timeS)).
			Str("message", message).Send()
	case statusCode >= 500 && statusCode <= 600:
		logs.Error().Stack().Str("method", r.Method).
			Int("status", statusCode).
			Str("url", r.URL.RequestURI()).
			Dur("elapsed time", time.Since(timeS)).
			Str("message", message).Send()

	}
}
