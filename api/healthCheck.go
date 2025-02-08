package main

import (
	"fmt"
	"net/http"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/utils"
)

// update this routes to perform it's duty like a bitch
func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.URL.RequestURI() != "/" {
		app.notFoundResponse(w, fmt.Errorf("this resource url cannot be found"))
		return
	}
	env, err := utils.ReadEnvVariable("APP_ENV")
	if err != nil {
		env = "DEVELOPMENT" // default environment set
	}
	app.writeResponse(w, http.StatusOK, toJson{"Version": "0.0.1", "App Environment": env})
}
