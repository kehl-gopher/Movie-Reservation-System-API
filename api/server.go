package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) Server() error {
	serv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", app.host, app.port),
		Handler:      app.routers(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}

	shutDownError := make(chan error)
	go func() {
		sigQuit := make(chan os.Signal, 1)

		signal.Notify(sigQuit, syscall.SIGTERM, syscall.SIGINT)
		s := <-sigQuit

		app.log.InfoLog(fmt.Sprintf("caught signal %s", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := serv.Shutdown(ctx)
		if err != nil {
			shutDownError <- err
		}
		app.log.InfoLog(fmt.Sprintf("Completing Background tasks %s", serv.Addr))
		app.wg.Wait()
		shutDownError <- nil
	}()
	app.log.InfoLog(fmt.Sprintf("Server connection successful on %s", serv.Addr))
	app.log.InfoLog("Database Connected Successfully")
	app.log.InfoLog("Redis Connected Successfully")

	err := serv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutDownError

	if err != nil {
		return err
	}

	app.log.InfoLog(fmt.Sprintf("Server shutting down %s...", serv.Addr))
	return nil
}
