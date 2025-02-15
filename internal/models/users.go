package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/data"
	"github.com/redis/go-redis/v9"
)

type UserDB struct {
	Db  *sql.DB
	Red *redis.Client
}

// create user from db
func (m *UserDB) CreateUser(user *data.Users) error {

	var uId int
	query := `INSERT INTO users (username, email, password, is_activated, is_admin)
	VALUES ($1, $2, $3, $4, $5) RETURNING ID
	`

	args := []interface{}{user.UserName, user.Email, user.Password.HashedPassword, user.IsActivated, user.IsAdmin}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.Db.QueryRowContext(ctx, query, args...).Scan(&uId)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicationValue
		default:
			return err
		}
	}

	// TODO insert user data to redis btw
	return nil
}

func (m *UserDB) GetByUserName(username string) (*data.Users, error) {

	var user data.Users
	// collect data from redis btw...

	query := `
	SELECT username, email, is_activated, is_admin
	FROM users
	WHERE username = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.Db.QueryRowContext(ctx, query, username).Scan(&user.Email, &user.IsActivated, &user.IsAdmin)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.NotFoundError
		default:
			return nil, err
		}
	}
	return &user, nil
}
