package data

import (
	"time"
	"unicode/utf8"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/validator"
)

type Movie struct {
	Title         string    `json:"title"`
	Synopsis      string    `json:"synopsis"`
	Profile_path  string    `json:"profile_path"`
	Backdrop_path string    `json:"backdrop_path"`
	Runtime       int32     `json:"runtime"`
	ReleaseDate   time.Time `json:"release_date"`
	Genres        []string  `json:"genres"`
	Status        string    `json:"status"`
}

// handle movie validation logic
func ValidateMovie(movie *Movie) {
	v := validator.NewValidator()

	// status
	v.CheckIsError(movie.Title != "", "title", "title field cannot be empty")
	v.CheckIsError(utf8.RuneCountInString(movie.Title) == 50, "title", "character too long for title field must be less than 50")

	//TODO: handle validation logic for images
	v.CheckIsError(len(movie.Genres) != 0, "genres", "genres field cannot be empty")
	v.CheckIsError(validator.CheckDuplicate[string](movie.Genres), "genres", "genres field cannot contain duplicate value")

	v.CheckIsError(movie.Status != "", "status", "status field cannot be empty")
	v.CheckIsError(v.CheckMovieStatus(movie.Status), "status", "status field can either be released or upcoming")
}
