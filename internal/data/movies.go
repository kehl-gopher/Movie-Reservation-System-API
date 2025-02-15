package data

import (
	"net/http"
	"unicode/utf8"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/utils"
	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/validator"
)

type MovieType interface {
	ReturnMovieObj(r *http.Request) (*Movie, *ErrorData)
}

type ErrorData struct {
	Err    error
	Verror map[string]interface{}
}

type Movie struct {
	ID            int      `json:"id" redis:"id"`
	Title         string   `json:"title" redis:"title"`
	Synopsis      string   `json:"synopsis" redis:"synopsis"`
	Profile_path  string   `json:"profile_path" redis:"profile_path"`
	Backdrop_path string   `json:"backdrop_path" redis:"backdrop_path"`
	Runtime       int32    `json:"runtime" redis:"runtime"`
	ReleaseDate   Dt       `json:"release_date" redis:"release_date"`
	Genres        []string `json:"genres" redis:"genres"`
	Status        string   `json:"status" redis:"status"`
}

// handle movie validation logic for user creating movies
func (movie *Movie) ValidateMovie(v *validator.ValidateData) {

	v.CheckIsError(movie.Title != "", "title", "title field cannot be empty")
	v.CheckIsError(utf8.RuneCountInString(movie.Title) == 50, "title", "character too long for title field must be less than 50")

	//TODO: handle validation logic for images
	v.CheckIsError(len(movie.Genres) != 0, "genres", "genres field cannot be empty")
	v.CheckIsError(validator.CheckDuplicate[string](movie.Genres), "genres", "genres field cannot contain duplicate value")

	v.CheckIsError(movie.Status != "", "status", "status field cannot be empty")
	v.CheckIsError(v.CheckMovieStatus(movie.Status), "status", "status field can either be released or upcoming")
}

// Handle logic for create movie... routes in the data level at least I hope
func (m *Movie) ReturnMovieObj(r *http.Request) (*Movie, *ErrorData) {
	movieStruct := m

	e := &ErrorData{Verror: make(map[string]interface{})}
	maxSize := 100 << 10 // 10 mb file size expected
	if err := r.ParseMultipartForm(int64(maxSize)); err != nil {
		e.Err = err
		return nil, e
	}
	formValues := make(map[string]interface{})
	for key, val := range r.Form {

		if key != "genres" {
			formValues[key] = r.FormValue(key)
		} else {
			formValues[key] = val
		}

	}

	for _, val := range []string{"backdrop_path", "profile_path"} {
		uploadPath := "uploads/movies"
		path, err := HandleImageFile(r, val, uploadPath)

		if err != nil {
			e.Err = err
			return nil, e
		}
		formValues[val] = path
	}
	err := utils.FillStruct(movieStruct, formValues)
	if err != nil {
		e.Err = err
		return nil, e
	}

	// validate movie value passed
	//
	v := validator.NewValidator()

	if movieStruct.ValidateMovie(v); v.CheckErrorExists() {
		e.Verror = v.Errors
		return nil, e
	}
	return movieStruct, nil
}

// handle user movie update
// TODO: fix error issue found here
func (m *Movie) UpdateMovieObj(r *http.Request, id int) (*Movie, *ErrorData) {
	var mov Movie

	movie, err := m.ReturnMovieObj(r)
	if err != nil {
		return nil, err
	}

	mov.Title = movie.Title
	mov.Synopsis = movie.Synopsis
	mov.Genres = movie.Genres
	mov.Runtime = movie.Runtime
	mov.Status = movie.Status
	mov.Profile_path = movie.Profile_path
	mov.Backdrop_path = movie.Backdrop_path

	return &mov, nil
}
