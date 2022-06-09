package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Buttercup struct {
	Service
}

type Service struct{}

func (*Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	w.Write([]byte("Not Implemented"))
}

func shutdown() (err error) {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(c)
}

func main() {
	var status int

	flag.Parse()

	server := http.Server{
		Addr:              ":8080",
		Handler:           new(Buttercup),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
	}

	go func(errchan chan error) {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errchan <- err
		}
	}(errchan)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	select {
	case v := <-err:
		status = 1
	case <-signals:

	}

	shutdown()
	os.Exit(status)
}
