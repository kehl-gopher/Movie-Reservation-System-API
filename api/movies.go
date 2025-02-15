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
		app.badErrorResponse(w, toJson{"error": "Expected form data, got JSON or unsupported format"})
		return
	}
	var movie data.MovieType = &data.Movie{}

	movieS, err := movie.ReturnMovieObj(r)
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
	id, err := getparams(r)

	if err != nil {
		app.notFoundResponse(w, err)
		return
	}
	movie, err := app.model.Movies.GetMovie(id)

	if err != nil {
		if errors.Is(data.NotFoundError, err) {
			app.notFoundResponse(w, err)
		} else {
			app.serverErrorResponse(w, err)
		}
		return
	}
	app.writeResponse(w, http.StatusOK, toJson{"movie": movie})
}

// TODO: implement user movie update using partial update cause m damn sure too lazy for real
func (app *application) updateMovieById(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		app.badErrorResponse(w, toJson{"error": "Expected form data, got JSON or unsupported format"})
		return
	}

	// id, errr := getparams(r)
	// if errr != nil {
	// 	app.notFoundResponse(w, errr)
	// }

	// silence errror atleast

}

// TODO: implemente user delete movie update
func (app *application) deleteMovieByID(w http.ResponseWriter, r *http.Request) {

}

// TODO: implement get all movies and paginations on user movies
func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {

}
