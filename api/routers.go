package main

import "net/http"

func (app *application) routers() http.Handler {
	srv := http.NewServeMux()

	// get app health check...
	srv.HandleFunc("/", app.healthCheck)

	// movie routes
	srv.HandleFunc("POST /movies", app.CreateMovieRoutes)

	// serve image files
	fs := http.FileServer(http.Dir("uploads"))
	srv.Handle("/images/", http.StripPrefix("/images/", fs))

	return app.recoverpanic(app.CORSMiddleWare(app.requestLogMiddleWare(srv)))
}
