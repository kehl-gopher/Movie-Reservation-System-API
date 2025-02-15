package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routers() http.Handler {
	srv := httprouter.New()

	// get app health check...
	srv.HandlerFunc(http.MethodGet, "/", app.healthCheck)

	// movie routes
	srv.HandlerFunc(http.MethodPost, "/api/movies", app.CreateMovieRoutes)
	srv.HandlerFunc(http.MethodGet, "/api/movies/:id", app.GetMovieById)

	// user routes
	srv.HandlerFunc(http.MethodPost, "/api/user", app.registerUser)

	// serve image files
	fs := http.FileServer(http.Dir("uploads"))
	srv.Handler(http.MethodGet, "/images/", http.StripPrefix("/images/", fs))

	return app.recoverpanic(app.CORSMiddleWare(app.requestLogMiddleWare(srv)))
}
