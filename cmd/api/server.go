package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"greenlight.hosseinnasiri.ir/internal/jsonlog"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})

		os.Exit(0)
	}()

	app.logger.PrintInfo("starting server ...", jsonlog.LoggerProperties{
		"env":  app.config.env,
		"addr": srv.Addr,
	})

	return srv.ListenAndServe()
}
