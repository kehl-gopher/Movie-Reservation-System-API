package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/data"
)

// contain all application dependency

func (app *application) CreateMovieRoutes(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		fmt.Println(r.Header.Get("Content-Type"))
		app.badErrorResponse(w, toJson{"error": "Expected form data, got JSON or unsupported format"})
		return
	}
	var movie data.MovieType = &data.Movie{}
	err := movie.CreateMovie(r)

	if err != nil {
		fmt.Println(err)
	}
}
