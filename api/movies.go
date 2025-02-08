package main

import (
	"errors"
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

	movieS, err := movie.CreateMovie(r)
	if err != nil {
		fmt.Println(err)
		switch {
		case errors.Is(&data.BadRequestError{}, err.Err):
			app.badErrorResponse(w, err)
		case errors.Is(&data.ValidationError{}, err.Err):
			app.validationErrorResponse(w, err)
			return
		case err.Verror != nil:
			app.validationErrorResponse(w, err.Verror)
		default:
			app.serverErrorResponse(w, err.Err)
		}
		return
	}

	errr := app.model.Movies.CreateMovie(movieS)

	if errr != nil {
		app.serverErrorResponse(w, errr)
		fmt.Println(errr)
		return
	}
	app.writeResponse(w, http.StatusCreated, toJson{"message": "Movie created successfully", "status": 201})
}

func (app *application) GetMovieById(w http.ResponseWriter, r *http.Request) {
	id, err := getparams(r.URL.Path)

	if err != nil {
		app.notFoundResponse(w, err)
		return
	}
	movie, err := app.model.Movies.GetMovie(id)

	if err != nil {
		if errors.Is(data.NotFoundError, err) {
			app.notFoundResponse(w, err)
			return
		}
		app.serverErrorResponse(w, err)
	}
	app.writeResponse(w, http.StatusOK, toJson{"movie": movie})
}
