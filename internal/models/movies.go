package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/data"
	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/logs"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	// insert into redis database
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
	_, err = m.Red.LPush(ctx, "movie_list", rArgs).Result()
	if err != nil {
		return err
	}

	return nil
}

// get individual movies by the provided id
// sorry my code is so messy i'lll refactor this godadamn code later
func (m *MovieDB) GetMovie(id int) (*data.Movie, error) {
	var movie data.Movie
	var mGenres struct {
		Genres string `redis:"genres"`
	}
	var date time.Time

	log := logs.Get(zerolog.InfoLevel)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// set redis  key
	m_key := fmt.Sprintf("movie_id_%d", id)

	// get genres value from redis first

	err := m.Red.HMGet(ctx, m_key, "genres").Scan(&mGenres)

	if err != nil {
		log.Err(err).Msg(err.Error())
		return nil, err
	}

	// get data from redis first...
	err = m.Red.HGetAll(ctx, m_key).Scan(&movie)
	if mGenres.Genres != "" && !reflect.DeepEqual(movie, data.Movie{}) {
		movie.Genres = strings.Split(mGenres.Genres, ",")
	}

	if err != nil {
		log.Err(err).Str("message", err.Error()).Send()
		return nil, err
	}

	// checks if there's a cache miss
	if reflect.DeepEqual(movie, data.Movie{}) {
		log.Err(err).Str("message", "Cache miss occur").Send()

		// get data from database and add to redis
		query := `SELECT id, title, synopsis, status,
				genre_names, profile_path, background_path, release_date
				FROM movies WHERE id = $1
		`
		err := m.DB.QueryRowContext(ctx, query, id).Scan(
			&movie.ID,
			&movie.Title,
			&movie.Synopsis,
			&movie.Status,
			pq.Array(&movie.Genres),
			&movie.Profile_path,
			&movie.Backdrop_path,
			&date,
		)
		movie.ReleaseDate = data.Dt(date.Format("2006-01-02"))
		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				log.Err(err).Msg("missing id")
				return nil, data.NotFoundError
			}
			log.Err(err).Msg("server error response")
			return nil, err
		}
		// add the data to redis...

		err = m.addDataToRedis(&movie, ctx, m_key)

		if err != nil {
			log.Err(err).Msg(err.Error())
			return nil, err
		}
	}

	// make movie genre arrays
	return &movie, nil

}

//add data to redis database

func (m *MovieDB) addDataToRedis(movie *data.Movie, ctx context.Context, key string) error {
	genres := strings.Join(movie.Genres, ",")
	rArgs := map[string]interface{}{
		"id":            movie.ID,
		"title":         movie.Title,
		"synopsis":      movie.Synopsis,
		"status":        movie.Status,
		"profile_path":  movie.Profile_path,
		"backdrop_path": movie.Backdrop_path,
		"genres":        genres,
		"release_date":  string(movie.ReleaseDate),
	}
	_, err := m.Red.HSet(ctx, key, rArgs).Result()
	if err != nil {
		log.Err(err).Str("message", err.Error()).Send()
		return err
	}
	return nil
}

// update user movie object...
func (m *MovieDB) UpdateMovieObj(movie *data.Movie) error {

	query := `
	UPDATE movies SET title = $2, synopsis = $3, status = $4,
	genre_names = $5, profile_path = $6, background_path = $7, release_date = $8
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		movie.ID,
		movie.Title,
		movie.Synopsis,
		movie.Status,
		movie.Profile_path,
		movie.Backdrop_path,
		pq.Array(movie.Genres),
		movie.ReleaseDate,
	}
	_, err := m.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	// pass updated data to redis

	err = m.updateDataInRedis(movie, ctx)
	return nil
}

// update data in redis so as to not cause data inconsistency
func (m *MovieDB) updateDataInRedis(movie *data.Movie, ctx context.Context) error {
	m_key := fmt.Sprintf("movie_id_%d", movie.ID)
	rArgs := []string{
		"id",
		"title",
		"synopsis",
		"status",
		"profile_path",
		"backdrop_path",
		"genres",
		"release_date",
	}
	err := m.Red.HDel(ctx, m_key, rArgs...).Err()

	if err != nil {
		return err
	}

	err = m.addDataToRedis(movie, ctx, m_key)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return nil
}
