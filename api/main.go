package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/logs"
	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/utils"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// to json a type for every app write
type toJson map[string]interface{}

type cors struct {
	originList []string
}
type config struct {
	port int
	host string
}
type application struct {
	config
	log            *logs.AppLogger
	trustedOrigins cors
}

const UploadDir string = "uploads/"

func main() {

	var conf config

	cors := cors{originList: []string{"localhost:3000",
		"localhost:3001",
		"localhost:8000",
	}} // specify your client origins here for production and development
	flag.IntVar(&conf.port, "port", 8000, "server port")
	flag.StringVar(&conf.host, "host", "127.0.0.1", "server connection host")

	// initialize app log
	logs := logs.NewAppLogger(zerolog.InfoLevel)

	// connect to database
	db, err := dbConnect()
	if err != nil {
		logs.FatalLog(err)
	}
	defer db.Close()

	// connect to redis client
	red, err := redConnection()
	if err != nil {
		logs.FatalLog(err)
	}
	defer red.Close()

	// application server initialze
	app := application{
		config:         conf,
		log:            logs,
		trustedOrigins: cors,
	}

	// handle file system to serve static file

	serv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", app.host, app.port),
		Handler:      app.routers(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}

	logs.InfoLog(fmt.Sprintf("Server connection successful on %s", serv.Addr))
	logs.InfoLog("Database Connected Successfully")
	logs.InfoLog("Redis Connected Successfully")
	err = serv.ListenAndServe()
	logs.FatalLog(err)
}

// database connection
func dbConnect() (*sql.DB, error) {
	env, err := utils.ReadEnvVariable("PG_SQL")
	if err != nil {
		return nil, fmt.Errorf("Database connection string missing in env")
	}

	db, err := sql.Open("postgres", env)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(13 * time.Minute)
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	defer db.Close()
	return db, nil
}

// redis connection
func redConnection() (*redis.Client, error) {
	// redis connection string
	addr, err := utils.ReadEnvVariable("RED_CONN")

	if err != nil {
		return nil, err
	}

	db, err := utils.ReadEnvVariable("DB")
	if err != nil {
		return nil, err
	}

	dbInt, err := strconv.Atoi(db)
	if err != nil {
		return nil, err
	}

	// password, err := utils.ReadEnvVariable("RED_PASSWORD")
	// if err != nil {
	// 	return nil, err
	// }
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       dbInt,
		Password: "",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ping, err := client.Ping(ctx).Result()

	if err != nil {
		return nil, err
	}
	fmt.Println(ping)
	return client, nil
}
