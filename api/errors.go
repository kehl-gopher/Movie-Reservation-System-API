package main

import "net/http"

// handle server error response >500 from the server
func (app *application) serverErrorResponse(w http.ResponseWriter, error error) {
	app.writeResponse(w, http.StatusInternalServerError, toJson{"error": error.Error()})
}

// handle 404 response from the server
func (app *application) notFoundResponse(w http.ResponseWriter) {
	app.writeResponse(w, http.StatusNotFound, toJson{"error": http.StatusText(http.StatusNotFound)})
}

// handle 400 response from the server
func (app *application) badErrorResponse(w http.ResponseWriter, message interface{}) {
	app.writeResponse(w, http.StatusBadRequest, toJson{"error": message})
}

func (app *application) validationErrorResponse(w http.ResponseWriter, message interface{}) {
	app.writeResponse(w, http.StatusUnprocessableEntity, toJson{"errors": message})
}
