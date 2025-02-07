package internal

import (
	"database/sql"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/models"
	"github.com/redis/go-redis/v9"
)

type AppModel struct {
	Movies *models.MovieDB
}

func InitAppModel(db *sql.DB, red *redis.Client) *AppModel {
	return &AppModel{Movies: &models.MovieDB{
		DB:  db,
		Red: red,
	}}
}
