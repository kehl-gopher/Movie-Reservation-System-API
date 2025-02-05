package main

import (
	"net/http"
	"time"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/utils"
)

type MiddleLogRequest struct {
	resp       http.ResponseWriter
	statusCode int
}

func (m *MiddleLogRequest) Write(msg []byte) (int, error) {
	return m.resp.Write(msg)
}

func (m *MiddleLogRequest) WriteHeader(statusCode int) {
	m.statusCode = statusCode
	m.resp.WriteHeader(statusCode)
}

func (m *MiddleLogRequest) Header() http.Header {
	return m.resp.Header()
}

func (app *application) requestLogMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t := time.Now()
		lrw := MiddleLogRequest{resp: w, statusCode: http.StatusOK}
		next.ServeHTTP(&lrw, r)
		// log user request this should only happen in development
		env, _ := utils.ReadEnvVariable("APP_ENV")
		if env == "DEVELOPMENT" {
			app.logRequest(lrw, r, t)
		}
	})
}

// Headers and CORS middleware
func (app *application) CORSMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle preflight options
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		origin := r.Header.Get("Origin")
		if origin != "" {

			for idx := range app.trustedOrigins {
				if origin != app.trustedOrigins[idx] {
					continue
				}
				w.Header().Set("Access-Control-Origin", origin)

				// handle preflight request
				if r.Method == http.MethodOptions || r.Header.Get("Access-Control-Request-Method") != "" {
					w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
					w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
					w.WriteHeader(http.StatusOK)
				}

				break
			}
		}

		next.ServeHTTP(w, r)
	})
}

// recover application panic error... and return 500 server response
func recoverpanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				// send a 500 Internal server error to the user
				// as a request
			}
		}()
		next.ServeHTTP(w, r)
	})
}
