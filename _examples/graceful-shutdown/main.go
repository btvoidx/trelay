package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/btvoidx/trelay"
)

func main() {
	// !
	// ! THIS EXAMPLE IS OUTDATED
	// !

	server := trelay.Server{
		Addr:    ":7777",
		Handler: trelay.Passthrough(":7878"),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("An error occured when starting the server: %s", err.Error())
		}
	}()

	defer func() {
		if err := server.Stop(); err != nil {
			log.Fatalf("An error occured when stopping the server: %s", err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
