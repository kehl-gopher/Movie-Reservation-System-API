package main

import "net/http"

// handle server error response >500 from the server
func (app *application) serverErrorResponse(w http.ResponseWriter, error interface{}) {
	app.writeResponse(w, http.StatusInternalServerError, toJson{"error": error})
}

// handle 404 response from the server
func (app *application) notFoundResponse(w http.ResponseWriter, err error) {
	app.writeResponse(w, http.StatusNotFound, toJson{"error": err.Error()})
}

// handle 400 response from the server
func (app *application) badErrorResponse(w http.ResponseWriter, message interface{}) {
	app.writeResponse(w, http.StatusBadRequest, toJson{"error": message})
}

func (app *application) validationErrorResponse(w http.ResponseWriter, message interface{}) {
	app.writeResponse(w, http.StatusUnprocessableEntity, toJson{"errors": message})
}
