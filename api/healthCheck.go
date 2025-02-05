package main

import "net/http"

// update this routes to perform it's duty like a bitch
func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Application healthcheck and versioning number incrementation"))
}
