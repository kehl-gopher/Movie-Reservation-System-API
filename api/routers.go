package main

import "net/http"

func (app *application) routers() http.Handler {
	srv := http.NewServeMux()
	srv.HandleFunc("GET /", app.healthCheck)

	return app.CORSMiddleWare(app.requestLogMiddleWare(srv))
}
