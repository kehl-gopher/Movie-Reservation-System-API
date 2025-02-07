package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/data"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type MovieDB struct {
	DB  *sql.DB
	red *redis.Client
}

func (m *MovieDB) CreateMovie(data *data.Movie) error {

	// INSERT INTO DATABASE
	query := `
	INSERT INTO movies (title, synopsis, status, profile_path background_path, genre_names, release_date)
	VALUES ($1, $2, $3, $4, $5,$6, $7)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{data.Title,
		data.Synopsis,
		data.Status,
		data.Profile_path,
		data.Backdrop_path,
		pq.Array(data.Genres),
		data.ReleaseDate,
	}
	_, err := m.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	// write to redis DB

	return nil
}
