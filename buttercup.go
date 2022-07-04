package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Buttercup struct {
	Service
}

type Service struct{}

func (*Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	w.Write([]byte("Not Implemented"))
}

func main() {

	var crt = flag.String("crt", "", "cert")
	var key = flag.String("key", "", "key")
	flag.Parse()

	var s http.Server = http.Server{
		Addr:    ":https",
		Handler: &Buttercup{},
	}

	done := func() <-chan struct{} {
		d := make(chan struct{})
		go func(c chan struct{}) {
			defer close(c)
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, os.Interrupt)
			<-signals
			if err := s.Shutdown(context.Background()); err != nil {
				log.Printf("HTTP server Shutdown: %v", err)
			}
		}(d)
		return d
	}

	if err := s.ListenAndServeTLS(*crt, *key); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
		return
	}
	<-done()
}
