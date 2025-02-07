package data

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/utils"
	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/validator"
)

type MovieType interface {
	CreateMovie(r *http.Request) error
}

type Movie struct {
	Title         string   `json:"title"`
	Synopsis      string   `json:"synopsis"`
	Profile_path  string   `json:"profile_path"`
	Backdrop_path string   `json:"backdrop_path"`
	Runtime       int32    `json:"runtime"`
	ReleaseDate   Dt       `json:"release_date"`
	Genres        []string `json:"genres"`
	Status        string   `json:"status"`
}

// handle movie validation logic
func (movie *Movie) ValidateMovie() {
	v := validator.NewValidator()

	v.CheckIsError(movie.Title != "", "title", "title field cannot be empty")
	v.CheckIsError(utf8.RuneCountInString(movie.Title) == 50, "title", "character too long for title field must be less than 50")

	//TODO: handle validation logic for images
	v.CheckIsError(len(movie.Genres) != 0, "genres", "genres field cannot be empty")
	v.CheckIsError(validator.CheckDuplicate[string](movie.Genres), "genres", "genres field cannot contain duplicate value")

	v.CheckIsError(movie.Status != "", "status", "status field cannot be empty")
	v.CheckIsError(v.CheckMovieStatus(movie.Status), "status", "status field can either be released or upcoming")
}

// create movie router...
func (m *Movie) CreateMovie(r *http.Request) error {
	movieStruct := &m
	maxSize := 100 << 10 // 10 mb file size expected
	if err := r.ParseMultipartForm(int64(maxSize)); err != nil {
		return err
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
			return err
		}
		formValues[val] = path
	}
	err := utils.FillStruct(movieStruct, formValues)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", *movieStruct)
	return nil
}
