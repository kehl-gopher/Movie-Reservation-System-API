package logs

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/utils"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

var Once sync.Once

var log zerolog.Logger

func Get(logLevel zerolog.Level) zerolog.Logger {
	// initialize the logger once
	Once.Do(func() {
		// only to be use in development environment
		var outPut io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		if env, err := utils.ReadEnvVariable("APP_ENV"); err == nil {
			if env != "DEVELOPMENT" {
				// handle for production logging
				fileLogger := &lumberjack.Logger{
					Filename:   "app-log.log",
					MaxSize:    5,
					MaxAge:     10,
					MaxBackups: 10,
					Compress:   true,
				}
				outPut = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
			}
		} else {
			panic(err)
		}

		log = zerolog.New(outPut).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Logger()
	})
	return log
}

type AppLogger struct {
	logLevel zerolog.Level
	log      zerolog.Logger
}

func NewAppLogger(level zerolog.Level) *AppLogger {
	return &AppLogger{logLevel: level, log: Get(level)}
}

func (a *AppLogger) InfoLog(message string) {
	a.logLevel = zerolog.InfoLevel
	a.printMessage(message)
}

func (a *AppLogger) ErrorLog(err error) {
	a.logLevel = zerolog.ErrorLevel
	a.printMessage(err)
}

func (a *AppLogger) FatalLog(err error) {
	a.logLevel = zerolog.FatalLevel
	a.printMessage(err)
}

func (a *AppLogger) DebugLog(message string) {
	a.logLevel = zerolog.DebugLevel
	a.printMessage(message)
}

// log error message
func (a *AppLogger) printMessage(message interface{}) {
	l := a.log.With().Timestamp().Logger()
	switch a.logLevel {
	case zerolog.InfoLevel:
		l.Info().Msgf("%s", message)
	case zerolog.DebugLevel:
		l.Debug().Msgf("%s", message)
	case zerolog.ErrorLevel:
		l.Error().Err(message.(error)).Msgf("%s", message.(error).Error())
	default:
		l.Fatal().Stack().Msgf("%s", message)
	}
}
