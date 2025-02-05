package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/logs"
	"github.com/rs/zerolog"
)

type config struct {
	port int
	host string
}
type application struct {
	config
	log *logs.AppLogger
}

func main() {

	var conf config

	flag.IntVar(&conf.port, "port", 8000, "server port")
	flag.StringVar(&conf.host, "host", "127.0.0.1", "server connection host")

	logs := logs.NewAppLogger(zerolog.InfoLevel)
	app := application{
		config: conf,
		log:    logs,
	}
	serv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", app.host, app.port),
		Handler:      app.routers(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	logs.InfoLog(fmt.Sprintf("Server connection successful on port %s", serv.Addr))
	err := serv.ListenAndServe()
	logs.FatalLog(err)
}
