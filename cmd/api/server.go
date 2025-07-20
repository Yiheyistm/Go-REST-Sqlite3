package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (app *application) server() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting server on %s", server.Addr)
	return server.ListenAndServe()
}
