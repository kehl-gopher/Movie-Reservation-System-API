package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/data"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type MovieDB struct {
	DB  *sql.DB
	Red *redis.Client
}

// create movie function for hanlding user creating movies
// in the database level while also performing the write-through cache
// invalidation logic...
func (m *MovieDB) CreateMovie(data *data.Movie) error {

	var movieID int
	// INSERT INTO DATABASE
	query := `
	INSERT INTO movies (title, synopsis, status, profile_path, background_path, genre_names, release_date)
	VALUES ($1, $2, $3, $4, $5,$6, $7) RETURNING id
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		data.Title,
		data.Synopsis,
		data.Status,
		data.Profile_path,
		data.Backdrop_path,
		pq.Array(data.Genres),
		data.ReleaseDate,
	}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movieID)

	if err != nil {
		return err
	}
	genres := strings.Join(data.Genres, ",")
	rArgs := map[string]interface{}{
		"id":            movieID,
		"title":         data.Title,
		"synopsis":      data.Synopsis,
		"status":        data.Status,
		"profile_path":  data.Profile_path,
		"backdrop_path": data.Backdrop_path,
		"genres":        genres,
		"release_date":  string(data.ReleaseDate),
	}
	fmt.Println(movieID)
	_, err = m.Red.LPush(ctx, "movie_list", rArgs).Result()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
