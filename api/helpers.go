package main

import (
	"net/http"
	"time"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/logs"
	"github.com/rs/zerolog"
)

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
		logs.Error().Str("method", r.Method).
			Int("status", statusCode).
			Str("url", r.URL.RequestURI()).
			Dur("elapsed time", time.Since(timeS)).
			Str("message", message).Send()

	}
}
